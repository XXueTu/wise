package tasks

import (
	"context"
	"errors"

	"github.com/zeromicro/go-zero/core/logx"

	"github.com/XXueTu/wise/internal/model"
	"github.com/XXueTu/wise/internal/svc"
	"github.com/XXueTu/wise/internal/types"
)

type ResumeTaskLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 恢复任务
func NewResumeTaskLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ResumeTaskLogic {
	return &ResumeTaskLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *ResumeTaskLogic) ResumeTask(req *types.ResumeTaskRequest) (resp *types.TaskOperationResponse, err error) {
	task, err := l.svcCtx.TasksModel.GetByTid(l.ctx, req.Tid)
	if err != nil {
		return nil, errors.New("任务不存在")
	}
	task.Status = model.TaskStatusInit
	if err := l.svcCtx.TasksModel.Update(l.ctx, task); err != nil {
		return nil, errors.New("恢复任务失败")
	}
	return resp, nil
}
