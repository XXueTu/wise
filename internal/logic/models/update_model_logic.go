package models

import (
	"context"
	"encoding/json"
	"errors"
	"time"

	"github.com/zeromicro/go-zero/core/logx"

	"github.com/XXueTu/wise/internal/svc"
	"github.com/XXueTu/wise/internal/types"
)

type UpdateModelLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 更新模型
func NewUpdateModelLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UpdateModelLogic {
	return &UpdateModelLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *UpdateModelLogic) UpdateModel(req *types.UpdateModelRequest) (resp *types.Model, err error) {
	tag := "[]"
	if len(req.Tag) > 0 {
		tagStr, _ := json.Marshal(req.Tag)
		tag = string(tagStr)
	}
	model, err := l.svcCtx.ModelsModel.Get(l.ctx, req.Id)
	if err != nil {
		return nil, errors.New("获取模型失败")
	}
	model.BaseUrl = req.BaseUrl
	model.Config = req.Config
	model.Type = req.Type
	model.ModelName = req.ModelName
	model.ModelRealName = req.ModelRealName
	model.Status = req.Status
	model.Tag = tag
	err = l.svcCtx.ModelsModel.Update(l.ctx, model)
	if err != nil {
		return nil, errors.New("更新模型失败")
	}
	resp = &types.Model{
		Id:            model.ID,
		BaseUrl:       model.BaseUrl,
		Config:        model.Config,
		Type:          model.Type,
		ModelName:     model.ModelName,
		ModelRealName: model.ModelRealName,
		Status:        model.Status,
		Tag:           req.Tag,
		CreatedAt:     model.CreatedAt.Format(time.DateTime),
		UpdatedAt:     model.UpdatedAt.Format(time.DateTime),
	}
	return resp, nil
}
