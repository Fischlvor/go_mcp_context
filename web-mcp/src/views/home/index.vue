<template>
  <div class="relative flex min-h-screen flex-col overflow-x-hidden bg-transparent antialiased">
    <!-- 顶部渐变背景 -->
    <div class="absolute inset-x-0 top-0 h-[260px] bg-gradient-to-b from-emerald-500/[0.15] to-transparent"></div>
    <!-- 底层背景色 -->
    <div class="fixed inset-0 -z-10 bg-stone-50"></div>

    <!-- 顶部 Header -->
    <AppHeader 
      :is-logged-in="isLoggedIn" 
      :user-email="userEmail" 
      :user-plan="userPlan"
      @add-docs="showAddDialog" 
      @sign-in="handleSignIn"
    />
    
    <!-- 主内容区 -->
    <main class="flex-grow pt-0">
      <div class="relative overflow-y-auto overflow-x-hidden pb-0">
        <div class="flex flex-col gap-4 px-4 pt-10 md:gap-[19px] md:pt-20">
          <!-- Hero Section -->
          <div class="mx-auto flex w-full max-w-[880px] flex-col gap-1">
            <h1 class="text-left text-lg font-semibold leading-[1.4] tracking-tight sm:text-xl md:text-2xl">
              <span class="text-emerald-600">Up-to-date Docs</span><br>
              <span class="text-stone-700">for LLMs and AI code editors</span>
            </h1>
            <p class="text-left text-sm text-stone-500 sm:text-base md:text-lg">
              Copy latest <span class="text-emerald-600">docs</span> &amp; <span class="text-emerald-600">code</span> — paste into <span class="text-emerald-600">Cursor</span>, <span class="text-emerald-600">Claude</span>, or other LLMs
            </p>
          </div>

          <!-- 搜索区域 -->
          <div class="mx-auto w-full max-w-[880px]">
            <div class="flex flex-col items-center gap-3 md:flex-row md:gap-4">
              <div class="relative w-full md:w-[460px]">
                <input
                  v-model="searchQuery"
                  type="text"
                  class="h-11 w-full rounded-xl border border-stone-400 bg-white px-4 pr-10 text-sm text-stone-800 shadow-md placeholder:text-stone-400 focus-within:ring-1 focus-within:ring-emerald-600 hover:border-emerald-600 focus:border-emerald-600 focus:outline-none sm:text-base md:h-[50px]"
                  placeholder="Search a library (e.g. Next, React)"
                  @input="handleSearch"
                />
              </div>
              <span class="text-sm font-normal text-stone-400">or</span>
              <a class="flex h-11 w-full items-center justify-center rounded-xl border border-stone-400 bg-white px-4 text-sm text-stone-600 shadow-md transition-colors hover:border-emerald-600 hover:text-emerald-600 focus:border-emerald-600 sm:text-base md:h-[50px] md:w-auto" href="#">Chat with Docs</a>
            </div>
          </div>

          <!-- Tabs + 表格区域 -->
          <div class="mx-auto mt-8 w-full max-w-[880px] md:mt-12">
            <!-- Tabs -->
            <div class="relative flex w-full items-end gap-0">
              <button 
                v-for="tab in tabs" 
                :key="tab.id"
                :class="['flex items-center font-medium gap-1 px-2 py-1.5 text-sm sm:gap-2 sm:px-4 sm:py-2 sm:text-base', activeTab === tab.id ? 'rounded-t-lg border border-stone-300 border-b-transparent text-stone-800' : 'border border-stone-300 border-l-transparent border-r-transparent border-t-transparent text-stone-500 hover:text-stone-600']"
                @click="activeTab = tab.id"
              >
                <svg v-if="tab.id === 'popular'" class="h-4 w-4 sm:h-5 sm:w-5" xmlns="http://www.w3.org/2000/svg" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="1.5" stroke-linecap="round" stroke-linejoin="round">
                  <path d="M12 17.75l-6.172 3.245l1.179 -6.873l-5 -4.867l6.9 -1l3.086 -6.253l3.086 6.253l6.9 1l-5 4.867l1.179 6.873z"></path>
                </svg>
                <svg v-else-if="tab.id === 'trending'" class="h-4 w-4 sm:h-5 sm:w-5" xmlns="http://www.w3.org/2000/svg" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="1.5" stroke-linecap="round" stroke-linejoin="round">
                  <path d="M3 17l6 -6l4 4l8 -8"></path>
                  <path d="M14 7l7 0l0 7"></path>
                </svg>
                <svg v-else class="h-4 w-4 sm:h-5 sm:w-5" xmlns="http://www.w3.org/2000/svg" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="1.5" stroke-linecap="round" stroke-linejoin="round">
                  <path d="M3 12a9 9 0 1 0 18 0a9 9 0 0 0 -18 0"></path>
                  <path d="M12 7v5l3 3"></path>
                </svg>
                {{ tab.label }}
              </button>
              <div class="flex-grow border-b border-stone-300"></div>
            </div>
          </div>

          <!-- 表格 -->
          <div class="flex justify-center overflow-x-auto">
            <div class="w-full max-w-[880px]">
              <div class="h-full min-w-[280px] sm:min-w-[600px]">
                <div class="table-container">
          <table class="library-table">
            <thead>
              <tr>
                <th class="col-name"></th>
                <th class="col-source">SOURCE</th>
                <th class="col-tokens">TOKENS</th>
                <th class="col-snippets">SNIPPETS</th>
                <th class="col-update">UPDATE</th>
                <th class="col-action"></th>
              </tr>
            </thead>
            <tbody>
              <tr v-if="libraries.length === 0 && !loading">
                <td colspan="6" class="empty-state">
                  <div class="empty-content">
                    <svg xmlns="http://www.w3.org/2000/svg" width="40" height="40" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="1" stroke-linecap="round" stroke-linejoin="round" class="empty-icon">
                      <rect width="16" height="20" x="4" y="2" rx="2" ry="2"></rect>
                      <path d="M9 22v-4h6v4"></path>
                      <path d="M8 6h.01"></path>
                      <path d="M16 6h.01"></path>
                      <path d="M12 6h.01"></path>
                      <path d="M12 10h.01"></path>
                      <path d="M12 14h.01"></path>
                      <path d="M16 10h.01"></path>
                      <path d="M16 14h.01"></path>
                      <path d="M8 10h.01"></path>
                      <path d="M8 14h.01"></path>
                    </svg>
                    <p class="empty-text">No libraries yet</p>
                    <p class="empty-hint">Add your first library to get started</p>
                  </div>
                </td>
              </tr>
              <tr 
                v-for="lib in libraries" 
                :key="lib.id" 
                class="table-row"
                @click="goToLibrary(lib.id)"
              >
                <td class="col-name">
                  <span class="lib-name">{{ lib.name }}</span>
                </td>
                <td class="col-source">
                  <div class="source-info">
                    <svg xmlns="http://www.w3.org/2000/svg" width="14" height="14" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round" class="source-icon">
                      <path d="M15 22v-4a4.8 4.8 0 0 0-1-3.5c3 0 6-2 6-5.5.08-1.25-.27-2.48-1-3.5.28-1.15.28-2.35 0-3.5 0 0-1 0-3 1.5-2.64-.5-5.36-.5-8 0C6 2 5 2 5 2c-.3 1.15-.3 2.35 0 3.5A5.403 5.403 0 0 0 4 9c0 3.5 3 5.5 6 5.5-.39.49-.68 1.05-.85 1.65-.17.6-.22 1.23-.15 1.85v4"></path>
                      <path d="M9 18c-4.51 2-5-2-7-2"></path>
                    </svg>
                    <span class="source-text">/{{ lib.name.toLowerCase().replace(/\s+/g, '-') }}/{{ lib.version || 'docs' }}</span>
                  </div>
                </td>
                <td class="col-tokens">{{ formatNumber(lib.chunk_count || 0) }}</td>
                <td class="col-snippets">{{ formatNumber(lib.document_count || 0) }}</td>
                <td class="col-update">{{ formatDate(lib.updated_at) }}</td>
                <td class="col-action">
                  <button class="arrow-btn" @click.stop="goToLibrary(lib.id)">
                    <svg xmlns="http://www.w3.org/2000/svg" width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round">
                      <circle cx="12" cy="12" r="10"></circle>
                      <path d="m10 8 4 4-4 4"></path>
                    </svg>
                  </button>
                </td>
              </tr>
            </tbody>
          </table>
                </div>
              </div>
            </div>
          </div>
        </div>
      </div>
    </main>

    <!-- Footer -->
    <AppFooter />

    <!-- 添加/编辑库对话框 -->
    <div v-if="dialogVisible" class="fixed inset-0 z-50 flex items-center justify-center bg-black/50 backdrop-blur-sm p-4" @click.self="dialogVisible = false">
      <div class="w-full max-w-md rounded-xl bg-white shadow-2xl">
        <div class="flex items-center justify-between border-b border-stone-200 px-6 py-4">
          <h3 class="text-lg font-semibold text-stone-900">{{ isEdit ? 'Edit Library' : 'Add New Library' }}</h3>
          <button class="rounded-lg p-1 text-stone-400 hover:bg-stone-100 hover:text-stone-600" @click="dialogVisible = false">
            <svg xmlns="http://www.w3.org/2000/svg" width="20" height="20" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round">
              <path d="M18 6 6 18"></path>
              <path d="m6 6 12 12"></path>
            </svg>
          </button>
        </div>
        <div class="space-y-4 px-6 py-4">
          <div>
            <label class="block text-sm font-medium text-stone-700 mb-1">Library Name</label>
            <input v-model="form.name" type="text" class="w-full h-10 px-3 rounded-lg border border-stone-300 text-sm focus:outline-none focus:border-emerald-600 focus:ring-1 focus:ring-emerald-600" placeholder="e.g. Vue.js" />
          </div>
          <div>
            <label class="block text-sm font-medium text-stone-700 mb-1">Version</label>
            <input v-model="form.version" type="text" class="w-full h-10 px-3 rounded-lg border border-stone-300 text-sm focus:outline-none focus:border-emerald-600 focus:ring-1 focus:ring-emerald-600" placeholder="e.g. 3.4.0" />
          </div>
          <div>
            <label class="block text-sm font-medium text-stone-700 mb-1">Description</label>
            <textarea v-model="form.description" rows="3" class="w-full px-3 py-2 rounded-lg border border-stone-300 text-sm resize-none focus:outline-none focus:border-emerald-600 focus:ring-1 focus:ring-emerald-600" placeholder="Brief description of the library"></textarea>
          </div>
        </div>
        <div class="flex justify-end gap-3 border-t border-stone-200 px-6 py-4">
          <button class="h-10 px-4 rounded-lg border border-stone-300 text-sm font-medium text-stone-700 hover:bg-stone-50" @click="dialogVisible = false">Cancel</button>
          <button class="h-10 px-4 rounded-lg bg-emerald-600 text-sm font-medium text-white hover:bg-emerald-700 disabled:opacity-50" @click="handleSubmit" :disabled="submitting">
            {{ isEdit ? 'Update' : 'Create' }}
          </button>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive, onMounted, onUnmounted } from 'vue'
