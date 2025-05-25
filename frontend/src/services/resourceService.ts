import api from './api'

export interface Resource {
  id: number
  url: string
  title: string
  content: string
  type: string
  tags: string[]      // 标签名称列表
  tag_uids: string[]  // 标签ID列表
  created_at: string
  updated_at: string
}

export interface ResourceQuery {
  page?: number
  page_size?: number
  type?: string
  keyword?: string
  tag_uids?: string[]  // 标签ID列表
}

export interface ResourceResponse {
  total: number
  resources: Resource[]
}

export interface CreateResourceRequest {
  url: string
  title: string
  content: string
  type: string
  tag_uids: string[]  // 标签ID列表
}

export interface UpdateResourceRequest {
  id: number
  url: string
  title: string
  content: string
  type: string
  tag_uids: string[]  // 标签ID列表
}

export const resourceService = {
  // 获取资源列表
  getResources: async (query: ResourceQuery): Promise<ResourceResponse> => {
    return api.post('/resources/list', query)
  },

  // 创建资源
  createResource: async (data: CreateResourceRequest): Promise<Resource> => {
    return api.post('/resources', data)
  },

  // 更新资源
  updateResource: async (id: number, data: Omit<UpdateResourceRequest, 'id'>): Promise<Resource> => {
    return api.put('/resources', { id, ...data })
  },

  // 删除资源
  deleteResource: async (id: number): Promise<void> => {
    return api.delete('/resources', { data: { id } })
  },

  // 获取单个资源
  getResource: async (id: number): Promise<Resource> => {
    return api.get(`/resources?id=${id}`)
  }
} 