package logic

import (
	"context"
	"fmt"
	"shorturl/rpc/transform/transform"
	"shorturl/rpc/transform/transformer"

	"shorturl/api/internal/svc"
	"shorturl/api/internal/types"

	"shorturl/wangjian-zero/core/logx"
)

type ShortenLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewShortenLogic(ctx context.Context, svcCtx *svc.ServiceContext) ShortenLogic {
	return ShortenLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}
func (l *ShortenLogic) Shorten(req types.ShortenReq) (*types.ShortenResp, error) {
	// manual code start
	resp, err := l.svcCtx.Transformer.Shorten(l.ctx, &transformer.ShortenReq{
		Url: req.Url,
	})
	if err != nil {
		return nil, err
	}

	resp1, err := l.svcCtx.Transformer.Expand(context.Background(), &transform.ExpandReq{Shorten: "1"})
	if err != nil {
		fmt.Printf("err:%v\n", err)
	}

	fmt.Printf("resp:%+v\n", resp1)

	return &types.ShortenResp{
		Shorten: resp.Shorten,
	}, nil
	// manual code stop
}
