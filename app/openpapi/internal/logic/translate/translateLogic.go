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
    wfi := l.svcCtx.WfItemModel
    // 先根据原来的语言查到词条
    wi, err := wfi.FindOneByNameLang(l.ctx, req.Name, req.From)
    if err != nil && err != model.ErrNotFound {
        l.Errorf("Translate.error: %+v", err)
        return nil, xerr.NewErrorWithMsg(fmt.Sprintf("查询词条出错: %s", req.Name))
    }
    if wi == nil {
        return nil, xerr.NewErrorWithMsg(fmt.Sprintf("未找到该词条: %s", req.Name))
    }
    // 根据词条的key, 和要翻译的语言
    res, err := wfi.FindOneByKeyLang(l.ctx, wi.Key, req.Target)
    if err != nil && err != model.ErrNotFound {
        l.Errorf("Translate.error: %+v", err)
        return nil, xerr.NewErrorWithMsg(fmt.Sprintf("查询词条出错: %s", req.Name))
    }
    if res == nil {
        return nil, xerr.NewErrorWithMsg(fmt.Sprintf("未找到该词条: %s对应的翻译内容", req.Name))
    }
    return &types.TranslateResp{res.Name}, nil
}
