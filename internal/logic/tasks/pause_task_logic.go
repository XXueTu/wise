package tasks

import (
	"context"
	"errors"

	"github.com/zeromicro/go-zero/core/logx"

	"github.com/XXueTu/wise/internal/model"
	"github.com/XXueTu/wise/internal/svc"
	"github.com/XXueTu/wise/internal/types"
)

type PauseTaskLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 暂停任务
func NewPauseTaskLogic(ctx context.Context, svcCtx *svc.ServiceContext) *PauseTaskLogic {
	return &PauseTaskLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *PauseTaskLogic) PauseTask(req *types.PauseTaskRequest) (resp *types.TaskOperationResponse, err error) {
	task, err := l.svcCtx.TasksModel.GetByTid(l.ctx, req.Tid)
	if err != nil {
		return nil, errors.New("任务不存在")
	}
	task.Status = model.TaskStatusCancelled
	if err := l.svcCtx.TasksModel.Update(l.ctx, task); err != nil {
		return nil, errors.New("暂停任务失败")
	}
	return
}
