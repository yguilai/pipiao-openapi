package trigger

import (
	"context"
	"github.com/yguilai/pipiao-openapi/app/openpapi/internal/biz/syncs"

	"github.com/yguilai/pipiao-openapi/app/openpapi/internal/svc"
	"github.com/yguilai/pipiao-openapi/app/openpapi/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type TriggerI18nLogic struct {
	logx.Logger
	ctx         context.Context
	svcCtx      *svc.ServiceContext
	i18nService *syncs.WfI18nItemService
}

func NewTriggerI18nLogic(ctx context.Context, svcCtx *svc.ServiceContext) *TriggerI18nLogic {
	return &TriggerI18nLogic{
		Logger:      logx.WithContext(ctx),
		ctx:         ctx,
		svcCtx:      svcCtx,
		i18nService: syncs.NewWfI18nItemService(svcCtx.Redis, svcCtx.WfI18nItemModel),
	}
}

func (l *TriggerI18nLogic) TriggerI18n() (resp *types.TranslateResp, err error) {
	url, sha, need := l.i18nService.NeedFetch(l.ctx)
	if !need {
		return nil, nil
	}
	err = l.i18nService.StartUpdate(l.ctx, url, sha)
	if err != nil {
		return nil, err
	}
	return &types.TranslateResp{Result: "1"}, nil
}
