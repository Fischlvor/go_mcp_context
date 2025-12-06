import { ref } from 'vue'

// 响应式 token 状态（独立模块，避免循环依赖）
export const accessToken = ref<string | null>(localStorage.getItem('access_token'))

// 设置 token（同时更新响应式变量和 localStorage）
export const setAccessToken = (token: string): void => {
  accessToken.value = token
  localStorage.setItem('access_token', token)
}

// 获取 token
export const getAccessToken = (): string | null => {
  return accessToken.value
}

// 清除 token
export const clearAccessToken = (): void => {
  accessToken.value = null
  localStorage.removeItem('access_token')
}

// 初始化 token（从 localStorage 恢复）
export const initAccessToken = (): void => {
  const savedToken = localStorage.getItem('access_token')
  if (savedToken) {
    accessToken.value = savedToken
  }
}
