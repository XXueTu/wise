package model

import (
	"context"

	"github.com/uptrace/bun"
	"github.com/zeromicro/go-zero/core/logx"
)

var _ TagGen = (*TagsModel)(nil)

type TagsModel struct {
	db *bun.DB
}

func NewTagsModel(db *bun.DB) *TagsModel {
	return &TagsModel{
		db: db,
	}
}

// TableName 返回表名
func (m *TagsModel) TableName() string {
	return "tags"
}

func (m *TagsModel) InitData() {
	tags := []*Tags{
		{
			Uid:         "default",
			Name:        "默认",
			Description: "默认标签",
			Color:       "red",
			Icon:        "icon",
		},
	}
	for _, tag := range tags {
		// 判断是否存在
		exist, err := m.db.NewSelect().Model((*Tags)(nil)).
			Where("uid = ?", tag.Uid).
			Exists(context.Background())
		if err != nil {
			logx.Error("InitData error", err)
			continue
		}
		if exist {
			continue
		}
		err = m.Create(context.Background(), tag)
		if err != nil {
			logx.Error("InitData error", err)
		}
	}
}

// Create 创建标签
func (m *TagsModel) Create(ctx context.Context, tag *Tags) error {
	_, err := m.db.NewInsert().Model(tag).Exec(ctx)
	if err != nil {
		logx.Error("Create error", err)
	}
	return err
}

// Update 更新标签
func (m *TagsModel) Update(ctx context.Context, tag *Tags) error {
	_, err := m.db.NewUpdate().
		Model(tag).
		WherePK().
		Exec(ctx)
	if err != nil {
		logx.Error("Update error", err)
	}
	return err
}

// Delete 删除模型
func (m *TagsModel) Delete(ctx context.Context, id int64) error {
	_, err := m.db.NewDelete().
		Model((*Tags)(nil)).
		Where("id = ?", id).
		Exec(ctx)
	if err != nil {
		logx.Error("Delete error", err)
	}
	return err
}

func (m *TagsModel) Get(ctx context.Context, id int64) (*Tags, error) {
	var tag Tags
	err := m.db.NewSelect().Model(&tag).Where("id = ?", id).Scan(ctx)
	return &tag, err
}

func (m *TagsModel) GetUid(ctx context.Context, uid string) (*Tags, error) {
	var tag Tags
	err := m.db.NewSelect().Model(&tag).Where("uid = ?", uid).Scan(ctx)
	return &tag, err
}

func (m *TagsModel) GetUids(ctx context.Context, uids []string) ([]*Tags, error) {
	var tags []*Tags
	err := m.db.NewSelect().Model(&tags).Where("uid IN (?)", bun.In(uids)).Scan(ctx)
	return tags, err
}

func (m *TagsModel) GetName(ctx context.Context, name string) (*Tags, error) {
	var tag Tags
	err := m.db.NewSelect().Model(&tag).Where("name = ?", name).Scan(ctx)
	return &tag, err
}

// TagsList 标签列表返回结构
type TagsList struct {
	Total int64   `json:"total"` // 总记录数
	List  []*Tags `json:"list"`  // 标签列表
}

// GetList 分页查询标签列表
func (m *TagsModel) GetList(ctx context.Context, page, size int64, name string) (*TagsList, error) {
	// 构建查询
	query := m.db.NewSelect().Model((*Tags)(nil))

	// 添加条件
	if name != "" {
		query = query.Where("name LIKE ?", "%"+name+"%")
	}

	// 获取总记录数
	total, err := query.Count(ctx)
	if err != nil {
		logx.Error("GetList total error", err)
		return nil, err
	}

	// 分页查询
	var tags []*Tags
	err = query.
		Order("created_at DESC").
		Offset(int((page-1)*size)).
		Limit(int(size)).
		Scan(ctx, &tags)
	if err != nil {
		logx.Error("GetList scan error", err)
		return nil, err
	}

	return &TagsList{
		Total: int64(total),
		List:  tags,
	}, nil
}
