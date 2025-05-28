package model

import (
	"context"

	"github.com/uptrace/bun"
	"github.com/zeromicro/go-zero/core/logx"
)

var _ TasksGen = (*TasksModel)(nil)

type TasksModel struct {
	db *bun.DB
}

func NewTasksModel(db *bun.DB) *TasksModel {
	return &TasksModel{
		db: db,
	}
}

// TableName 返回表名
func (m *TasksModel) TableName() string {
	return "tasks"
}

func (m *TasksModel) InitData() {

}

// Create 创建标签
func (m *TasksModel) Create(ctx context.Context, task *Tasks) error {
	_, err := m.db.NewInsert().Model(task).Exec(ctx)
	if err != nil {
		logx.Error("Create error", err)
	}
	return err
}

// Update 更新标签
func (m *TasksModel) Update(ctx context.Context, task *Tasks) error {
	_, err := m.db.NewUpdate().
		Model(task).
		WherePK().
		Exec(ctx)
	if err != nil {
		logx.Error("Update error", err)
	}
	return err
}

// Delete 删除模型
func (m *TasksModel) Delete(ctx context.Context, id int64) error {
	_, err := m.db.NewDelete().
		Model((*Tasks)(nil)).
		Where("id = ?", id).
		Exec(ctx)
	if err != nil {
		logx.Error("Delete error", err)
	}
	return err
}

func (m *TasksModel) Get(ctx context.Context, id int64) (*Tasks, error) {
	var task Tasks
	err := m.db.NewSelect().Model(&task).Where("id = ?", id).Scan(ctx)
	return &task, err
}

func (m *TasksModel) GetStatus(ctx context.Context, status string) ([]*Tasks, error) {
	var tasks []*Tasks
	err := m.db.NewSelect().Model(&tasks).Where("status = ?", status).Scan(ctx)
	return tasks, err
}

func (m *TasksModel) GetStatusLimit(ctx context.Context, status string, limit int) ([]*Tasks, error) {
	var tasks []*Tasks
	err := m.db.NewSelect().Model(&tasks).Where("status = ?", status).Limit(limit).Scan(ctx)
	return tasks, err
}

func (m *TasksModel) GetByTid(ctx context.Context, tid string) (*Tasks, error) {
	var task Tasks
	err := m.db.NewSelect().Model(&task).Where("tid = ?", tid).Scan(ctx)
	return &task, err
}

// TagsList 标签列表返回结构
type TasksList struct {
	Total int64    `json:"total"` // 总记录数
	List  []*Tasks `json:"list"`  // 标签列表
}

func (m *TasksModel) GetPage(ctx context.Context, page int64, pageSize int64, name string, status string, types string) (*TasksList, error) {
	// 构建查询
	query := m.db.NewSelect().Model((*Tasks)(nil))

	// 添加条件
	if name != "" {
		query = query.Where("name LIKE ?", "%"+name+"%")
	}
	if status != "" {
		query = query.Where("status = ?", status)
	}
	if types != "" {
		query = query.Where("types = ?", types)
	}

	// 获取总记录数
	total, err := query.Count(ctx)
	if err != nil {
		logx.Error("GetPage total error", err)
		return nil, err
	}

	// 分页查询
	var tasks []*Tasks
	err = query.
		Order("created_at DESC").
		Offset(int((page-1)*pageSize)).
		Limit(int(pageSize)).
		Scan(ctx, &tasks)
	if err != nil {
		logx.Error("GetPage scan error", err)
		return nil, err
	}

	return &TasksList{
		Total: int64(total),
		List:  tasks,
	}, nil
}

// UpdateState 更新任务状态
func (m *TasksModel) UpdateState(ctx context.Context, tid string, state string, result string) error {
	task, err := m.GetByTid(ctx, tid)
	if err != nil {
		return err
	}
	task.CurrentState = state
	task.Result = result
	return m.Update(ctx, task)
}

// UpdateStateAndStep 更新任务状态和步骤
func (m *TasksModel) UpdateStateAndStep(ctx context.Context, tid string, state string, step int64, result string) error {
	task, err := m.GetByTid(ctx, tid)
	if err != nil {
		return err
	}
	task.CurrentState = state
	task.CurrentStep = step
	task.Result = result
	return m.Update(ctx, task)
}

// UpdateStatus 更新任务状态
func (m *TasksModel) UpdateStatus(ctx context.Context, tid string, status string, error string) error {
	task, err := m.GetByTid(ctx, tid)
	if err != nil {
		return err
	}
	task.Status = status
	task.Error = error
	return m.Update(ctx, task)
}
