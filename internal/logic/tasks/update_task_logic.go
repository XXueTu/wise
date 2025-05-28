package tasks

import (
	"context"
	"errors"

	"github.com/zeromicro/go-zero/core/logx"

	"github.com/XXueTu/wise/internal/svc"
	"github.com/XXueTu/wise/internal/types"
)

type UpdateTaskLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 更新任务
func NewUpdateTaskLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UpdateTaskLogic {
	return &UpdateTaskLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *UpdateTaskLogic) UpdateTask(req *types.UpdateTaskRequest) (resp *types.UpdateTaskResponse, err error) {
	task, err := l.svcCtx.TasksModel.GetByTid(l.ctx, req.Tid)
	if err != nil {
		return nil, errors.New("任务不存在")
	}
	task.Name = req.Name
	if err := l.svcCtx.TasksModel.Update(l.ctx, task); err != nil {
		return nil, errors.New("更新任务失败")
	}
	return
}
