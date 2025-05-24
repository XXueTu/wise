package model

import (
	"context"
	"time"

	"github.com/uptrace/bun"
	"github.com/zeromicro/go-zero/core/logx"
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

type ResourceModel struct {
	db *bun.DB
}

func NewResourceModel(db *bun.DB) *ResourceModel {
	return &ResourceModel{
		db: db,
	}
}

// ResourceList 资源列表返回结构
type ResourceList struct {
	Total int64       `json:"total"` // 总记录数
	List  []*Resource `json:"list"`  // 资源列表
}

// TableName 返回表名
func (r *ResourceModel) TableName() string {
	return "resources"
}

// Create 创建资源
func (r *ResourceModel) Create(ctx context.Context, resource *Resource) error {
	_, err := r.db.NewInsert().Model(resource).Exec(ctx)
	if err != nil {
		logx.Error("Create error", err)
	}
	return err
}

// GetByURL 根据URL获取资源
func (r *ResourceModel) GetByURL(ctx context.Context, url string) (*Resource, error) {
	resource := new(Resource)
	err := r.db.NewSelect().
		Model(resource).
		Where("url = ?", url).
		Scan(ctx)
	if err != nil {
		logx.Error("GetByURL error", err)
		return nil, err
	}
	return resource, nil
}

// Update 更新资源
func (r *ResourceModel) Update(ctx context.Context, resource *Resource) error {
	_, err := r.db.NewUpdate().
		Model(resource).
		WherePK().
		Exec(ctx)
	if err != nil {
		logx.Error("Update error", err)
	}
	return err
}

// Delete 删除资源
func (r *ResourceModel) Delete(ctx context.Context, id int64) error {
	_, err := r.db.NewDelete().
		Model((*Resource)(nil)).
		Where("id = ?", id).
		Exec(ctx)
	if err != nil {
		logx.Error("Delete error", err)
	}
	return err
}

// GetList 分页查询资源列表
func (r *ResourceModel) GetList(ctx context.Context, page, size int, resourceType, title string) (*ResourceList, error) {
	// 构建查询
	query := r.db.NewSelect().Model((*Resource)(nil))

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
		logx.Error("GetList total error", err)
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
		logx.Error("GetList scan error", err)
		return nil, err
	}

	return &ResourceList{
		Total: int64(total),
		List:  resources,
	}, nil
}

func (r *ResourceModel) Get(ctx context.Context, id int64) (*Resource, error) {
	var resource Resource
	err := r.db.NewSelect().Model(&resource).Where("id = ?", id).Scan(ctx)
	if err != nil {
		logx.Error("Get error", err)
	}
	return &resource, err
}
