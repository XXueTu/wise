import api from './api'

export interface Model {
  id: number
  name: string
  tag: string
  status: string
  createdAt: string
}

export interface ModelQuery {
  page: number
  pageSize: number
  name?: string
  tag?: string
  status?: string
}

export interface ModelResponse {
  items: Model[]
  total: number
}

export const modelService = {
  // 获取模型列表
  getModels: async (query: ModelQuery): Promise<ModelResponse> => {
    const params = new URLSearchParams()
    Object.entries(query).forEach(([key, value]) => {
      if (value !== undefined && value !== '') {
        params.append(key, String(value))
      }
    })
    return api.get(`/models?${params.toString()}`)
  },

  // 创建模型
  createModel: async (data: Omit<Model, 'id' | 'createdAt'>): Promise<Model> => {
    return api.post('/models', data)
  },

  // 更新模型
  updateModel: async (id: number, data: Partial<Model>): Promise<Model> => {
    return api.put(`/models/${id}`, data)
  },

  // 删除模型
  deleteModel: async (id: number): Promise<void> => {
    return api.delete(`/models/${id}`)
  }
} 