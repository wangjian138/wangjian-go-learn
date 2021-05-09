package svc

import (
	"shorturl/api/internal/config"
	"shorturl/rpc/transform/transformer"
	"shorturl/wangjian-zero/zrpc"
)

type ServiceContext struct {
	Config      config.Config
	Transformer transformer.Transformer // manual code
}

func NewServiceContext(c config.Config) *ServiceContext {
	return &ServiceContext{
		Config:      c,
		Transformer: transformer.NewTransformer(zrpc.MustNewClient(c.Transform)), // manual code
	}
}
