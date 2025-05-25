import api from './api'

export interface Tag {
  uid: string
  name: string
  description: string
  color: string
  icon: string
  created_at: string
  updated_at: string
}

export interface TagQuery {
  page?: number
  page_size?: number
  name?: string
}

export interface TagResponse {
  total: number
  list: Tag[]
}

export interface CreateTagRequest {
  name: string
  description: string
  color: string
  icon: string
}

export interface UpdateTagRequest {
  name: string
  description: string
  color: string
  icon: string
}

export const tagService = {
  // 获取标签列表
  getTags: async (query: TagQuery): Promise<TagResponse> => {
    const params = new URLSearchParams()
    Object.entries(query).forEach(([key, value]) => {
      if (value !== undefined && value !== '') {
        params.append(key, String(value))
      }
    })
    return api.get(`/tags?${params.toString()}`)
  },

  // 创建标签
  createTag: async (data: CreateTagRequest): Promise<Tag> => {
    return api.post('/tag', data)
  },

  // 更新标签
  updateTag: async (uid: string, data: UpdateTagRequest): Promise<Tag> => {
    return api.put('/tag', { uid, ...data })
  },

  // 删除标签
  deleteTag: async (uid: string): Promise<void> => {
    return api.delete('/tag', { data: { uid } })
  },

  // 获取标签详情
  getTag: async (uid: string): Promise<Tag> => {
    return api.get(`/tag?uid=${uid}`)
  }
} 