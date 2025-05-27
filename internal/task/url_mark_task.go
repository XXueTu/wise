package task

import (
	"context"
	"encoding/json"

	"github.com/google/uuid"
	"github.com/zeromicro/go-zero/core/logx"

	"github.com/XXueTu/wise/internal/model"
	"github.com/XXueTu/wise/internal/svc"
)

/*
校验URL参数
读取内容
拆分内容
标记内容 (tag,标题,描述)

构建拆分内容向量
构建全文内容索引
*/

const (
	UrlMarkTaskTypes = "url_mark"
)

type UrlMarkState struct {
	Code string // 状态代码
	Name string // 状态名称
}

var UrlMarkStates = struct {
	Start  UrlMarkState
	Check  UrlMarkState
	Read   UrlMarkState
	Split  UrlMarkState
	Mark   UrlMarkState
	Vector UrlMarkState
	Index  UrlMarkState
}{
	Start:  UrlMarkState{Code: "start", Name: "开始"},
	Check:  UrlMarkState{Code: "check", Name: "检查"},
	Read:   UrlMarkState{Code: "read", Name: "读取"},
	Split:  UrlMarkState{Code: "split", Name: "拆分"},
	Mark:   UrlMarkState{Code: "mark", Name: "标记"},
	Vector: UrlMarkState{Code: "vector", Name: "向量化"},
	Index:  UrlMarkState{Code: "index", Name: "索引"},
}

// UrlMarkTaskArgs URL标记任务参数
type UrlMarkTaskArgs struct {
	Url string `json:"url"`
	Tid string `json:"tid"`
}

var _ UrlMarkTaskFsm = (*UrlMarkTask)(nil)

// UrlMarkTask URL标记任务状态机
type UrlMarkTask struct {
	svc      *svc.ServiceContext
	stateMap map[string]State[UrlMarkTaskArgs, UrlMarkTask]
	states   []UrlMarkState
}

var urlMarkTaskInstance *UrlMarkTask

// NewUrlMarkTask 创建URL标记任务状态机
func NewUrlMarkTask(svc *svc.ServiceContext) *UrlMarkTask {
	states := []UrlMarkState{
		{
			Code: "～",
			Name: "～",
		},
		UrlMarkStates.Start,
		UrlMarkStates.Check,
		UrlMarkStates.Read,
		UrlMarkStates.Split,
		UrlMarkStates.Mark,
		UrlMarkStates.Vector,
		UrlMarkStates.Index,
	}
	task := &UrlMarkTask{
		svc:      svc,
		stateMap: make(map[string]State[UrlMarkTaskArgs, UrlMarkTask]),
		states:   states,
	}
	// 注册状态
	task.stateMap[UrlMarkStates.Start.Code] = task.Start
	task.stateMap[UrlMarkStates.Check.Code] = task.Check
	task.stateMap[UrlMarkStates.Read.Code] = task.Read
	task.stateMap[UrlMarkStates.Split.Code] = task.Split
	task.stateMap[UrlMarkStates.Mark.Code] = task.Mark
	task.stateMap[UrlMarkStates.Vector.Code] = task.Vector
	task.stateMap[UrlMarkStates.Index.Code] = task.Index
	urlMarkTaskInstance = task
	return task
}

func (t *UrlMarkTask) Run(ctx context.Context, tid string, params string) {
	current := t.stateMap[UrlMarkStates.Start.Code]
	var args UrlMarkTaskArgs
	var err error
	if err := json.Unmarshal([]byte(params), &args); err != nil {
		logx.Errorf("解析任务参数失败: %v", err)
		return
	}
	args.Tid = tid
	// 查询任务执行计划

	for {
		args, current, err = current(ctx, args)
		if err != nil {
			return
		}
		if current == nil {
			break
		}
	}
}

func CreateUrlMarkTask(ctx context.Context, args UrlMarkTaskArgs) error {
	argsJson, err := json.Marshal(args)
	if err != nil {
		return err
	}
	tid := uuid.New().String()
	// 创建任务
	err = urlMarkTaskInstance.svc.TasksModel.Create(ctx, &model.Tasks{
		Tid:          tid,
		Name:         "链接标记任务",
		Types:        UrlMarkTaskTypes,
		Status:       model.TaskStatusInit,
		CurrentState: UrlMarkStates.Start.Code,
		TotalSteps:   7,
		CurrentStep:  1,
		Params:       string(argsJson),
		Result:       "{}",
		Duration:     0,
		Error:        "{}",
		Extend:       "{}",
	})
	if err != nil {
		return err
	}
	// 批量创建计划
	plans := make([]model.TaskPlans, 0)
	for i := 1; i < len(urlMarkTaskInstance.states); i++ {
		plans = append(plans, model.TaskPlans{
			Tid:      tid,
			Pid:      uuid.New().String(),
			Types:    urlMarkTaskInstance.states[i].Code,
			Name:     urlMarkTaskInstance.states[i].Name,
			Index:    int64(i),
			Status:   model.TaskPlanStatusInit,
			Params:   string(argsJson),
			Result:   "{}",
			Duration: 0,
			Error:    "{}",
		})
	}
	err = urlMarkTaskInstance.svc.TaskPlansModel.CreateBatch(ctx, plans)
	if err != nil {
		return err
	}
	return nil
}

