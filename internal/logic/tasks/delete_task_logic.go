package tasks

import (
	"context"
	"errors"

	"github.com/zeromicro/go-zero/core/logx"

	"github.com/XXueTu/wise/internal/svc"
	"github.com/XXueTu/wise/internal/types"
)

type DeleteTaskLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 删除任务
func NewDeleteTaskLogic(ctx context.Context, svcCtx *svc.ServiceContext) *DeleteTaskLogic {
	return &DeleteTaskLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *DeleteTaskLogic) DeleteTask(req *types.DeleteTaskRequest) (resp *types.DeleteTaskResponse, err error) {
	task, err := l.svcCtx.TasksModel.GetByTid(l.ctx, req.Tid)
	if err != nil {
		return nil, errors.New("任务不存在")
	}
	if err := l.svcCtx.TasksModel.Delete(l.ctx, task.ID); err != nil {
		return nil, errors.New("删除任务失败")
	}
	return
}
