package tags

import (
	"context"
	"errors"

	"github.com/zeromicro/go-zero/core/logx"

	"github.com/XXueTu/wise/internal/svc"
	"github.com/XXueTu/wise/internal/types"
)

type ListTagLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 获取标签列表
func NewListTagLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ListTagLogic {
	return &ListTagLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *ListTagLogic) ListTag(req *types.ListTagRequest) (resp *types.ListTagResponse, err error) {
	tagList, err := l.svcCtx.TagsModel.GetList(l.ctx, req.Page, req.PageSize, req.Name)
	if err != nil {
		return nil, errors.New("查询标签失败")
	}
	var list []types.TagResponse
	for _, tag := range tagList.List {
		list = append(list, types.TagResponse{
			Uid:         tag.Uid,
			Name:        tag.Name,
			Description: tag.Description,
			Color:       tag.Color,
			Icon:        tag.Icon,
		})
	}
	resp = &types.ListTagResponse{
		Total: tagList.Total,
		List:  list,
	}
	return
}
