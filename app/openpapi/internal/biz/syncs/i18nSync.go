package syncs

import (
	"context"
	"fmt"
	"github.com/yguilai/pipiao-openapi/app/model"
	"github.com/yguilai/pipiao-openapi/common/maps"
	"github.com/yguilai/pipiao-openapi/common/xerr"
	"github.com/zeromicro/go-zero/core/logx"
	"github.com/zeromicro/go-zero/core/stores/redis"
)

const (
	warframeI18nUpdateTaskKey    = "wf:sync_unique_key:i18n_update"
	warframeI18nUpdateTaskExpire = 60 * 10
)

var needSaveLang = []string{"zh"}

type wfI18nItemService struct {
	redis *redis.Redis
	model.WfI18nItemModel
}

func NewWfI18nItemService(r *redis.Redis, m model.WfI18nItemModel) SyncService {
	return &wfI18nItemService{
		redis:           r,
		WfI18nItemModel: m,
	}
}

func (s *wfI18nItemService) NeedFetch(ctx context.Context) (string, string, bool) {
	return NeedFetch(ctx, s.redis, WfI18n)
}

func (s *wfI18nItemService) StartUpdate(ctx context.Context, downloadUrl, sha string) error {
	lock := redis.NewRedisLock(s.redis, warframeI18nUpdateTaskKey)
	lock.SetExpire(warframeI18nUpdateTaskExpire)
	ok, err := lock.AcquireCtx(ctx)
	if err != nil {
		return err
	}
	if !ok {
		return xerr.NewErrorWithMsg("更新任务正在进行中, 请稍后再试~")
	}
	defer lock.Release()
	i18nMap := make(map[string]I18nEntry, 0)
	err = FetchData(downloadUrl, &i18nMap)
	if err != nil {
		return xerr.NewErrorWithFormat("获取数据失败: %+v", err)
	}
	if len(i18nMap) == 0 {
		logx.WithContext(ctx).Infof("没有需要更新的数据\n")
		return nil
	}

	entries := maps.Entries(i18nMap)
	syncPool := NewSyncPool[maps.Entry[string, I18nEntry]](1)
	err = syncPool.SyncAll(ctx, entries, s.doI18nUpdate)
	if err != nil {
		return err
	}

	_ = UpdateLastSHA(ctx, s.redis, WfI18n, sha, warframeSyncSHARedisExpire)
	logx.Errorf("更新上次Warframe字典文件SHA失败: %+v", err)
	return nil
}

func (s *wfI18nItemService) FindOld(ctx context.Context, n *model.WfI18nItem) (*model.WfI18nItem, error) {
	return s.FindOneByUniqueNameLang(ctx, n.UniqueName, n.Lang)
}

func (s *wfI18nItemService) Modify(ctx context.Context, o, n *model.WfI18nItem) error {
	n.Id = o.Id
	return s.Update(ctx, n)
}

func (s *wfI18nItemService) NeedModify(o, n *model.WfI18nItem) bool {
	newKey := fmt.Sprintf("%s%s%s", n.UniqueName, n.Name, n.Lang)
	oldKey := fmt.Sprintf("%s%s%s", o.UniqueName, o.Name, o.Lang)
	return NeedUpdate(oldKey, newKey)
}

func (s *wfI18nItemService) newCommonSync() *CommonSync[model.WfI18nItem] {
	return &CommonSync[model.WfI18nItem]{
		FindOld:    s.FindOld,
		Insert:     s.Insert,
		Modify:     s.Modify,
		NeedModify: s.NeedModify,
	}
}

func (s *wfI18nItemService) doI18nUpdate(ctx context.Context, e *maps.Entry[string, I18nEntry], errch chan<- error) {
	for _, lang := range needSaveLang {
		item := e.Value[lang]
		var name string
		if item.SystemName != "" {
			name = item.SystemName
		}
		i18n := &model.WfI18nItem{
			UniqueName: e.Key,
			Lang:       lang,
			Name:       name,
		}
		DoUpdate(ctx, s.newCommonSync(), i18n, errch)
	}
}
