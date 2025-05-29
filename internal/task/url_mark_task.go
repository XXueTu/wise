package task

import (
	"context"
	"encoding/json"
	"errors"
	"time"

	"github.com/zeromicro/go-zero/core/logx"

	"github.com/XXueTu/wise/internal/model"
	"github.com/XXueTu/wise/internal/svc"
	"github.com/XXueTu/wise/pkg/spiders"
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

// UrlMarkState URL标记任务状态
type UrlMarkState struct {
	Index int64  // 状态索引
	Code  string // 状态代码
	Name  string // 状态名称
}

type UrlMarkTaskArgs struct {
	Url string `json:"url"`
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
	Start:  UrlMarkState{Index: 1, Code: "start", Name: "开始"},
	Check:  UrlMarkState{Index: 2, Code: "check", Name: "检查"},
	Read:   UrlMarkState{Index: 3, Code: "read", Name: "读取"},
	Split:  UrlMarkState{Index: 4, Code: "split", Name: "拆分"},
	Mark:   UrlMarkState{Index: 5, Code: "mark", Name: "标记"},
	Vector: UrlMarkState{Index: 6, Code: "vector", Name: "向量化"},
	Index:  UrlMarkState{Index: 7, Code: "index", Name: "索引"},
}

// UrlMarkTask URL标记任务
type UrlMarkTask struct {
	*StateMachine[*UrlMarkTask]
	svc *svc.ServiceContext
}

// NewUrlMarkTask 创建URL标记任务
func NewUrlMarkTask(svc *svc.ServiceContext) *UrlMarkTask {
	task := &UrlMarkTask{
		StateMachine: NewStateMachine[*UrlMarkTask](svc),
		svc:          svc,
	}

	// 注册状态
	task.RegisterState(UrlMarkStates.Start.Code, UrlMarkStates.Start.Name, task.Start)
	task.RegisterState(UrlMarkStates.Check.Code, UrlMarkStates.Check.Name, task.Check)
	task.RegisterState(UrlMarkStates.Read.Code, UrlMarkStates.Read.Name, task.Read)
	task.RegisterState(UrlMarkStates.Split.Code, UrlMarkStates.Split.Name, task.Split)
	task.RegisterState(UrlMarkStates.Mark.Code, UrlMarkStates.Mark.Name, task.Mark)
	task.RegisterState(UrlMarkStates.Vector.Code, UrlMarkStates.Vector.Name, task.Vector)
	task.RegisterState(UrlMarkStates.Index.Code, UrlMarkStates.Index.Name, task.Index)

	return task
}

// Execute 执行任务
func (t *UrlMarkTask) Execute(ctx context.Context, task *model.Tasks) error {
	return t.Run(ctx, task.Tid, task.Params, task.CurrentState)
}

// GetTaskType 获取任务类型
func (t *UrlMarkTask) GetTaskType() string {
	return UrlMarkTaskTypes
}

// GetMaxRetries 获取最大重试次数
func (t *UrlMarkTask) GetMaxRetries() int {
	return 3
}

// GetRetryInterval 获取重试间隔
func (t *UrlMarkTask) GetRetryInterval() time.Duration {
	return 5 * time.Second
}

// CreateTaskPlans 创建任务计划
func (t *UrlMarkTask) CreateTaskPlans(ctx context.Context, task *model.Tasks) error {
	plans := make([]model.TaskPlans, 0)
	for i, state := range t.states {
		plans = append(plans, model.TaskPlans{
			Tid:      task.Tid,
			Pid:      task.Tid + "_" + state.Code,
			Types:    state.Code,
			Name:     state.Name,
			Index:    int64(i),
			Status:   model.TaskPlanStatusInit,
			Params:   task.Params,
			Result:   "{}",
			Duration: 0,
			Error:    "{}",
		})
	}
	return t.svc.TaskPlansModel.CreateBatch(ctx, plans)
}

// CreateUrlMarkTask 创建URL标记任务
func CreateUrlMarkTask(ctx context.Context, svc *svc.ServiceContext, args Args) error {
	task := NewUrlMarkTask(svc)
	return task.CreateTask(ctx, args, "链接标记任务", UrlMarkTaskTypes)
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
	Start(ctx context.Context, args Args) (Args, State[*UrlMarkTask], error)
	Check(ctx context.Context, args Args) (Args, State[*UrlMarkTask], error)
	Read(ctx context.Context, args Args) (Args, State[*UrlMarkTask], error)
	Split(ctx context.Context, args Args) (Args, State[*UrlMarkTask], error)
	Mark(ctx context.Context, args Args) (Args, State[*UrlMarkTask], error)
	Vector(ctx context.Context, args Args) (Args, State[*UrlMarkTask], error)
	Index(ctx context.Context, args Args) (Args, State[*UrlMarkTask], error)
}

// Start 开始状态
func (t *UrlMarkTask) Start(ctx context.Context, args Args) (Args, State[*UrlMarkTask], error) {
	logx.Info("url mark task start")
	args.Result = args.Params
	return args, t.Check, nil
}

type UrlMarkTaskCheckArgs struct {
	Url   string `json:"url"`
	Types string `json:"types"`
}

// Check 检查状态
func (t *UrlMarkTask) Check(ctx context.Context, args Args) (Args, State[*UrlMarkTask], error) {
	logx.Info("url mark task check")
	startArgs := UrlMarkTaskArgs{}
	checkArgs := UrlMarkTaskCheckArgs{}
	_ = json.Unmarshal([]byte(args.Params), &startArgs)

	types := spiders.NewPattern().GetPatternTypes(startArgs.Url)
	if types == "unknown" {
		checkArgs.Types = "unknown"
		result, _ := json.Marshal(checkArgs)
		args.Result = string(result)
		return args, nil, errors.New("不支持的链接类型")
	}

	checkArgs.Url = startArgs.Url
	checkArgs.Types = types
	result, _ := json.Marshal(checkArgs)
	args.Result = string(result)
	return args, t.Read, nil
}

type UrlMarkTaskReadArgs struct {
	Title   string `json:"title"`
	Content string `json:"content"`
}

// Read 读取状态
func (t *UrlMarkTask) Read(ctx context.Context, args Args) (Args, State[*UrlMarkTask], error) {
	logx.Info("url mark task read")
	checkArgs := UrlMarkTaskCheckArgs{}
	readArgs := UrlMarkTaskReadArgs{}
	_ = json.Unmarshal([]byte(args.Params), &checkArgs)

	// 检查是否是微信公众号链接
	title, content, err := spiders.NewPattern().GetPattern(checkArgs.Url)
	if err != nil {
		return args, t.Read, err
	}

	err = t.svc.ResourceModel.Create(ctx, &model.Resource{
		URL:     checkArgs.Url,
		Title:   title,
		Content: content,
		Type:    checkArgs.Types,
	})
	if err != nil {
		return args, t.Read, err
	}

	readArgs.Title = title
	readArgs.Content = content
	result, _ := json.Marshal(readArgs)
	args.Result = string(result)
	return args, t.Split, nil
}

// Split 拆分状态
func (t *UrlMarkTask) Split(ctx context.Context, args Args) (Args, State[*UrlMarkTask], error) {
	logx.Info("url mark task split")
	// TODO: 实现内容拆分逻辑
	return args, t.Mark, nil
}

// Mark 标记状态
func (t *UrlMarkTask) Mark(ctx context.Context, args Args) (Args, State[*UrlMarkTask], error) {
	logx.Info("url mark task mark")
	// TODO: 实现内容标记逻辑
	return args, t.Vector, nil
}

// Vector 向量化状态
func (t *UrlMarkTask) Vector(ctx context.Context, args Args) (Args, State[*UrlMarkTask], error) {
	logx.Info("url mark task vector")
	// TODO: 实现向量化逻辑
	return args, t.Index, nil
}

// Index 索引状态
func (t *UrlMarkTask) Index(ctx context.Context, args Args) (Args, State[*UrlMarkTask], error) {
	logx.Info("url mark task index")
	// TODO: 实现索引逻辑
	return args, nil, nil
}
