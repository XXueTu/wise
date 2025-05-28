import { Button } from "@/components/ui/button"
import {
  Dialog,
  DialogContent,
  DialogDescription,
  DialogFooter,
  DialogHeader,
  DialogTitle
} from "@/components/ui/dialog"
import { Input } from "@/components/ui/input"
import {
  Pagination,
  PaginationContent,
  PaginationItem,
  PaginationLink,
  PaginationNext,
  PaginationPrevious,
} from "@/components/ui/pagination"
import { Progress } from "@/components/ui/progress"
import {
  Select,
  SelectContent,
  SelectItem,
  SelectTrigger,
  SelectValue,
} from "@/components/ui/select"
import {
  Table,
  TableBody,
  TableCell,
  TableHead,
  TableHeader,
  TableRow,
} from "@/components/ui/table"
import { Task, TaskPlanDetail, TaskVisualization, taskService } from "@/services/taskService"
import { Check, Copy } from "lucide-react"
import { useEffect, useRef, useState } from "react"

// 长文本展示组件
function LongTextDisplay({ 
  title, 
  content, 
  isError = false,
  isJson = false 
}: { 
  title: string
  content: string
  isError?: boolean
  isJson?: boolean
}) {
  const [isCopied, setIsCopied] = useState(false)
  const [formattedContent, setFormattedContent] = useState("")

  useEffect(() => {
    if (isJson) {
      try {
        const parsed = JSON.parse(content)
        setFormattedContent(JSON.stringify(parsed, null, 2))
      } catch {
        setFormattedContent(content)
      }
    } else {
      setFormattedContent(content)
    }
  }, [content, isJson])

  const handleCopy = async () => {
    await navigator.clipboard.writeText(content)
    setIsCopied(true)
    setTimeout(() => setIsCopied(false), 2000)
  }

  const bgColor = isError ? 'bg-red-50' : 'bg-gray-50'
  const textColor = isError ? 'text-red-600' : 'text-gray-900'

  return (
    <div className="space-y-2">
      <div className="flex items-center justify-between">
        <h4 className="text-sm font-medium">{title}</h4>
        <Button
          variant="ghost"
          size="sm"
          className="h-8 px-2"
          onClick={handleCopy}
        >
          {isCopied ? (
            <Check className="h-4 w-4 text-green-500" />
          ) : (
            <Copy className="h-4 w-4" />
          )}
        </Button>
      </div>
      <div className={`${bgColor} p-4 rounded-lg`}>
        <pre className={`text-sm whitespace-pre-wrap ${textColor} max-h-[500px] overflow-y-auto`}>
          {formattedContent}
        </pre>
      </div>
    </div>
  )
}

