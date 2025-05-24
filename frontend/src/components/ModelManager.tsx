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
  const [searchTag, setSearchTag] = useState("")
  const [searchStatus, setSearchStatus] = useState("")
  const [isDialogOpen, setIsDialogOpen] = useState(false)
  const [editingModel, setEditingModel] = useState<Model | null>(null)

  const loadData = async () => {
    try {
      const data = await modelService.getModels({
        page: currentPage,
        pageSize,
        name: searchName,
        tag: searchTag,
        status: searchStatus
      })
      setModels(data.items)
      setTotal(data.total)
    } catch (error) {
      console.error("加载模型失败:", error)
    }
  }

  useEffect(() => {
    loadData()
  }, [currentPage, pageSize, searchName, searchTag, searchStatus])

  const handleSearch = () => {
    setCurrentPage(1)
    loadData()
  }

  const handleReset = () => {
    setSearchName("")
    setSearchTag("")
    setSearchStatus("")
    setCurrentPage(1)
    loadData()
  }

  const handleAdd = () => {
    setEditingModel(null)
    setIsDialogOpen(true)
  }

  const handleEdit = (model: Model) => {
    setEditingModel(model)
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
      name: formData.get("name") as string,
      tag: formData.get("tag") as string,
      status: formData.get("status") as string,
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
            placeholder="模型类型"
            value={searchTag}
            onChange={(e) => setSearchTag(e.target.value)}
            className="w-[200px]"
          />
          <Input
            placeholder="状态"
            value={searchStatus}
            onChange={(e) => setSearchStatus(e.target.value)}
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
            <TableHead>名称</TableHead>
            <TableHead>标签</TableHead>
            <TableHead>状态</TableHead>
            <TableHead>创建时间</TableHead>
            <TableHead>操作</TableHead>
          </TableRow>
        </TableHeader>
        <TableBody>
          {models && models.length > 0 ? (
            models.map((model) => (
              <TableRow key={model.id}>
                <TableCell>{model.id}</TableCell>
                <TableCell>{model.name}</TableCell>
                <TableCell>{model.tag}</TableCell>
                <TableCell>{model.status}</TableCell>
                <TableCell>{model.createdAt}</TableCell>
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
              <TableCell colSpan={6} className="text-center">
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
                <label htmlFor="name">名称</label>
                <Input
                  id="name"
                  name="name"
                  defaultValue={editingModel?.name}
                  required
                />
              </div>
              <div className="grid gap-2">
                <label htmlFor="tag">标签</label>
                <Input
                  id="tag"
                  name="tag"
                  defaultValue={editingModel?.tag}
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