package consul

import (
	"github.com/hashicorp/consul/api"
	"shorturl/wangjian-zero/core/consul"
)

var ConsulConn *consul.ZeroConsul

func NewConsulByAddr(addr string) {

	ConsulConn = consul.NewConsulClient(addr)
}

//注册ip
func RegisterAddr(registration *api.AgentServiceRegistration) {
	err := ConsulConn.Agent().ServiceRegister(registration)
	if err != nil {
		panic(err)
	}
}

//反注册ip
func DeRegisterAddr(serviceId string) {
	err := ConsulConn.Agent().ServiceDeregister(serviceId)
	if err != nil {
		panic(err)
	}
}
