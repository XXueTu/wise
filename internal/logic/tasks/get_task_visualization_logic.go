package tasks

import (
	"context"
	"errors"
	"time"

	"github.com/zeromicro/go-zero/core/logx"

	"github.com/XXueTu/wise/internal/svc"
	"github.com/XXueTu/wise/internal/types"
)

type GetTaskVisualizationLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 获取任务可视化信息
func NewGetTaskVisualizationLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetTaskVisualizationLogic {
	return &GetTaskVisualizationLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GetTaskVisualizationLogic) GetTaskVisualization(req *types.GetTaskVisualizationRequest) (resp *types.TaskVisualizationResponse, err error) {
	task, err := l.svcCtx.TasksModel.GetByTid(l.ctx, req.Tid)
	if err != nil {
		return nil, errors.New("任务不存在")
	}
	plans, err := l.svcCtx.TaskPlansModel.GetByTid(l.ctx, task.Tid)
	if err != nil {
		return nil, errors.New("获取任务计划失败")
	}
	taskPlanDetails := make([]types.TaskPlanDetail, 0)
	for _, plan := range plans {
		taskPlanDetails = append(taskPlanDetails, types.TaskPlanDetail{
			Pid:       plan.Pid,
			Name:      plan.Name,
			Index:     plan.Index,
			Status:    plan.Status,
			Params:    plan.Params,
			Result:    plan.Result,
			Duration:  plan.Duration,
			Error:     plan.Error,
			CreatedAt: plan.CreatedAt.Format(time.DateTime),
			UpdatedAt: plan.UpdatedAt.Format(time.DateTime),
		})
	}
	resp = &types.TaskVisualizationResponse{
		Tid:          task.Tid,
		Name:         task.Name,
		Types:        task.Types,
		Status:       task.Status,
		CurrentState: task.CurrentState,
		TotalSteps:   task.TotalSteps,
		CurrentStep:  task.CurrentStep,
		Plans:        taskPlanDetails,
		CreatedAt:    task.CreatedAt.Format(time.DateTime),
		UpdatedAt:    task.UpdatedAt.Format(time.DateTime),
	}
	return resp, nil
}
