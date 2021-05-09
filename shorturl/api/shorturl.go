package main

import (
	"flag"
	"fmt"
	"shorturl/api/arg"
	"shorturl/api/lib"
	"shorturl/wangjian-zero/rest"

	"shorturl/api/internal/config"
	"shorturl/api/internal/handler"
	"shorturl/api/internal/svc"
	_ "shorturl/wangjian-zero/core/proc"
)

func main() {
	flag.Parse()

	//进行实例化
	lib.Init()

	c := config.GetConfigByServiceName(*arg.ServiceName)

	ctx := svc.NewServiceContext(c)
	server := rest.MustNewServer(c.RestConf)
	defer server.Stop()

	handler.RegisterHandlers(server, ctx)

	fmt.Printf("Starting server at %s:%d...\n", c.Host, c.Port)
	server.Start()
}
