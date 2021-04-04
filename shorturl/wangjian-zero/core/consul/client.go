package consul

import (
	"fmt"
	"github.com/hashicorp/consul/api"
	"github.com/hashicorp/go-cleanhttp"
)

type ZeroConsul struct {
	*api.Client
}

func NewConsulClient(addr string) *ZeroConsul {
	conf := &api.Config{
		Address:   addr,
		Transport: cleanhttp.DefaultPooledTransport(),
	}

	var err error
	client, err := api.NewClient(conf)

	if err != nil {
		panic("consul创建失败：" + err.Error())
	}

	//_, err = client.Status().Peers()
	//if err != nil {
	//	panic("consul连接失败：" + err.Error())
	//}

	return &ZeroConsul{client}
}

func (z *ZeroConsul) GetConfigByServiceName(serviceName string) (string, error) {
	v, err := z.getKV(fmt.Sprintf("/config/%v", serviceName), nil)

	if err != nil {
		panic(err)
	}

	return string(v), nil
}

func (z *ZeroConsul) getKV(k string, q *api.QueryOptions) (v []byte, err error) {
	pair, meta, err := z.KV().Get(k, q)
	if err != nil {
		err = fmt.Errorf("检查[%s]配置出错：%s", k, err)
		return
	}

	if pair == nil {
		err = fmt.Errorf("[%s]配置不存在", k)
		return
	}

	v = pair.Value
	if q == nil {
		return
	}

	q.WaitIndex = meta.LastIndex
	return
}
