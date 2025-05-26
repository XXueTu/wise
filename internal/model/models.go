package model

import (
	"context"
	"fmt"
	"os"

	"github.com/uptrace/bun"
	"github.com/zeromicro/go-zero/core/logx"
)

var _ ModelsGen = (*ModelsModel)(nil)

func NewModelsModel(db *bun.DB) *ModelsModel {
	return &ModelsModel{
		db: db,
	}
}

type ModelsModel struct {
	db *bun.DB
}

// TableName 返回表名
func (m *ModelsModel) TableName() string {
	return "models"
}

func (m *ModelsModel) InitData() {
	apiKey := os.Getenv("DEFAULT_API_KEY")
	models := []*Models{
		{
			BaseUrl:       "https://ark.cn-beijing.volces.com/api/v3",
			Config:        fmt.Sprintf("{\"apiKey\":\"%s\"}", apiKey),
			Type:          "doubao",
			ModelName:     "豆包1.5",
			ModelRealName: "doubao-1-5-pro-32k-250115",
			Status:        "active",
			Tag:           "function",
		},
	}
	for _, model := range models {
		// 判断是否存在
		exist, err := m.db.NewSelect().Model((*Models)(nil)).
			Where("base_url = ?", model.BaseUrl).
			Where("model_name = ?", model.ModelName).
			Where("model_real_name = ?", model.ModelRealName).
			Exists(context.Background())
		if err != nil {
			logx.Error("InitData error", err)
			continue
		}
		if exist {
			continue
		}
		err = m.Create(context.Background(), model)
		if err != nil {
			logx.Error("InitData error", err)
		}
	}
}

// Create 创建模型
func (m *ModelsModel) Create(ctx context.Context, model *Models) error {
	_, err := m.db.NewInsert().Model(model).Exec(ctx)
	if err != nil {
		logx.Error("Create error", err)
	}
	return err
}

// Update 更新模型
func (m *ModelsModel) Update(ctx context.Context, model *Models) error {
	_, err := m.db.NewUpdate().
		Model(model).
		WherePK().
		Exec(ctx)
	if err != nil {
		logx.Error("Update error", err)
	}
	return err
}

// Delete 删除模型
func (m *ModelsModel) Delete(ctx context.Context, id int64) error {
	_, err := m.db.NewDelete().
		Model((*Models)(nil)).
		Where("id = ?", id).
		Exec(ctx)
	if err != nil {
		logx.Error("Delete error", err)
	}
	return err
}

func (m *ModelsModel) Get(ctx context.Context, id int64) (*Models, error) {
	var model Models
	err := m.db.NewSelect().Model(&model).Where("id = ?", id).Scan(ctx)
	return &model, err
}

// ModelsList 模型列表返回结构
type ModelsList struct {
	Total int64     `json:"total"` // 总记录数
	List  []*Models `json:"list"`  // 模型列表
}

// GetList 分页查询模型列表
func (m *ModelsModel) GetList(ctx context.Context, page, size int64, modelType string, tag []string, status, modelName string) (*ModelsList, error) {
	// 构建查询
	query := m.db.NewSelect().Model((*Models)(nil))

	// 添加条件
	if modelType != "" {
		query = query.Where("type = ?", modelType)
	}
	if len(tag) > 0 {
		for _, t := range tag {
			query = query.Where("tag LIKE ?", "%"+t+"%")
		}
	}
	if status != "" {
		query = query.Where("status = ?", status)
	}
	if modelName != "" {
		query = query.Where("model_name LIKE ?", "%"+modelName+"%")
	}

	// 获取总记录数
	total, err := query.Count(ctx)
	if err != nil {
		logx.Error("GetList total error", err)
		return nil, err
	}

	// 分页查询
	var models []*Models
	err = query.
		Order("created_at DESC").
		Offset(int((page-1)*size)).
		Limit(int(size)).
		Scan(ctx, &models)
	if err != nil {
		logx.Error("GetList scan error", err)
		return nil, err
	}

	return &ModelsList{
		Total: int64(total),
		List:  models,
	}, nil
}
