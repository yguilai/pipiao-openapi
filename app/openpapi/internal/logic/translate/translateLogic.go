package translate

import (
	"context"
	"fmt"
	"github.com/yguilai/pipiao-openapi/app/model"
	"github.com/yguilai/pipiao-openapi/app/openpapi/internal/svc"
	"github.com/yguilai/pipiao-openapi/app/openpapi/internal/types"
	"github.com/yguilai/pipiao-openapi/common/xerr"

	"github.com/zeromicro/go-zero/core/logx"
)

type TranslateLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewTranslateLogic(ctx context.Context, svcCtx *svc.ServiceContext) *TranslateLogic {
	return &TranslateLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *TranslateLogic) Translate(req *types.TranslateReq) (*types.TranslateResp, error) {
	i18nItem, err := l.svcCtx.WfI18nItemModel.FindOneByNameLang(l.ctx, req.Name, "zh")
	if err != nil && err != model.ErrNotFound {
		l.Errorf("Translate.error: %+v", err)
		return nil, xerr.NewErrorWithMsg(fmt.Sprintf("查询词条出错: %s", req.Name))
	}
	if i18nItem == nil {
		return nil, xerr.NewErrorWithMsg(fmt.Sprintf("未找到该词条: %s", req.Name))
	}
	entry, err := l.svcCtx.WfEntryModel.FindOneByUniqueName(l.ctx, i18nItem.UniqueName)
	if err != nil && err != model.ErrNotFound {
		l.Errorf("Translate.error: %+v", err)
		return nil, xerr.NewErrorWithMsg(fmt.Sprintf("查询词条出错: %s", req.Name))
	}
	if entry == nil {
		return nil, xerr.NewErrorWithMsg(fmt.Sprintf("未找到该词条: %s对应的翻译内容", req.Name))
	}
	return &types.TranslateResp{Result: entry.Name}, nil
}
