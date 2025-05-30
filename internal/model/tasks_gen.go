package model

import (
	"context"
	"time"

	"github.com/uptrace/bun"
)

// Tasks 任务
type Tasks struct {
	bun.BaseModel `bun:"table:tasks,alias:t"`
	ID            int64     `bun:"id,pk,autoincrement" json:"id"`
	Tid           string    `bun:"tid,notnull" json:"tid"`                     // 任务唯一标识
	Name          string    `bun:"name,notnull" json:"name"`                   // 任务名称
	Types         string    `bun:"types,notnull" json:"types"`                 // 任务类型
	Status        string    `bun:"status,notnull" json:"status"`               // 任务状态
	CurrentState  string    `bun:"current_state,notnull" json:"current_state"` // 当前状态机
	TotalSteps    int64     `bun:"total_steps,notnull" json:"total_steps"`     // 总步骤
	CurrentStep   int64     `bun:"current_step,notnull" json:"current_step"`   // 当前步骤
	RetryCount    int64     `bun:"retry_count,notnull" json:"retry_count"`     // 重试次数
	Params        string    `bun:"params,notnull" json:"params"`               // 任务参数
	Result        string    `bun:"result,notnull" json:"result"`               // 任务结果
	Duration      int64     `bun:"duration,notnull" json:"duration"`           // 任务耗时 ms
	Error         string    `bun:"error,notnull" json:"error"`                 // 任务错误
	Extend        string    `bun:"extend,notnull" json:"extend"`               // 扩展字段
	CreatedAt     time.Time `bun:"created_at,notnull,default:current_timestamp" json:"created_at"`
	UpdatedAt     time.Time `bun:"updated_at,notnull,default:current_timestamp" json:"updated_at"`
}

type TasksGen interface {
	TableName() string
	InitData()
	Create(ctx context.Context, tag *Tasks) error
	Update(ctx context.Context, tag *Tasks) error
	Delete(ctx context.Context, id int64) error
	Get(ctx context.Context, id int64) (*Tasks, error)
	GetByTid(ctx context.Context, tid string) (*Tasks, error)
	GetStatus(ctx context.Context, status string) ([]*Tasks, error)
	GetStatusLimit(ctx context.Context, status string, limit int) ([]*Tasks, error)
	GetPage(ctx context.Context, page int64, pageSize int64, name string, status string, types string) (*TasksList, error)
	UpdateState(ctx context.Context, tid string, state string, result string) error
	UpdateStateAndStep(ctx context.Context, tid string, state string, step int64, result string) error
	UpdateStatus(ctx context.Context, tid string, status string, error string) error
}

const (
	TaskStatusInit      = "init"      // 初始状态
	TaskStatusRunning   = "running"   // 运行中
	TaskStatusSuccess   = "success"   // 成功
	TaskStatusFailed    = "failed"    // 失败
	TaskStatusRetry     = "retry"     // 重试中
	TaskStatusCancelled = "cancelled" // 已取消
)

func (m *Tasks) BeforeInsert(ctx context.Context, query *bun.InsertQuery) error {
	m.CreatedAt = time.Now()
	return nil
}

func (m *Tasks) BeforeUpdate(ctx context.Context, query *bun.UpdateQuery) error {
	m.UpdatedAt = time.Now()
	return nil
}
