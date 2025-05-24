package resources

import (
	"context"
	"errors"

	"github.com/zeromicro/go-zero/core/logx"

	"github.com/XXueTu/wise/internal/svc"
	"github.com/XXueTu/wise/internal/types"
)

type UpdateResourceLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 更新资源
func NewUpdateResourceLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UpdateResourceLogic {
	return &UpdateResourceLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *UpdateResourceLogic) UpdateResource(req *types.UpdateResourceRequest) (resp *types.Resource, err error) {
	resource, err := l.svcCtx.ResourceModel.Get(l.ctx, req.Id)
	if err != nil {
		return nil, errors.New("获取资源失败")
	}
	resource.Title = req.Title
	resource.Content = req.Content
	resource.Type = req.Type
	err = l.svcCtx.ResourceModel.Update(l.ctx, resource)
	if err != nil {
		return nil, errors.New("更新资源失败")
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
