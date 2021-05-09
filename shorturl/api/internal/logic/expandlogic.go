package logic

import (
	"context"
	"shorturl/rpc/transform/transformer"

	"shorturl/api/internal/svc"
	"shorturl/api/internal/types"

	"shorturl/wangjian-zero/core/logx"
)

type ExpandLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewExpandLogic(ctx context.Context, svcCtx *svc.ServiceContext) ExpandLogic {
	return ExpandLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *ExpandLogic) Expand(req types.ExpandReq) (*types.ExpandResp, error) {
	// manual code start
	resp, err := l.svcCtx.Transformer.Expand(l.ctx, &transformer.ExpandReq{
		Shorten: req.Shorten,
	})
	if err != nil {
		return nil, err
	}

	return &types.ExpandResp{
		Url: resp.Url,
	}, nil
	// manual code stop
}
