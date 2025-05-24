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

type ListModelLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 分页查询模型列表
func NewListModelLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ListModelLogic {
	return &ListModelLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *ListModelLogic) ListModel(req *types.ListModelRequest) (resp *types.ListModelResponse, err error) {

	models, err := l.svcCtx.ModelsModel.GetList(l.ctx,
		req.Page, req.PageSize,
		req.Type,
		req.Tag, req.Status,
		req.Keyword)
	if err != nil {
		return nil, errors.New("获取模型列表失败")
	}
	resp = &types.ListModelResponse{
		Total:  models.Total,
		Models: make([]types.Model, len(models.List)),
	}
	for i, model := range models.List {
		tag := []string{}
		json.Unmarshal([]byte(model.Tag), &tag)
		resp.Models[i] = types.Model{
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
	}
	return resp, nil
}
