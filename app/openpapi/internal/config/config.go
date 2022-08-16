package config

import (
    "github.com/zeromicro/go-zero/core/stores/cache"
    "github.com/zeromicro/go-zero/core/stores/redis"
    "github.com/zeromicro/go-zero/rest"
)

type Config struct {
    rest.RestConf
    DataSource string
    Cache      cache.CacheConf
    Redis      redis.RedisConf
}
