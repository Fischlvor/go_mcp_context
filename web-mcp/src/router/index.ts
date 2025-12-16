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
    path: '/libraries/:id',
    name: 'library-detail',
    component: () => import('@/views/library/detail.vue'),
    meta: { title: 'Library Detail' }
  },
  {
    path: '/libraries/:id/admin',
    name: 'library-admin',
    component: () => import('@/views/library/admin.vue'),
    meta: { title: 'Library Admin' }
  },
  {
    path: '/libraries/:id/:version',
    name: 'library-version',
    component: () => import('@/views/library/detail.vue'),
    meta: { title: 'Library Version' }
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
