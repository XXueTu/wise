package tags

import (
	"context"
	"errors"

	"github.com/zeromicro/go-zero/core/logx"

	"github.com/XXueTu/wise/internal/svc"
	"github.com/XXueTu/wise/internal/types"
)

type UpdateTagLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 更新标签
func NewUpdateTagLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UpdateTagLogic {
	return &UpdateTagLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *UpdateTagLogic) UpdateTag(req *types.UpdateTagRequest) (resp *types.UpdateTagResponse, err error) {
	// 查找是否存在
	tag, err := l.svcCtx.TagsModel.GetUid(l.ctx, req.Uid)
	if err != nil {
		return nil, errors.New("查询标签失败")
	}
	if tag == nil {
		return nil, errors.New("标签不存在")
	}
	tag.Name = req.Name
	tag.Description = req.Description
	tag.Color = req.Color
	tag.Icon = req.Icon
	err = l.svcCtx.TagsModel.Update(l.ctx, tag)
	if err != nil {
		return nil, errors.New("更新标签失败")
	}
	resp = &types.UpdateTagResponse{
		Uid:         tag.Uid,
		Name:        tag.Name,
		Description: tag.Description,
		Color:       tag.Color,
		Icon:        tag.Icon,
	}
	return
}
