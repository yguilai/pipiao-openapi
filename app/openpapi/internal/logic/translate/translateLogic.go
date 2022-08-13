package translate

import (
    "context"
    "github.com/yguilai/pipiao-openapi/app/openpapi/internal/svc"
    "github.com/yguilai/pipiao-openapi/app/openpapi/internal/types"

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

func (l *TranslateLogic) Translate(req *types.TranslateReq) (resp *types.TranslateResp, err error) {
    // todo: add your logic here and delete this line

    return
}