export function TaskManager() {
  const [tasks, setTasks] = useState<Task[]>([])
  const [currentPage, setCurrentPage] = useState(1)
  const [pageSize] = useState(10)
  const [total, setTotal] = useState(0)
  const [searchName, setSearchName] = useState("")
  const [searchType, setSearchType] = useState("all")
  const [searchStatus, setSearchStatus] = useState("all")
  const [isDialogOpen, setIsDialogOpen] = useState(false)
  const [isVisualizationDialogOpen, setIsVisualizationDialogOpen] = useState(false)
  const [editingTask, setEditingTask] = useState<Task | null>(null)
  const [visualizationData, setVisualizationData] = useState<TaskVisualization | null>(null)
  const [selectedPlan, setSelectedPlan] = useState<TaskPlanDetail | null>(null)
  
  // 用于存储自动刷新的定时器
  const pageRefreshTimer = useRef<NodeJS.Timeout | null>(null)
  const visualizationRefreshTimer = useRef<NodeJS.Timeout | null>(null)

  const loadData = async () => {
    try {
      const data = await taskService.getTasks({
        page: currentPage,
        page_size: pageSize,
        name: searchName || undefined,
        types: searchType === 'all' ? undefined : searchType,
        status: searchStatus === 'all' ? undefined : searchStatus
      })
      setTasks(data.list || [])
      setTotal(data.total)
    } catch (error) {
      console.error("加载任务失败:", error)
    }
  }

  // 启动页面自动刷新
  const startPageRefresh = () => {
    if (pageRefreshTimer.current) {
      clearInterval(pageRefreshTimer.current)
    }
    pageRefreshTimer.current = setInterval(loadData, 5000)
  }

  // 停止页面自动刷新
  const stopPageRefresh = () => {
    if (pageRefreshTimer.current) {
      clearInterval(pageRefreshTimer.current)
      pageRefreshTimer.current = null
    }
  }

  // 启动可视化数据自动刷新
  const startVisualizationRefresh = (tid: string) => {
    if (visualizationRefreshTimer.current) {
      clearInterval(visualizationRefreshTimer.current)
    }
    visualizationRefreshTimer.current = setInterval(async () => {
      try {
        const data = await taskService.getTaskVisualization(tid)
        setVisualizationData(data)
        // 如果当前选中的计划存在，更新其数据
        if (selectedPlan) {
          const updatedPlan = data.plans.find(p => p.pid === selectedPlan.pid)
          if (updatedPlan) {
            setSelectedPlan(updatedPlan)
          }
        }
      } catch (error) {
        console.error("刷新任务可视化信息失败:", error)
      }
    }, 2000)
  }

  // 停止可视化数据自动刷新
  const stopVisualizationRefresh = () => {
    if (visualizationRefreshTimer.current) {
      clearInterval(visualizationRefreshTimer.current)
      visualizationRefreshTimer.current = null
    }
  }

  useEffect(() => {
    loadData()
    startPageRefresh()

    // 组件卸载时清理定时器
    return () => {
      stopPageRefresh()
      stopVisualizationRefresh()
    }
  }, [currentPage, pageSize, searchName, searchType, searchStatus])

  const handleSearch = () => {
    setCurrentPage(1)
    loadData()
  }

  const handleReset = () => {
    setSearchName("")
    setSearchType("all")
    setSearchStatus("all")
    setCurrentPage(1)
    loadData()
  }

  const handleAdd = () => {
    setEditingTask(null)
    setIsDialogOpen(true)
  }

  const handleEdit = (task: Task) => {
    setEditingTask(task)
    setIsDialogOpen(true)
  }

  const handleDelete = async (tid: string) => {
    if (window.confirm("确定要删除这个任务吗？")) {
      try {
        await taskService.deleteTask(tid)
        loadData()
      } catch (error) {
        console.error("删除任务失败:", error)
      }
    }
  }

  const handleRetry = async (tid: string) => {
    try {
      await taskService.retryTask(tid)
      loadData()
    } catch (error) {
      console.error("重试任务失败:", error)
    }
  }

  const handlePause = async (tid: string) => {
    try {
      await taskService.pauseTask(tid)
      loadData()
    } catch (error) {
      console.error("暂停任务失败:", error)
    }
  }

  const handleResume = async (tid: string) => {
    try {
      await taskService.resumeTask(tid)
      loadData()
    } catch (error) {
      console.error("恢复任务失败:", error)
    }
  }

  const handleCancel = async (tid: string) => {
    try {
      await taskService.cancelTask(tid)
      loadData()
    } catch (error) {
      console.error("取消任务失败:", error)
    }
  }

  const handleViewVisualization = async (tid: string) => {
    try {
      const data = await taskService.getTaskVisualization(tid)
      setVisualizationData(data)
      // 默认选中第一个计划
      if (data.plans && data.plans.length > 0) {
        setSelectedPlan(data.plans[0])
      }
      setIsVisualizationDialogOpen(true)
      startVisualizationRefresh(tid)
    } catch (error) {
      console.error("获取任务可视化信息失败:", error)
    }
  }

  // 监听对话框关闭
  useEffect(() => {
    if (!isVisualizationDialogOpen) {
      stopVisualizationRefresh()
      setSelectedPlan(null)
      setVisualizationData(null)
    }
  }, [isVisualizationDialogOpen])

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault()
    const formData = new FormData(e.target as HTMLFormElement)
    const data = {
      name: formData.get("name") as string,
      types: formData.get("types") as string,
      params: formData.get("params") as string,
      total_steps: Number(formData.get("total_steps")),
      current_state: formData.get("current_state") as string,
      status: formData.get("status") as string || editingTask?.status || "init"
    }

    try {
      if (editingTask) {
        await taskService.updateTask(editingTask.tid, data)
      } else {
        await taskService.createTask(data)
      }
      setIsDialogOpen(false)
      loadData()
    } catch (error) {
      console.error("保存任务失败:", error)
    }
  }

  return (
    <div className="container mx-auto py-6">
      <div className="flex justify-between items-center mb-6">
        <div className="flex gap-4">
          <Input
            placeholder="任务名称"
            value={searchName}
            onChange={(e) => setSearchName(e.target.value)}
            className="w-[200px]"
          />
          <Select value={searchType} onValueChange={setSearchType}>
            <SelectTrigger className="w-[180px]">
              <SelectValue placeholder="任务类型" />
            </SelectTrigger>
            <SelectContent>
              <SelectItem value="all">全部</SelectItem>
              <SelectItem value="url_mark">URL标记</SelectItem>
              <SelectItem value="url_parse">URL解析</SelectItem>
            </SelectContent>
          </Select>
          <Select value={searchStatus} onValueChange={setSearchStatus}>
            <SelectTrigger className="w-[180px]">
              <SelectValue placeholder="任务状态" />
            </SelectTrigger>
            <SelectContent>
              <SelectItem value="all">全部</SelectItem>
              <SelectItem value="pending">待处理</SelectItem>
              <SelectItem value="running">运行中</SelectItem>
              <SelectItem value="completed">已完成</SelectItem>
              <SelectItem value="failed">失败</SelectItem>
              <SelectItem value="cancelled">已取消</SelectItem>
            </SelectContent>
          </Select>
          <Button onClick={handleSearch}>搜索</Button>
          <Button variant="outline" onClick={handleReset}>
            重置
          </Button>
        </div>
        <Button onClick={handleAdd}>添加任务</Button>
      </div>

      <Table>
        <TableHeader>
          <TableRow>
            <TableHead>TID</TableHead>
            <TableHead>名称</TableHead>
            <TableHead>类型</TableHead>
            <TableHead>状态</TableHead>
            <TableHead>当前状态</TableHead>
            <TableHead>进度</TableHead>
            <TableHead>创建时间</TableHead>
            <TableHead>操作</TableHead>
          </TableRow>
        </TableHeader>
        <TableBody>
          {tasks && tasks.length > 0 ? (
            tasks.map((task) => (
              <TableRow key={task.tid}>
                <TableCell>{task.tid}</TableCell>
                <TableCell>{task.name}</TableCell>
                <TableCell>{task.types}</TableCell>
                <TableCell>{task.status}</TableCell>
                <TableCell>{task.current_state}</TableCell>
                <TableCell>
                  <div className="space-y-1">
                    <Progress 
                      value={(task.current_step / task.total_steps) * 100} 
                      className="h-2"
                    />
                    <div className="text-xs text-gray-500">
                      {task.current_step}/{task.total_steps}
                    </div>
                  </div>
                </TableCell>
                <TableCell>{task.created_at}</TableCell>
                <TableCell>
                  <div className="flex gap-2">
                    <Button
                      variant="ghost"
                      size="sm"
                      onClick={() => handleViewVisualization(task.tid)}
                    >
                      查看
                    </Button>
                    <Button
                      variant="ghost"
                      size="sm"
                      onClick={() => handleEdit(task)}
                    >
                      编辑
                    </Button>
                    {task.status === 'failed' && (
                      <Button
                        variant="ghost"
                        size="sm"
                        onClick={() => handleRetry(task.tid)}
                      >
                        重试
                      </Button>
                    )}
                    {task.status === 'running' && (
                      <>
                        <Button
                          variant="ghost"
                          size="sm"
                          onClick={() => handlePause(task.tid)}
                        >
                          暂停
                        </Button>
                        <Button
                          variant="ghost"
                          size="sm"
                          onClick={() => handleCancel(task.tid)}
                        >
                          取消
                        </Button>
                      </>
                    )}
                    {task.status === 'paused' && (
                      <Button
                        variant="ghost"
                        size="sm"
                        onClick={() => handleResume(task.tid)}
                      >
                        恢复
                      </Button>
                    )}
                    <Button
                      variant="ghost"
                      size="sm"
                      onClick={() => handleDelete(task.tid)}
                    >
                      删除
                    </Button>
                  </div>
                </TableCell>
              </TableRow>
            ))
          ) : (
            <TableRow>
              <TableCell colSpan={8} className="text-center">
                暂无数据
              </TableCell>
            </TableRow>
          )}
        </TableBody>
      </Table>

      {tasks && tasks.length > 0 && (
        <div className="mt-4">
          <Pagination>
            <PaginationContent>
              <PaginationItem>
                <PaginationPrevious
                  onClick={() => setCurrentPage((p) => Math.max(1, p - 1))}
                  disabled={currentPage === 1}
                />
              </PaginationItem>
              {Array.from({ length: Math.ceil(total / pageSize) }, (_, i) => i + 1).map(
                (page) => (
                  <PaginationItem key={page}>
                    <PaginationLink
                      isActive={currentPage === page}
                      onClick={() => setCurrentPage(page)}
                    >
                      {page}
                    </PaginationLink>
                  </PaginationItem>
                )
              )}
              <PaginationItem>
                <PaginationNext
                  onClick={() =>
                    setCurrentPage((p) => Math.min(Math.ceil(total / pageSize), p + 1))
                  }
                  disabled={currentPage === Math.ceil(total / pageSize)}
                />
              </PaginationItem>
            </PaginationContent>
          </Pagination>
        </div>
      )}

      <Dialog open={isDialogOpen} onOpenChange={setIsDialogOpen}>
        <DialogContent>
          <DialogHeader>
            <DialogTitle>
              {editingTask ? "编辑任务" : "添加任务"}
            </DialogTitle>
            <DialogDescription>
              请填写任务信息
            </DialogDescription>
          </DialogHeader>
          <form onSubmit={handleSubmit}>
            <div className="grid gap-4 py-4">
              <div className="grid gap-2">
                <label htmlFor="name">名称</label>
                <Input
                  id="name"
                  name="name"
                  defaultValue={editingTask?.name}
                  required
                />
              </div>
              <div className="grid gap-2">
                <label htmlFor="types">类型</label>
                <Select name="types" defaultValue={editingTask?.types}>
                  <SelectTrigger>
                    <SelectValue placeholder="选择任务类型" />
                  </SelectTrigger>
                  <SelectContent>
                    <SelectItem value="url_mark">URL标记</SelectItem>
                    <SelectItem value="url_parse">URL解析</SelectItem>
                  </SelectContent>
                </Select>
              </div>
              <div className="grid gap-2">
                <label htmlFor="params">参数</label>
                <Input
                  id="params"
                  name="params"
                  defaultValue={editingTask?.params}
                  required
                />
              </div>
              <div className="grid gap-2">
                <label htmlFor="total_steps">总步骤</label>
                <Input
                  id="total_steps"
                  name="total_steps"
                  type="number"
                  defaultValue={editingTask?.total_steps}
                  required
                />
              </div>
              <div className="grid gap-2">
                <label htmlFor="current_state">当前状态</label>
                <Input
                  id="current_state"
                  name="current_state"
                  defaultValue={editingTask?.current_state}
                  required
                />
              </div>
              <div className="grid gap-2">
                <label htmlFor="status">状态</label>
                <Select name="status" defaultValue={editingTask?.status || "init"}>
                  <SelectTrigger>
                    <SelectValue placeholder="选择任务状态" />
                  </SelectTrigger>
                  <SelectContent>
                    <SelectItem value="init">初始状态</SelectItem>
                    <SelectItem value="running">运行中</SelectItem>
                    <SelectItem value="success">成功</SelectItem>
                    <SelectItem value="failed">失败</SelectItem>
                    <SelectItem value="retry">重试中</SelectItem>
                    <SelectItem value="cancelled">已取消</SelectItem>
                  </SelectContent>
                </Select>
              </div>
            </div>
            <DialogFooter>
              <Button type="submit">保存</Button>
            </DialogFooter>
          </form>
        </DialogContent>
      </Dialog>

      <Dialog open={isVisualizationDialogOpen} onOpenChange={setIsVisualizationDialogOpen}>
        <DialogContent className="max-w-6xl h-[90vh] p-0 flex flex-col">
          <div className="flex flex-col h-full">
            <div className="px-6 py-4 border-b flex-shrink-0">
              <div className="flex items-center justify-between pr-8">
                <div>
                  <h2 className="text-lg font-semibold">任务可视化</h2>
                </div>
                {visualizationData && (
                  <div className="flex items-center gap-4">
                    <div>
                      <span className="text-sm text-gray-500">任务名称：</span>
                      <span className="font-medium">{visualizationData.name}</span>
                    </div>
                    <div>
                      <span className="text-sm text-gray-500">当前状态：</span>
                      <span className={`px-2 py-1 rounded text-sm ${
                        visualizationData.status === 'completed' ? 'bg-green-100 text-green-800' :
                        visualizationData.status === 'failed' ? 'bg-red-100 text-red-800' :
                        visualizationData.status === 'running' ? 'bg-blue-100 text-blue-800' :
                        'bg-gray-100 text-gray-800'
                      }`}>
                        {visualizationData.status}
                      </span>
                    </div>
                    <div>
                      <span className="text-sm text-gray-500">进度：</span>
                      <span className="font-medium">{visualizationData.current_step}/{visualizationData.total_steps}</span>
                    </div>
                  </div>
                )}
              </div>
            </div>
            {visualizationData && (
              <div className="flex flex-1 min-h-0">
                {/* 左侧时间线 */}
                <div className="w-1/3 border-r p-4 overflow-y-auto">
                  <div className="space-y-4">
                    {visualizationData.plans.map((plan, index) => (
                      <div
                        key={plan.pid}
                        className={`cursor-pointer p-4 rounded-lg transition-colors ${
                          plan.status === 'completed' ? 'bg-green-50 hover:bg-green-100' :
                          plan.status === 'failed' ? 'bg-red-50 hover:bg-red-100' :
                          plan.status === 'running' ? 'bg-blue-50 hover:bg-blue-100' :
                          'bg-gray-50 hover:bg-gray-100'
                        }`}
                        onClick={() => setSelectedPlan(plan)}
                      >
                        <div className="flex items-center gap-3">
                          <div className={`w-2 h-2 rounded-full ${
                            plan.status === 'completed' ? 'bg-green-500' :
                            plan.status === 'failed' ? 'bg-red-500' :
                            plan.status === 'running' ? 'bg-blue-500' :
                            'bg-gray-500'
                          }`} />
                          <div className="flex-1">
                            <p className="font-medium">{plan.name}</p>
                            <p className="text-sm text-gray-500">
                              {new Date(plan.created_at).toLocaleString()}
                            </p>
                          </div>
                          <span className={`px-2 py-1 rounded text-sm ${
                            plan.status === 'completed' ? 'bg-green-100 text-green-800' :
                            plan.status === 'failed' ? 'bg-red-100 text-red-800' :
                            plan.status === 'running' ? 'bg-blue-100 text-blue-800' :
                            'bg-gray-100 text-gray-800'
                          }`}>
                            {plan.status}
                          </span>
                        </div>
                      </div>
                    ))}
                  </div>
                </div>

                {/* 右侧详情 */}
                <div className="flex-1 p-4 overflow-y-auto">
                  {selectedPlan ? (
                    <div className="space-y-6">
                      <div className="border-b pb-4">
                        <h3 className="text-lg font-medium mb-2">计划详情</h3>
                        <div className="grid grid-cols-2 gap-4">
                          <div>
                            <p className="text-sm text-gray-500">计划名称</p>
                            <p className="font-medium">{selectedPlan.name}</p>
                          </div>
                          <div>
                            <p className="text-sm text-gray-500">计划索引</p>
                            <p className="font-medium">{selectedPlan.index}</p>
                          </div>
                          <div>
                            <p className="text-sm text-gray-500">状态</p>
                            <span className={`px-2 py-1 rounded text-sm ${
                              selectedPlan.status === 'completed' ? 'bg-green-100 text-green-800' :
                              selectedPlan.status === 'failed' ? 'bg-red-100 text-red-800' :
                              selectedPlan.status === 'running' ? 'bg-blue-100 text-blue-800' :
                              'bg-gray-100 text-gray-800'
                            }`}>
                              {selectedPlan.status}
                            </span>
                          </div>
                          <div>
                            <p className="text-sm text-gray-500">耗时</p>
                            <p className="font-medium">{selectedPlan.duration}ms</p>
                          </div>
                        </div>
                      </div>

                      <div className="space-y-4">
                        <LongTextDisplay 
                          title="参数" 
                          content={selectedPlan.params} 
                          isJson={true}
                        />

                        <LongTextDisplay 
                          title="结果" 
                          content={selectedPlan.result} 
                          isJson={true}
                        />

                        {selectedPlan.error && (
                          <LongTextDisplay 
                            title="错误信息" 
                            content={selectedPlan.error} 
                            isError={true}
                          />
                        )}

                        <div className="grid grid-cols-2 gap-4">
                          <div>
                            <p className="text-sm text-gray-500">创建时间</p>
                            <p className="font-medium">{new Date(selectedPlan.created_at).toLocaleString()}</p>
                          </div>
                          <div>
                            <p className="text-sm text-gray-500">更新时间</p>
                            <p className="font-medium">{new Date(selectedPlan.updated_at).toLocaleString()}</p>
                          </div>
                        </div>
                      </div>
                    </div>
                  ) : (
                    <div className="h-full flex items-center justify-center text-gray-500">
                      请从左侧选择一个执行计划查看详情
                    </div>
                  )}
                </div>
              </div>
            )}
          </div>
        </DialogContent>
      </Dialog>
    </div>
  )
} 