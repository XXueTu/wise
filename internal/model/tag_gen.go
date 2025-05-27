package model

import (
	"context"
	"time"

	"github.com/uptrace/bun"
)

// Tags 标签集合
type Tags struct {
	bun.BaseModel `bun:"table:tags,alias:t"`

	ID          int64     `bun:"id,pk,autoincrement" json:"id"`
	Uid         string    `bun:"uid,notnull" json:"uid"`                 // 标签唯一标识
	Name        string    `bun:"name,notnull" json:"name"`               // 标签名称
	Description string    `bun:"description,notnull" json:"description"` // 标签描述
	Color       string    `bun:"color,notnull" json:"color"`             // 标签颜色
	Icon        string    `bun:"icon,notnull" json:"icon"`               // 标签图标
	CreatedAt   time.Time `bun:"created_at,notnull,default:current_timestamp" json:"created_at"`
	UpdatedAt   time.Time `bun:"updated_at,notnull,default:current_timestamp" json:"updated_at"`
}

type TagGen interface {
	TableName() string
	InitData()
	Create(ctx context.Context, tag *Tags) error
	Update(ctx context.Context, tag *Tags) error
	Delete(ctx context.Context, id int64) error
	Get(ctx context.Context, id int64) (*Tags, error)
	GetUid(ctx context.Context, uid string) (*Tags, error)
	GetUids(ctx context.Context, uids []string) ([]*Tags, error)
	GetName(ctx context.Context, name string) (*Tags, error)
	GetList(ctx context.Context, page, size int64, name string) (*TagsList, error)
	FindBatchByNames(ctx context.Context, names []string) ([]*Tags, error)
	CreateBatch(ctx context.Context, tags []Tags) error
}

// BeforeCreate 创建前的钩子
func (m *Tags) BeforeCreate(ctx context.Context) (context.Context, error) {
	now := time.Now()
	m.CreatedAt = now
	m.UpdatedAt = now
	return ctx, nil
}

// BeforeUpdate 更新前的钩子
func (m *Tags) BeforeUpdate(ctx context.Context) (context.Context, error) {
	m.UpdatedAt = time.Now()
	return ctx, nil
}
