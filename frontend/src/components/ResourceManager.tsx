import { TagSelector } from "@/components/TagSelector"
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
import { Resource, resourceService } from "@/services/resourceService"
import { Check, ChevronDown, Copy, ExternalLink, FileText, Globe, Tag } from "lucide-react"
import { useEffect, useRef, useState } from "react"

// Tooltip 组件
function Tooltip({ content, children }: { content: string, children: React.ReactNode }) {
  const [show, setShow] = useState(false)
  const [position, setPosition] = useState({ top: 0, left: 0 })
  const tooltipRef = useRef<HTMLDivElement>(null)

  const handleMouseEnter = (e: React.MouseEvent) => {
    const rect = e.currentTarget.getBoundingClientRect()
    setPosition({
      top: rect.bottom + window.scrollY + 5,
      left: rect.left + window.scrollX
    })
    setShow(true)
  }

  return (
    <>
      <span
        onMouseEnter={handleMouseEnter}
        onMouseLeave={() => setShow(false)}
        className="relative inline-block"
      >
        {children}
      </span>
      {show && (
        <div
          ref={tooltipRef}
          className="fixed z-[9999] bg-black text-white text-xs rounded shadow-lg px-3 py-2 max-w-xl break-all whitespace-pre-wrap"
          style={{
            top: `${position.top}px`,
            left: `${position.left}px`,
            transform: 'translateX(-50%)'
          }}
        >
          {content}
        </div>
      )}
    </>
  )
}

// 长文本展示组件（带Tooltip，单行省略）
function LongTextDisplay({ 
  content,
  maxLength = 30,
  icon,
  isLink = false
}: { 
  content: string
  maxLength?: number
  icon?: React.ReactNode
  isLink?: boolean
}) {
  const [isCopied, setIsCopied] = useState(false)

  const handleCopy = async () => {
    await navigator.clipboard.writeText(content)
    setIsCopied(true)
    setTimeout(() => setIsCopied(false), 2000)
  }

  const isOverflow = content.length > maxLength
  const displayContent = isOverflow ? content.slice(0, maxLength) + '...' : content

  return (
    <div className="group flex items-center gap-2 min-w-0">
      {icon && (
        <div className="flex-shrink-0 text-gray-400">{icon}</div>
      )}
      <span className="flex-1 min-w-0 relative">
        {isOverflow ? (
          <Tooltip content={content}>
            {isLink ? (
              <a
                href={content}
                target="_blank"
                rel="noopener noreferrer"
                className="text-blue-600 hover:text-blue-800 hover:underline flex items-center gap-1 truncate max-w-[180px]"
                style={{maxWidth: 180, display: 'inline-block', verticalAlign: 'bottom'}}
              >
                {displayContent}
                <ExternalLink className="h-3 w-3" />
              </a>
            ) : (
              <span className="truncate max-w-[180px] inline-block align-bottom cursor-pointer">{displayContent}</span>
            )}
          </Tooltip>
        ) : (
          isLink ? (
            <a
              href={content}
              target="_blank"
              rel="noopener noreferrer"
              className="text-blue-600 hover:text-blue-800 hover:underline flex items-center gap-1 truncate max-w-[180px]"
              style={{maxWidth: 180, display: 'inline-block', verticalAlign: 'bottom'}}
            >
              {displayContent}
              <ExternalLink className="h-3 w-3" />
            </a>
          ) : (
            <span className="truncate max-w-[180px] inline-block align-bottom">{displayContent}</span>
          )
        )}
      </span>
      <Button
        variant="ghost"
        size="sm"
        className="h-7 px-2"
        onClick={handleCopy}
      >
        {isCopied ? (
          <Check className="h-3 w-3 text-green-500" />
        ) : (
          <Copy className="h-3 w-3" />
        )}
      </Button>
    </div>
  )
}

