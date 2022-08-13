package handler

import (
    "context"
    "github.com/yguilai/pipiao-openapi/app/openpapi/internal/logic/cron"
    "github.com/yguilai/pipiao-openapi/app/openpapi/internal/svc"
    "github.com/zeromicro/go-zero/core/service"
)

func RegisterJobs(group *service.ServiceGroup, svcCtx *svc.ServiceContext) {
    group.Add(cron.NewWfDictSyncTask(context.Background(), svcCtx))
}
