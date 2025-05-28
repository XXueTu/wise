package task

import (
	"context"
	"encoding/json"
	"errors"
	"sync"
	"time"

	"github.com/google/uuid"
	"github.com/zeromicro/go-zero/core/logx"
	"golang.org/x/sync/semaphore"

	"github.com/XXueTu/wise/internal/model"
	"github.com/XXueTu/wise/internal/svc"
)

// State 状态机状态函数
type State[F any, T any] func(ctx context.Context, args F) (F, State[F, T], error)

// StateInfo 状态信息
type StateInfo struct {
	Code string // 状态代码
	Name string // 状态名称
}

// TaskExecutor 任务执行器接口
type TaskExecutor interface {
	// Execute 执行任务
	Execute(ctx context.Context, task *model.Tasks) error
	// GetTaskType 获取任务类型
	GetTaskType() string
	// GetMaxRetries 获取最大重试次数
	GetMaxRetries() int
	// GetRetryInterval 获取重试间隔
	GetRetryInterval() time.Duration
	// CreateTaskPlans 创建任务计划
	CreateTaskPlans(ctx context.Context, task *model.Tasks) error
}

// StateMachine 状态机基类
type StateMachine[F any, T any] struct {
	svc      *svc.ServiceContext
	stateMap map[string]State[F, T]
	states   []StateInfo
}

// NewStateMachine 创建状态机
func NewStateMachine[F any, T any](svc *svc.ServiceContext) *StateMachine[F, T] {
	return &StateMachine[F, T]{
		svc:      svc,
		stateMap: make(map[string]State[F, T]),
		states:   make([]StateInfo, 0),
	}
}

// RegisterState 注册状态
func (m *StateMachine[F, T]) RegisterState(code string, name string, state State[F, T]) {
	m.stateMap[code] = state
	m.states = append(m.states, StateInfo{
		Code: code,
		Name: name,
	})
}

// Run 运行状态机
func (m *StateMachine[F, T]) Run(ctx context.Context, tid string, params string, startState string) error {
	current := m.stateMap[startState]
	if current == nil {
		return errors.New("初始状态不存在: " + startState)
	}

	stateCode := startState
	var args F
	var err error

	// 解析参数
	if err := json.Unmarshal([]byte(params), &args); err != nil {
		logx.Errorf("解析任务参数失败: %v", err)
		return errors.New("解析任务参数失败" + err.Error())
	}

	// 查询任务执行计划
	plans, err := m.svc.TaskPlansModel.GetByTid(ctx, tid)
	if err != nil {
		logx.Errorf("查询任务执行计划失败: %v", err)
		return errors.New("查询任务执行计划失败" + err.Error())
	}
	if len(plans) == 0 {
		logx.Errorf("任务执行计划为空: %v", tid)
		return errors.New("任务执行计划为空")
	}

	plansMap := make(map[string]*model.TaskPlans)
	for _, plan := range plans {
		plansMap[plan.Types] = plan
		logx.Infof("任务执行计划: types: %s, tid: %s", plan.Types, plan.Tid)
	}

	// 执行状态机
	for {
		logx.Infof("执行状态机: %s", stateCode)
		startTime := time.Now()

		// 更新任务状态
		plan := plansMap[stateCode]
		if plan == nil {
			logx.Errorf("当前状态不存在: %s", stateCode)
			return errors.New("当前状态不存在" + stateCode)
		}

		// 更新计划状态为运行中
		plan.Status = model.TaskPlanStatusRunning
		if err := m.svc.TaskPlansModel.Update(ctx, plan); err != nil {
			logx.Errorf("更新任务计划状态失败: %v", err)
			return err
		}

		// 执行状态
		args, current, err = current(ctx, args)
		if err != nil {
			logx.Errorf("执行状态失败: %+v", err)
			// 更新计划状态为失败
			plan.Status = model.TaskPlanStatusFailed
			plan.Error = err.Error()
			if err := m.svc.TaskPlansModel.Update(ctx, plan); err != nil {
				logx.Errorf("更新任务计划状态失败: %v", err)
			}
			return errors.New("执行状态失败" + err.Error())
		}

		// 更新计划状态为成功
		plan.Status = model.TaskPlanStatusSuccess
		plan.Duration = time.Since(startTime).Milliseconds()
		if err := m.svc.TaskPlansModel.Update(ctx, plan); err != nil {
			logx.Errorf("更新任务计划状态失败: %v", err)
			return err
		}

		// 更新任务状态
		if err := m.svc.TasksModel.UpdateState(ctx, tid, plan.Types, "{}"); err != nil {
			logx.Errorf("更新任务状态失败: %v", err)
			return err
		}

		// 检查是否结束
		if current == nil {
			logx.Infof("任务执行结束: %s", tid)
			break
		}

		// 获取下一个状态
		nextState := ""
		if argsWithNext, ok := any(args).(interface{ GetNextState() string }); ok {
			nextState = argsWithNext.GetNextState()
		}

		// 验证下一个状态
		if nextState == "" {
			logx.Infof("没有下一个状态，任务执行结束: %s", tid)
			break
		}

		// 验证状态是否存在
		if _, exists := m.stateMap[nextState]; !exists {
			logx.Errorf("下一个状态不存在: %s", nextState)
			return errors.New("下一个状态不存在" + nextState)
		}

		// 验证状态是否在计划中
		if _, exists := plansMap[nextState]; !exists {
			logx.Errorf("下一个状态不在计划中: %s", nextState)
			return errors.New("下一个状态不在计划中" + nextState)
		}

		stateCode = nextState
	}

	return nil
}

