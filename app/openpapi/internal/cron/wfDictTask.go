package cron

import (
    "context"
    "github.com/robfig/cron/v3"
    "github.com/yguilai/pipiao-openapi/app/openpapi/internal/biz"
    "github.com/yguilai/pipiao-openapi/app/openpapi/internal/svc"
    "github.com/yguilai/pipiao-openapi/common/xcron"
    "github.com/zeromicro/go-zero/core/logx"
)

type WfDictSyncTask struct {
    logx.Logger
    ctx        context.Context
    svcCtx     *svc.ServiceContext
    task       *cron.Cron
    dictHelper *biz.WfDictService
}

func NewWfDictSyncTask(ctx context.Context, svcCtx *svc.ServiceContext) *WfDictSyncTask {
    lg := logx.WithContext(ctx)
    return &WfDictSyncTask{
        ctx:    ctx,
        svcCtx: svcCtx,
        Logger: lg,
        task: cron.New(
            cron.WithParser(cron.NewParser(cron.SecondOptional|cron.Minute|cron.Hour|cron.Dom|cron.Month|cron.Dow|cron.Descriptor)),
            cron.WithLogger(xcron.NewLogger(lg, svcCtx.IsDev())),
        ),
        dictHelper: biz.NewWfDictService(svcCtx.Redis),
    }
}

func (t *WfDictSyncTask) Start() {
    // 每天凌晨执行一次拉取
    spec := "0 0 0 */1 * ?"
    entryId, err := t.task.AddFunc(spec, t.fetchWfDictFromGithub)
    if err != nil {
        panic(err)
    }
    t.task.Start()
    t.Infof("warframe dict sync task started, entryId: %d", entryId)
}

func (t *WfDictSyncTask) Stop() {
    t.Infof("warframe dict sync task stop\n")
    t.task.Stop()
}

func (t *WfDictSyncTask) fetchWfDictFromGithub() {
    downURL, isNeed := t.dictHelper.NeedFetch(t.ctx)
    if !isNeed || downURL == "" {
        t.Infof("本次任务无需更新词典")
        return
    }
}
