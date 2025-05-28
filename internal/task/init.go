package task

import (
	"github.com/XXueTu/wise/internal/svc"
)

var globalScheduler *TaskScheduler

// InitScheduler 初始化任务调度器
func InitScheduler(svc *svc.ServiceContext) {
	if globalScheduler != nil {
		return
	}

	globalScheduler = NewTaskScheduler(svc)

	// 注册URL标记任务执行器
	urlMarkTask := NewUrlMarkTask(svc)
	globalScheduler.RegisterTaskExecutor(urlMarkTask)

	// 启动调度器
	globalScheduler.Start()
}

// GetScheduler 获取全局调度器实例
func GetScheduler() *TaskScheduler {
	return globalScheduler
}
