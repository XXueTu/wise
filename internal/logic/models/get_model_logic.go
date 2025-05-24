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

type GetModelLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 获取单个模型
func NewGetModelLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetModelLogic {
	return &GetModelLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GetModelLogic) GetModel(req *types.GetModelRequest) (resp *types.Model, err error) {
	model, err := l.svcCtx.ModelsModel.Get(l.ctx, req.Id)
	if err != nil {
		return nil, errors.New("获取模型失败")
	}
	tag := []string{}
	json.Unmarshal([]byte(model.Tag), &tag)
	resp = &types.Model{
		Id:            model.ID,
		BaseUrl:       model.BaseUrl,
		Config:        model.Config,
		Type:          model.Type,
		ModelName:     model.ModelName,
		ModelRealName: model.ModelRealName,
		Status:        model.Status,
		Tag:           tag,
		CreatedAt:     model.CreatedAt.Format(time.DateTime),
		UpdatedAt:     model.UpdatedAt.Format(time.DateTime),
	}
	return resp, nil
}
