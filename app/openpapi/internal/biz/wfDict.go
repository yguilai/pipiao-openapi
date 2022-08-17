package biz

import (
    "context"
    "encoding/json"
    "github.com/google/go-github/github"
    "github.com/panjf2000/ants/v2"
    "github.com/yguilai/pipiao-openapi/app/model"
    "github.com/yguilai/pipiao-openapi/common/xerr"
    "github.com/yguilai/pipiao-openapi/common/xgithub"
    "github.com/zeromicro/go-zero/core/hash"
    "github.com/zeromicro/go-zero/core/logx"
    "github.com/zeromicro/go-zero/core/stores/redis"
    "github.com/zeromicro/go-zero/core/threading"
    "io"
    "net/http"
    "strings"
    "sync"
)

const (
    warframeDictOwner           = "WFCD"
    warframeDictRepo            = "warframe-items"
    warframeDictPath            = "/data/json/i18n.json"
    warframeDictSyncSHARedisKey = "wf:dict:last_sha"
    // 定时任务每天跑一次, 缓存过期时间比24h长就行, 这里就先设为48h
    warframeDictSyncSHARedisExpire = 86400 * 2
    // warframeDictUpdateTaskKey 词典更新redis唯一key, 用于全局分布式锁, 避免同时跑太多任务
    warframeDictUpdateTaskKey = "wf:dict:update_unique_key"
    // warframeDictUpdateTaskExpire key过期时间, 10分钟
    warframeDictUpdateTaskExpire = 60 * 10
)

var needSaveLang = []string{"zh", "uk"}

type WfDictService struct {
    redis    *redis.Redis
    client   *github.Client
    wfiModel model.WfItemModel
    pool     *ants.Pool
}

func NewWfDictService(r *redis.Redis, m model.WfItemModel) *WfDictService {
    pool, _ := ants.NewPool(3)
    return &WfDictService{
        redis:    r,
        wfiModel: m,
        client:   xgithub.NewSimpleClient(),
        pool:     pool,
    }
}

// NeedFetch 判断是否需要拉取wf词条
func (s *WfDictService) NeedFetch(ctx context.Context) (string, string, bool) {
    ctt, _, _, err := s.client.Repositories.GetContents(ctx, warframeDictOwner, warframeDictRepo, warframeDictPath, nil)
    if err != nil {
        logx.Errorf("获取Wf词条仓库信息出错: %+v", err)
        return "", "", false
    }
    // TODO: 这里其实存在一致性问题, 可以改用lua脚本来校验, 保证原子性, 目前只是单机部署, 问题不大.
    //       可惜go这个redis操作库不能像java里一样支持针对某个key开启事务 不然就简单了
    // 获取上次拉取的
    lastSHA, err := s.redis.Get(warframeDictSyncSHARedisKey)
    if err != nil {
        logx.Errorf("获取上次拉取SHA信息出错: %+v", err)
        return "", "", false
    }

    // lastSHA不为空说明不是第一次拉取, lastSHA == *ctt.SHA说明词条还未更新
    if lastSHA != "" && lastSHA == *ctt.SHA {
        return "", "", false
    }
    return *ctt.DownloadURL, *ctt.SHA, true
}

func (s *WfDictService) StartUpdateTask(ctx context.Context, downloadUrl, sha string) error {
    lock := redis.NewRedisLock(s.redis, warframeDictUpdateTaskKey)
    lock.SetExpire(warframeDictUpdateTaskExpire)
    ok, err := lock.AcquireCtx(ctx)
    if err != nil {
        return err
    }
    if !ok {
        return xerr.NewErrorWithMsg("更新任务正在进行中, 请稍后再试~")
    }
    defer lock.Release()
    data, err := fetchDictData(downloadUrl)
    if err != nil {
        return err
    }
    if len(data) == 0 {
        logx.WithContext(ctx).Infof("没有需要更新的数据\n")
        return nil
    }

    errch := make(chan error, 1)
    threading.GoSafe(func() {
        defer close(errch)
        var wg sync.WaitGroup

        for k, entry := range data {
            miArr := generateWfItem(k, entry, needSaveLang)
            for _, mi := range miArr {
                wg.Add(1)
                err := s.pool.Submit(func() {
                    defer wg.Done()
                    doUpdate(ctx, s.wfiModel, &mi, errch)
                })
                if err != nil {
                    errch <- err
                }
            }
        }
        wg.Wait()
    })

    for e := range errch {
        if e != nil {
            return e
        }
    }
    err = s.updateLastSHA(ctx, sha)
    logx.Errorf("更新上次Warframe字典文件SHA失败: %+v", err)
    return nil
}

func (s *WfDictService) updateLastSHA(ctx context.Context, sha string) error {
    return s.redis.SetexCtx(ctx, warframeDictSyncSHARedisKey, sha, warframeDictSyncSHARedisExpire)
}

func doUpdate(ctx context.Context, wim model.WfItemModel, mi *model.WfItem, errch chan<- error) {
    wfItem, err := wim.FindOneByKeyLang(ctx, mi.Key, mi.Lang)
    if err != nil && err != model.ErrNotFound {
        errch <- err
        return
    }
    // 说明是新增的
    if wfItem == nil {
        res, err := wim.Insert(ctx, mi)
        if err != nil {
            errch <- err
            return
        }
        if rows, _ := res.RowsAffected(); rows == 0 {
            errch <- xerr.NewErrorWithMsg("插入失败")
            return
        }
        return
    }
    // 不为空则是已存在, 判断一下是否需要更新
    if !needUpdate(wfItem, mi) {
        logx.WithContext(ctx).Infof("词条%s-%s-%s无需更新", mi.Key, mi.Lang, mi.Name)
        return
    }
    mi.Id = wfItem.Id
    err = wim.Update(ctx, mi)
    if err != nil {
        errch <- err
        return
    }
}

func fetchDictData(downloadUrl string) (map[string]I18nEntry, error) {
    resp, err := http.Get(downloadUrl)
    if err != nil {
        return nil, err
    }
    defer resp.Body.Close()
    bytes, err := io.ReadAll(resp.Body)
    if err != nil {
        return nil, err
    }
    m := make(map[string]I18nEntry)
    err = json.Unmarshal(bytes, &m)
    if err != nil {
        panic(err)
    }
    return m, nil
}

func generateWfItem(key string, entry I18nEntry, needSaveLang []string) []model.WfItem {
    r := make([]model.WfItem, len(needSaveLang))
    for _, lang := range needSaveLang {
        item := entry[lang]
        mi := model.WfItem{
            Key:  key,
            Lang: lang,
            Name: item.Name,
        }
        switch item.Description.(type) {
        case string:
            mi.Description = item.Description.(string)
        case []string:
            ds := item.Description.([]string)
            mi.Description = strings.Join(ds, ";")
        default:
            mi.Description = ""
        }
        r = append(r, mi)
    }
    return r
}

func needUpdate(old *model.WfItem, n *model.WfItem) bool {
    oldKey := hash.Md5Hex([]byte(old.Key + old.Lang + old.Name + old.Description))
    newKey := hash.Md5Hex([]byte(n.Key + n.Lang + n.Name + n.Description))
    return oldKey != newKey
}
