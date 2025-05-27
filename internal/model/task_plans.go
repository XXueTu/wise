package model

import (
	"context"

	"github.com/uptrace/bun"
	"github.com/zeromicro/go-zero/core/logx"
)

var _ TaskPlansGen = (*TaskPlansModel)(nil)

type TaskPlansModel struct {
	db *bun.DB
}

func NewTaskPlansModel(db *bun.DB) *TaskPlansModel {
	return &TaskPlansModel{
		db: db,
	}
}

// TableName 返回表名
func (m *TaskPlansModel) TableName() string {
	return "task_plans"
}

func (m *TaskPlansModel) InitData() {

}

// Create 创建标签
func (m *TaskPlansModel) Create(ctx context.Context, taskPlans *TaskPlans) error {
	_, err := m.db.NewInsert().Model(taskPlans).Exec(ctx)
	if err != nil {
		logx.Error("Create error", err)
	}
	return err
}

// CreateBatch 批量创建标签
func (m *TaskPlansModel) CreateBatch(ctx context.Context, taskPlans []TaskPlans) error {
	_, err := m.db.NewInsert().Model(&taskPlans).Exec(ctx)
	if err != nil {
		logx.Error("CreateBatch error", err)
	}
	return err
}

// Update 更新标签
func (m *TaskPlansModel) Update(ctx context.Context, taskPlans *TaskPlans) error {
	_, err := m.db.NewUpdate().
		Model(taskPlans).
		WherePK().
		Exec(ctx)
	if err != nil {
		logx.Error("Update error", err)
	}
	return err
}

// Delete 删除模型
func (m *TaskPlansModel) Delete(ctx context.Context, id int64) error {
	_, err := m.db.NewDelete().
		Model((*TaskPlans)(nil)).
		Where("id = ?", id).
		Exec(ctx)
	if err != nil {
		logx.Error("Delete error", err)
	}
	return err
}

func (m *TaskPlansModel) Get(ctx context.Context, id int64) (*TaskPlans, error) {
	var taskPlans TaskPlans
	err := m.db.NewSelect().Model(&taskPlans).Where("id = ?", id).Scan(ctx)
	return &taskPlans, err
}

func (m *TaskPlansModel) GetInitByTid(ctx context.Context, tid string) ([]*TaskPlans, error) {
	var taskPlans []*TaskPlans
	err := m.db.NewSelect().Model(&taskPlans).Where("tid = ?", tid).Where("status = ?", TaskPlanStatusInit).Scan(ctx)
	return taskPlans, err
}