import { useRouter } from 'vue-router'
import AppHeader from '@/components/AppHeader.vue'
import AppFooter from '@/components/AppFooter.vue'
import { useUser } from '@/stores/user'
// TODO: 替换为 shadcn-vue 的 toast 组件
const ElMessage = {
  success: (msg: string) => alert(msg),
  warning: (msg: string) => alert(msg),
  error: (msg: string) => alert(msg),
}
const ElMessageBox = {
  confirm: (msg: string, title: string, options: any) => Promise.resolve(confirm(msg) ? 'confirm' : Promise.reject()),
}
import { getLibraries, createLibrary, updateLibrary, deleteLibrary } from '@/api/library'
import type { Library } from '@/api/library'

const router = useRouter()

// 用户状态
const { isLoggedIn, userEmail, userPlan, initUserState, redirectToSSO } = useUser()

// 初始化用户状态
onMounted(() => {
  initUserState()
})

// 登录处理：跳转到 SSO 登录
const handleSignIn = () => {
  redirectToSSO()
}

// Tabs 配置
const tabs = [
  { id: 'popular', label: 'Popular' },
  { id: 'trending', label: 'Trending' },
  { id: 'recent', label: 'Recent' }
]
const activeTab = ref('popular')

// 数据状态
const loading = ref(false)
const libraries = ref<Library[]>([])
const page = ref(1)
const pageSize = ref(20)
const total = ref(0)
const searchQuery = ref('')