// CreateTask 创建任务
func (m *StateMachine[F, T]) CreateTask(ctx context.Context, args F, taskName string, taskType string) error {
	argsJson, err := json.Marshal(args)
	if err != nil {
		return err
	}

	tid := uuid.New().String()
	// 创建任务
	err = m.svc.TasksModel.Create(ctx, &model.Tasks{
		Tid:          tid,
		Name:         taskName,
		Types:        taskType,
		Status:       model.TaskStatusInit,
		CurrentState: m.states[0].Code,
		TotalSteps:   int64(len(m.states)),
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
	for i, state := range m.states {
		plans = append(plans, model.TaskPlans{
			Tid:      tid,
			Pid:      uuid.New().String(),
			Types:    state.Code,
			Name:     state.Name,
			Index:    int64(i),
			Status:   model.TaskPlanStatusInit,
			Params:   string(argsJson),
			Result:   "{}",
			Duration: 0,
			Error:    "{}",
		})
	}
	return m.svc.TaskPlansModel.CreateBatch(ctx, plans)
}

// TaskScheduler 任务调度器
type TaskScheduler struct {
	svc          *svc.ServiceContext
	workerPool   *semaphore.Weighted
	maxWorkers   int64
	taskTimeout  time.Duration
	scanInterval time.Duration
	stopChan     chan struct{}
	wg           sync.WaitGroup
	taskCtxs     sync.Map                // 存储任务上下文，用于取消任务
	taskRegistry map[string]TaskExecutor // 任务执行器注册表
}

// NewTaskScheduler 创建任务调度器
func NewTaskScheduler(svc *svc.ServiceContext) *TaskScheduler {
	return &TaskScheduler{
		svc:          svc,
		maxWorkers:   int64(svc.Config.Task.PoolSize),
		taskTimeout:  20 * time.Minute,
		scanInterval: 10 * time.Second,
		stopChan:     make(chan struct{}),
		workerPool:   semaphore.NewWeighted(int64(svc.Config.Task.PoolSize)),
		taskRegistry: make(map[string]TaskExecutor),
	}
}

// RegisterTaskExecutor 注册任务执行器
func (s *TaskScheduler) RegisterTaskExecutor(executor TaskExecutor) {
	s.taskRegistry[executor.GetTaskType()] = executor
}

// Start 启动任务调度器
func (s *TaskScheduler) Start() {
	go s.startScheduler()
}

// Stop 停止任务调度器
func (s *TaskScheduler) Stop() {
	close(s.stopChan)
	s.wg.Wait()
}

// CancelTask 取消任务
func (s *TaskScheduler) CancelTask(tid string) error {
	if ctx, ok := s.taskCtxs.Load(tid); ok {
		if cancel, ok := ctx.(context.CancelFunc); ok {
			cancel()
		}
	}

	task, err := s.svc.TasksModel.GetByTid(context.Background(), tid)
	if err != nil {
		return err
	}

	task.Status = model.TaskStatusCancelled
	return s.svc.TasksModel.Update(context.Background(), task)
}

func (s *TaskScheduler) startScheduler() {
	ticker := time.NewTicker(s.scanInterval)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			s.scanAndExecuteTasks()
		case <-s.stopChan:
			return
		}
	}
}

