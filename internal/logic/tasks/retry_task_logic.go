package tasks

import (
	"context"
	"errors"

	"github.com/zeromicro/go-zero/core/logx"

	"github.com/XXueTu/wise/internal/model"
	"github.com/XXueTu/wise/internal/svc"
	"github.com/XXueTu/wise/internal/types"
)

type RetryTaskLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 重试任务
func NewRetryTaskLogic(ctx context.Context, svcCtx *svc.ServiceContext) *RetryTaskLogic {
	return &RetryTaskLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *RetryTaskLogic) RetryTask(req *types.RetryTaskRequest) (resp *types.TaskOperationResponse, err error) {
	task, err := l.svcCtx.TasksModel.GetByTid(l.ctx, req.Tid)
	if err != nil {
		return nil, errors.New("任务不存在")
	}
	task.Status = model.TaskStatusRunning
	task.CurrentStep = 0
	if err := l.svcCtx.TasksModel.Update(l.ctx, task); err != nil {
		return nil, errors.New("重试任务失败")
	}
	return
}
