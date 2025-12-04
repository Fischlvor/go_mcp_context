import { createRouter, createWebHistory } from 'vue-router'

const routes = [
  {
    path: '/',
    name: 'home',
    component: () => import('@/views/home/index.vue'),
    meta: { title: 'MCP Context' }
  },
  {
    path: '/dashboard',
    name: 'dashboard',
    component: () => import('@/views/dashboard/index.vue'),
    meta: { title: 'Dashboard' }
  },
  {
    path: '/libraries/:id/documents',
    name: 'documents',
    component: () => import('@/views/document/index.vue'),
    meta: { title: '文档管理' }
  },
  {
    path: '/search',
    name: 'search',
    component: () => import('@/views/search/index.vue'),
    meta: { title: '搜索测试' }
  },
  {
    path: '/sso-callback',
    name: 'sso-callback',
    component: () => import('@/views/SSOCallback.vue'),
    meta: { title: '登录中...' }
  }
]

const router = createRouter({
  history: createWebHistory(import.meta.env.BASE_URL),
  routes
})

export default router
