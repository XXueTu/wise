package resources

import (
	"context"
	"errors"
	"time"

	"github.com/zeromicro/go-zero/core/logx"

	"github.com/XXueTu/wise/internal/svc"
	"github.com/XXueTu/wise/internal/types"
)

type ListResourceLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 分页查询资源列表
func NewListResourceLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ListResourceLogic {
	return &ListResourceLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *ListResourceLogic) ListResource(req *types.ListResourceRequest) (resp *types.ListResourceResponse, err error) {
	logx.Infof("ListResourceLogic: %+v", req)
	resources, err := l.svcCtx.ResourceModel.GetList(l.ctx, int(req.Page), int(req.PageSize), req.Type, req.Keyword)
	if err != nil {
		return nil, errors.New("获取资源列表失败")
	}
	resp = &types.ListResourceResponse{
		Total:     resources.Total,
		Resources: make([]types.Resource, len(resources.List)),
	}
	for i, resource := range resources.List {
		resp.Resources[i] = types.Resource{
			Id:        resource.ID,
			URL:       resource.URL,
			Title:     resource.Title,
			Content:   resource.Content,
			Type:      resource.Type,
			CreatedAt: resource.CreatedAt.Format(time.DateTime),
			UpdatedAt: resource.UpdatedAt.Format(time.DateTime),
		}
	}
	return resp, nil
}