// Start 开始状态
func (t *UrlMarkTask) Start(ctx context.Context, args UrlMarkTaskArgs) (UrlMarkTaskArgs, State[UrlMarkTaskArgs, UrlMarkTask], error) {
	logx.Info("url mark task start")
	return args, t.Check, nil
}

// Check 检查状态
func (t *UrlMarkTask) Check(ctx context.Context, args UrlMarkTaskArgs) (UrlMarkTaskArgs, State[UrlMarkTaskArgs, UrlMarkTask], error) {
	logx.Info("url mark task check")
	// TODO: 实现URL检查逻辑
	return args, t.Read, nil
}

// Read 读取状态
func (t *UrlMarkTask) Read(ctx context.Context, args UrlMarkTaskArgs) (UrlMarkTaskArgs, State[UrlMarkTaskArgs, UrlMarkTask], error) {
	logx.Info("url mark task read")
	// TODO: 实现内容读取逻辑
	return args, t.Split, nil
}

// Split 拆分状态
func (t *UrlMarkTask) Split(ctx context.Context, args UrlMarkTaskArgs) (UrlMarkTaskArgs, State[UrlMarkTaskArgs, UrlMarkTask], error) {
	logx.Info("url mark task split")
	// TODO: 实现内容拆分逻辑
	return args, t.Mark, nil
}

// Mark 标记状态
func (t *UrlMarkTask) Mark(ctx context.Context, args UrlMarkTaskArgs) (UrlMarkTaskArgs, State[UrlMarkTaskArgs, UrlMarkTask], error) {
	logx.Info("url mark task mark")
	// TODO: 实现内容标记逻辑
	return args, t.Vector, nil
}

// Vector 向量化状态
func (t *UrlMarkTask) Vector(ctx context.Context, args UrlMarkTaskArgs) (UrlMarkTaskArgs, State[UrlMarkTaskArgs, UrlMarkTask], error) {
	logx.Info("url mark task vector")
	// TODO: 实现向量化逻辑
	return args, t.Index, nil
}

// Index 索引状态
func (t *UrlMarkTask) Index(ctx context.Context, args UrlMarkTaskArgs) (UrlMarkTaskArgs, State[UrlMarkTaskArgs, UrlMarkTask], error) {
	logx.Info("url mark task index")
	// TODO: 实现索引逻辑
	return args, nil, nil
}

// 定义任务 index 和函数的枚举
type UrlMarkTaskIndex int

const (
	UrlMarkTaskIndexStart UrlMarkTaskIndex = iota
	UrlMarkTaskIndexCheck
	UrlMarkTaskIndexRead
	UrlMarkTaskIndexSplit
	UrlMarkTaskIndexMark
	UrlMarkTaskIndexVector
	UrlMarkTaskIndexIndex
)

type UrlMarkTaskFsm interface {
	Start(ctx context.Context, args UrlMarkTaskArgs) (UrlMarkTaskArgs, State[UrlMarkTaskArgs, UrlMarkTask], error)
	Check(ctx context.Context, args UrlMarkTaskArgs) (UrlMarkTaskArgs, State[UrlMarkTaskArgs, UrlMarkTask], error)
	Read(ctx context.Context, args UrlMarkTaskArgs) (UrlMarkTaskArgs, State[UrlMarkTaskArgs, UrlMarkTask], error)
	Split(ctx context.Context, args UrlMarkTaskArgs) (UrlMarkTaskArgs, State[UrlMarkTaskArgs, UrlMarkTask], error)
	Mark(ctx context.Context, args UrlMarkTaskArgs) (UrlMarkTaskArgs, State[UrlMarkTaskArgs, UrlMarkTask], error)
	Vector(ctx context.Context, args UrlMarkTaskArgs) (UrlMarkTaskArgs, State[UrlMarkTaskArgs, UrlMarkTask], error)
	Index(ctx context.Context, args UrlMarkTaskArgs) (UrlMarkTaskArgs, State[UrlMarkTaskArgs, UrlMarkTask], error)
}
