package tasks

import (
	"context"
	"time"

	"github.com/zeromicro/go-zero/core/logx"

	"github.com/XXueTu/wise/internal/svc"
	"github.com/XXueTu/wise/internal/types"
)

type ListTaskLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 获取任务列表
func NewListTaskLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ListTaskLogic {
	return &ListTaskLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *ListTaskLogic) ListTask(req *types.ListTaskRequest) (resp *types.ListTaskResponse, err error) {

	taskList, err := l.svcCtx.TasksModel.GetPage(l.ctx, req.Page, req.PageSize, req.Name, req.Status, req.Types)
	if err != nil {
		return nil, err
	}
	response := make([]types.TaskResponse, 0)
	for _, task := range taskList.List {
		response = append(response, types.TaskResponse{
			Tid:          task.Tid,
			Name:         task.Name,
			Types:        task.Types,
			Status:       task.Status,
			CurrentState: task.CurrentState,
			TotalSteps:   task.TotalSteps,
			CurrentStep:  task.CurrentStep,
			RetryCount:   task.RetryCount,
			Params:       task.Params,
			Result:       task.Result,
			Error:        task.Error,
			Extend:       task.Extend,
			CreatedAt:    task.CreatedAt.Format(time.DateTime),
			UpdatedAt:    task.UpdatedAt.Format(time.DateTime),
		})
	}
	return &types.ListTaskResponse{
		Total: taskList.Total,
		List:  response,
	}, nil
}