// 对话框状态
const dialogVisible = ref(false)
const isEdit = ref(false)
const editId = ref<number | null>(null)
const submitting = ref(false)

// 下拉菜单状态
const openDropdownId = ref<number | null>(null)

const form = reactive({
  name: '',
  version: '',
  description: ''
})

let searchTimer: ReturnType<typeof setTimeout> | null = null

const fetchLibraries = async () => {
  loading.value = true
  try {
    const res = await getLibraries({
      name: searchQuery.value || undefined,
      page: page.value,
      page_size: pageSize.value
    })
    if (res.code === 0) {
      libraries.value = res.data.list || []
      total.value = res.data.total
    }
  } finally {
    loading.value = false
  }
}

const handleSearch = () => {
  if (searchTimer) clearTimeout(searchTimer)
  searchTimer = setTimeout(() => {
    page.value = 1
    fetchLibraries()
  }, 300)
}

const clearSearch = () => {
  searchQuery.value = ''
  page.value = 1
  fetchLibraries()
}

const showAddDialog = () => {
  isEdit.value = false
  editId.value = null
  form.name = ''
  form.version = ''
  form.description = ''
  dialogVisible.value = true
}

const handleEdit = (lib: Library) => {
  isEdit.value = true
  editId.value = lib.id
  form.name = lib.name
  form.version = lib.version
  form.description = lib.description
  openDropdownId.value = null
  dialogVisible.value = true
}

