package xerr

import "fmt"

type CodeError struct {
    errCode uint32
    errMsg  string
}

func (e *CodeError) GetErrCode() uint32 {
    return e.errCode
}

func (e *CodeError) GetErrMsg() string {
    return e.errMsg
}

func (e *CodeError) Error() string {
    return fmt.Sprintf("ErrCode:%d, ErrMsg:%s", e.errCode, e.errMsg)
}

func NewCodeError(errCode uint32, errMsg string) *CodeError {
    return &CodeError{errCode: errCode, errMsg: errMsg}
}

func NewErrorWithMsg(errMsg string) *CodeError {
    return &CodeError{errCode: ServerErrorCode, errMsg: errMsg}
}
