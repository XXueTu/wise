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
import { Resource, resourceService } from "@/services/resourceService"
import { useEffect, useState } from "react"

export function ResourceManager() {
  const [resources, setResources] = useState<Resource[]>([])
  const [currentPage, setCurrentPage] = useState(1)
  const [pageSize] = useState(10)
  const [total, setTotal] = useState(0)
  const [searchName, setSearchName] = useState("")
  const [searchType, setSearchType] = useState("")
  const [searchStatus, setSearchStatus] = useState("")
  const [isDialogOpen, setIsDialogOpen] = useState(false)
  const [editingResource, setEditingResource] = useState<Resource | null>(null)

  const loadData = async () => {
    try {
      const data = await resourceService.getResources({
        page: currentPage,
        pageSize,
        name: searchName,
        type: searchType,
        status: searchStatus
      })
      setResources(data.items)
      setTotal(data.total)
    } catch (error) {
      console.error("加载资源失败:", error)
    }
  }

  useEffect(() => {
    loadData()
  }, [currentPage, pageSize, searchName, searchType, searchStatus])

  const handleSearch = () => {
    setCurrentPage(1)
    loadData()
  }

  const handleReset = () => {
    setSearchName("")
    setSearchType("")
    setSearchStatus("")
    setCurrentPage(1)
    loadData()
  }

  const handleAdd = () => {
    setEditingResource(null)
    setIsDialogOpen(true)
  }

  const handleEdit = (resource: Resource) => {
    setEditingResource(resource)
    setIsDialogOpen(true)
  }

  const handleDelete = async (id: number) => {
    if (window.confirm("确定要删除这个资源吗？")) {
      try {
        await resourceService.deleteResource(id)
        loadData()
      } catch (error) {
        console.error("删除资源失败:", error)
      }
    }
  }

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault()
    const formData = new FormData(e.target as HTMLFormElement)
    const data = {
      name: formData.get("name") as string,
      type: formData.get("type") as string,
      url: formData.get("url") as string,
      status: formData.get("status") as string,
    }

    try {
      if (editingResource) {
        await resourceService.updateResource(editingResource.id, data)
      } else {
        await resourceService.createResource(data)
      }
      setIsDialogOpen(false)
      loadData()
    } catch (error) {
      console.error("保存资源失败:", error)
    }
  }

  return (
    <div className="container mx-auto py-6">
      <div className="flex justify-between items-center mb-6">
        <div className="flex gap-4">
          <Input
            placeholder="资源名称"
            value={searchName}
            onChange={(e) => setSearchName(e.target.value)}
            className="w-[200px]"
          />
          <Input
            placeholder="资源类型"
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
          <Button onClick={handleSearch}>搜索</Button>
          <Button variant="outline" onClick={handleReset}>
            重置
          </Button>
        </div>
        <Button onClick={handleAdd}>添加资源</Button>
      </div>

      <Table>
        <TableHeader>
          <TableRow>
            <TableHead>ID</TableHead>
            <TableHead>名称</TableHead>
            <TableHead>类型</TableHead>
            <TableHead>URL</TableHead>
            <TableHead>状态</TableHead>
            <TableHead>创建时间</TableHead>
            <TableHead>操作</TableHead>
          </TableRow>
        </TableHeader>
        <TableBody>
          {resources && resources.length > 0 ? (
            resources.map((resource) => (
              <TableRow key={resource.id}>
                <TableCell>{resource.id}</TableCell>
                <TableCell>{resource.name}</TableCell>
                <TableCell>{resource.type}</TableCell>
                <TableCell>{resource.url}</TableCell>
                <TableCell>{resource.status}</TableCell>
                <TableCell>{resource.createdAt}</TableCell>
                <TableCell>
                  <Button
                    variant="ghost"
                    size="sm"
                    onClick={() => handleEdit(resource)}
                  >
                    编辑
                  </Button>
                  <Button
                    variant="ghost"
                    size="sm"
                    onClick={() => handleDelete(resource.id)}
                  >
                    删除
                  </Button>
                </TableCell>
              </TableRow>
            ))
          ) : (
            <TableRow>
              <TableCell colSpan={7} className="text-center">
                暂无数据
              </TableCell>
            </TableRow>
          )}
        </TableBody>
      </Table>

      {resources && resources.length > 0 && (
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
              {editingResource ? "编辑资源" : "添加资源"}
            </DialogTitle>
            <DialogDescription>
              请填写资源信息
            </DialogDescription>
          </DialogHeader>
          <form onSubmit={handleSubmit}>
            <div className="grid gap-4 py-4">
              <div className="grid gap-2">
                <label htmlFor="name">名称</label>
                <Input
                  id="name"
                  name="name"
                  defaultValue={editingResource?.name}
                  required
                />
              </div>
              <div className="grid gap-2">
                <label htmlFor="type">类型</label>
                <Input
                  id="type"
                  name="type"
                  defaultValue={editingResource?.type}
                  required
                />
              </div>
              <div className="grid gap-2">
                <label htmlFor="url">URL</label>
                <Input
                  id="url"
                  name="url"
                  defaultValue={editingResource?.url}
                  required
                />
              </div>
              <div className="grid gap-2">
                <label htmlFor="status">状态</label>
                <Input
                  id="status"
                  name="status"
                  defaultValue={editingResource?.status}
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