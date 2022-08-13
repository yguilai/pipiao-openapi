// Code generated by goctl. DO NOT EDIT.
package handler

import (
    "github.com/yguilai/pipiao-openapi/app/openpapi/internal/handler/translate"
    "github.com/yguilai/pipiao-openapi/app/openpapi/internal/svc"
    "net/http"

    "github.com/zeromicro/go-zero/rest"
)

func RegisterHandlers(server *rest.Server, serverCtx *svc.ServiceContext) {
	server.AddRoutes(
		rest.WithMiddlewares(
			[]rest.Middleware{serverCtx.Oauth},
			[]rest.Route{
				{
					Method:  http.MethodGet,
					Path:    "/translate/:from",
					Handler: translate.TranslateHandler(serverCtx),
				},
			}...,
		),
		rest.WithPrefix("/wf/v1"),
	)
}
