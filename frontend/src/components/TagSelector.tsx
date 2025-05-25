import { cn } from "@/lib/utils"
import { Tag, tagService } from "@/services/tagService"
import { Check } from "lucide-react"
import { useEffect, useState } from "react"
import Select from "react-select"

interface TagOption {
  value: string
  label: string
}

interface TagSelectorProps {
  value: string[]
  onChange: (value: string[]) => void
  maxDisplayedTags?: number
  placeholder?: string
}

export function TagSelector({ value, onChange, maxDisplayedTags = 4, placeholder = "选择标签..." }: TagSelectorProps) {
  const [inputValue, setInputValue] = useState("")
  const [tags, setTags] = useState<Tag[]>([])
  const [isLoading, setIsLoading] = useState(false)
  const [currentPage, setCurrentPage] = useState(1)
  const pageSize = 1000

  const loadTags = async (keyword: string) => {
    try {
      setIsLoading(true)
      const data = await tagService.getTags({
        page: currentPage,
        page_size: pageSize,
        name: keyword || undefined
      })
      const newTags = data.list || []
      setTags(prevTags => {
        const uniqueTags = newTags.filter(
          newTag => !prevTags.some(prevTag => prevTag.uid === newTag.uid)
        )
        return [...prevTags, ...uniqueTags]
      })
    } catch (error) {
      console.error("加载标签失败:", error)
    } finally {
      setIsLoading(false)
    }
  }

  useEffect(() => {
    loadTags("")
  }, [])

  const handleInputChange = (newValue: string) => {
    setInputValue(newValue)
    setCurrentPage(1)
    setTags([])
    loadTags(newValue)
  }

  const handleMenuScrollToBottom = () => {
    loadTags(inputValue)
    setCurrentPage(currentPage + 1)
  }

  const formatOptionLabel = (option: TagOption) => (
    <div className="flex items-center gap-2">
      <div className="w-4 h-4 flex items-center justify-center">
        <Check
          className={cn(
            "h-4 w-4",
            value.includes(option.value) ? "opacity-100 text-primary" : "opacity-0"
          )}
        />
      </div>
      <span>{option.label}</span>
    </div>
  )

  const getTagWidth = (label: string) => {
    const chineseChars = label.match(/[\u4e00-\u9fa5]/g)?.length || 0
    const englishChars = label.length - chineseChars
    const width = 60 + chineseChars * 14 + englishChars * 8
    return Math.min(width, 300)
  }

  const options = tags.map(tag => ({
    value: tag.uid,
    label: tag.name
  }))

  return (
    <Select
      isMulti
      value={options.filter(option => value.includes(option.value))}
      options={options}
      onChange={(newValue) => onChange(newValue.map(v => v.value))}
      onInputChange={handleInputChange}
      onMenuScrollToBottom={handleMenuScrollToBottom}
      isLoading={isLoading}
      formatOptionLabel={formatOptionLabel}
      placeholder={placeholder}
      noOptionsMessage={() => "没有找到标签"}
      loadingMessage={() => "加载中..."}
      styles={{
        control: (base) => ({
          ...base,
          minHeight: "40px",
          maxHeight: "80px",
          overflowY: "auto"
        }),
        valueContainer: (base) => ({
          ...base,
          flexWrap: "wrap",
          gap: "4px",
          maxWidth: `${maxDisplayedTags * 100}px`
        }),
        multiValue: (base) => ({
          ...base,
          margin: 0,
          backgroundColor: "transparent"
        }),
        multiValueLabel: (base, state) => ({
          ...base,
          padding: "2px 8px",
          backgroundColor: "hsl(var(--primary))",
          color: "hsl(var(--primary-foreground))",
          borderRadius: "4px",
          maxWidth: `${getTagWidth(state.data.label)}px`,
          overflow: "hidden",
          textOverflow: "ellipsis",
          whiteSpace: "nowrap"
        }),
        multiValueRemove: (base) => ({
          ...base,
          display: "none"
        }),
        menu: (base) => ({
          ...base,
          minWidth: "400px"
        })
      }}
    />
  )
}

// 添加全局动画样式
const style = document.createElement('style')
style.textContent = `
  @keyframes slideDown {
    from {
      opacity: 0;
      transform: translateY(-10px);
    }
    to {
      opacity: 1;
      transform: translateY(0);
    }
  }
`
document.head.appendChild(style) 