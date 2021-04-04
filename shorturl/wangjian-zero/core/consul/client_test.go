package consul

import (
	"context"
	"testing"
)

var (
	ctx = context.Background()
)

func TestNewConsulClient(t *testing.T) {
	zConsul := NewConsulClient("http://localhost:8500")
	conf, err := zConsul.GetConfigByServiceName("wangjian-zero")
	if err != nil {
		t.Errorf("err:%v", err)
		return
	}

	t.Logf("conf:%v", string(conf))
}
