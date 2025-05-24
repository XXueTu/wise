import api from './api';

export interface Stats {
  totalResources: number
  totalModels: number
  resourceTypes: { name: string; value: number }[]
  modelTypes: { name: string; value: number }[]
  dailyStats: { date: string; resources: number; models: number }[]
}

export const statsService = {
  // 获取统计数据
  getStats: async (): Promise<Stats> => {
    return api.get('/stats')
  }
} 