package consul

import (
	"github.com/hashicorp/consul/api"
	"testing"
)

func TestMain(m *testing.M) {
	NewConsulByAddr("localhost:8500")
}

func TestRegisterAddr(t *testing.T) {
	registration := &api.AgentServiceRegistration{
		Name:    "wangjian-zero",
		Address: "localhost",
		Port:    8888,
	}
	defer func() {
		if errs := recover(); errs != nil {
			t.Errorf("errs:%v", errs)
		}
	}()
	RegisterAddr(registration)
}
