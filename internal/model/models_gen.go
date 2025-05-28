package model

import (
	"context"
	"time"

	"github.com/uptrace/bun"
)

// Models 模型集合
type Models struct {
	bun.BaseModel `bun:"table:models,alias:m"`

	ID            int64     `bun:"id,pk,autoincrement" json:"id"`
	BaseUrl       string    `bun:"base_url,notnull" json:"base_url"`               // 资源链接 // https://ark.cn-beijing.volces.com/api/v3
	Config        string    `bun:"config,notnull" json:"config"`                   // 资源配置 {"apiKey":"9567f3a1-7e2e-4fa7-a8db-5a7ee0926"}
	Type          string    `bun:"type,notnull" json:"type"`                       // 资源类型（如：doubao, ollama等）
	ModelName     string    `bun:"model_name,notnull" json:"model_name"`           // 模型名称 豆包1.5
	ModelRealName string    `bun:"model_real_name,notnull" json:"model_real_name"` // 模型名称 doubao-1-5-pro-32k-250115
	Status        string    `bun:"status,notnull" json:"status"`                   // 状态（如：active,inactive,available,not_available）
	Tag           string    `bun:"tag,notnull" json:"tag"`                         // 标签（如：doubao, ollama等）
	CreatedAt     time.Time `bun:"created_at,notnull,default:current_timestamp" json:"created_at"`
	UpdatedAt     time.Time `bun:"updated_at,notnull,default:current_timestamp" json:"updated_at"`
}

type ModelsGen interface {
	TableName() string
	InitData()
	Create(ctx context.Context, model *Models) error
	Update(ctx context.Context, model *Models) error
	Delete(ctx context.Context, id int64) error
	Get(ctx context.Context, id int64) (*Models, error)
	GetList(ctx context.Context, page, size int64, modelType string, tag []string, status, modelName string) (*ModelsList, error)
}

func (m *Models) BeforeInsert(ctx context.Context, query *bun.InsertQuery) error {
	m.CreatedAt = time.Now()
	return nil
}

func (m *Models) BeforeUpdate(ctx context.Context, query *bun.UpdateQuery) error {
	m.UpdatedAt = time.Now()
	return nil
}
