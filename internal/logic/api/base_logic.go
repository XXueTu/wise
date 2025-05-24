package api

import (
	"context"

	"github.com/XXueTu/wise/internal/svc"
	"github.com/XXueTu/wise/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type BaseLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// API基础
func NewBaseLogic(ctx context.Context, svcCtx *svc.ServiceContext) *BaseLogic {
	return &BaseLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *BaseLogic) Base(req *types.BaseRequest) (resp *types.BaseResponse, err error) {
	// todo: add your logic here and delete this line

	return
}
