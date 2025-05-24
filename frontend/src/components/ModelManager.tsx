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
import {
  Table,
  TableBody,
  TableCell,
  TableHead,
  TableHeader,
  TableRow,
} from "@/components/ui/table"
import { Model, modelService } from "@/services/modelService"
import { useEffect, useState } from "react"

export function ModelManager() {
  const [models, setModels] = useState<Model[]>([])
  const [currentPage, setCurrentPage] = useState(1)
  const [pageSize] = useState(10)
  const [total, setTotal] = useState(0)
  const [searchName, setSearchName] = useState("")
  const [searchType, setSearchType] = useState("")
  const [searchStatus, setSearchStatus] = useState("")
  const [searchTagInput, setSearchTagInput] = useState("")
  const [searchTags, setSearchTags] = useState<string[]>([])
  const [isDialogOpen, setIsDialogOpen] = useState(false)
  const [editingModel, setEditingModel] = useState<Model | null>(null)
  const [tagInput, setTagInput] = useState("")

  const loadData = async () => {
    try {
      const data = await modelService.getModels({
        page: currentPage,
        page_size: pageSize,
        type: searchType,
        status: searchStatus,
        tag: searchTags,
        keyword: searchName
      })
      setModels(data.models)
      setTotal(data.total)
    } catch (error) {
      console.error("加载模型失败:", error)
    }
  }

  useEffect(() => {
    loadData()
  }, [currentPage, pageSize, searchName, searchType, searchStatus, searchTags])

  const handleSearch = () => {
    setCurrentPage(1)
    loadData()
  }

  const handleReset = () => {
    setSearchName("")
    setSearchType("")
    setSearchStatus("")
    setSearchTagInput("")
    setSearchTags([])
    setCurrentPage(1)
    loadData()
  }

  const handleAdd = () => {
    setEditingModel(null)
    setTagInput("")
    setIsDialogOpen(true)
  }

  const handleEdit = (model: Model) => {
    setEditingModel(model)
    setTagInput(model.tag.join(", "))
    setIsDialogOpen(true)
  }

  const handleDelete = async (id: number) => {
    if (window.confirm("确定要删除这个模型吗？")) {
      try {
        await modelService.deleteModel(id)
        loadData()
      } catch (error) {
        console.error("删除模型失败:", error)
      }
    }
  }

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault()
    const formData = new FormData(e.target as HTMLFormElement)
    const data = {
      base_url: formData.get("base_url") as string,
      config: formData.get("config") as string,
      type: formData.get("type") as string,
      model_name: formData.get("model_name") as string,
      model_real_name: formData.get("model_real_name") as string,
      status: formData.get("status") as string,
      tag: tagInput.split(",").map(tag => tag.trim()).filter(tag => tag !== "")
    }

    try {
      if (editingModel) {
        await modelService.updateModel(editingModel.id, data)
      } else {
        await modelService.createModel(data)
      }
      setIsDialogOpen(false)
      loadData()
    } catch (error) {
      console.error("保存模型失败:", error)
    }
  }

  const handleSearchTagChange = (value: string) => {
    setSearchTagInput(value)
    const tags = value.split(",")
      .map(tag => tag.trim())
      .filter(tag => tag !== "")
    setSearchTags(tags)
  }

  return (
    <div className="container mx-auto py-6">
      <div className="flex justify-between items-center mb-6">
        <div className="flex gap-4">
          <Input
            placeholder="模型名称"
            value={searchName}
            onChange={(e) => setSearchName(e.target.value)}
            className="w-[200px]"
          />
          <Input
            placeholder="类型"
            value={searchType}
            onChange={(e) => setSearchType(e.target.value)}
            className="w-[200px]"
          />
          <Input
            placeholder="状态"
            value={searchStatus}
            onChange={(e) => setSearchStatus(e.target.value)}
            className="w-[200px]"
          />
          <Input
            placeholder="标签（用逗号分隔）"
            value={searchTagInput}
            onChange={(e) => handleSearchTagChange(e.target.value)}
            className="w-[200px]"
          />
          <Button onClick={handleSearch}>搜索</Button>
          <Button variant="outline" onClick={handleReset}>
            重置
          </Button>
        </div>
        <Button onClick={handleAdd}>添加模型</Button>
      </div>

      <Table>
        <TableHeader>
          <TableRow>
            <TableHead>ID</TableHead>
            <TableHead>基础URL</TableHead>
            <TableHead>模型名称</TableHead>
            <TableHead>真实名称</TableHead>
            <TableHead>类型</TableHead>
            <TableHead>状态</TableHead>
            <TableHead>标签</TableHead>
            <TableHead>创建时间</TableHead>
            <TableHead>操作</TableHead>
          </TableRow>
        </TableHeader>
        <TableBody>
          {models && models.length > 0 ? (
            models.map((model) => (
              <TableRow key={model.id}>
                <TableCell>{model.id}</TableCell>
                <TableCell>{model.base_url}</TableCell>
                <TableCell>{model.model_name}</TableCell>
                <TableCell>{model.model_real_name}</TableCell>
                <TableCell>{model.type}</TableCell>
                <TableCell>{model.status}</TableCell>
                <TableCell>
                  <div className="flex flex-wrap gap-1">
                    {model.tag.map((tag, index) => (
                      <span
                        key={index}
                        className="px-2 py-1 bg-gray-100 rounded-full text-sm"
                      >
                        {tag}
                      </span>
                    ))}
                  </div>
                </TableCell>
                <TableCell>{model.created_at}</TableCell>
                <TableCell>
                  <Button
                    variant="ghost"
                    size="sm"
                    onClick={() => handleEdit(model)}
                  >
                    编辑
                  </Button>
                  <Button
                    variant="ghost"
                    size="sm"
                    onClick={() => handleDelete(model.id)}
                  >
                    删除
                  </Button>
                </TableCell>
              </TableRow>
            ))
          ) : (
            <TableRow>
              <TableCell colSpan={9} className="text-center">
                暂无数据
              </TableCell>
            </TableRow>
          )}
        </TableBody>
      </Table>

      {models && models.length > 0 && (
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
              {editingModel ? "编辑模型" : "添加模型"}
            </DialogTitle>
            <DialogDescription>
              请填写模型信息
            </DialogDescription>
          </DialogHeader>
          <form onSubmit={handleSubmit}>
            <div className="grid gap-4 py-4">
              <div className="grid gap-2">
                <label htmlFor="base_url">基础URL</label>
                <Input
                  id="base_url"
                  name="base_url"
                  defaultValue={editingModel?.base_url}
                  required
                />
              </div>
              <div className="grid gap-2">
                <label htmlFor="config">配置信息</label>
                <Input
                  id="config"
                  name="config"
                  defaultValue={editingModel?.config}
                  required
                />
              </div>
              <div className="grid gap-2">
                <label htmlFor="type">类型</label>
                <Input
                  id="type"
                  name="type"
                  defaultValue={editingModel?.type}
                  required
                />
              </div>
              <div className="grid gap-2">
                <label htmlFor="model_name">模型名称</label>
                <Input
                  id="model_name"
                  name="model_name"
                  defaultValue={editingModel?.model_name}
                  required
                />
              </div>
              <div className="grid gap-2">
                <label htmlFor="model_real_name">真实名称</label>
                <Input
                  id="model_real_name"
                  name="model_real_name"
                  defaultValue={editingModel?.model_real_name}
                  required
                />
              </div>
              <div className="grid gap-2">
                <label htmlFor="status">状态</label>
                <Input
                  id="status"
                  name="status"
                  defaultValue={editingModel?.status}
                  required
                />
              </div>
              <div className="grid gap-2">
                <label htmlFor="tag">标签（用逗号分隔）</label>
                <Input
                  id="tag"
                  name="tag"
                  value={tagInput}
                  onChange={(e) => setTagInput(e.target.value)}
                  placeholder="例如：标签1, 标签2, 标签3"
                  required
                />
              </div>
            </div>
            <DialogFooter>
              <Button type="submit">保存</Button>
            </DialogFooter>
          </form>
        </DialogContent>
      </Dialog>
    </div>
  )
} 