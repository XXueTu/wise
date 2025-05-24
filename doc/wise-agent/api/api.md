### 1. "API基础"

1. route definition

- Url: /wise/api/base
- Method: POST
- Request: `BaseRequest`
- Response: `BaseResponse`

2. request definition



```golang
type BaseRequest struct {
	Id string `json:"id"` // 主键
	Name string `json:"name"` // 名称
}
```


3. response definition



```golang
type BaseResponse struct {
	Result string `json:"result"` // 结果
}
```

### 2. "URL链接识别"

1. route definition

- Url: /wise/api/url
- Method: POST
- Request: `URLRequest`
- Response: `URLResponse`

2. request definition



```golang
type URLRequest struct {
	URL string `json:"url"` // URL链接
}
```


3. response definition



```golang
type URLResponse struct {
	Tag []string `json:"tag"` // 标签
	Title string `json:"title"` // 标题
	Description string `json:"description"` // 描述
	Link string `json:"link"` // 链接
}
```

