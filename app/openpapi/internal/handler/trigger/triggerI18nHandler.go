package trigger

import (
	"github.com/yguilai/pipiao-openapi/common/response"
	"net/http"

	"github.com/yguilai/pipiao-openapi/app/openpapi/internal/logic/trigger"
	"github.com/yguilai/pipiao-openapi/app/openpapi/internal/svc"
)

func TriggerI18nHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		l := trigger.NewTriggerI18nLogic(r.Context(), svcCtx)
		resp, err := l.TriggerI18n()
		response.HttpResult(r, w, resp, err)
	}
}
