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

type WfI18nSyncTask struct {
	logx.Logger
	ctx        context.Context
	svcCtx     *svc.ServiceContext
	task       *cron.Cron
	syncHelper syncs.SyncService
}

func NewWfI18nSyncTask(ctx context.Context, svcCtx *svc.ServiceContext) *WfI18nSyncTask {
	lg := logx.WithContext(ctx)
	return &WfI18nSyncTask{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: lg,
		task: cron.New(
			cron.WithParser(cron.NewParser(cron.SecondOptional|cron.Minute|cron.Hour|cron.Dom|cron.Month|cron.Dow|cron.Descriptor)),
			cron.WithLogger(xcron.NewLogger(lg, svcCtx.IsDev())),
		),
		syncHelper: syncs.NewWfI18nItemService(svcCtx.Redis, svcCtx.WfI18nItemModel),
	}
}

func (t *WfI18nSyncTask) Start() {
	// 每天凌晨执行一次拉取
	spec := "0 0 0 */1 * ?"
	entryId, err := t.task.AddFunc(spec, t.fetchWfDictFromGithub)
	if err != nil {
		panic(err)
	}
	t.task.Start()
	t.Infof("warframe dict sync task started, entryId: %d", entryId)
}

func (t *WfI18nSyncTask) Stop() {
	t.Infof("warframe dict sync task stop\n")
	t.task.Stop()
}

func (t *WfI18nSyncTask) fetchWfDictFromGithub() {
	ctx := context.WithValue(t.ctx, "trace", id.SimpleUUID())
	downURL, sha, isNeed := t.syncHelper.NeedFetch(ctx)
	if !isNeed || downURL == "" {
		t.Infof("本次任务无需更新I18n")
		return
	}
	err := t.syncHelper.StartUpdate(ctx, downURL, sha)
	if err != nil {
		t.Errorf("i18n词条更新失败: %+v", err)
		return
	}
}
