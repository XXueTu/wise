package tags

import (
	"context"
	"errors"
	"time"

	"github.com/google/uuid"
	"github.com/zeromicro/go-zero/core/logx"

	"github.com/XXueTu/wise/internal/model"
	"github.com/XXueTu/wise/internal/svc"
	"github.com/XXueTu/wise/internal/types"
)

type CreateTagLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 创建标签
func NewCreateTagLogic(ctx context.Context, svcCtx *svc.ServiceContext) *CreateTagLogic {
	return &CreateTagLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *CreateTagLogic) CreateTag(req *types.CreateTagRequest) (resp *types.CreateTagResponse, err error) {
	// 查找是否存在
	tag, err := l.svcCtx.TagsModel.GetName(l.ctx, req.Name)
	if err != nil {
		return nil, errors.New("查询标签失败")
	}
	if tag != nil {
		return nil, errors.New("标签已经存在")
	}
	uid := uuid.New().String()
	tag = &model.Tags{
		Uid:         uid,
		Name:        req.Name,
		Description: req.Description,
		Color:       req.Color,
		Icon:        req.Icon,
	}
	err = l.svcCtx.TagsModel.Create(l.ctx, tag)
	if err != nil {
		return nil, errors.New("创建标签失败")
	}

	resp = &types.CreateTagResponse{
		Id:          tag.ID,
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
