import axios from 'axios'
import type { AxiosRequestConfig, InternalAxiosRequestConfig, AxiosError, AxiosResponse } from 'axios'
import { ElMessage } from 'element-plus'
import { getDeviceId } from './deviceId'
import { getAccessToken, setAccessToken, clearAccessToken } from './token'

const service = axios.create({
  baseURL: import.meta.env.VITE_BASE_API,
  timeout: 30000
})

export interface ApiResponse<T> {
  code: number
  msg: string
  data: T
}

// 清除认证信息
const clearAuth = (): void => {
  clearAccessToken()
  localStorage.removeItem('user')
}

// 请求拦截器
service.interceptors.request.use(
  (config: AxiosRequestConfig) => {
    const token = getAccessToken()
    config.headers = {
      'Content-Type': 'application/json',
      'X-Device-Id': getDeviceId(),
      ...config.headers
    }
    // 添加 Authorization header
    if (token) {
      config.headers['Authorization'] = `Bearer ${token}`
    }
    return config as InternalAxiosRequestConfig
  },
  (error: AxiosError) => {
    ElMessage.error({
      showClose: true,
      message: error.message,
      type: 'error'
    })
    return Promise.reject(error)
  }
)

// 响应拦截器
service.interceptors.response.use(
  (response: AxiosResponse) => {
    // 检查是否有新的 token（自动刷新）
    const newToken = response.headers['X-New-Access-Token']
    if (newToken) {
      console.log('✓ Token 已自动刷新')
      setAccessToken(newToken)
    }

    // 检查业务状态码
    if (response.data.code !== 0) {
      // 检查是否需要重新登录
      if (response.data.data?.reload) {
        ElMessage.error('登录已过期，请重新登录')
        clearAuth()
        // 跳转到首页
        window.location.href = '/'
        return Promise.reject(new Error(response.data.msg))
      }
      
      // 业务错误：显示错误信息并 reject
      ElMessage.error(response.data.msg)
      return Promise.reject(response.data)
    }
    
    // 成功响应：返回 data 部分（如果存在则返回 data.data，否则返回整个 data）
    return response.data.data !== undefined ? response.data.data : response.data
  },
  (error: AxiosError<ApiResponse<any>>) => {
    // HTTP 错误处理
    if (error.response) {
      const { status, data } = error.response
      
      if (status === 401 || status === 403) {
        // 未授权或禁止访问
        if (data?.data?.reload) {
          ElMessage.error('登录已过期，请重新登录')
          clearAuth()
          window.location.href = '/'
        } else {
          ElMessage.error(data?.msg || '无权限访问')
        }
      } else {
        ElMessage.error(data?.msg || error.message || '请求失败')
      }
    } else {
      ElMessage.error(error.message || '网络错误')
    }
    return Promise.reject(error)
  }
)

export default service
