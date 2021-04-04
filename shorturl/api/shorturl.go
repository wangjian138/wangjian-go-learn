package main

import (
	"flag"
	"fmt"
	"shorturl/api/lib/consul"
	"shorturl/wangjian-zero/rest"

	"shorturl/api/internal/config"
	"shorturl/api/internal/handler"
	"shorturl/api/internal/svc"
)

var configFile = flag.String("f", "etc/shorturl-api.yaml", "the config file")
var serviceName = flag.String("serviceName", "wangjian-zero", "the service name")
var ConsulAddr = flag.String("consulAddr", "localhost:8500", "the consul adde")

func main() {
	flag.Parse()

	//todo 抽象化
	consul.NewConsulByAddr(*ConsulAddr)
	c := config.GetConfigByServiceName(*serviceName)

	ctx := svc.NewServiceContext(c)
	server := rest.MustNewServer(c.RestConf)
	defer server.Stop()

	handler.RegisterHandlers(server, ctx)

	fmt.Printf("Starting server at %s:%d...\n", c.Host, c.Port)
	server.Start()
}
