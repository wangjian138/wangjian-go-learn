package internal

import "shorturl/wangjian-zero/grpc"

//WithStreamServerInterceptors uses given server stream interceptors.
func WithStreamServerInterceptors(interceptors ...grpc.StreamServerInterceptor) grpc.ServerOption {
	return grpc.ChainStreamInterceptor(interceptors...)
}

//
// WithUnaryServerInterceptors uses given server unary interceptors.
func WithUnaryServerInterceptors(interceptors ...grpc.UnaryServerInterceptor) grpc.ServerOption {
	return grpc.ChainUnaryInterceptor(interceptors...)
}
