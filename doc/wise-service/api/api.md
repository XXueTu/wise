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

