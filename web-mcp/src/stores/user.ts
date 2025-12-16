import { ref, computed } from 'vue'
import request from '@/utils/request'
import { accessToken, setAccessToken, clearAccessToken, initAccessToken } from '@/utils/token'

// 用户信息接口
export interface User {
  email: string
  plan: 'FREE' | 'PRO' | 'TEAM'
  team?: string
}

// 用户状态
const user = ref<User | null>(null)

// 初始化标记，确保只初始化一次
let initPromise: Promise<void> | null = null

// 是否已登录
export const isLoggedIn = computed(() => user.value !== null && accessToken.value !== null)

// 获取用户信息
export const currentUser = computed(() => user.value)

// 获取用户邮箱
export const userEmail = computed(() => user.value?.email || '')

// 获取用户计划
export const userPlan = computed(() => user.value?.plan || 'FREE')

// 获取当前团队
export const currentTeam = computed(() => user.value?.team || 'Personal')

// 获取 access_token
export const getAccessToken = computed(() => accessToken.value)

// 登录
export function login(userData: User, token?: string) {
  user.value = userData
  if (token) {
    setAccessToken(token)
  }
  localStorage.setItem('user', JSON.stringify(userData))
}

// 登出
export async function logout() {
  // 调用后端登出接口
  try {
    await request.post('/auth/logout')
  } catch (e) {
    console.error('登出请求失败:', e)
  }
  
  user.value = null
  clearAccessToken()
  localStorage.removeItem('user')
  initPromise = null // 重置初始化标记
}

// 初始化：从 localStorage 恢复用户状态，并从 SSO 获取最新用户信息
export async function initUserState() {
  // 如果已经在初始化中，等待初始化完成
  if (initPromise) {
    return initPromise
  }
  
  // 创建初始化 Promise
  initPromise = (async () => {
    initAccessToken()
    
    if (accessToken.value) {
      // 从 SSO 获取最新用户信息
      await fetchUserInfo()
    }
  })()
  
  return initPromise
}

// 从后端获取用户信息
export async function fetchUserInfo() {
  if (!accessToken.value) return
  
  try {
    const res: any = await request.get('/user/info')
    if (res.code === 0 && res.data) {
      const userInfo = res.data
      user.value = {
        email: userInfo.email || userInfo.nickname || '',
        plan: 'FREE',
        team: userInfo.nickname || 'Personal'
      }
      localStorage.setItem('user', JSON.stringify(user.value))
    } else {
      // token 无效，清除登录状态
      console.error('获取用户信息失败:', res.msg)
      user.value = null
      clearAccessToken()
      localStorage.removeItem('user')
    }
  } catch (e) {
    console.error('获取用户信息失败:', e)
    // 网络错误时不清除状态，使用缓存的用户信息
    const savedUser = localStorage.getItem('user')
    if (savedUser) {
      try {
        user.value = JSON.parse(savedUser)
      } catch (err) {
        // ignore
      }
    }
  }
}

// 切换团队
export function switchTeam(teamName: string) {
  if (user.value) {
    user.value = { ...user.value, team: teamName }
    localStorage.setItem('user', JSON.stringify(user.value))
  }
}

// 跳转到 SSO 登录
export async function redirectToSSO(returnUrl?: string) {
  try {
    const redirectUri = encodeURIComponent(window.location.origin + '/sso-callback')
    const returnPath = encodeURIComponent(returnUrl || window.location.pathname)
    
    const res: any = await request.get(`/auth/sso_login_url?redirect_uri=${redirectUri}&return_url=${returnPath}`)
    
    if (res.code === 0 && res.data.sso_login_url) {
      window.location.href = res.data.sso_login_url
    } else {
      console.error('获取 SSO 登录地址失败:', res.msg)
      alert('获取登录地址失败，请稍后重试')
    }
  } catch (error) {
    console.error('获取 SSO 登录 URL 失败:', error)
    alert('登录服务异常，请稍后重试')
  }
}

// 导出 useUser composable
export function useUser() {
  return {
    user: currentUser,
    isLoggedIn,
    userEmail,
    userPlan,
    currentTeam,
    accessToken: getAccessToken,
    login,
    logout,
    initUserState,
    fetchUserInfo,
    switchTeam,
    redirectToSSO
  }
}