// 内容预览组件（2行省略，无Tooltip，仅弹窗）
function ContentPreview({ content }: { content: string }) {
  const [isDialogOpen, setIsDialogOpen] = useState(false)
  const [isCopied, setIsCopied] = useState(false)

  const handleCopy = async () => {
    await navigator.clipboard.writeText(content)
    setIsCopied(true)
    setTimeout(() => setIsCopied(false), 2000)
  }

  // 2行省略
  return (
    <>
      <div className="group relative min-w-0">
        <div className="flex items-start gap-2 min-w-0">
          <div className="flex-1 min-w-0">
            <div
              className="text-sm text-gray-600 bg-gray-50 rounded-md p-2 cursor-pointer line-clamp-2 max-h-[3.2em] min-h-[2.4em] break-all"
              style={{ display: '-webkit-box', WebkitLineClamp: 2, WebkitBoxOrient: 'vertical', overflow: 'hidden' }}
              onClick={() => setIsDialogOpen(true)}
              title={content.length > 0 ? '点击查看全部内容' : ''}
            >
              {content}
            </div>
          </div>
          <div className="flex-shrink-0 opacity-0 group-hover:opacity-100 transition-opacity flex items-center gap-1">
            <Button
              variant="ghost"
              size="sm"
              className="h-7 px-2"
              onClick={handleCopy}
            >
              {isCopied ? (
                <Check className="h-3 w-3 text-green-500" />
              ) : (
                <Copy className="h-3 w-3" />
              )}
            </Button>
            {content.length > 50 && (
              <Button
                variant="ghost"
                size="sm"
                className="h-7 px-2"
                onClick={() => setIsDialogOpen(true)}
              >
                <ChevronDown className="h-3 w-3" />
              </Button>
            )}
          </div>
        </div>
      </div>
      <Dialog open={isDialogOpen} onOpenChange={setIsDialogOpen}>
        <DialogContent className="max-w-4xl">
          <div className="flex flex-col gap-4">
            <div className="flex justify-between items-center">
              <div className="font-semibold text-base">内容详情</div>
              <Button
                variant="ghost"
                size="sm"
                className="h-7 px-2"
                onClick={handleCopy}
              >
                {isCopied ? (
                  <Check className="h-4 w-4 text-green-500" />
                ) : (
                  <Copy className="h-4 w-4" />
                )}
                <span className="ml-1 text-xs">复制</span>
              </Button>
            </div>
            <div className="bg-gray-50 rounded-md p-4 max-h-[60vh] overflow-y-auto text-sm text-gray-800 break-all whitespace-pre-wrap">
              {content}
            </div>
          </div>
        </DialogContent>
      </Dialog>
    </>
  )
}

