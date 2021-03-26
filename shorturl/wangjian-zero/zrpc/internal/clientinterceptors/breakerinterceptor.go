package clientinterceptors

import (
	"context"
	"path"

	"google.golang.org/grpc"
	"shorturl/go-zero/core/breaker"
	"shorturl/go-zero/zrpc/internal/codes"
)

// BreakerInterceptor is an interceptor that acts as a circuit breaker.
func BreakerInterceptor(ctx context.Context, method string, req, reply interface{},
	cc *grpc.ClientConn, invoker grpc.UnaryInvoker, opts ...grpc.CallOption) error {
	breakerName := path.Join(cc.Target(), method)
	return breaker.DoWithAcceptable(breakerName, func() error {
		return invoker(ctx, method, req, reply, cc, opts...)
	}, codes.Acceptable)
}
