package model

import (
	"context"
	"time"

	"github.com/uptrace/bun"
)

// Resource 资源模型
type Resource struct {
	bun.BaseModel `bun:"table:resources,alias:r"`

	ID        int64     `bun:"id,pk,autoincrement" json:"id"`
	URL       string    `bun:"url,notnull" json:"url"`         // 资源URL
	Title     string    `bun:"title,notnull" json:"title"`     // 资源标题
	Content   string    `bun:"content,notnull" json:"content"` // 资源内容
	Type      string    `bun:"type,notnull" json:"type"`       // 资源类型（如：wechat, zhihu等）
	Tags      string    `bun:"tags,notnull" json:"tags"`       // 资源标签
	CreatedAt time.Time `bun:"created_at,notnull,default:current_timestamp" json:"created_at"`
	UpdatedAt time.Time `bun:"updated_at,notnull,default:current_timestamp" json:"updated_at"`
}

type ResourceGen interface {
	TableName() string
	InitData()
	Create(ctx context.Context, resource *Resource) error
	Update(ctx context.Context, resource *Resource) error
	Delete(ctx context.Context, id int64) error
	Get(ctx context.Context, id int64) (*Resource, error)
	GetByURL(ctx context.Context, url string) (*Resource, error)
	GetList(ctx context.Context, page, size int, resourceType, title string, tagUids []string) (*ResourceList, error)
}

func (m *Resource) BeforeInsert(ctx context.Context, query *bun.InsertQuery) error {
	m.CreatedAt = time.Now()
	return nil
}

func (m *Resource) BeforeUpdate(ctx context.Context, query *bun.UpdateQuery) error {
	m.UpdatedAt = time.Now()
	return nil
}
