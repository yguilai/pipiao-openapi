package biz

import (
    "context"
    "fmt"
    "github.com/google/go-github/github"
    "github.com/yguilai/pipiao-openapi/common/xgithub"
    "github.com/zeromicro/go-zero/core/logx"
    "github.com/zeromicro/go-zero/core/stores/redis"
)

const (
    warframeDictOwner           = "WFCD"
    warframeDictRepo            = "warframe-items"
    warframeDictPath            = "/data/json/i18n.json"
    warframeDictSyncSHARedisKey = "wf:dict:last_sha"
    // 定时任务每天跑一次, 缓存过期时间比24h长就行, 这里就先设为48h
    warframeDictSyncSHARedisExpire = 86400 * 2
)

type WfDictService struct {
    redis  *redis.Redis
    client *github.Client
}

func NewWfDictService(r *redis.Redis) *WfDictService {
    return &WfDictService{
        redis:  r,
        client: xgithub.NewSimpleClient(),
    }
}

// NeedFetch 判断是否需要拉取wf词条
func (s *WfDictService) NeedFetch(ctx context.Context) (string, bool) {
    ctt, _, _, err := s.client.Repositories.GetContents(ctx, warframeDictOwner, warframeDictRepo, warframeDictPath, nil)
    if err != nil {
        logx.Errorf("获取Wf词条仓库信息出错: %+v", err)
        return "", false
    }
    // TODO: 这里其实存在一致性问题, 可以改用lua脚本来校验, 保证原子性, 目前只是单机部署, 问题不大.
    //       可惜go这个redis操作库不能像java里一样支持针对某个key开启事务 不然就简单了
    // 获取上次拉取的
    lastSHA, err := s.redis.Get(warframeDictSyncSHARedisKey)
    if err != nil {
        logx.Errorf("获取上次拉取SHA信息出错: %+v", err)
        return "", false
    }

    // lastSHA不为空说明不是第一次拉取, lastSHA == *ctt.SHA说明词条还未更新
    if lastSHA != "" && lastSHA == *ctt.SHA {
        return "", false
    }
    err = s.redis.SetexCtx(ctx, warframeDictSyncSHARedisKey, *ctt.SHA, warframeDictSyncSHARedisExpire)
    if err != nil {
        // 这里就只告警一下吧
        logx.Alert(fmt.Sprintf("保存本次拉取的SHA失败: %+v", err))
        return *ctt.DownloadURL, true
    }
    return *ctt.DownloadURL, true
}
