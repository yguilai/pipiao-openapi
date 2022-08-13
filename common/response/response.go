package response

import "github.com/yguilai/pipiao-openapi/common/xerr"

type Response struct {
    Code    uint32      `json:"code"`
    Message string      `json:"message"`
    Data    interface{} `json:"data,omitempty"`
}

type NullJson struct {
}

func success(data interface{}) *Response {
    return &Response{
        Code:    200,
        Message: "ok",
        Data:    data,
    }
}

func failed(code uint32, msg string) *Response {
    return &Response{Code: code, Message: msg}
}

func failedWithErr(err *xerr.CodeError) *Response {
    return &Response{
        Code:    err.GetErrCode(),
        Message: err.GetErrMsg(),
    }
}
