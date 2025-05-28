package tasks

import (
	"context"
	"errors"
	"time"

	"github.com/zeromicro/go-zero/core/logx"

	"github.com/XXueTu/wise/internal/svc"
	"github.com/XXueTu/wise/internal/types"
)

type GetTaskLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 获取任务详情
func NewGetTaskLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetTaskLogic {
	return &GetTaskLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GetTaskLogic) GetTask(req *types.GetTaskRequest) (resp *types.TaskResponse, err error) {
	task, err := l.svcCtx.TasksModel.GetByTid(l.ctx, req.Tid)
	if err != nil {
		return nil, errors.New("任务不存在")
	}
	resp = &types.TaskResponse{
		Tid:          task.Tid,
		Name:         task.Name,
		Types:        task.Types,
		Params:       task.Params,
		TotalSteps:   task.TotalSteps,
		CurrentState: task.CurrentState,
		Status:       task.Status,
		CurrentStep:  task.CurrentStep,
		RetryCount:   task.RetryCount,
		Result:       task.Result,
		Error:        task.Error,
		Extend:       task.Extend,
		CreatedAt:    task.CreatedAt.Format(time.DateTime),
		UpdatedAt:    task.UpdatedAt.Format(time.DateTime),
	}
	return resp, nil
}
