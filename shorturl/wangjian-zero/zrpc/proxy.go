package zrpc

import (
	"context"
	"sync"

	"google.golang.org/grpc"
	"shorturl/wangjian-zero/core/syncx"
	"shorturl/wangjian-zero/zrpc/internal"
	"shorturl/wangjian-zero/zrpc/internal/auth"
)

// A RpcProxy is a rpc proxy.
type RpcProxy struct {
	backend     string
	clients     map[string]Client
	options     []internal.ClientOption
	sharedCalls syncx.SharedCalls
	lock        sync.Mutex
}

// NewProxy returns a RpcProxy.
func NewProxy(backend string, opts ...internal.ClientOption) *RpcProxy {
	return &RpcProxy{
		backend:     backend,
		clients:     make(map[string]Client),
		options:     opts,
		sharedCalls: syncx.NewSharedCalls(),
	}
}

// TakeConn returns a grpc.ClientConn.
func (p *RpcProxy) TakeConn(ctx context.Context) (*grpc.ClientConn, error) {
	cred := auth.ParseCredential(ctx)
	key := cred.App + "/" + cred.Token
	val, err := p.sharedCalls.Do(key, func() (interface{}, error) {
		p.lock.Lock()
		client, ok := p.clients[key]
		p.lock.Unlock()
		if ok {
			return client, nil
		}

		opts := append(p.options, WithDialOption(grpc.WithPerRPCCredentials(&auth.Credential{
			App:   cred.App,
			Token: cred.Token,
		})))
		client, err := NewClientWithTarget(p.backend, opts...)
		if err != nil {
			return nil, err
		}

		p.lock.Lock()
		p.clients[key] = client
		p.lock.Unlock()
		return client, nil
	})
	if err != nil {
		return nil, err
	}

	return val.(Client).Conn(), nil
}
