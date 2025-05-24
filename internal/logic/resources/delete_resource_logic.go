package resources

import (
	"context"
	"errors"

	"github.com/zeromicro/go-zero/core/logx"

	"github.com/XXueTu/wise/internal/svc"
	"github.com/XXueTu/wise/internal/types"
)

type DeleteResourceLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 删除资源
func NewDeleteResourceLogic(ctx context.Context, svcCtx *svc.ServiceContext) *DeleteResourceLogic {
	return &DeleteResourceLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *DeleteResourceLogic) DeleteResource(req *types.DeleteResourceRequest) (resp *types.Resource, err error) {
	err = l.svcCtx.ResourceModel.Delete(l.ctx, req.Id)
	if err != nil {
		return nil, errors.New("删除资源失败")
	}
	return
}
