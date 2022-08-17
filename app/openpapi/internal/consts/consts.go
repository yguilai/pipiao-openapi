package consts

const (
    // WarframeDictUpdateTaskKey 词典更新redis唯一key, 用于全局分布式锁, 避免同时跑太多任务
    WarframeDictUpdateTaskKey = "wf:dict:update_unique_key"
    // WarframeDictUpdateTaskExpire key过期时间, 10分钟, 也就是说每十分钟仅能更新一次, 比较文件很大
    WarframeDictUpdateTaskExpire = 60 * 10
)