func (s *TaskScheduler) scanAndExecuteTasks() {
	ctx := context.Background()

	// 获取所有未完成的任务
	tasks, err := s.svc.TasksModel.GetStatusLimit(ctx, model.TaskStatusInit, int(s.maxWorkers))
	if err != nil {
		logx.Errorf("获取任务列表失败: %v", err)
		return
	}
	if len(tasks) == 0 {
		// 查询重试任务
		tasks, err = s.svc.TasksModel.GetStatusLimit(ctx, model.TaskStatusRetry, int(s.maxWorkers))
		if err != nil {
			logx.Errorf("获取重试任务列表失败: %v", err)
			return
		}
		if len(tasks) == 0 {
			return
		}
	}

	for _, task := range tasks {
		if !s.workerPool.TryAcquire(1) {
			logx.Debugf("工作协程池已满，等待下次调度")
			return
		}

		s.wg.Add(1)
		go func(t *model.Tasks) {
			defer s.wg.Done()
			defer s.workerPool.Release(1)
			s.executeTask(t)
		}(task)
	}
}

func (s *TaskScheduler) executeTask(task *model.Tasks) {
	ctx, cancel := context.WithTimeout(context.Background(), s.taskTimeout)
	s.taskCtxs.Store(task.Tid, cancel)
	defer func() {
		cancel()
		s.taskCtxs.Delete(task.Tid)
	}()

	logx.Infof("执行任务: %v", task.Tid)
	s.svc.TasksModel.UpdateStatus(ctx, task.Tid, model.TaskStatusRunning, "{}")

	executor, ok := s.taskRegistry[task.Types]
	if !ok {
		err := "未找到任务执行器"
		logx.Errorf("执行任务失败: %v", err)
		s.svc.TasksModel.UpdateStatus(ctx, task.Tid, model.TaskStatusFailed, err)
		return
	}

	// 执行任务
	err := executor.Execute(ctx, task)
	if err != nil {
		// 检查是否需要重试
		if task.CurrentStep < int64(executor.GetMaxRetries()) {
			task.Status = model.TaskStatusRetry
			task.CurrentStep++
			task.Error = err.Error()
			s.svc.TasksModel.Update(ctx, task)
			// 等待重试间隔后重新调度
			time.Sleep(executor.GetRetryInterval())
			logx.Infof("重试任务: %v", task.Tid)
			return
		}

		errMap := make(map[string]string)
		errMap["error"] = err.Error()
		s.svc.TasksModel.UpdateStatus(ctx, task.Tid, model.TaskStatusFailed, err.Error())
		logx.Errorf("执行任务失败: %v", err)
		return
	}

	s.svc.TasksModel.UpdateStatus(ctx, task.Tid, model.TaskStatusSuccess, "{}")
	logx.Infof("执行任务成功: %v", task.Tid)
}

// InitTaskScheduler 初始化任务调度器
func InitTaskScheduler(svc *svc.ServiceContext) *TaskScheduler {
	scheduler := NewTaskScheduler(svc)

	// 注册URL标记任务执行器
	urlMarkTask := NewUrlMarkTask(svc)
	scheduler.RegisterTaskExecutor(urlMarkTask)

	// 启动调度器
	scheduler.Start()

	return scheduler
}
