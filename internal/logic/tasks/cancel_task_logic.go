package tasks

import (
	"context"

	"github.com/zeromicro/go-zero/core/logx"

	"github.com/XXueTu/wise/internal/model"
	"github.com/XXueTu/wise/internal/svc"
	"github.com/XXueTu/wise/internal/types"
)

type CancelTaskLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 取消任务
func NewCancelTaskLogic(ctx context.Context, svcCtx *svc.ServiceContext) *CancelTaskLogic {
	return &CancelTaskLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *CancelTaskLogic) CancelTask(req *types.CancelTaskRequest) (resp *types.TaskOperationResponse, err error) {
	// todo: add your logic here and delete this line
	task, err := l.svcCtx.TasksModel.GetByTid(l.ctx, req.Tid)
	if err != nil {
		return nil, err
	}
	task.Status = model.TaskStatusCancelled
	if err := l.svcCtx.TasksModel.Update(l.ctx, task); err != nil {
		return nil, err
	}
	return
}
