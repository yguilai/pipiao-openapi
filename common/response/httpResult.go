package response

import (
    "github.com/pkg/errors"
    "github.com/yguilai/pipiao-openapi/common/xerr"
    "github.com/zeromicro/go-zero/core/logx"
    "github.com/zeromicro/go-zero/rest/httpx"
    "google.golang.org/grpc/status"
    "net/http"
)

func HttpResult(r *http.Request, w http.ResponseWriter, data interface{}, err error) {
    if err == nil {
        r := success(data)
        httpx.WriteJson(w, http.StatusOK, r)
        return
    }
    var respError *xerr.CodeError
    causeErr := errors.Cause(err)
    if e, ok := causeErr.(*xerr.CodeError); ok {
        respError = e
    } else if grpcStatus, ok := status.FromError(causeErr); ok {
        respError = xerr.NewCodeError(uint32(grpcStatus.Code()), grpcStatus.Message())
    } else {
        respError = xerr.ServerError
    }
    logx.WithContext(r.Context()).Errorf("[API-ERROR]: %+v", err)
    httpx.WriteJson(w, http.StatusOK, failedWithErr(respError))
}

func AuthorizationFailedResult(r *http.Request, w http.ResponseWriter, err error) {
    if err == nil {
        httpx.WriteJson(w, http.StatusForbidden, failedWithErr(xerr.ForbiddenError))
        return
    }
    var respError *xerr.CodeError
    causeErr := errors.Cause(err)
    if e, ok := causeErr.(*xerr.CodeError); ok {
        respError = e
    } else {
        respError = xerr.ServerError
    }
    logx.WithContext(r.Context()).Errorf("[API-ERROR]: %+v", err)
    httpx.WriteJson(w, http.StatusForbidden, failedWithErr(respError))
}
