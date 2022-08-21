package trigger

import (
	"context"
	"github.com/yguilai/pipiao-openapi/app/openpapi/internal/biz/syncs"

	"github.com/yguilai/pipiao-openapi/app/openpapi/internal/svc"
	"github.com/yguilai/pipiao-openapi/app/openpapi/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type TriggerEntryLogic struct {
	logx.Logger
	ctx          context.Context
	svcCtx       *svc.ServiceContext
	entryService syncs.SyncService
}

func NewTriggerEntryLogic(ctx context.Context, svcCtx *svc.ServiceContext) *TriggerEntryLogic {
	return &TriggerEntryLogic{
		Logger:       logx.WithContext(ctx),
		ctx:          ctx,
		svcCtx:       svcCtx,
		entryService: syncs.NewWfEntrySyncService(svcCtx.Redis, svcCtx.WfEntryModel),
	}
}

func (l *TriggerEntryLogic) TriggerEntry() (resp *types.TranslateResp, err error) {
	url, sha, need := l.entryService.NeedFetch(l.ctx)
	if !need {
		return nil, nil
	}
	err = l.entryService.StartUpdate(l.ctx, url, sha)
	if err != nil {
		return nil, err
	}
	return &types.TranslateResp{Result: "1"}, nil
}
