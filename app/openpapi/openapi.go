package main

import (
    "flag"
    "fmt"
    "github.com/yguilai/pipiao-openapi/app/openpapi/internal/config"
    "github.com/yguilai/pipiao-openapi/app/openpapi/internal/handler"
    "github.com/yguilai/pipiao-openapi/app/openpapi/internal/svc"
    "github.com/zeromicro/go-zero/core/conf"
    "github.com/zeromicro/go-zero/core/service"
    "github.com/zeromicro/go-zero/core/threading"
    "github.com/zeromicro/go-zero/rest"
)

var configFile = flag.String("f", "etc/%s/openapi.yaml", "the config file")
var mode = flag.String("p", "dev", "the environment")

func main() {
    flag.Parse()

    var c config.Config
    cfgPath := fmt.Sprintf(*configFile, *mode)
    conf.MustLoad(cfgPath, &c)

    server := rest.MustNewServer(c.RestConf)
    defer server.Stop()
    group := service.NewServiceGroup()
    defer group.Stop()

    ctx := svc.NewServiceContext(c)
    handler.RegisterHandlers(server, ctx)
    handler.RegisterJobs(group, ctx)

    threading.GoSafe(func() {
        fmt.Println("Starting jobs")
        group.Start()
    })
    fmt.Printf("Starting server at %s:%d...\n", c.Host, c.Port)
    server.Start()
}
