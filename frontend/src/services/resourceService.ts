import api from './api'

export interface Resource {
  id: number
  name: string
  type: string
  url: string
  status: string
  createdAt: string
}

export interface ResourceQuery {
  page: number
  pageSize: number
  name?: string
  type?: string
  status?: string
}

export interface ResourceResponse {
  items: Resource[]
  total: number
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
    return api.get(`/resources?${params.toString()}`)
  },

  // 创建资源
  createResource: async (data: Omit<Resource, 'id' | 'createdAt'>): Promise<Resource> => {
    return api.post('/resources', data)
  },

  // 更新资源
  updateResource: async (id: number, data: Partial<Resource>): Promise<Resource> => {
    return api.put(`/resources/${id}`, data)
  },

  // 删除资源
  deleteResource: async (id: number): Promise<void> => {
    return api.delete(`/resources/${id}`)
  }
} 