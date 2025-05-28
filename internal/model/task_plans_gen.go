package model

import (
	"context"
	"time"

	"github.com/uptrace/bun"
)

// Task 任务
type TaskPlans struct {
	bun.BaseModel `bun:"table:task_plans,alias:p"`

	ID        int64     `bun:"id,pk,autoincrement" json:"id"`
	Tid       string    `bun:"tid,notnull" json:"tid"`           // 任务唯一标识
	Pid       string    `bun:"pid,notnull" json:"pid"`           // 任务计划唯一标识
	Types     string    `bun:"types,notnull" json:"types"`       // 任务类型
	Name      string    `bun:"name,notnull" json:"name"`         // 任务名称
	Index     int64     `bun:"index,notnull" json:"index"`       // 任务计划索引
	Status    string    `bun:"status,notnull" json:"status"`     // 任务状态 init,running,success,failed,cancelled
	Params    string    `bun:"params,notnull" json:"params"`     // 任务参数
	Result    string    `bun:"result,notnull" json:"result"`     // 任务结果
	Duration  int64     `bun:"duration,notnull" json:"duration"` // 任务耗时 ms
	Error     string    `bun:"error,notnull" json:"error"`       // 任务错误
	CreatedAt time.Time `bun:"created_at,notnull,default:current_timestamp" json:"created_at"`
	UpdatedAt time.Time `bun:"updated_at,notnull,default:current_timestamp" json:"updated_at"`
}

type TaskPlansGen interface {
	TableName() string
	InitData()
	Create(ctx context.Context, tag *TaskPlans) error
	CreateBatch(ctx context.Context, tag []TaskPlans) error
	Update(ctx context.Context, tag *TaskPlans) error
	Delete(ctx context.Context, id int64) error
	Get(ctx context.Context, id int64) (*TaskPlans, error)
	GetInitByTid(ctx context.Context, tid string) ([]*TaskPlans, error)
	GetByTid(ctx context.Context, tid string) ([]*TaskPlans, error)
}

const (
	TaskPlanStatusInit      = "init"
	TaskPlanStatusRunning   = "running"
	TaskPlanStatusSuccess   = "success"
	TaskPlanStatusFailed    = "failed"
	TaskPlanStatusCancelled = "cancelled"
)

// BeforeCreate 创建前的钩子
func (m *TaskPlans) BeforeCreate(ctx context.Context) (context.Context, error) {
	now := time.Now()
	m.CreatedAt = now
	m.UpdatedAt = now
	return ctx, nil
}

// BeforeUpdate 更新前的钩子
func (m *TaskPlans) BeforeUpdate(ctx context.Context) (context.Context, error) {
	m.UpdatedAt = time.Now()
	return ctx, nil
}
