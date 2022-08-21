package handler

import (
	"context"
	"github.com/yguilai/pipiao-openapi/app/openpapi/internal/cron"
	"github.com/yguilai/pipiao-openapi/app/openpapi/internal/svc"
	"github.com/zeromicro/go-zero/core/service"
)

func RegisterJobs(group *service.ServiceGroup, svcCtx *svc.ServiceContext) {
	group.Add(cron.NewWfEntrySyncTask(context.Background(), svcCtx))
}
