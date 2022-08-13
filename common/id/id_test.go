package id

import (
    "github.com/zeromicro/go-zero/core/logx"
    "testing"
)

func TestNew(t *testing.T) {
    id := New("cli")
    logx.Info(id)
    //cli_1000000000
    //cli_3361977894
    //cli_1953844557
    //cli_3812958319
}
