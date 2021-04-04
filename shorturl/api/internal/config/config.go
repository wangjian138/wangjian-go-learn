package config

import (
	"shorturl/api/lib/consul"
	"shorturl/wangjian-zero/core/mapping"
	"shorturl/wangjian-zero/rest"
)

type Config struct {
	rest.RestConf
}

//从consul获取 后面进行抽象化 支持etcd等
func GetConfigByServiceName(serviceName string) Config {
	v, err := consul.ConsulConn.GetConfigByServiceName(serviceName)
	if err != nil {
		panic(err)
	}
	conf := Config{}

	err = mapping.UnmarshalYamlBytes([]byte(v), &conf)
	if err != nil {
		panic(err)
	}

	return conf
}
