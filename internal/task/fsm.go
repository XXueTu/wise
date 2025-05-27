package task

import (
	"context"
	"sync"
	"time"

	"github.com/zeromicro/go-zero/core/logx"
	"golang.org/x/sync/semaphore"

	"github.com/XXueTu/wise/internal/model"
	"github.com/XXueTu/wise/internal/svc"
)

// TaskActuator 任务执行器
type TaskActuator struct {
	poolSize     int
	Svc          *svc.ServiceContext
	workerPool   *semaphore.Weighted
	maxWorkers   int64
	taskTimeout  time.Duration
	scanInterval time.Duration
	stopChan     chan struct{}
	wg           sync.WaitGroup
	taskCtxs     sync.Map // 存储任务上下文，用于取消任务
	// 状态机实例
	urlMarkTaskInstance *UrlMarkTask
}

// NewTaskActuator 创建任务执行器
func NewTaskActuator(svc *svc.ServiceContext) *TaskActuator {
	// 注册状态机

	return &TaskActuator{
		poolSize:            svc.Config.Task.PoolSize,
		Svc:                 svc,
		maxWorkers:          int64(svc.Config.Task.PoolSize),
		taskTimeout:         20 * time.Minute,
		scanInterval:        10 * time.Second,
		stopChan:            make(chan struct{}),
		workerPool:          semaphore.NewWeighted(int64(svc.Config.Task.PoolSize)),
		urlMarkTaskInstance: NewUrlMarkTask(svc),
	}
}

// CancelTask 取消任务
func (a *TaskActuator) CancelTask(tid string) error {
	if ctx, ok := a.taskCtxs.Load(tid); ok {
		if cancel, ok := ctx.(context.CancelFunc); ok {
			cancel()
		}
	}

	task, err := a.Svc.TasksModel.GetByTid(context.Background(), tid)
	if err != nil {
		return err
	}

	task.Status = model.TaskStatusCancelled
	return a.Svc.TasksModel.Update(context.Background(), task)
}

// Start 启动任务执行器
func (a *TaskActuator) Start() {
	go a.startScheduler()
}

// Stop 停止任务执行器
func (a *TaskActuator) Stop() {
	close(a.stopChan)
	a.wg.Wait()
}

func (a *TaskActuator) startScheduler() {
	ticker := time.NewTicker(a.scanInterval)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			a.scanAndExecuteTasks()
		case <-a.stopChan:
			return
		}
	}
}

func (a *TaskActuator) scanAndExecuteTasks() {
	ctx := context.Background()

	// 获取所有未完成的任务
	tasks, err := a.Svc.TasksModel.GetStatusLimit(ctx, model.TaskStatusInit, a.poolSize)
	if err != nil {
		logx.Errorf("获取任务列表失败: %v", err)
		return
	}
	if len(tasks) == 0 {
		// 查询重试任务
		tasks, err = a.Svc.TasksModel.GetStatusLimit(ctx, model.TaskStatusRetry, a.poolSize)
		if err != nil {
			logx.Errorf("获取重试任务列表失败: %v", err)
			return
		}
		if len(tasks) == 0 {
			return
		}
	}

	for _, task := range tasks {
		if !a.workerPool.TryAcquire(1) {
			logx.Debugf("工作协程池已满，等待下次调度")
			return
		}

		a.wg.Add(1)
		go func(t *model.Tasks) {
			defer a.wg.Done()
			defer a.workerPool.Release(1)
			a.executeTask(t)
		}(task)
	}
}

func (a *TaskActuator) executeTask(task *model.Tasks) {
	ctx, cancel := context.WithTimeout(context.Background(), a.taskTimeout)
	a.taskCtxs.Store(task.Tid, cancel)
	defer func() {
		cancel()
		a.taskCtxs.Delete(task.Tid)
	}()

	if err := a.Execute(ctx, task); err != nil {
		logx.Errorf("执行任务失败: %v", err)
	}
}

// Execute 执行状态机
func (a *TaskActuator) Execute(ctx context.Context, task *model.Tasks) error {
	switch task.Types {
	case UrlMarkTaskTypes:
		a.urlMarkTaskInstance.Run(ctx, task.Tid, task.Params)
	}
	return nil
}

// State 状态机状态函数
type State[F any, T any] func(ctx context.Context, args F) (F, State[F, T], error)
