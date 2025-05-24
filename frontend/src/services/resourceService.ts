import api from './api'

export interface Resource {
  id: number
  url: string
  title: string
  content: string
  type: string
  created_at: string
  updated_at: string
}

export interface ResourceQuery {
  page: number
  page_size: number
  type?: string
  keyword?: string
}

export interface ResourceResponse {
  total: number
  resources: Resource[]
}

export const resourceService = {
  // 获取资源列表
  getResources: async (query: ResourceQuery): Promise<ResourceResponse> => {
    const params = new URLSearchParams()
    Object.entries(query).forEach(([key, value]) => {
      if (value !== undefined && value !== '') {
        params.append(key, String(value))
      }
    })
    return api.get(`/resources/list?${params.toString()}`)
  },

  // 创建资源
  createResource: async (data: Omit<Resource, 'id' | 'created_at' | 'updated_at'>): Promise<Resource> => {
    return api.post('/resources', data)
  },

  // 更新资源
  updateResource: async (id: number, data: Partial<Resource>): Promise<Resource> => {
    return api.put('/resources', { id, ...data })
  },

  // 删除资源
  deleteResource: async (id: number): Promise<Resource> => {
    return api.delete('/resources', { data: { id } })
  },

  // 获取单个资源
  getResource: async (id: number): Promise<Resource> => {
    return api.get(`/resources?id=${id}`)
  }
} 