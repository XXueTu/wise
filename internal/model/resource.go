package model

import (
	"context"
	"time"

	"github.com/uptrace/bun"
)

type Resource struct {
	bun.BaseModel `bun:"table:resource"`

	ID          int64     `bun:",pk,autoincrement"`
	URL         string    `bun:"url"`
	Title       string    `bun:"title"`
	Description string    `bun:"description"`
	Tag         string    `bun:"tag"`
	CreatedAt   time.Time `bun:"created_at"`
	UpdatedAt   time.Time `bun:"updated_at"`
}

func ResourceModel(id int64, url string, title string, description string, tag string) *Resource {

	return &Resource{
		ID:          id,
		URL:         url,
		Title:       title,
		Description: description,
		Tag:         tag,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}
}

func (r *Resource) TableName() string {
	return "resource"
}

func (r *Resource) Create(ctx context.Context, resource *Resource) error {
	resource.ID = 0
	resource.CreatedAt = time.Now()
	resource.UpdatedAt = time.Now()
	_, err := sqliteDB.NewInsert().Model(resource).Exec(ctx)
	return err
}

func (r *Resource) Get(ctx context.Context, id int64) (*Resource, error) {
	var resource Resource
	err := sqliteDB.NewSelect().Model(&resource).Where("id = ?", id).Scan(ctx)
	return &resource, err
}

func (r *Resource) Update(ctx context.Context, resource *Resource) error {
	_, err := sqliteDB.NewUpdate().Model(resource).Where("id = ?", resource.ID).Exec(ctx)
	return err
}

func (r *Resource) Delete(ctx context.Context, id int64) error {
	_, err := sqliteDB.NewDelete().Model(r).Where("id = ?", id).Exec(ctx)
	return err
}
