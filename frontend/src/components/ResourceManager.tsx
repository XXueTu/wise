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
  const [searchTitle, setSearchTitle] = useState("")
  const [searchType, setSearchType] = useState("")
  const [isDialogOpen, setIsDialogOpen] = useState(false)
  const [editingResource, setEditingResource] = useState<Resource | null>(null)

  const loadData = async () => {
    try {
      const data = await resourceService.getResources({
        page: currentPage,
        page_size: pageSize,
        type: searchType,
        keyword: searchTitle
      })
      setResources(data.resources)
      setTotal(data.total)
    } catch (error) {
      console.error("加载资源失败:", error)
    }
  }

  useEffect(() => {
    loadData()
  }, [currentPage, pageSize, searchTitle, searchType])

  const handleSearch = () => {
    setCurrentPage(1)
    loadData()
  }

  const handleReset = () => {
    setSearchTitle("")
    setSearchType("")
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
      url: formData.get("url") as string,
      title: formData.get("title") as string,
      content: formData.get("content") as string,
      type: formData.get("type") as string,
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
            placeholder="标题"
            value={searchTitle}
            onChange={(e) => setSearchTitle(e.target.value)}
            className="w-[200px]"
          />
          <Input
            placeholder="类型"
            value={searchType}
            onChange={(e) => setSearchType(e.target.value)}
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
            <TableHead>URL</TableHead>
            <TableHead>标题</TableHead>
            <TableHead>内容</TableHead>
            <TableHead>类型</TableHead>
            <TableHead>创建时间</TableHead>
            <TableHead>操作</TableHead>
          </TableRow>
        </TableHeader>
        <TableBody>
          {resources && resources.length > 0 ? (
            resources.map((resource) => (
              <TableRow key={resource.id}>
                <TableCell>{resource.id}</TableCell>
                <TableCell>{resource.url}</TableCell>
                <TableCell>{resource.title}</TableCell>
                <TableCell>{resource.content}</TableCell>
                <TableCell>{resource.type}</TableCell>
                <TableCell>{resource.created_at}</TableCell>
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
                <label htmlFor="url">URL</label>
                <Input
                  id="url"
                  name="url"
                  defaultValue={editingResource?.url}
                  required
                />
              </div>
              <div className="grid gap-2">
                <label htmlFor="title">标题</label>
                <Input
                  id="title"
                  name="title"
                  defaultValue={editingResource?.title}
                  required
                />
              </div>
              <div className="grid gap-2">
                <label htmlFor="content">内容</label>
                <Input
                  id="content"
                  name="content"
                  defaultValue={editingResource?.content}
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