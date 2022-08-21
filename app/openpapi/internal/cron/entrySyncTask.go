package cron

import (
	"context"
	"github.com/robfig/cron/v3"
	"github.com/yguilai/pipiao-openapi/app/openpapi/internal/biz/syncs"
	"github.com/yguilai/pipiao-openapi/app/openpapi/internal/svc"
	"github.com/yguilai/pipiao-openapi/common/id"
	"github.com/yguilai/pipiao-openapi/common/xcron"
	"github.com/zeromicro/go-zero/core/logx"
)

type WfEntrySyncTask struct {
	logx.Logger
	ctx        context.Context
	svcCtx     *svc.ServiceContext
	task       *cron.Cron
	syncHelper syncs.SyncService
}

func NewWfEntrySyncTask(ctx context.Context, svcCtx *svc.ServiceContext) *WfEntrySyncTask {
	lg := logx.WithContext(ctx)
	return &WfEntrySyncTask{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: lg,
		task: cron.New(
			cron.WithParser(cron.NewParser(cron.SecondOptional|cron.Minute|cron.Hour|cron.Dom|cron.Month|cron.Dow|cron.Descriptor)),
			cron.WithLogger(xcron.NewLogger(lg, svcCtx.IsDev())),
		),
		syncHelper: syncs.NewWfEntrySyncService(svcCtx.Redis, svcCtx.WfEntryModel),
	}
}

func (t *WfEntrySyncTask) Start() {
	// 每天凌晨执行一次拉取
	spec := "0 0 0 */1 * ?"
	entryId, err := t.task.AddFunc(spec, t.fetchWfDictFromGithub)
	if err != nil {
		panic(err)
	}
	t.task.Start()
	t.Infof("warframe dict sync task started, entryId: %d", entryId)
}

func (t *WfEntrySyncTask) Stop() {
	t.Infof("warframe dict sync task stop\n")
	t.task.Stop()
}

func (t *WfEntrySyncTask) fetchWfDictFromGithub() {
	ctx := context.WithValue(t.ctx, "trace", id.SimpleUUID())
	downURL, sha, isNeed := t.syncHelper.NeedFetch(ctx)
	if !isNeed || downURL == "" {
		t.Infof("本次任务无需更新词典")
		return
	}
	err := t.syncHelper.StartUpdate(ctx, downURL, sha)
	if err != nil {
		t.Errorf("词条更新失败: %+v", err)
		return
	}
}
