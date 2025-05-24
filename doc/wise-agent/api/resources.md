### 1. "创建资源"

1. route definition

- Url: /wise/api/resources
- Method: POST
- Request: `CreateResourceRequest`
- Response: `Resource`

2. request definition



```golang
type CreateResourceRequest struct {
	URL string `json:"url"` // URL链接
	Title string `json:"title"` // 标题
	Content string `json:"content"` // 内容
	Type string `json:"type"` // 类型
}
```


3. response definition



```golang
type Resource struct {
	Id int64 `json:"id"` // 主键
	URL string `json:"url"` // URL链接
	Title string `json:"title"` // 标题
	Content string `json:"content"` // 内容
	Type string `json:"type"` // 类型
	CreatedAt string `json:"created_at"` // 创建时间
	UpdatedAt string `json:"updated_at"` // 更新时间
}
```

### 2. "更新资源"

1. route definition

- Url: /wise/api/resources
- Method: PUT
- Request: `UpdateResourceRequest`
- Response: `Resource`

2. request definition



```golang
type UpdateResourceRequest struct {
	Id int64 `json:"id"` // 主键
	URL string `json:"url"` // URL链接
	Title string `json:"title"` // 标题
	Content string `json:"content"` // 内容
	Type string `json:"type"` // 类型
}
```


3. response definition



```golang
type Resource struct {
	Id int64 `json:"id"` // 主键
	URL string `json:"url"` // URL链接
	Title string `json:"title"` // 标题
	Content string `json:"content"` // 内容
	Type string `json:"type"` // 类型
	CreatedAt string `json:"created_at"` // 创建时间
	UpdatedAt string `json:"updated_at"` // 更新时间
}
```

### 3. "删除资源"

1. route definition

- Url: /wise/api/resources
- Method: DELETE
- Request: `DeleteResourceRequest`
- Response: `Resource`

2. request definition



```golang
type DeleteResourceRequest struct {
	Id int64 `json:"id"` // 主键
}
```


3. response definition



```golang
type Resource struct {
	Id int64 `json:"id"` // 主键
	URL string `json:"url"` // URL链接
	Title string `json:"title"` // 标题
	Content string `json:"content"` // 内容
	Type string `json:"type"` // 类型
	CreatedAt string `json:"created_at"` // 创建时间
	UpdatedAt string `json:"updated_at"` // 更新时间
}
```

### 4. "获取单个资源"

1. route definition

- Url: /wise/api/resources
- Method: GET
- Request: `GetResourceRequest`
- Response: `Resource`

2. request definition



```golang
type GetResourceRequest struct {
	Id int64 `json:"id"` // 主键
}
```


3. response definition



```golang
type Resource struct {
	Id int64 `json:"id"` // 主键
	URL string `json:"url"` // URL链接
	Title string `json:"title"` // 标题
	Content string `json:"content"` // 内容
	Type string `json:"type"` // 类型
	CreatedAt string `json:"created_at"` // 创建时间
	UpdatedAt string `json:"updated_at"` // 更新时间
}
```

### 5. "分页查询资源列表"

1. route definition

- Url: /wise/api/resources/list
- Method: GET
- Request: `ListResourceRequest`
- Response: `ListResourceResponse`

2. request definition



```golang
type ListResourceRequest struct {
	Page int64 `json:"page"` // 页码
	PageSize int64 `json:"page_size"` // 每页数量
	Type string `json:"type"` // 类型（可选）
	Keyword string `json:"keyword"` // 关键词（可选）
}
```


3. response definition



```golang
type ListResourceResponse struct {
	Total int64 `json:"total"` // 总数
	Resources []Resource `json:"resources"` // 资源列表
}
```

