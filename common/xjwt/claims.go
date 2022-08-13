package xjwt

import "github.com/golang-jwt/jwt/v4"

type OpenapiClaims struct {
    jwt.RegisteredClaims
    AppId string
}

type IdentityClaims struct {
    jwt.RegisteredClaims
    AccountId string
    SessionId string
}
