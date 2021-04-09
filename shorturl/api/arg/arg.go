package arg

import (
	"flag"
	"shorturl/wangjian-zero/core/netx"
)

var configFile = flag.String("f", "etc/shorturl-api.yaml", "the config file")
var ServiceName = flag.String("ServiceName", "wangjian-zero", "the service name")
var ConsulAddr = flag.String("consulAddr", "localhost:8500", "the consul adde")
var OutAddr = flag.String("OutAddr", "", "the out adde")
var OutAddrPort = flag.Int("OutAddrPort", 8888, "the out adde")

func init() {
	//判断是否是对位的地址
	if *OutAddr == "" {
		*OutAddr = netx.InternalIp()
	}
}
