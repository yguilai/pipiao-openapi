package svc

import (
    "github.com/yguilai/pipiao-openapi/app/model"
    "github.com/yguilai/pipiao-openapi/app/openpapi/internal/config"
    "github.com/yguilai/pipiao-openapi/app/openpapi/internal/middleware"
    "github.com/zeromicro/go-zero/core/stores/redis"
    "github.com/zeromicro/go-zero/core/stores/sqlx"
    "github.com/zeromicro/go-zero/rest"
)

type ServiceContext struct {
    Config       config.Config
    Oauth        rest.Middleware
    OpenApiModel model.OpenapiAuthModel
    WfItemModel  model.WfItemModel
    Redis        *redis.Redis
}

func (c *ServiceContext) IsDev() bool {
    return c.Config.Mode == "dev"
}

func (c *ServiceContext) IsPro() bool {
    return c.Config.Mode == "pro"
}

func NewServiceContext(c config.Config) *ServiceContext {
    mysql := sqlx.NewMysql(c.DataSource)
    authModel := model.NewOpenapiAuthModel(mysql, c.Cache)

    return &ServiceContext{
        Config:       c,
        Oauth:        middleware.NewOauthMiddleware(authModel).Handle,
        OpenApiModel: authModel,
        WfItemModel:  model.NewWfItemModel(mysql, c.Cache),
        Redis:        c.Redis.NewRedis(),
    }
}
