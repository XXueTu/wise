package models

import (
	"context"
	"errors"

	"github.com/zeromicro/go-zero/core/logx"

	"github.com/XXueTu/wise/internal/svc"
	"github.com/XXueTu/wise/internal/types"
)

type DeleteModelLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 删除模型
func NewDeleteModelLogic(ctx context.Context, svcCtx *svc.ServiceContext) *DeleteModelLogic {
	return &DeleteModelLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *DeleteModelLogic) DeleteModel(req *types.DeleteModelRequest) (resp *types.Model, err error) {
	err = l.svcCtx.ModelsModel.Delete(l.ctx, req.Id)
	if err != nil {
		return nil, errors.New("删除模型失败")
	}
	resp = &types.Model{
		Id: req.Id,
	}
	return resp, nil
}
