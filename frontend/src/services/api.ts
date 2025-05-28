import axios from 'axios'

// 获取当前访问地址
const getBaseUrl = () => {
  const { protocol, hostname, port } = window.location
  // 如果当前是开发环境（localhost），使用默认地址
  if (hostname === 'localhost' || hostname === '127.0.0.1') {
    return 'http://127.0.0.1:8888/wise/api'
  }
  // 生产环境使用当前访问地址
  return `${protocol}//${hostname}${port ? `:${port}` : ''}/wise/api`
}

// 创建 axios 实例
const api = axios.create({
  baseURL: getBaseUrl(),
  timeout: 10000,
  headers: {
    'Content-Type': 'application/json',
    'Accept': 'application/json'
  },
  withCredentials: false
})

// 请求拦截器
api.interceptors.request.use(
  (config) => {
    // 打印请求信息
    console.log('Request URL:', config.url)
    console.log('Full URL:', `${config.baseURL}${config.url}`)
    console.log('Request Method:', config.method)
    console.log('Request Data:', config.data)
    console.log('Request Params:', config.params)
    console.log('Request Headers:', config.headers)
    return config
  },
  (error) => {
    return Promise.reject(error)
  }
)

// 响应拦截器
api.interceptors.response.use(
  (response) => {
    const { code, msg, data } = response.data
    
    // 打印响应信息
    console.log('Response:', response.data)
    console.log('Response Headers:', response.headers)
    
    // 处理业务错误
    if (code !== 0) {
      return Promise.reject(new Error(msg || '请求失败'))
    }
    
    // 返回数据部分
    return data
  },
  (error) => {
    // 打印错误信息
    console.error('API Error:', error.response || error)
    console.error('Error Headers:', error.response?.headers)
    
    // 处理 HTTP 错误
    if (error.response) {
      const { status, data } = error.response
      return Promise.reject(new Error(data?.msg || `HTTP Error: ${status}`))
    }
    
    // 处理网络错误
    return Promise.reject(new Error(error.message || '网络错误'))
  }
)

export default api 