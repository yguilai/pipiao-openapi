package syncs

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/yguilai/pipiao-openapi/app/model"
	"github.com/yguilai/pipiao-openapi/common/collect"
	"github.com/yguilai/pipiao-openapi/common/xerr"
	"github.com/zeromicro/go-zero/core/logx"
	"github.com/zeromicro/go-zero/core/stores/redis"
)

const (
	// 定时任务每天跑一次, 缓存过期时间比24h长就行, 这里就先设为48h
	warframeDictSyncSHARedisExpire = 86400 * 2
	// warframeEntryUpdateTaskKey 词典更新redis唯一key, 用于全局分布式锁, 避免同时跑太多任务
	warframeEntryUpdateTaskKey = "wf:sync_unique_key:entry_update"
	// warframeEntryUpdateTaskExpire key过期时间, 10分钟
	warframeEntryUpdateTaskExpire = 60 * 10
)

type WfEntrySyncService struct {
	redis *redis.Redis
	model.WfEntryModel
}

func NewWfEntrySyncService(r *redis.Redis, m model.WfEntryModel) *WfEntrySyncService {
	return &WfEntrySyncService{
		redis:        r,
		WfEntryModel: m,
	}
}

// NeedFetch 判断是否需要拉取wf词条
func (s *WfEntrySyncService) NeedFetch(ctx context.Context) (string, string, bool) {
	return NeedFetch(ctx, s.redis, WfAllEntry)
}

func (s *WfEntrySyncService) StartUpdate(ctx context.Context, downloadUrl, sha string) error {
	lock := redis.NewRedisLock(s.redis, warframeEntryUpdateTaskKey)
	lock.SetExpire(warframeEntryUpdateTaskExpire)
	ok, err := lock.AcquireCtx(ctx)
	if err != nil {
		return err
	}
	if !ok {
		return xerr.NewErrorWithMsg("更新任务正在进行中, 请稍后再试~")
	}
	defer lock.Release()
	entries := make(collect.Slice[WfEntry], 0)
	err = FetchData(downloadUrl, &entries)
	if err != nil {
		return xerr.NewErrorWithFormat("获取数据失败: %+v", err)
	}
	if len(entries) == 0 {
		logx.WithContext(ctx).Infof("没有需要更新的数据\n")
		return nil
	}

	syncPool := NewSyncPool[collect.Slice[WfEntry], WfEntry](1)
	err = syncPool.SyncAll(ctx, entries, s.doEntryUpdate)
	if err != nil {
		return err
	}

	_ = UpdateLastSHA(ctx, s.redis, WfAllEntry, sha, warframeDictSyncSHARedisExpire)
	logx.Errorf("更新上次Warframe字典文件SHA失败: %+v", err)
	return nil
}

func (s *WfEntrySyncService) FindOld(ctx context.Context, newEntry interface{}) (interface{}, error) {
	n := newEntry.(*model.WfEntry)
	return s.FindOneByUniqueName(ctx, n.UniqueName)
}

func (s *WfEntrySyncService) Add(ctx context.Context, e interface{}) (sql.Result, error) {
	entry := e.(*model.WfEntry)
	return s.Insert(ctx, entry)
}

func (s *WfEntrySyncService) Modify(ctx context.Context, oldEntry, newEntry interface{}) error {
	n := newEntry.(*model.WfEntry)
	o := oldEntry.(*model.WfEntry)
	n.Id = o.Id
	return s.Update(ctx, n)
}

func (s *WfEntrySyncService) NeedModify(oldEntry, newEntry interface{}) bool {
	if oldEntry == nil {
		return true
	}
	o := oldEntry.(*model.WfEntry)
	n := newEntry.(*model.WfEntry)
	newKey := fmt.Sprintf("%s%s%s%d", n.UniqueName, n.Name, n.Category, n.Tradable)
	oldKey := fmt.Sprintf("%s%s%s%d", o.UniqueName, o.Name, o.Category, o.Tradable)
	return NeedUpdate(oldKey, newKey)
}

func (s *WfEntrySyncService) doEntryUpdate(ctx context.Context, e *WfEntry, errch chan<- error) {
	if e.UniqueName == "" {
		return
	}
	newEntry := &model.WfEntry{
		UniqueName: e.UniqueName,
		Category:   e.Category,
		Name:       e.Name,
		Tradable:   transferTradable(e.Tradable),
	}
	DoUpdate(ctx, s, newEntry, errch)
}

func transferTradable(tradable bool) int64 {
	if tradable {
		return 1
	}
	return 0
}
