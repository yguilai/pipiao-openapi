package trigger

import (
	"github.com/yguilai/pipiao-openapi/common/response"
	"net/http"

	"github.com/yguilai/pipiao-openapi/app/openpapi/internal/logic/trigger"
	"github.com/yguilai/pipiao-openapi/app/openpapi/internal/svc"
)

func TriggerEntryHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		l := trigger.NewTriggerEntryLogic(r.Context(), svcCtx)
		resp, err := l.TriggerEntry()
		response.HttpResult(r, w, resp, err)
	}
}
