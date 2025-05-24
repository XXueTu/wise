package models

import (
	"context"
	"errors"

	"github.com/zeromicro/go-zero/core/logx"

	"github.com/XXueTu/wise/internal/model"
	"github.com/XXueTu/wise/internal/svc"
	"github.com/XXueTu/wise/internal/types"
)

type CreateModelLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 创建模型
func NewCreateModelLogic(ctx context.Context, svcCtx *svc.ServiceContext) *CreateModelLogic {
	return &CreateModelLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *CreateModelLogic) CreateModel(req *types.CreateModelRequest) (resp *types.Model, err error) {
	models := model.Models{
		BaseUrl:       req.BaseUrl,
		Config:        req.Config,
		Type:          req.Type,
		ModelName:     req.ModelName,
		ModelRealName: req.ModelRealName,
	}
	err = l.svcCtx.ModelsModel.Create(l.ctx, &models)
	if err != nil {
		return nil, errors.New("创建模型失败")
	}
	resp = &types.Model{
		Id:            models.ID,
		BaseUrl:       models.BaseUrl,
		Config:        models.Config,
		Type:          models.Type,
		ModelName:     models.ModelName,
		ModelRealName: models.ModelRealName,
	}
	return resp, nil
}
