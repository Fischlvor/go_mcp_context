<template>
  <div class="relative flex min-h-screen flex-col overflow-x-hidden bg-stone-50 antialiased">
    <!-- 顶部 Header -->
    <AppHeader 
      :is-logged-in="isLoggedIn" 
      :user-email="userEmail" 
      :user-plan="userPlan"
      @add-docs="showAddDialog" 
      @sign-in="handleSignIn"
    />
    
    <!-- 主内容区 -->
    <main class="flex flex-col gap-4 px-4 pt-10 sm:px-6 md:gap-[19px] md:pt-20">
      <!-- 1. Hero Section -->
      <div class="mx-auto flex w-full max-w-[880px] flex-col gap-1">
        <h1 class="text-left text-lg font-semibold leading-[1.4] tracking-tight sm:text-xl md:text-2xl">
          <span class="text-emerald-600">Up-to-date Docs</span><br>
          <span class="text-stone-700">for LLMs and AI code editors</span>
        </h1>
        <p class="text-left text-sm text-stone-500 sm:text-base md:text-lg">
          Copy latest docs &amp; code — paste into Cursor, Claude, or other LLMs
        </p>
      </div>

      <!-- 2. 搜索区域 -->
      <div class="mx-auto w-full max-w-[880px]">
        <div class="flex flex-col items-center gap-3 md:flex-row md:gap-4">
          <div class="relative w-full md:w-[460px]">
            <input
              v-model="searchQuery"
              type="text"
              aria-label="Search for a library"
              class="h-11 w-full rounded-xl border border-stone-400 bg-white px-4 pr-10 text-sm text-stone-800 shadow-md placeholder:text-stone-400 focus-within:ring-1 focus-within:ring-emerald-600 hover:border-emerald-600 focus:border-emerald-600 focus:outline-none sm:text-base md:h-[50px]"
              placeholder="Search a library (e.g. Next, React)"
              @input="handleSearch"
            />
          </div>
          <span class="text-sm font-normal text-stone-400">or</span>
          <a class="flex h-11 w-full items-center justify-center rounded-xl border border-stone-400 bg-white px-4 text-sm text-stone-600 shadow-md transition-colors hover:border-emerald-600 hover:text-emerald-600 focus:border-emerald-600 sm:text-base md:h-[50px] md:w-auto" href="#">Chat with Docs</a>
        </div>
      </div>

      <!-- 3. 表格区域（Tabs + Table + 底部栏） -->
      <div>
        <!-- Tabs -->
        <div class="mx-auto mt-8 w-full max-w-[880px] md:mt-12">
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

        <!-- Table -->
        <div class="flex justify-center overflow-x-auto">
          <div class="w-full max-w-[880px]">
            <div class="h-full min-w-[280px] sm:min-w-[600px]">
              <table class="w-full table-fixed border-b border-stone-200">
                <thead class="top-0 z-10 border-b border-stone-200">
                  <tr>
                    <th class="w-auto px-2 py-3 sm:w-[150px] sm:px-4"></th>
                    <th class="hidden w-[230px] px-2 py-3 text-left text-sm font-normal uppercase leading-none text-stone-400 sm:table-cell sm:px-4">SOURCE</th>
                    <th class="hidden w-[80px] px-2 py-3 text-right text-sm font-normal uppercase leading-none text-stone-400 sm:table-cell sm:px-4">TOKENS</th>
                    <th class="w-[80px] px-2 py-3 text-right text-sm font-normal uppercase leading-none text-stone-400 sm:w-[100px] sm:px-4">SNIPPETS</th>
                    <th class="hidden w-[115px] px-2 py-3 text-right text-sm font-normal uppercase leading-none text-stone-400 sm:table-cell sm:px-4">UPDATE</th>
                    <th class="hidden w-[30px] px-1 py-3 text-center text-sm font-normal uppercase leading-none text-stone-400 sm:table-cell"></th>
                  </tr>
                </thead>
                <tbody class="divide-y divide-stone-200">
                  <tr v-if="libraries.length === 0 && !loading">
                    <td colspan="6" class="py-16 text-center">
                      <div class="flex flex-col items-center gap-2">
                        <svg xmlns="http://www.w3.org/2000/svg" width="40" height="40" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="1" stroke-linecap="round" stroke-linejoin="round" class="text-stone-300">
                          <rect width="16" height="20" x="4" y="2" rx="2" ry="2"></rect>
                          <path d="M9 22v-4h6v4"></path>
                          <path d="M8 6h.01"></path>
                          <path d="M16 6h.01"></path>
                          <path d="M12 6h.01"></path>
                        </svg>
                        <p class="text-sm font-medium text-stone-500">No libraries yet</p>
                        <p class="text-sm text-stone-400">Add your first library to get started</p>
                      </div>
                    </td>
                  </tr>
                  <tr 
                    v-for="lib in libraries" 
                    :key="lib.id" 
                    class="group cursor-pointer transition-colors hover:bg-white"
                    @click="goToLibrary(lib.id)"
                  >
                    <td class="h-11 px-2 align-middle sm:px-4">
                      <a :title="lib.name" class="block max-w-[auto] truncate text-base font-semibold leading-tight text-emerald-600 hover:text-emerald-500 hover:underline">{{ lib.name }}</a>
                    </td>
                    <td class="hidden h-11 px-2 text-left align-middle text-base font-normal slashed-zero tabular-nums text-stone-800 sm:table-cell sm:px-4 sm:leading-tight">
                      <div class="flex flex-row items-center gap-2">
                        <!-- GitHub Icon -->
                        <span v-if="lib.source_type === 'github'" class="inline-flex flex-shrink-0 items-center justify-center" style="width: 20px; height: 20px;">
                          <svg role="img" viewBox="0 0 16 16" width="18" height="18" xmlns="http://www.w3.org/2000/svg" class="flex-shrink-0">
                            <path fill="currentColor" d="M8 0C3.58 0 0 3.58 0 8c0 3.54 2.29 6.53 5.47 7.59.4.07.55-.17.55-.38 0-.19-.01-.82-.01-1.49-2.01.37-2.53-.49-2.69-.94-.09-.23-.48-.94-.82-1.13-.28-.15-.68-.52-.01-.53.63-.01 1.08.58 1.23.82.72 1.21 1.87.87 2.33.66.07-.52.28-.87.51-1.07-1.78-.2-3.64-.89-3.64-3.95 0-.87.31-1.59.82-2.15-.08-.2-.36-1.02.08-2.12 0 0 .67-.21 2.2.82a7.65 7.65 0 0 1 2-.27c.68.003 1.36.092 2 .27 1.53-1.04 2.2-.82 2.2-.82.44 1.1.16 1.92.08 2.12.51.56.82 1.27.82 2.15 0 3.07-1.87 3.75-3.65 3.95.29.25.54.73.54 1.48 0 1.07-.01 1.93-.01 2.2 0 .21.15.46.55.38A8.013 8.013 0 0 0 16 8c0-4.42-3.58-8-8-8z"></path>
                          </svg>
                        </span>
                        <!-- Website Icon -->
                        <span v-else-if="lib.source_type === 'website'" class="inline-flex flex-shrink-0 items-center justify-center" style="width: 20px; height: 20px;">
                          <svg xmlns="http://www.w3.org/2000/svg" width="18" height="18" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round" class="flex-shrink-0">
                            <circle cx="12" cy="12" r="10"></circle>
                            <path d="M2 12h20"></path>
                            <path d="M12 2a15.3 15.3 0 0 1 4 10 15.3 15.3 0 0 1 -4 10 15.3 15.3 0 0 1 -4 -10 15.3 15.3 0 0 1 4 -10z"></path>
                          </svg>
                        </span>
                        <!-- Local Upload Icon -->
                        <span v-else class="inline-flex flex-shrink-0 items-center justify-center" style="width: 20px; height: 20px;">
                          <svg xmlns="http://www.w3.org/2000/svg" width="18" height="18" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round" class="flex-shrink-0">
                            <path d="M21 15v4a2 2 0 0 1 -2 2H5a2 2 0 0 1 -2 -2v-4"></path>
                            <polyline points="17 8 12 3 7 8"></polyline>
                            <line x1="12" y1="3" x2="12" y2="15"></line>
                          </svg>
                        </span>
                        <!-- Source URL or Label -->
                        <span class="truncate text-stone-800 hover:text-stone-600">
                          <template v-if="lib.source_type === 'github'">
                            {{ lib.source_url || 'GitHub' }}
                          </template>
                          <template v-else-if="lib.source_type === 'website'">
                            {{ lib.source_url || 'Website' }}
                          </template>
                          <template v-else>
                            Local Upload
                          </template>
                        </span>
                      </div>
                    </td>
                    <td class="hidden h-11 px-2 text-right align-middle text-base font-normal slashed-zero tabular-nums leading-tight text-stone-800 sm:table-cell sm:px-4 sm:leading-none">{{ formatNumber(lib.token_count || 0) }}</td>
                    <td class="h-11 px-2 text-right align-middle text-base font-normal slashed-zero tabular-nums leading-tight text-stone-800 sm:px-4 sm:leading-none">{{ formatNumber(lib.chunk_count || 0) }}</td>
                    <td class="hidden h-11 px-2 text-right align-middle text-base font-normal slashed-zero tabular-nums leading-tight text-stone-800 sm:table-cell sm:px-4 sm:leading-none">{{ formatDate(lib.updated_at) }}</td>
                    <td class="hidden h-11 px-1 text-center align-middle sm:table-cell">
                      <div class="flex items-center justify-center">
                        <div class="relative inline-flex items-center justify-center">
                          <svg xmlns="http://www.w3.org/2000/svg" width="24" height="24" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round" class="h-5 w-5 cursor-pointer text-emerald-600">
                            <path d="M12 12m-9 0a9 9 0 1 0 18 0a9 9 0 1 0 -18 0"></path>
                            <path d="M9 12l2 2l4 -4"></path>
                          </svg>
                        </div>
                      </div>
                    </td>
                  </tr>
                </tbody>
              </table>
            </div>
          </div>
        </div>

        <!-- 底部统计栏 -->
        <div class="mx-auto flex w-full max-w-[880px] items-center justify-between bg-stone-100 px-4 py-3 md:px-4">
          <span class="text-sm font-normal uppercase text-stone-400">{{ total }} Libraries</span>
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
            <label class="block text-sm font-medium text-stone-700 mb-1">Description</label>
            <textarea v-model="form.description" rows="3" class="w-full px-3 py-2 rounded-lg border border-stone-300 text-sm resize-none focus:outline-none focus:border-emerald-600 focus:ring-1 focus:ring-emerald-600" placeholder="Brief description of the library"></textarea>
          </div>
          <div>
            <label class="block text-sm font-medium text-stone-700 mb-1">Source Type</label>
            <select v-model="form.source_type" class="w-full h-10 px-3 rounded-lg border border-stone-300 text-sm focus:outline-none focus:border-emerald-600 focus:ring-1 focus:ring-emerald-600">
              <option value="local">Local Upload</option>
              <option value="github">GitHub</option>
              <option value="website">Website</option>
            </select>
          </div>
          <div v-if="form.source_type !== 'local'">
            <label class="block text-sm font-medium text-stone-700 mb-1">Source URL</label>
            <input v-model="form.source_url" type="text" class="w-full h-10 px-3 rounded-lg border border-stone-300 text-sm focus:outline-none focus:border-emerald-600 focus:ring-1 focus:ring-emerald-600" :placeholder="form.source_type === 'github' ? 'e.g. vuejs/docs' : 'e.g. vuejs.org/guide'" />
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

