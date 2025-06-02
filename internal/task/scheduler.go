package task

import (
	"context"
	"sync"
	"time"

	"github.com/zeromicro/go-zero/core/logx"
	"golang.org/x/sync/semaphore"

	"github.com/XXueTu/wise/internal/model"
	"github.com/XXueTu/wise/internal/svc"
	"github.com/XXueTu/wise/pkg/agent/url_analyse.go"
)

// CreateTask 创建任务
func CreateTask(ctx context.Context, svc *svc.ServiceContext, args string, taskName string, taskType string, totalStep int64) error {
	tid := model.GenUid()
	// 创建任务
	err := svc.TasksModel.Create(ctx, &model.Tasks{
		Tid:          tid,
		Name:         taskName,
		Types:        taskType,
		Status:       model.TaskStatusInit,
		CurrentState: "start",
		RetryCount:   0,
		TotalSteps:   totalStep,
		CurrentStep:  0,
		Params:       args,
		Result:       "{}",
		Duration:     0,
		Error:        "{}",
		Extend:       "{}",
	})
	if err != nil {
		return err
	}
	return nil
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
	taskCtxs     sync.Map // 存储任务上下文，用于取消任务
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
	}
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
	_ = s.svc.TasksModel.UpdateStatus(ctx, task.Tid, model.TaskStatusRunning, "{}")

	// 执行任务
	err := url_analyse.RunUrlAnalyseAgent(task.Tid, task.Params)
	if err != nil {
		// 检查是否需要重试
		if task.RetryCount < 2 {
			task.Status = model.TaskStatusRetry
			task.RetryCount++
			task.Error = err.Error()
			_ = s.svc.TasksModel.Update(ctx, task)
			// 等待重试间隔后重新调度
			time.Sleep(5 * time.Second)
			logx.Infof("重试任务: %v", task.Tid)
			return
		}
		_ = s.svc.TasksModel.UpdateStatus(ctx, task.Tid, model.TaskStatusFailed, err.Error())
		logx.Errorf("执行任务失败: %v", err)
		return
	}
	_ = s.svc.TasksModel.UpdateStatus(ctx, task.Tid, model.TaskStatusSuccess, "{}")
	logx.Infof("执行任务成功: %v", task.Tid)
}
