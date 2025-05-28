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
    PaginationTotal,
} from "@/components/ui/pagination"
import {
    Table,
    TableBody,
    TableCell,
    TableHead,
    TableHeader,
    TableRow,
} from "@/components/ui/table"
import { Tag, tagService } from "@/services/tagService"
import { useEffect, useState } from "react"

export function TagManager() {
  const [tags, setTags] = useState<Tag[]>([])
  const [currentPage, setCurrentPage] = useState(1)
  const [pageSize] = useState(10)
  const [total, setTotal] = useState(0)
  const [searchName, setSearchName] = useState("")
  const [isDialogOpen, setIsDialogOpen] = useState(false)
  const [editingTag, setEditingTag] = useState<Tag | null>(null)

  const loadData = async () => {
    try {
      const data = await tagService.getTags({
        page: currentPage,
        page_size: pageSize,
        name: searchName || undefined
      })
      setTags(data.list || [])
      setTotal(data.total)
    } catch (error) {
      console.error("加载标签失败:", error)
    }
  }

  useEffect(() => {
    loadData()
  }, [currentPage, pageSize, searchName])

  const handleSearch = () => {
    setCurrentPage(1)
    loadData()
  }

  const handleReset = () => {
    setSearchName("")
    setCurrentPage(1)
    loadData()
  }

  const handleAdd = () => {
    setEditingTag(null)
    setIsDialogOpen(true)
  }

  const handleEdit = (tag: Tag) => {
    setEditingTag(tag)
    setIsDialogOpen(true)
  }

  const handleDelete = async (uid: string) => {
    if (window.confirm("确定要删除这个标签吗？")) {
      try {
        await tagService.deleteTag(uid)
        loadData()
      } catch (error) {
        console.error("删除标签失败:", error)
      }
    }
  }

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault()
    const formData = new FormData(e.target as HTMLFormElement)
    const data = {
      name: formData.get("name") as string,
      description: formData.get("description") as string,
      color: formData.get("color") as string,
      icon: formData.get("icon") as string,
    }

    try {
      if (editingTag) {
        await tagService.updateTag(editingTag.uid, data)
      } else {
        await tagService.createTag(data)
      }
      setIsDialogOpen(false)
      loadData()
    } catch (error) {
      console.error("保存标签失败:", error)
    }
  }

  return (
    <div className="container mx-auto py-6">
      <div className="flex justify-between items-center mb-6">
        <div className="flex gap-4">
          <Input
            placeholder="标签名称"
            value={searchName}
            onChange={(e) => setSearchName(e.target.value)}
            className="w-[200px]"
          />
          <Button onClick={handleSearch}>搜索</Button>
          <Button variant="outline" onClick={handleReset}>
            重置
          </Button>
        </div>
        <Button onClick={handleAdd}>添加标签</Button>
      </div>

      <Table>
        <TableHeader>
          <TableRow>
            <TableHead>UID</TableHead>
            <TableHead>名称</TableHead>
            <TableHead>描述</TableHead>
            <TableHead>颜色</TableHead>
            <TableHead>图标</TableHead>
            <TableHead>创建时间</TableHead>
            <TableHead>操作</TableHead>
          </TableRow>
        </TableHeader>
        <TableBody>
          {tags && tags.length > 0 ? (
            tags.map((tag) => (
              <TableRow key={tag.uid}>
                <TableCell>{tag.uid}</TableCell>
                <TableCell>{tag.name}</TableCell>
                <TableCell>{tag.description}</TableCell>
                <TableCell>
                  <div className="flex items-center gap-2">
                    <span
                      className="w-4 h-4 rounded-full"
                      style={{ backgroundColor: tag.color }}
                    />
                    {tag.color}
                  </div>
                </TableCell>
                <TableCell>{tag.icon}</TableCell>
                <TableCell>{tag.created_at}</TableCell>
                <TableCell>
                  <Button
                    variant="ghost"
                    size="sm"
                    onClick={() => handleEdit(tag)}
                  >
                    编辑
                  </Button>
                  <Button
                    variant="ghost"
                    size="sm"
                    onClick={() => handleDelete(tag.uid)}
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

      {tags && tags.length > 0 && (
        <div className="mt-4">
          <Pagination>
            <PaginationTotal total={total} />
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
              {editingTag ? "编辑标签" : "添加标签"}
            </DialogTitle>
            <DialogDescription>
              请填写标签信息
            </DialogDescription>
          </DialogHeader>
          <form onSubmit={handleSubmit}>
            <div className="grid gap-4 py-4">
              <div className="grid gap-2">
                <label htmlFor="name">名称</label>
                <Input
                  id="name"
                  name="name"
                  defaultValue={editingTag?.name}
                  required
                />
              </div>
              <div className="grid gap-2">
                <label htmlFor="description">描述</label>
                <Input
                  id="description"
                  name="description"
                  defaultValue={editingTag?.description}
                />
              </div>
              <div className="grid gap-2">
                <label htmlFor="color">颜色</label>
                <Input
                  id="color"
                  name="color"
                  type="color"
                  defaultValue={editingTag?.color}
                  required
                />
              </div>
              <div className="grid gap-2">
                <label htmlFor="icon">图标</label>
                <Input
                  id="icon"
                  name="icon"
                  defaultValue={editingTag?.icon}
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