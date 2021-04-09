package consul

import (
	"fmt"
	"github.com/hashicorp/consul/api"
)

// HealthState 获取存活的健康服务列表
func HealthState(q *api.QueryOptions) (api.HealthChecks, error) {
	return state(api.HealthPassing, q)
}

// AnyState 获取服务列表
func AnyState(q *api.QueryOptions) (api.HealthChecks, error) {
	return state(api.HealthAny, q)
}

func state(state string, q *api.QueryOptions) (api.HealthChecks, error) {
	checks, meta, err := ConsulConn.Health().State(state, q)
	if err != nil {
		return nil, fmt.Errorf("检查服务列表出错：%s", err)
	}

	q.WaitIndex = meta.LastIndex
	return checks, nil
}
