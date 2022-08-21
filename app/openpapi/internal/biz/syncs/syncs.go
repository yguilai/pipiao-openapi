package syncs

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/yguilai/pipiao-openapi/app/model"
	"github.com/yguilai/pipiao-openapi/common/xerr"
	"github.com/yguilai/pipiao-openapi/common/xgithub"
	"github.com/zeromicro/go-zero/core/hash"
	"github.com/zeromicro/go-zero/core/logx"
	"github.com/zeromicro/go-zero/core/stores/redis"
	"io"
	"net/http"
)

const (
	WarframeDictOwner       = "WFCD"
	WarframeDictRepo        = "warframe-items"
	warframeSyncSHARedisKey = "wf:sync:last_sha:"
)

// SyncType 同步类型
type SyncType int

const (
	WfAllEntry SyncType = iota + 1
	WfI18n
)

func getSyncPathByType(t SyncType) string {
	switch t {
	case WfAllEntry:
		return "/data/json/All.json"
	case WfI18n:
		return "/data/json/i18n.json"
	}
	panic("unsupported sync type")
}

// NeedFetch 通用判断文件是否已更新
func NeedFetch(ctx context.Context, store *redis.Redis, syncType SyncType) (string, string, bool) {
	client := xgithub.NewSimpleClient()
	ctt, _, _, err := client.Repositories.GetContents(ctx, WarframeDictOwner, WarframeDictRepo, getSyncPathByType(syncType), nil)
	if err != nil {
		logx.Errorf("获取Wf词条仓库信息出错: %+v", err)
		return "", "", false
	}
	// TODO: 这里其实存在一致性问题, 可以改用lua脚本来校验, 保证原子性, 目前只是单机部署, 问题不大.
	//       可惜go这个redis操作库不能像java里一样支持针对某个key开启事务 不然就简单了
	// 获取上次拉取的
	lastSHA, err := store.Get(genSyncTypeKey(syncType))
	if err != nil {
		logx.Errorf("获取上次拉取SHA信息出错: %+v", err)
		return "", "", false
	}

	// lastSHA不为空说明不是第一次拉取, lastSHA == *ctt.SHA说明词条还未更新
	if lastSHA != "" && lastSHA == *ctt.SHA {
		return "", "", false
	}
	return *ctt.DownloadURL, *ctt.SHA, true
}

func UpdateLastSHA(ctx context.Context, store *redis.Redis, syncType SyncType, sha string, expire int) error {
	return store.SetexCtx(ctx, genSyncTypeKey(syncType), sha, expire)
}

func genSyncTypeKey(syncType SyncType) string {
	return fmt.Sprintf("%s%d", warframeSyncSHARedisKey, syncType)
}

func FetchData(downloadUrl string, v interface{}) error {
	client := &http.Client{Timeout: 0}
	resp, err := client.Get(downloadUrl)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	bytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return err
	}
	err = json.Unmarshal(bytes, v)
	if err != nil {
		return err
	}
	return nil
}

func NeedUpdate(oldKey string, newKey string) bool {
	return hash.Md5Hex([]byte(oldKey)) != hash.Md5Hex([]byte(newKey))
}

func DoUpdate[T any](ctx context.Context, m *CommonSync[T], newEntry *T, errch chan<- error) {
	oldEntry, err := m.FindOld(ctx, newEntry)
	if err != nil && err != model.ErrNotFound {
		errch <- err
		return
	}
	// 说明是新增的
	if err == model.ErrNotFound {
		res, err := m.Insert(ctx, newEntry)
		if err != nil {
			errch <- err
			return
		}
		if rows, _ := res.RowsAffected(); rows == 0 {
			errch <- xerr.NewErrorWithMsg("插入失败")
			return
		}
		return
	}

	// 不为空则是已存在, 判断一下是否需要更新
	if !m.NeedModify(oldEntry, newEntry) {
		logx.WithContext(ctx).Infof("entry: %+v, don't need to update", newEntry)
		return
	}
	err = m.Modify(ctx, oldEntry, newEntry)
	if err != nil {
		errch <- err
		return
	}
}
