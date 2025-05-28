package tasks

import (
	"context"
	"encoding/json"
	"errors"

	"github.com/zeromicro/go-zero/core/logx"

	"github.com/XXueTu/wise/internal/svc"
	"github.com/XXueTu/wise/internal/task"
	"github.com/XXueTu/wise/internal/types"
)

type CreateTaskLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 创建任务
func NewCreateTaskLogic(ctx context.Context, svcCtx *svc.ServiceContext) *CreateTaskLogic {
	return &CreateTaskLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *CreateTaskLogic) CreateTask(req *types.CreateTaskRequest) (resp *types.CreateTaskResponse, err error) {
	taskArgs := task.UrlMarkTaskArgs{}
	err = json.Unmarshal([]byte(req.Params), &taskArgs)
	if err != nil {
		return nil, err
	}
	err = task.CreateUrlMarkTask(l.ctx, l.svcCtx, taskArgs)
	if err != nil {
		return nil, errors.New("创建任务失败")
	}
	return
}
