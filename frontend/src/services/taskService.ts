import api from './api'

export interface Task {
  tid: string
  name: string
  types: string
  status: string
  current_state: string
  total_steps: number
  current_step: number
  params: string
  result: string
  duration: number
  error: string
  extend: string
  created_at: string
  updated_at: string
}

export interface TaskQuery {
  page?: number
  page_size?: number
  name?: string
  types?: string
  status?: string
}

export interface TaskResponse {
  total: number
  list: Task[]
}

export interface CreateTaskRequest {
  name: string
  types: string
  params: string
  total_steps: number
  current_state: string
}

export interface UpdateTaskRequest {
  name?: string
  status?: string
  current_state?: string
  current_step?: number
  result?: string
  error?: string
  extend?: string
}

export interface TaskPlanDetail {
  pid: string
  name: string
  index: number
  status: string
  params: string
  result: string
  duration: number
  error: string
  created_at: string
  updated_at: string
}

export interface TaskVisualization {
  tid: string
  name: string
  types: string
  status: string
  current_state: string
  total_steps: number
  current_step: number
  plans: TaskPlanDetail[]
  created_at: string
  updated_at: string
}

export const taskService = {
  // 获取任务列表
  getTasks: async (query: TaskQuery): Promise<TaskResponse> => {
    const params = new URLSearchParams()
    Object.entries(query).forEach(([key, value]) => {
      if (value !== undefined && value !== '') {
        params.append(key, String(value))
      }
    })
    return api.get(`/tasks?${params.toString()}`)
  },

  // 创建任务
  createTask: async (data: CreateTaskRequest): Promise<Task> => {
    return api.post('/task', data)
  },

  // 更新任务
  updateTask: async (tid: string, data: UpdateTaskRequest): Promise<Task> => {
    return api.put('/task', { tid, ...data })
  },

  // 删除任务
  deleteTask: async (tid: string): Promise<void> => {
    return api.delete('/task', { data: { tid } })
  },

  // 获取任务详情
  getTask: async (tid: string): Promise<Task> => {
    return api.get(`/task?tid=${tid}`)
  },

  // 重试任务
  retryTask: async (tid: string): Promise<{ result: string }> => {
    return api.post('/task/retry', { tid })
  },

  // 暂停任务
  pauseTask: async (tid: string): Promise<{ result: string }> => {
    return api.post('/task/pause', { tid })
  },

  // 恢复任务
  resumeTask: async (tid: string): Promise<{ result: string }> => {
    return api.post('/task/resume', { tid })
  },

  // 取消任务
  cancelTask: async (tid: string): Promise<{ result: string }> => {
    return api.post('/task/cancel', { tid })
  },

  // 获取任务可视化信息
  getTaskVisualization: async (tid: string): Promise<TaskVisualization> => {
    return api.get(`/task/visualization?tid=${tid}`)
  }
} 