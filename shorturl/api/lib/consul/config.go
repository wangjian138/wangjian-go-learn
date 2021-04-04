package consul

import (
	"shorturl/wangjian-zero/core/consul"
)

var ConsulConn *consul.ZeroConsul

func NewConsulByAddr(addr string) {

	ConsulConn = consul.NewConsulClient(addr)
}
