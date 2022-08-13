package xjwt

import (
    "fmt"
    "github.com/golang-jwt/jwt/v4"
    "testing"
    "time"
)

func TestSign(t *testing.T) {
    c := OpenapiClaims{
        RegisteredClaims: jwt.RegisteredClaims{
            Issuer:    "ygl",
            ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Second * 86400)),
            NotBefore: jwt.NewNumericDate(time.Now()),
            IssuedAt:  jwt.NewNumericDate(time.Now()),
        },
        AppId: "cli_1000000000",
    }
    sign, err := Sign("d00cbf369b424fbb96abe3640f64a637", c)
    if err != nil {
        panic(err)
    }
    fmt.Println(sign)
}
