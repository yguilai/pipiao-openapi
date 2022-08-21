package translate

import (
    "github.com/yguilai/pipiao-openapi/common/response"
    "net/http"

    "github.com/yguilai/pipiao-openapi/app/openpapi/internal/logic/translate"
    "github.com/yguilai/pipiao-openapi/app/openpapi/internal/svc"
)

func TriggerHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
    return func(w http.ResponseWriter, r *http.Request) {
        l := translate.NewTriggerLogic(r.Context(), svcCtx)
        resp, err := l.Trigger()
        response.HttpResult(r, w, resp, err)
    }
}
