package resources

import (
	"context"
	"errors"
	"strings"

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
	tagUids := strings.Join(req.TagUids, ",")
	resource.Tags = tagUids
	err = l.svcCtx.ResourceModel.Update(l.ctx, resource)
	if err != nil {
		return nil, errors.New("更新资源失败")
	}
	// 获取标签
	tagList, err := l.svcCtx.TagsModel.GetUids(l.ctx, req.TagUids)
	if err != nil {
		return nil, errors.New("获取标签失败")
	}
	var tags []string
	for _, tag := range tagList {
		tags = append(tags, tag.Name)
	}
	resp = &types.Resource{
		Id:      resource.ID,
		URL:     resource.URL,
		Title:   resource.Title,
		Content: resource.Content,
		Type:    resource.Type,
		Tags:    tags,
	}
	return resp, nil
}
