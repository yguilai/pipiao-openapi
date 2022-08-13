package xjwt

import "github.com/golang-jwt/jwt/v4"

type Claims interface {
    OpenapiClaims | IdentityClaims
}

type Parser[T Claims] struct {
}

func (p *Parser[T]) Parse(token, secret string) (*T, error) {
    return p.doParse(token, secret)
}

func (p *Parser[T]) doParse(token, secret string) (*T, error) {
    var claims T
    var tokens *jwt.Token
    var err error
    switch any(claims).(type) {
    case IdentityClaims:
        tokens, err = jwt.ParseWithClaims(token, &IdentityClaims{}, func(token *jwt.Token) (interface{}, error) {
            return []byte(secret), nil
        }, jwt.WithJSONNumber())
        break
    case OpenapiClaims:
        tokens, err = jwt.ParseWithClaims(token, &OpenapiClaims{}, func(token *jwt.Token) (interface{}, error) {
            return []byte(secret), nil
        }, jwt.WithJSONNumber())
        break
    }
    if err != nil {
        return &claims, err
    }
    if err := tokens.Claims.Valid(); err != nil {
        return &claims, err
    }
    return any(tokens.Claims).(*T), nil
}

func NewParser[T Claims]() *Parser[T] {
    return &Parser[T]{}
}

func Sign(secret string, claims jwt.Claims) (string, error) {
    token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
    return token.SignedString([]byte(secret))
}
