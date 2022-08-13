package middleware

import (
    "github.com/yguilai/pipiao-openapi/app/model"
    "github.com/yguilai/pipiao-openapi/common/response"
    "github.com/yguilai/pipiao-openapi/common/xerr"
    "github.com/yguilai/pipiao-openapi/common/xjwt"
    "github.com/zeromicro/go-zero/core/logx"
    "net/http"
    "strings"
)

type OauthMiddleware struct {
    OpenApiModel model.OpenapiAuthModel
    logx.Logger
}

const (
    AuthorizationHeader = "Authorization"
    AppIdHeader         = "x-open-appid"
    AuthorizationPrefix = "Bearer "
)

func NewOauthMiddleware(m model.OpenapiAuthModel) *OauthMiddleware {
    return &OauthMiddleware{OpenApiModel: m}
}

func (m *OauthMiddleware) Handle(next http.HandlerFunc) http.HandlerFunc {
    return func(w http.ResponseWriter, r *http.Request) {
        m.Logger = logx.WithContext(r.Context())
        // token e.g: Bearer xxxxxx
        token, err := m.verifyToken(r.Header.Get(AuthorizationHeader))
        if err != nil {
            response.AuthorizationFailedResult(r, w, err)
            return
        }
        auth, err := m.verifyAppIdAndGetAuth(r)
        if err != nil {
            response.AuthorizationFailedResult(r, w, err)
            return
        }

        parser := xjwt.NewParser[xjwt.OpenapiClaims]()
        claims, err := parser.Parse(token, auth.AppKey)
        if err != nil {
            m.Errorf("token parse error: %+v", err)
            response.AuthorizationFailedResult(r, w, xerr.TokenParseError)
            return
        }
        if err := claims.Valid(); err != nil {
            m.Errorf("token valid error: %+v", err)
            response.AuthorizationFailedResult(r, w, xerr.TokenExpireError)
            return
        }
        next(w, r)
    }
}

func (m *OauthMiddleware) verifyToken(token string) (string, error) {
    // no authorization token 403
    if token == "" {
        return "", xerr.NewErrorWithMsg("token不能为空")
    }
    if !strings.HasPrefix(token, AuthorizationPrefix) {
        return "", xerr.NewErrorWithMsg("token格式错误")
    }
    return strings.ReplaceAll(token, AuthorizationPrefix, ""), nil
}

func (m *OauthMiddleware) verifyAppIdAndGetAuth(r *http.Request) (*model.OpenapiAuth, error) {
    appId := r.Header.Get(AppIdHeader)
    if appId == "" {
        return nil, xerr.NewErrorWithMsg("appId不能为空")
    }
    // 自带缓存
    auth, err := m.OpenApiModel.FindOneByAppId(r.Context(), appId)
    if err != nil {
        logx.WithContext(r.Context()).Errorf("verifyAppId error: %+v", err)
        if err == model.ErrNotFound {
            return nil, xerr.NewErrorWithMsg("未授权的appId")
        }
        return nil, err
    }
    return auth, nil
}
