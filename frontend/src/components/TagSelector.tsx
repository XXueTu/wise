import { cn } from "@/lib/utils"
import { Tag, tagService } from "@/services/tagService"
import { Check } from "lucide-react"
import { useEffect, useState } from "react"
import Select from 'react-select'

interface TagSelectorProps {
  value: string[]
  onChange: (value: string[]) => void
  placeholder?: string
  maxDisplayedTags?: number
}

interface TagOption {
  value: string
  label: string
  icon?: string
  color?: string
}

export function TagSelector({
  value,
  onChange,
  placeholder = "选择标签...",
  maxDisplayedTags
}: TagSelectorProps) {
  const [inputValue, setInputValue] = useState("")
  const [tags, setTags] = useState<Tag[]>([])
  const [isLoading, setIsLoading] = useState(false)
  const [hasMore, setHasMore] = useState(true)
  const [page, setPage] = useState(1)
  const pageSize = 1000

  useEffect(() => {
    loadTags("")
  }, [])

  const loadTags = async (searchText: string) => {
    try {
      setIsLoading(true)
      const response = await tagService.getTags({
        page: 1,
        page_size: 1000,
        name: searchText || undefined
      })
      setTags(response.list || [])
    } catch (error) {
      console.error("加载标签失败:", error)
      setTags([])
    } finally {
      setIsLoading(false)
    }
  }

  const options: TagOption[] = tags.map(tag => ({
    value: tag.uid,
    label: tag.name,
    icon: tag.icon,
    color: tag.color
  }))

  const selectedOptions = options.filter(option => value.includes(option.value))

  const handleInputChange = (newValue: string) => {
    setInputValue(newValue)
    loadTags(newValue)
  }

  const handleMenuScrollToBottom = () => {
    if (hasMore) {
      loadTags(inputValue)
      setPage(page + 1)
    }
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
    const baseWidth = 60 // 基础宽度
    const charWidth = 14 // 每个中文字符的宽度
    const maxWidth = 300 // 增加最大宽度
    // 计算中文字符数量
    const chineseChars = (label.match(/[\u4e00-\u9fa5]/g) || []).length
    // 计算英文字符数量
    const englishChars = label.length - chineseChars
    // 中文字符宽度更大，英文字符宽度较小
    const width = baseWidth + (chineseChars * charWidth) + (englishChars * 8)
    return Math.min(width, maxWidth)
  }

  return (
    <Select
      isMulti
      value={value.map(uid => {
        const tag = tags.find(t => t.uid === uid)
        return {
          value: uid,
          label: tag?.name || uid,
          color: tag?.color,
          icon: tag?.icon
        }
      })}
      onChange={(newValue) => {
        onChange(newValue.map(v => v.value))
      }}
      onInputChange={handleInputChange}
      onMenuScrollToBottom={handleMenuScrollToBottom}
      options={tags.map(tag => ({
        value: tag.uid,
        label: tag.name,
        color: tag.color,
        icon: tag.icon
      }))}
      isLoading={isLoading}
      placeholder={placeholder}
      noOptionsMessage={() => "没有找到标签"}
      loadingMessage={() => "加载中..."}
      formatOptionLabel={formatOptionLabel}
      styles={{
        control: (base) => ({
          ...base,
          minHeight: '36px',
          backgroundColor: 'white',
          borderColor: 'hsl(var(--input))',
          transition: 'all 0.2s ease',
          '&:hover': {
            borderColor: 'hsl(var(--input))',
            boxShadow: '0 2px 4px rgba(0,0,0,0.05)'
          }
        }),
        valueContainer: (base) => ({
          ...base,
          flexWrap: 'wrap',
          gap: '4px',
          padding: '2px 8px'
        }),
        multiValue: (base, state) => ({
          ...base,
          backgroundColor: 'hsl(var(--secondary))',
          maxWidth: maxDisplayedTags ? `${getTagWidth(state.data.label)}px` : 'none',
          margin: '0',
          borderRadius: '4px',
          transition: 'all 0.2s ease',
          transform: 'scale(1)',
          '&:hover': {
            transform: 'scale(1.02)',
            boxShadow: '0 2px 4px rgba(0,0,0,0.1)'
          }
        }),
        multiValueLabel: (base) => ({
          ...base,
          color: 'hsl(var(--secondary-foreground))',
          padding: '2px 8px',
          fontSize: '14px',
          whiteSpace: 'nowrap',
          overflow: 'hidden',
          textOverflow: 'ellipsis',
          transition: 'all 0.2s ease'
        }),
        multiValueRemove: (base) => ({
          ...base,
          color: 'hsl(var(--secondary-foreground))',
          padding: '0 4px',
          transition: 'all 0.2s ease',
          ':hover': {
            backgroundColor: 'hsl(var(--secondary-hover))',
            color: 'hsl(var(--secondary-foreground))',
            transform: 'scale(1.1)'
          }
        }),
        menu: (base) => ({
          ...base,
          backgroundColor: 'white',
          border: '1px solid hsl(var(--border))',
          borderRadius: '6px',
          boxShadow: '0 2px 10px rgba(0, 0, 0, 0.1)',
          minWidth: '400px',
          animation: 'slideDown 0.2s ease'
        }),
        option: (base, state) => ({
          ...base,
          backgroundColor: state.isFocused ? 'hsl(var(--accent))' : 'white',
          color: state.isFocused ? 'hsl(var(--accent-foreground))' : 'hsl(var(--foreground))',
          transition: 'all 0.2s ease',
          ':active': {
            backgroundColor: 'hsl(var(--accent))',
            color: 'hsl(var(--accent-foreground))',
            transform: 'scale(0.98)'
          }
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