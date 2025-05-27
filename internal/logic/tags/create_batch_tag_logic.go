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

type CreateBatchTagLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 批量创建标签
func NewCreateBatchTagLogic(ctx context.Context, svcCtx *svc.ServiceContext) *CreateBatchTagLogic {
	return &CreateBatchTagLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *CreateBatchTagLogic) CreateBatchTag(req *types.CreateBatchTagRequest) (resp *types.CreateBatchTagResponse, err error) {
	if len(req.Tags) == 0 {
		return nil, errors.New("tags is empty")
	}
	// 批量查询标签是否存在
	tagNames := make([]string, 0)
	for _, tag := range req.Tags {
		tagNames = append(tagNames, tag.Name)
	}
	existTags, err := l.svcCtx.TagsModel.FindBatchByNames(l.ctx, tagNames)
	if err != nil {
		return nil, err
	}
	// 判断哪些标签不存在
	notExistTags := make([]model.Tags, 0)
	existTagMap := make(map[string]struct{})
	for _, existTag := range existTags {
		existTagMap[existTag.Name] = struct{}{}
	}
	for _, tag := range req.Tags {
		if _, exist := existTagMap[tag.Name]; !exist {
			notExistTags = append(notExistTags, model.Tags{
				Name:        tag.Name,
				Uid:         uuid.New().String(),
				Description: tag.Description,
				Color:       tag.Color,
				Icon:        tag.Icon,
			})
		}
	}
	// 批量创建标签
	if len(notExistTags) > 0 {
		err = l.svcCtx.TagsModel.CreateBatch(l.ctx, notExistTags)
		if err != nil {
			return nil, err
		}
	}
	// 查询添加的标签
	addedTags, err := l.svcCtx.TagsModel.FindBatchByNames(l.ctx, tagNames)
	if err != nil {
		return nil, err
	}

	for _, tag := range addedTags {
		tagNames = append(tagNames, tag.Name)
	}
	createTagResponses := make([]types.CreateTagResponse, len(req.Tags))
	for i, tag := range addedTags {
		createTagResponses[i] = types.CreateTagResponse{
			Id:          tag.ID,
			Uid:         tag.Uid,
			Name:        tag.Name,
			Description: tag.Description,
			Color:       tag.Color,
			Icon:        tag.Icon,
			CreatedAt:   tag.CreatedAt.Format(time.DateTime),
			UpdatedAt:   tag.UpdatedAt.Format(time.DateTime),
		}
	}
	// 返回结果
	resp = &types.CreateBatchTagResponse{
		CreatedTotal:       int64(len(notExistTags)),
		ExistedTotal:       int64(len(existTags)),
		CreateTagResponses: createTagResponses,
	}
	return resp, nil
}