export function ResourceManager() {
  const [resources, setResources] = useState<Resource[]>([])
  const [currentPage, setCurrentPage] = useState(1)
  const [pageSize] = useState(10)
  const [total, setTotal] = useState(0)
  const [searchTitle, setSearchTitle] = useState("")
  const [searchType, setSearchType] = useState("")
  const [searchTagUids, setSearchTagUids] = useState<string[]>([])
  const [isDialogOpen, setIsDialogOpen] = useState(false)
  const [editingResource, setEditingResource] = useState<Resource | null>(null)
  const [selectedTagUids, setSelectedTagUids] = useState<string[]>([])

  const loadData = async () => {
    try {
      const data = await resourceService.getResources({
        page: currentPage,
        page_size: pageSize,
        type: searchType || undefined,
        keyword: searchTitle || undefined,
        tag_uids: searchTagUids.length > 0 ? searchTagUids : undefined
      })
      setResources(data.resources)
      setTotal(data.total)
    } catch (error) {
      console.error("加载资源失败:", error)
    }
  }

  useEffect(() => {
    loadData()
  }, [currentPage, pageSize, searchTitle, searchType, searchTagUids])

  const handleSearch = () => {
    setCurrentPage(1)
    loadData()
  }

  const handleReset = () => {
    setSearchTitle("")
    setSearchType("")
    setSearchTagUids([])
    setCurrentPage(1)
    loadData()
  }

  const handleAdd = () => {
    setEditingResource(null)
    setSelectedTagUids([])
    setIsDialogOpen(true)
  }

  const handleEdit = (resource: Resource) => {
    setEditingResource(resource)
    setSelectedTagUids(resource.tag_uids || [])
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
      tag_uids: selectedTagUids,
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
          <div className="w-[300px]">
            <TagSelector
              value={searchTagUids}
              onChange={setSearchTagUids}
              placeholder="选择标签..."
              maxDisplayedTags={4}
            />
          </div>
          <Button onClick={handleSearch}>搜索</Button>
          <Button variant="outline" onClick={handleReset}>
            重置
          </Button>
        </div>
        <Button onClick={handleAdd}>添加资源</Button>
      </div>

      <div className="rounded-lg border bg-white shadow-sm">
        <Table>
          <TableHeader>
            <TableRow className="hover:bg-transparent">
              <TableHead className="w-[80px]">ID</TableHead>
              <TableHead className="w-[180px]">URL</TableHead>
              <TableHead className="w-[150px]">标题</TableHead>
              <TableHead className="max-w-[220px] w-[220px]">内容</TableHead>
              <TableHead className="w-[120px] text-center">类型</TableHead>
              <TableHead className="w-[180px]">标签</TableHead>
              <TableHead className="w-[180px]">创建时间</TableHead>
              <TableHead className="w-[120px]">操作</TableHead>
            </TableRow>
          </TableHeader>
          <TableBody>
            {resources && resources.length > 0 ? (
              resources.map((resource) => (
                <TableRow key={resource.id} className="hover:bg-gray-50">
                  <TableCell className="font-mono text-sm text-gray-500">{resource.id}</TableCell>
                  <TableCell>
                    <LongTextDisplay 
                      content={resource.url} 
                      maxLength={10}
                      icon={<Globe className="h-4 w-4" />}
                      isLink={true}
                    />
                  </TableCell>
                  <TableCell>
                    <LongTextDisplay 
                      content={resource.title} 
                      maxLength={10}
                      icon={<FileText className="h-4 w-4" />}
                    />
                  </TableCell>
                  <TableCell className="max-w-[220px] w-[220px]">
                    <ContentPreview content={resource.content} />
                  </TableCell>
                  <TableCell className="w-[120px] text-center">
                    <span className="px-3 py-1 bg-gray-100 rounded text-sm inline-block">
                      {resource.type}
                    </span>
                  </TableCell>
                  <TableCell className="w-[180px]">
                    <div className="flex flex-wrap gap-1">
                      {resource.tags?.map((tag, index) => (
                        <span 
                          key={index} 
                          className="inline-flex items-center gap-1 px-2 py-1 bg-blue-50 text-blue-700 rounded text-sm"
                        >
                          <Tag className="h-3 w-3" />
                          {tag}
                        </span>
                      ))}
                    </div>
                  </TableCell>
                  <TableCell className="text-sm text-gray-500">
                    {resource.created_at}
                  </TableCell>
                  <TableCell>
                    <div className="flex gap-2">
                      <Button
                        variant="ghost"
                        size="sm"
                        onClick={() => handleEdit(resource)}
                        className="h-8 hover:bg-gray-100"
                      >
                        编辑
                      </Button>
                      <Button
                        variant="ghost"
                        size="sm"
                        onClick={() => handleDelete(resource.id)}
                        className="h-8 hover:bg-red-50 hover:text-red-600"
                      >
                        删除
                      </Button>
                    </div>
                  </TableCell>
                </TableRow>
              ))
            ) : (
              <TableRow>
                <TableCell colSpan={8} className="text-center py-8 text-gray-500">
                  暂无数据
                </TableCell>
              </TableRow>
            )}
          </TableBody>
        </Table>
      </div>

      {resources && resources.length > 0 && (
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
              <div className="grid gap-2">
                <label htmlFor="tag_uids">标签</label>
                <TagSelector
                  value={selectedTagUids}
                  onChange={setSelectedTagUids}
                  placeholder="选择标签..."
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