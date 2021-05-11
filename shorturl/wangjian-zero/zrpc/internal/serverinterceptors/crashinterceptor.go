package serverinterceptors

import (
	"context"
	"runtime/debug"

	"shorturl/wangjian-zero/core/logx"
	"shorturl/wangjian-zero/grpc"
	"shorturl/wangjian-zero/grpc/codes"
	"shorturl/wangjian-zero/grpc/status"
)

// StreamCrashInterceptor catches panics in processing stream requests and recovers.
func StreamCrashInterceptor(srv interface{}, stream grpc.ServerStream, info *grpc.StreamServerInfo,
	handler grpc.StreamHandler) (err error) {
	defer handleCrash(func(r interface{}) {
		err = toPanicError(r)
	})

	return handler(srv, stream)
}

// UnaryCrashInterceptor catches panics in processing unary requests and recovers.
func UnaryCrashInterceptor() grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo,
		handler grpc.UnaryHandler) (resp interface{}, err error) {
		defer handleCrash(func(r interface{}) {
			err = toPanicError(r)
		})

		return handler(ctx, req)
	}
}

func handleCrash(handler func(interface{})) {
	if r := recover(); r != nil {
		handler(r)
	}
}

func toPanicError(r interface{}) error {
	logx.Errorf("%+v %s", r, debug.Stack())
	return status.Errorf(codes.Internal, "panic: %v", r)
}
