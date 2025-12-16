/**
 * SSE (Server-Sent Events) 统一工具
 * 类似于 axios 的请求拦截器，提供统一的 SSE 连接管理
 */

import { getDeviceId } from './deviceId'

// SSE 事件数据类型
export interface SSEEvent<T = any> {
  type: string
  data: T
}

// SSE 配置选项
export interface SSEOptions {
  // 事件处理器
  onMessage?: (event: SSEEvent) => void
  onError?: (error: Error) => void
  onComplete?: () => void
  // 超时时间（毫秒）
  timeout?: number
}

// 获取 token
const getToken = (): string | null => {
  return localStorage.getItem('access_token')
}

// 获取完整 URL
const getFullUrl = (url: string): string => {
  const baseUrl = import.meta.env.VITE_BASE_API || ''
  return url.startsWith('http') ? url : `${baseUrl}${url}`
}

/**
 * 创建 SSE 连接（GET 请求）
 */
export const createSSE = (url: string, options: SSEOptions = {}): EventSource => {
  const fullUrl = getFullUrl(url)
  const eventSource = new EventSource(fullUrl)

  eventSource.onmessage = (event) => {
    try {
      const data = JSON.parse(event.data)
      options.onMessage?.({ type: data.stage || data.type || 'message', data })
      
      // 检查是否完成
      if (data.stage === 'completed' || data.stage === 'failed' || data.type === 'completed' || data.type === 'failed') {
        eventSource.close()
        options.onComplete?.()
      }
    } catch (e) {
      console.error('SSE parse error:', e)
    }
  }

  eventSource.onerror = (error) => {
    console.error('SSE error:', error)
    eventSource.close()
    options.onError?.(new Error('SSE connection error'))
  }

  return eventSource
}

/**
 * 创建 SSE 连接（POST 请求，使用 fetch + ReadableStream）
 * 用于需要发送请求体的场景（如文件上传）
 */
export const createSSEPost = async <T = any>(
  url: string,
  data: FormData | Record<string, any>,
  options: SSEOptions = {}
): Promise<void> => {
  const fullUrl = getFullUrl(url)
  const token = getToken()
  
  const headers: Record<string, string> = {
    'X-Device-Id': getDeviceId(),
  }
  
  if (token) {
    headers['Authorization'] = `Bearer ${token}`
  }

  // 如果不是 FormData，设置 Content-Type
  if (!(data instanceof FormData)) {
    headers['Content-Type'] = 'application/json'
  }

  try {
    const response = await fetch(fullUrl, {
      method: 'POST',
      headers,
      body: data instanceof FormData ? data : JSON.stringify(data),
    })

    // HTTP 错误处理
    if (!response.ok) {
      let errorMessage = `HTTP ${response.status}: ${response.statusText}`
      
      // 尝试解析响应体获取更详细的错误信息
      try {
        const contentType = response.headers.get('content-type')
        if (contentType?.includes('application/json')) {
          const errorData = await response.json()
          errorMessage = errorData.msg || errorData.message || errorMessage
        } else {
          const text = await response.text()
          if (text) errorMessage = text
        }
      } catch (e) {
        // 忽略解析错误，使用默认错误信息
      }
      
      const error = new Error(errorMessage)
      ;(error as any).status = response.status
      throw error
    }

    const reader = response.body?.getReader()
    if (!reader) {
      throw new Error('No response body')
    }

    const decoder = new TextDecoder()
    let buffer = ''

    while (true) {
      const { done, value } = await reader.read()
      
      if (done) {
        options.onComplete?.()
        break
      }

      buffer += decoder.decode(value, { stream: true })
      
      // 解析 SSE 数据（格式：data: {...}\n\n）
      const lines = buffer.split('\n\n')
      buffer = lines.pop() || '' // 保留未完成的部分

      for (const line of lines) {
        if (line.startsWith('data: ')) {
          try {
            const jsonStr = line.slice(6) // 移除 "data: " 前缀
            const response = JSON.parse(jsonStr)
            
            // 统一格式：{ code: 0, msg: "...", data: { stage, progress, message, ... } }
            const { code, msg, data } = response
            
            // 检查错误（code !== 0）
            if (code !== 0) {
              const error = new Error(msg || 'Upload failed')
              ;(error as any).code = code
              options.onError?.(error)
              return
            }
            
            // 检查失败事件
            if (data?.stage === 'failed') {
              const error = new Error(data.message || msg || 'Processing failed')
              options.onError?.(error)
              return
            }
            
            options.onMessage?.({ 
              type: data?.stage || 'message', 
              data: data 
            })

            // 检查是否完成
            if (data?.stage === 'completed') {
              options.onComplete?.()
              return
            }
          } catch (e) {
            console.error('SSE parse error:', e, line)
          }
        }
      }
    }
  } catch (error) {
    console.error('SSE POST error:', error)
    options.onError?.(error as Error)
  }
}

/**
 * 上传文件并监听 SSE 状态
 */
export const uploadWithSSE = (
  url: string,
  file: File,
  extraData: Record<string, string> = {},
  options: SSEOptions = {}
): Promise<void> => {
  const formData = new FormData()
  formData.append('file', file)
  
  // 添加额外数据
  Object.entries(extraData).forEach(([key, value]) => {
    formData.append(key, value)
  })

  return createSSEPost(url, formData, options)
}
