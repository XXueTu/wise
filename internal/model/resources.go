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
	CreatedAt time.Time `bun:"created_at,notnull,default:current_timestamp" json:"created_at"`
	UpdatedAt time.Time `bun:"updated_at,notnull,default:current_timestamp" json:"updated_at"`
}

// ResourceList 资源列表返回结构
type ResourceList struct {
	Total int64       `json:"total"` // 总记录数
	List  []*Resource `json:"list"`  // 资源列表
}

// TableName 返回表名
func (r *Resource) TableName() string {
	return "resources"
}

// BeforeCreate 创建前的钩子
func (r *Resource) BeforeCreate(ctx context.Context) (context.Context, error) {
	now := time.Now()
	r.CreatedAt = now
	r.UpdatedAt = now
	return ctx, nil
}

// BeforeUpdate 更新前的钩子
func (r *Resource) BeforeUpdate(ctx context.Context) (context.Context, error) {
	r.UpdatedAt = time.Now()
	return ctx, nil
}

// Create 创建资源
func (r *Resource) Create(ctx context.Context, resource *Resource) error {
	_, err := sqliteDB.NewInsert().Model(resource).Exec(ctx)
	return err
}

// GetByURL 根据URL获取资源
func (r *Resource) GetByURL(ctx context.Context, url string) (*Resource, error) {
	resource := new(Resource)
	err := sqliteDB.NewSelect().
		Model(resource).
		Where("url = ?", url).
		Scan(ctx)
	if err != nil {
		return nil, err
	}
	return resource, nil
}

// Update 更新资源
func (r *Resource) Update(ctx context.Context, resource *Resource) error {
	_, err := sqliteDB.NewUpdate().
		Model(resource).
		WherePK().
		Exec(ctx)
	return err
}

// Delete 删除资源
func (r *Resource) Delete(ctx context.Context, id int64) error {
	_, err := sqliteDB.NewDelete().
		Model((*Resource)(nil)).
		Where("id = ?", id).
		Exec(ctx)
	return err
}

// NewResource 创建新的资源实例
func NewResource(url, title, content, resourceType string) *Resource {
	return &Resource{
		URL:     url,
		Title:   title,
		Content: content,
		Type:    resourceType,
	}
}

// GetList 分页查询资源列表
func (r *Resource) GetList(ctx context.Context, page, size int, title, resourceType string) (*ResourceList, error) {
	// 构建查询
	query := sqliteDB.NewSelect().Model((*Resource)(nil))

	// 添加条件
	if title != "" {
		query = query.Where("title LIKE ?", "%"+title+"%")
	}
	if resourceType != "" {
		query = query.Where("type = ?", resourceType)
	}

	// 获取总记录数
	total, err := query.Count(ctx)
	if err != nil {
		return nil, err
	}

	// 分页查询
	var resources []*Resource
	err = query.
		Order("created_at DESC").
		Offset((page-1)*size).
		Limit(size).
		Scan(ctx, &resources)
	if err != nil {
		return nil, err
	}

	return &ResourceList{
		Total: int64(total),
		List:  resources,
	}, nil
}

func (r *Resource) Get(ctx context.Context, id int64) (*Resource, error) {
	var resource Resource
	err := sqliteDB.NewSelect().Model(&resource).Where("id = ?", id).Scan(ctx)
	return &resource, err
}
