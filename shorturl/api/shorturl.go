package main

import (
	"flag"
	"fmt"
	"shorturl/wangjian-zero/rest"

	"shorturl/api/internal/config"
	"shorturl/api/internal/handler"
	"shorturl/api/internal/svc"

	"shorturl/wangjian-zero/core/conf"
)

var configFile = flag.String("f", "etc/shorturl-api.yaml", "the config file")
var serviceName = flag.String("serviceName", "wangjian-zero", "the service name")

func main() {
	flag.Parse()

	var c config.Config
	conf.MustLoad(*configFile, &c)

	ctx := svc.NewServiceContext(c)
	server := rest.MustNewServer(c.RestConf)
	defer server.Stop()

	handler.RegisterHandlers(server, ctx)

	fmt.Printf("Starting server at %s:%d...\n", c.Host, c.Port)
	server.Start()
}
