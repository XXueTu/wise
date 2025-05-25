package tags

import (
	"context"
	"errors"
	"time"

	"github.com/zeromicro/go-zero/core/logx"

	"github.com/XXueTu/wise/internal/svc"
	"github.com/XXueTu/wise/internal/types"
)

type GetTagLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 获取标签详情
func NewGetTagLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetTagLogic {
	return &GetTagLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GetTagLogic) GetTag(req *types.GetTagRequest) (resp *types.TagResponse, err error) {
	// 查找是否存在
	tag, err := l.svcCtx.TagsModel.GetUid(l.ctx, req.Uid)
	if err != nil {
		return nil, errors.New("查询标签失败")
	}
	if tag == nil {
		return nil, errors.New("标签不存在")
	}
	resp = &types.TagResponse{
		Uid:         tag.Uid,
		Name:        tag.Name,
		Description: tag.Description,
		Color:       tag.Color,
		Icon:        tag.Icon,
		CreatedAt:   tag.CreatedAt.Format(time.DateTime),
		UpdatedAt:   tag.UpdatedAt.Format(time.DateTime),
	}
	return
}
