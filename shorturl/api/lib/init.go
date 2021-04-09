package lib

import (
	"fmt"
	"github.com/hashicorp/consul/api"
	"shorturl/api/arg"
	"shorturl/api/lib/consul"
	zeroConsul "shorturl/wangjian-zero/core/consul"
	"strings"
)

func Init() {
	consul.NewConsulByAddr(*arg.ConsulAddr)

	//注册服务
	registration := &api.AgentServiceRegistration{
		ID:      fmt.Sprintf("%v-%v", *arg.OutAddr, *arg.OutAddrPort),
		Name:    *arg.ServiceName,
		Address: *arg.OutAddr,
		Port:    *arg.OutAddrPort,
		Check: &api.AgentServiceCheck{
			Interval: "4s",
			Timeout:  "5s",
			HTTP:     "http://127.0.0.1:8500",
		},
	}
	consul.RegisterAddr(registration)

	//获取服务
	q := &api.QueryOptions{RequireConsistent: true, WaitIndex: 0}
	checks, err := consul.AnyState(q)
	if err != nil {
		panic(err)
	}

	for _, check := range checks {
		if check.ServiceID == "" {
			continue
		}

		serviceSplit := strings.Split(check.ServiceID, "-")
		if len(serviceSplit) != 2 {
			//进行反解析
			consul.DeRegisterAddr(check.ServiceID)
			continue
		}

		serviceMap := consul.ConsulConn.ServiceMap[check.ServiceName]
		serviceMap = append(serviceMap, zeroConsul.ZeroConsulService{Host: serviceSplit[0], Port: serviceSplit[1]})
		consul.ConsulConn.ServiceMap[check.ServiceName] = serviceMap
	}
}
