package mock

import (
	"context"

	"shorturl/wangjian-zero/grpc/codes"
	"shorturl/wangjian-zero/grpc/status"
)

type DepositServer struct {
}

func (*DepositServer) Deposit(ctx context.Context, req *DepositRequest) (*DepositResponse, error) {
	if req.GetAmount() < 0 {
		return nil, status.Errorf(codes.InvalidArgument, "cannot deposit %v", req.GetAmount())
	}

	return &DepositResponse{Ok: true}, nil
}
