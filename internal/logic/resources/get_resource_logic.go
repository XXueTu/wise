package resources

import (
	"context"
	"errors"

	"github.com/zeromicro/go-zero/core/logx"

	"github.com/XXueTu/wise/internal/svc"
	"github.com/XXueTu/wise/internal/types"
)

type GetResourceLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 获取单个资源
func NewGetResourceLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetResourceLogic {
	return &GetResourceLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GetResourceLogic) GetResource(req *types.GetResourceRequest) (resp *types.Resource, err error) {
	resource, err := l.svcCtx.ResourceModel.Get(l.ctx, req.Id)
	if err != nil {
		return nil, errors.New("获取资源失败")
	}
	resp = &types.Resource{
		Id:      resource.ID,
		URL:     resource.URL,
		Title:   resource.Title,
		Content: resource.Content,
		Type:    resource.Type,
	}
	return resp, nil
}
