import api from './api'

export interface Model {
  id: number
  base_url: string
  config: string
  type: string
  model_name: string
  model_real_name: string
  status: string
  tag: string[]
  created_at: string
  updated_at: string
}

export interface ModelQuery {
  page: number
  page_size: number
  type?: string
  status?: string
  tag?: string[]
  keyword?: string
}

export interface ModelResponse {
  total: number
  models: Model[]
}

export const modelService = {
  // 获取模型列表
  getModels: async (query: ModelQuery): Promise<ModelResponse> => {
    const { page, page_size, ...rest } = query
    const params = new URLSearchParams()
    params.append('page', String(page))
    params.append('page_size', String(page_size))
    
    // 将其他参数放在请求体中
    return api.post('/models/list', {
      ...rest,
      page,
      page_size
    })
  },

  // 创建模型
  createModel: async (data: Omit<Model, 'id' | 'created_at' | 'updated_at'>): Promise<Model> => {
    return api.post('/models', data)
  },

  // 更新模型
  updateModel: async (id: number, data: Partial<Model>): Promise<Model> => {
    return api.put('/models', { id, ...data })
  },

  // 删除模型
  deleteModel: async (id: number): Promise<Model> => {
    return api.delete('/models', { data: { id } })
  },

  // 获取单个模型
  getModel: async (id: number): Promise<Model> => {
    return api.get(`/models?id=${id}`)
  }
} 