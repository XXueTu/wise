package resources

import (
	"context"
	"errors"

	"github.com/zeromicro/go-zero/core/logx"

	"github.com/XXueTu/wise/internal/model"
	"github.com/XXueTu/wise/internal/svc"
	"github.com/XXueTu/wise/internal/types"
)

type CreateResourceLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 创建资源
func NewCreateResourceLogic(ctx context.Context, svcCtx *svc.ServiceContext) *CreateResourceLogic {
	return &CreateResourceLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *CreateResourceLogic) CreateResource(req *types.CreateResourceRequest) (resp *types.Resource, err error) {
	resourceModel := &model.Resource{
		URL:     req.URL,
		Title:   req.Title,
		Content: req.Content,
		Type:    req.Type,
	}
	err = l.svcCtx.ResourceModel.Create(l.ctx, resourceModel)
	if err != nil {
		return nil, errors.New("创建资源失败")
	}
	resp = &types.Resource{
		Id:      resourceModel.ID,
		URL:     resourceModel.URL,
		Title:   resourceModel.Title,
		Content: resourceModel.Content,
		Type:    resourceModel.Type,
	}
	return resp, nil
}