const ElMessage = {
  success: (msg: string) => console.log('✓', msg),
  warning: (msg: string) => alert(msg),
  error: (msg: string) => alert(msg),
}
const ElMessageBox = {
  confirm: (msg: string, title: string, options: any) => Promise.resolve(confirm(msg) ? 'confirm' : Promise.reject()),
}
import { getLibraries, createLibrary, updateLibrary, deleteLibrary } from '@/api/library'
import type { LibraryListItem } from '@/api/library'

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
const libraries = ref<LibraryListItem[]>([])
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
  description: '',
  source_type: 'local',
  source_url: ''
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
  form.description = ''
  form.source_type = 'local'
  form.source_url = ''
  dialogVisible.value = true
}

const handleEdit = (lib: LibraryListItem) => {
  isEdit.value = true
  editId.value = lib.id
  form.name = lib.name
  form.description = ''  // 列表项没有 description，需要从详情获取
  form.source_type = lib.source_type
  form.source_url = lib.source_url
  openDropdownId.value = null
  dialogVisible.value = true
}

const handleDelete = (lib: LibraryListItem) => {
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
  if (!form.name) {
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
  router.push(`/libraries/${id}`)
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
  
  // 如果时间戳无效或是未来时间，显示 'now'
  if (isNaN(date.getTime()) || date > now) {
    return 'just now'
  }
  
  const diff = now.getTime() - date.getTime()
  const minutes = Math.floor(diff / (1000 * 60))
  const hours = Math.floor(diff / (1000 * 60 * 60))
  const days = Math.floor(diff / (1000 * 60 * 60 * 24))
  const weeks = Math.floor(days / 7)
  const months = Math.floor(days / 30)
  const years = Math.floor(days / 365)
  
  // Context7 风格：简洁的数字 + 时间单位
  if (minutes < 1) return 'just now'
  if (minutes < 60) return `${minutes} minute${minutes > 1 ? 's' : ''}`
  if (hours < 24) return `${hours} hour${hours > 1 ? 's' : ''}`
  if (days < 7) return `${days} day${days > 1 ? 's' : ''}`
  if (weeks < 4) return `${weeks} week${weeks > 1 ? 's' : ''}`
  if (months < 12) return `${months} month${months > 1 ? 's' : ''}`
  return `${years} year${years > 1 ? 's' : ''}`
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
/* 页面无需额外样式，全部使用 Tailwind CSS */
</style>
