import { Card, CardContent, CardHeader, CardTitle } from "@/components/ui/card"
import { useEffect, useState } from "react"
import {
    CartesianGrid,
    Cell,
    Line,
    LineChart,
    Pie,
    PieChart,
    ResponsiveContainer,
    Tooltip,
    XAxis,
    YAxis,
} from "recharts"

interface Stats {
  totalResources: number
  totalModels: number
  resourceTypes: { name: string; value: number }[]
  modelTypes: { name: string; value: number }[]
  dailyStats: { date: string; resources: number; models: number }[]
}

// 模拟数据
const mockStats: Stats = {
  totalResources: 156,
  totalModels: 42,
  resourceTypes: [
    { name: "文本", value: 45 },
    { name: "图片", value: 35 },
    { name: "音频", value: 25 },
    { name: "视频", value: 15 },
  ],
  modelTypes: [
    { name: "文本模型", value: 20 },
    { name: "图像模型", value: 12 },
    { name: "音频模型", value: 6 },
    { name: "视频模型", value: 4 },
  ],
  dailyStats: [
    { date: "2024-03-01", resources: 10, models: 2 },
    { date: "2024-03-02", resources: 15, models: 3 },
    { date: "2024-03-03", resources: 12, models: 4 },
    { date: "2024-03-04", resources: 18, models: 5 },
    { date: "2024-03-05", resources: 20, models: 6 },
    { date: "2024-03-06", resources: 25, models: 7 },
    { date: "2024-03-07", resources: 30, models: 8 },
  ],
}

export function Dashboard() {
  const [stats, setStats] = useState<Stats>(mockStats)
  const [loading, setLoading] = useState(true)

  useEffect(() => {
    const fetchStats = async () => {
      try {
        setLoading(true)
        const response = await fetch("/api/stats")
        const data = await response.json()
        setStats(data)
      } catch (error) {
        console.error("加载统计数据失败:", error)
        // 如果 API 调用失败，使用模拟数据
        setStats(mockStats)
      } finally {
        setLoading(false)
      }
    }

    fetchStats()
  }, [])

  const COLORS = ["#0088FE", "#00C49F", "#FFBB28", "#FF8042"]

  if (loading) {
    return (
      <div className="container mx-auto py-10">
        <div className="flex items-center justify-center h-[400px]">
          <div className="text-lg">加载中...</div>
        </div>
      </div>
    )
  }

  return (
    <div className="container mx-auto py-6">
      <div className="grid gap-4 md:grid-cols-2 lg:grid-cols-4 mb-8">
        <Card>
          <CardHeader className="flex flex-row items-center justify-between space-y-0 pb-2">
            <CardTitle className="text-sm font-medium">资源总数</CardTitle>
          </CardHeader>
          <CardContent>
            <div className="text-2xl font-bold">{stats?.totalResources || 0}</div>
          </CardContent>
        </Card>
        <Card>
          <CardHeader className="flex flex-row items-center justify-between space-y-0 pb-2">
            <CardTitle className="text-sm font-medium">模型总数</CardTitle>
          </CardHeader>
          <CardContent>
            <div className="text-2xl font-bold">{stats?.totalModels || 0}</div>
          </CardContent>
        </Card>
      </div>

      <div className="grid gap-4 md:grid-cols-2 mb-8">
        <Card>
          <CardHeader>
            <CardTitle>近7天趋势</CardTitle>
          </CardHeader>
          <CardContent>
            <div className="h-[300px]">
              <ResponsiveContainer width="100%" height="100%">
                <LineChart data={stats?.dailyStats || []}>
                  <CartesianGrid strokeDasharray="3 3" />
                  <XAxis dataKey="date" />
                  <YAxis />
                  <Tooltip />
                  <Line
                    type="monotone"
                    dataKey="resources"
                    stroke="#8884d8"
                    name="资源"
                  />
                  <Line
                    type="monotone"
                    dataKey="models"
                    stroke="#82ca9d"
                    name="模型"
                  />
                </LineChart>
              </ResponsiveContainer>
            </div>
          </CardContent>
        </Card>

        <Card>
          <CardHeader>
            <CardTitle>资源类型分布</CardTitle>
          </CardHeader>
          <CardContent>
            <div className="h-[300px]">
              <ResponsiveContainer width="100%" height="100%">
                <PieChart>
                  <Pie
                    data={stats?.resourceTypes || []}
                    dataKey="value"
                    nameKey="name"
                    cx="50%"
                    cy="50%"
                    outerRadius={80}
                    label
                  >
                    {(stats?.resourceTypes || []).map((_, index) => (
                      <Cell
                        key={`cell-${index}`}
                        fill={COLORS[index % COLORS.length]}
                      />
                    ))}
                  </Pie>
                  <Tooltip />
                </PieChart>
              </ResponsiveContainer>
            </div>
          </CardContent>
        </Card>
      </div>
    </div>
  )
} 