const handleDelete = (lib: Library) => {
  openDropdownId.value = null
  ElMessageBox.confirm(
    `Are you sure to delete "${lib.name}"?`,
    'Delete Library',
    {
      confirmButtonText: 'Delete',
      cancelButtonText: 'Cancel',
      type: 'warning'
    }
  ).then(async () => {
    await deleteLibrary(lib.id)
    ElMessage.success('Library deleted')
    fetchLibraries()
  }).catch(() => {})
}

const handleSubmit = async () => {
  if (!form.name || !form.version) {
    ElMessage.warning('Please fill in required fields')
    return
  }
  
  submitting.value = true
  try {
    if (isEdit.value && editId.value) {
      await updateLibrary(editId.value, form)
      ElMessage.success('Library updated')
    } else {
      await createLibrary(form)
      ElMessage.success('Library created')
    }
    dialogVisible.value = false
    fetchLibraries()
  } finally {
    submitting.value = false
  }
}

const toggleDropdown = (id: number) => {
  openDropdownId.value = openDropdownId.value === id ? null : id
}

const closeDropdown = (e: MouseEvent) => {
  const target = e.target as HTMLElement
  if (!target.closest('.action-dropdown')) {
    openDropdownId.value = null
  }
}

const goToLibrary = (id: number) => {
  router.push(`/libraries/${id}/documents`)
}

const formatNumber = (num: number) => {
  if (num >= 1000000) return (num / 1000000).toFixed(1) + 'M'
  if (num >= 1000) return (num / 1000).toFixed(1) + 'K'
  return num.toString()
}

