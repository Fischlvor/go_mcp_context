<template>
  <div class="flex min-h-screen items-center justify-center bg-stone-50">
    <div class="flex flex-col items-center gap-4">
      <!-- Loading Spinner -->
      <div class="h-12 w-12 animate-spin rounded-full border-4 border-emerald-200 border-t-emerald-600"></div>
      <div class="text-lg font-medium text-stone-600">正在登录中...</div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { onMounted } from 'vue'
import { useRouter } from 'vue-router'
import { useUser, fetchUserInfo } from '@/stores/user'
import { handleSSOCallback } from '@/api/auth'

const router = useRouter()
const { login } = useUser()

onMounted(async () => {
  try {
    // 获取 URL 参数
    const urlParams = new URLSearchParams(window.location.search)
    const code = urlParams.get('code')
    const state = urlParams.get('state')

    if (!code) {
      console.error('登录失败：未获取到授权码')
      alert('登录失败：未获取到授权码')
      router.push('/')
      return
    }

    // 解析 state 获取 return_url
    let returnUrl = '/'
    if (state) {
      try {
        const stateObj = JSON.parse(decodeURIComponent(state))
        returnUrl = stateObj.return_url || '/'
      } catch (e) {
        console.error('解析 state 失败:', e)
      }
    }

    // 用 code 向后端换取 token（不需要手动编码，axios 会自动编码）
    const redirectUri = window.location.origin + '/sso-callback'
    const res = await handleSSOCallback(code, redirectUri)

    // 保存 token 到 store
    login({
      email: '',
      plan: 'FREE',
      team: 'Personal'
    }, res.access_token)

    // 获取用户信息
    await fetchUserInfo()

    console.log('登录成功')
    
    // 跳转到 return_url 指定的页面
    router.push(returnUrl)
  } catch (error) {
    console.error('SSO 回调处理失败:', error)
    alert('登录处理失败，请重试')
    router.push('/')
  }
})
</script>
