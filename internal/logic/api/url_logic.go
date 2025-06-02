package api

import (
	"context"

	"github.com/zeromicro/go-zero/core/logx"

	"github.com/XXueTu/wise/internal/svc"
	"github.com/XXueTu/wise/internal/types"
)

type UrlLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// URL链接识别
func NewUrlLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UrlLogic {
	return &UrlLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *UrlLogic) Url(req *types.URLRequest) (resp *types.URLResponse, err error) {
	return resp, nil
}