const formatDate = (dateStr: string) => {
  if (!dateStr) return '-'
  const date = new Date(dateStr)
  const now = new Date()
  const diff = now.getTime() - date.getTime()
  const days = Math.floor(diff / (1000 * 60 * 60 * 24))
  
  if (days === 0) return 'today'
  if (days === 1) return '1 day'
  if (days < 7) return `${days} days`
  if (days < 30) return `${Math.floor(days / 7)} weeks`
  if (days < 365) return `${Math.floor(days / 30)} months`
  return `${Math.floor(days / 365)} years`
}

onMounted(() => {
  fetchLibraries()
  document.addEventListener('click', closeDropdown)
})

onUnmounted(() => {
  document.removeEventListener('click', closeDropdown)
})
</script>

<style scoped>
/* ========== 表格样式 (保留，因为表格内容复杂) ========== */

.table-container {
  overflow-x: auto;
}

.library-table {
  width: 100%;
  min-width: 600px;
  border-collapse: collapse;
  table-layout: fixed;
}

.library-table thead tr {
  border-bottom: 1px solid #e7e5e4;
}

.library-table th {
  padding: 10px 16px;
  font-size: 11px;
  font-weight: 500;
  text-transform: uppercase;
  color: #a8a29e;
  text-align: left;
  letter-spacing: 0.3px;
}

.library-table th.col-name {
  width: 120px;
}

.library-table th.col-source {
  width: auto;
}

.library-table th.col-tokens,
.library-table th.col-snippets {
  width: 80px;
  text-align: right;
}

.library-table th.col-update {
  width: 80px;
  text-align: right;
}

.library-table th.col-action {
  width: 36px;
  text-align: center;
}

.library-table tbody tr {
  transition: background-color 0.15s;
  border-bottom: 1px solid #f5f5f4;
}

.library-table tbody tr:last-child {
  border-bottom: none;
}

.library-table tbody tr:hover {
  background: #fafaf9;
}

.library-table tbody tr.table-row {
  cursor: pointer;
}

.library-table td {
  padding: 12px 16px;
  font-size: 14px;
  color: #1c1917;
  vertical-align: middle;
}

.library-table td.col-tokens,
.library-table td.col-snippets,
.library-table td.col-update {
  text-align: right;
  color: #78716c;
  font-size: 14px;
}

.lib-name {
  font-size: 15px;
  font-weight: 600;
  color: #059669;
  transition: text-decoration 0.15s;
}

.lib-name:hover {
  text-decoration: underline;
  text-underline-offset: 2px;
}

.source-info {
  display: flex;
  align-items: center;
  gap: 8px;
}

.source-icon {
  flex-shrink: 0;
  color: #78716c;
}

.source-text {
  color: #57534e;
  font-size: 14px;
  font-family: ui-monospace, SFMono-Regular, Menlo, Monaco, Consolas, monospace;
}

/* ========== 箭头按钮 ========== */
.arrow-btn {
  display: flex;
  align-items: center;
  justify-content: center;
  width: 28px;
  height: 28px;
  border: none;
  border-radius: 50%;
  background: transparent;
  color: #059669;
  cursor: pointer;
  transition: background-color 0.15s;
}

.arrow-btn:hover {
  background: #ecfdf5;
}


/* ========== 空状态 ========== */
.empty-state {
  padding: 60px 20px;
  text-align: center;
}

.empty-content {
  display: flex;
  flex-direction: column;
  align-items: center;
  gap: 4px;
}

.empty-icon {
  color: #d6d3d1;
  margin-bottom: 12px;
  width: 48px;
  height: 48px;
}

.empty-text {
  font-size: 14px;
  font-weight: 500;
  color: #57534e;
  margin: 0;
}

.empty-hint {
  font-size: 13px;
  color: #a8a29e;
  margin: 0;
}

</style>
