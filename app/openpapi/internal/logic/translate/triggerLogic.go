package translate

import (
	"context"
	"github.com/yguilai/pipiao-openapi/app/openpapi/internal/biz/syncs"

	"github.com/yguilai/pipiao-openapi/app/openpapi/internal/svc"
	"github.com/yguilai/pipiao-openapi/app/openpapi/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type TriggerLogic struct {
	logx.Logger
	ctx         context.Context
	svcCtx      *svc.ServiceContext
	dictService *syncs.WfEntrySyncService
}

func NewTriggerLogic(ctx context.Context, svcCtx *svc.ServiceContext) *TriggerLogic {
	return &TriggerLogic{
		Logger:      logx.WithContext(ctx),
		ctx:         ctx,
		svcCtx:      svcCtx,
		dictService: syncs.NewWfEntrySyncService(svcCtx.Redis, svcCtx.WfEntryModel),
	}
}

func (l *TriggerLogic) Trigger() (resp *types.TranslateResp, err error) {
	url, sha, need := l.dictService.NeedFetch(l.ctx)
	if !need {
		return nil, nil
	}
	err = l.dictService.StartUpdate(l.ctx, url, sha)
	if err != nil {
		return nil, err
	}
	return &types.TranslateResp{Result: "1"}, nil
}
