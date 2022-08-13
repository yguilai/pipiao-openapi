package xerr

const OK uint32 = 200

const (
    // ServerErrorCode 服务端通用错误码
    ServerErrorCode  uint32 = 100001
    ParamRequestCode uint32 = 100002
    TokenExpireCode  uint32 = 100003
    TokenFailedCode  uint32 = 100004
)

var (
    ServerError       = &CodeError{errCode: ServerErrorCode, errMsg: "服务器开小差了, 稍后再来试一试"}
    ForbiddenError    = &CodeError{errCode: TokenFailedCode, errMsg: "服务器开小差了, 稍后再来试一试"}
    ParamRequestError = &CodeError{errCode: ParamRequestCode, errMsg: "请求参数有误"}
    TokenParseError   = &CodeError{errCode: TokenFailedCode, errMsg: "token解析失败"}
    TokenExpireError  = &CodeError{errCode: TokenExpireCode, errMsg: "token已失效"}
)
