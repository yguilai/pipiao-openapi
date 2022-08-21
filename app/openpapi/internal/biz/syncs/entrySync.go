package syncs

import (
	"context"
	"fmt"
	"github.com/yguilai/pipiao-openapi/app/model"
	"github.com/yguilai/pipiao-openapi/common/xerr"
	"github.com/zeromicro/go-zero/core/logx"
	"github.com/zeromicro/go-zero/core/stores/redis"
)

const (
	// 定时任务每天跑一次, 缓存过期时间比24h长就行, 这里就先设为48h
	warframeSyncSHARedisExpire = 86400 * 2
	// warframeEntryUpdateTaskKey 词典更新redis唯一key, 用于全局分布式锁, 避免同时跑太多任务
	warframeEntryUpdateTaskKey = "wf:sync_unique_key:entry_update"
	// warframeEntryUpdateTaskExpire key过期时间, 10分钟
	warframeEntryUpdateTaskExpire = 60 * 10
)

type wfEntrySyncService struct {
	redis *redis.Redis
	model.WfEntryModel
}

func NewWfEntrySyncService(r *redis.Redis, m model.WfEntryModel) SyncService {
	return &wfEntrySyncService{
		redis:        r,
		WfEntryModel: m,
	}
}

// NeedFetch 判断是否需要拉取wf词条
func (s *wfEntrySyncService) NeedFetch(ctx context.Context) (string, string, bool) {
	return NeedFetch(ctx, s.redis, WfAllEntry)
}

func (s *wfEntrySyncService) StartUpdate(ctx context.Context, downloadUrl, sha string) error {
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
	entries := make([]WfEntry, 0)
	err = FetchData(downloadUrl, &entries)
	if err != nil {
		return xerr.NewErrorWithFormat("获取数据失败: %+v", err)
	}
	if len(entries) == 0 {
		logx.WithContext(ctx).Infof("没有需要更新的数据\n")
		return nil
	}

	syncPool := NewSyncPool[WfEntry](1)
	err = syncPool.SyncAll(ctx, entries, s.doEntryUpdate)
	if err != nil {
		return err
	}

	_ = UpdateLastSHA(ctx, s.redis, WfAllEntry, sha, warframeSyncSHARedisExpire)
	logx.Errorf("更新上次Warframe字典文件SHA失败: %+v", err)
	return nil
}

func (s *wfEntrySyncService) FindOld(ctx context.Context, newEntry *model.WfEntry) (*model.WfEntry, error) {
	return s.FindOneByUniqueName(ctx, newEntry.UniqueName)
}

func (s *wfEntrySyncService) Modify(ctx context.Context, oldEntry, newEntry *model.WfEntry) error {
	newEntry.Id = oldEntry.Id
	return s.Update(ctx, newEntry)
}

func (s *wfEntrySyncService) NeedModify(o, n *model.WfEntry) bool {
	if o == nil || n == nil {
		return true
	}
	newKey := fmt.Sprintf("%s%s%s%d", n.UniqueName, n.Name, n.Category, n.Tradable)
	oldKey := fmt.Sprintf("%s%s%s%d", o.UniqueName, o.Name, o.Category, o.Tradable)
	return NeedUpdate(oldKey, newKey)
}

func (s *wfEntrySyncService) newCommonSync() *CommonSync[model.WfEntry] {
	return &CommonSync[model.WfEntry]{
		FindOld:    s.FindOld,
		Insert:     s.Insert,
		Modify:     s.Modify,
		NeedModify: s.NeedModify,
	}
}

func (s *wfEntrySyncService) doEntryUpdate(ctx context.Context, e *WfEntry, errch chan<- error) {
	if e.UniqueName == "" {
		return
	}
	newEntry := &model.WfEntry{
		UniqueName: e.UniqueName,
		Category:   e.Category,
		Name:       e.Name,
		Tradable:   transferTradable(e.Tradable),
	}
	DoUpdate[model.WfEntry](ctx, s.newCommonSync(), newEntry, errch)
}

func transferTradable(tradable bool) int64 {
	if tradable {
		return 1
	}
	return 0
}
