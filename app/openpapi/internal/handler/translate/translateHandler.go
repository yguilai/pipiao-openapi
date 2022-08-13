package translate

import (
    "github.com/yguilai/pipiao-openapi/app/openpapi/internal/logic/translate"
    "github.com/yguilai/pipiao-openapi/app/openpapi/internal/svc"
    "github.com/yguilai/pipiao-openapi/app/openpapi/internal/types"
    "github.com/yguilai/pipiao-openapi/common/response"
    "net/http"

    "github.com/zeromicro/go-zero/rest/httpx"
)

func TranslateHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
    return func(w http.ResponseWriter, r *http.Request) {
        var req types.TranslateReq
        if err := httpx.Parse(r, &req); err != nil {
            httpx.Error(w, err)
            return
        }

        l := translate.NewTranslateLogic(r.Context(), svcCtx)
        resp, err := l.Translate(&req)
        response.HttpResult(r, w, resp, err)
    }
}
