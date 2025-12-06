<template>
  <div class="relative flex min-h-screen flex-col overflow-x-hidden bg-stone-50 antialiased">
    <!-- 顶部 Header -->
    <AppHeader 
      :is-logged-in="isLoggedIn" 
      :user-email="userEmail" 
      :user-plan="userPlan"
      @sign-in="handleSignIn"
    />

    <!-- 主内容区 -->
    <main class="flex-grow pt-0">
      <div class="mx-auto flex w-full max-w-[880px] flex-col items-center justify-center px-0">
        <div class="mx-auto flex w-full max-w-[880px] flex-col px-4 pt-10 lg:px-0">
          
          <!-- 库信息卡片 -->
          <div class="w-full rounded-3xl border-2 border-emerald-600 bg-white p-6 shadow-sm sm:p-8">
            <div class="flex flex-col space-y-5">
              <!-- 标题行 -->
              <div class="flex items-start justify-between gap-4">
                <div class="flex min-w-0 flex-1 flex-col gap-1">
                  <h2 class="flex items-center gap-2 text-xl font-semibold leading-[100%] tracking-[0%] text-stone-800">
                    {{ library.name }}
                  </h2>
                  <div class="w-fit max-w-full">
                    <span class="block overflow-hidden text-ellipsis whitespace-nowrap text-base font-normal leading-normal text-stone-500">
                      Version: {{ library.version }}
                    </span>
                  </div>
                  <h3 class="text-base font-normal leading-[140%] text-stone-500">
                    {{ library.description || 'No description' }}
                  </h3>
                </div>
                <!-- Manage 按钮 -->
                <div class="relative inline-flex">
                  <router-link 
                    :to="`/libraries/${libraryId}/admin`"
                    class="flex h-8 items-center justify-center gap-1.5 rounded-lg border border-stone-300 text-base text-stone-500 transition hover:border-stone-400 px-3 py-2"
                  >
                    <svg xmlns="http://www.w3.org/2000/svg" width="24" height="24" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="1.5" stroke-linecap="round" stroke-linejoin="round" class="h-5 w-5 text-stone-500">
                      <path d="M10.325 4.317c.426 -1.756 2.924 -1.756 3.35 0a1.724 1.724 0 0 0 2.573 1.066c1.543 -.94 3.31 .826 2.37 2.37a1.724 1.724 0 0 0 1.065 2.572c1.756 .426 1.756 2.924 0 3.35a1.724 1.724 0 0 0 -1.066 2.573c.94 1.543 -.826 3.31 -2.37 2.37a1.724 1.724 0 0 0 -2.572 1.065c-.426 1.756 -2.924 1.756 -3.35 0a1.724 1.724 0 0 0 -2.573 -1.066c-1.543 .94 -3.31 -.826 -2.37 -2.37a1.724 1.724 0 0 0 -1.065 -2.572c-1.756 -.426 -1.756 -2.924 0 -3.35a1.724 1.724 0 0 0 1.066 -2.573c-.94 -1.543 .826 -3.31 2.37 -2.37c1 .608 2.296 .07 2.572 -1.065z"></path>
                      <path d="M9 12a3 3 0 1 0 6 0a3 3 0 0 0 -6 0"></path>
                    </svg>
                    <span>Manage</span>
                  </router-link>
                </div>
              </div>

              <!-- 状态标签 -->
              <div class="flex flex-col-reverse gap-4 sm:flex-row sm:flex-wrap sm:items-start sm:justify-between">
                <div class="flex flex-wrap gap-2 text-sm sm:gap-1">
                  <div class="flex items-center gap-1 rounded-md bg-emerald-50 px-2 py-1">
                    <svg xmlns="http://www.w3.org/2000/svg" width="24" height="24" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="1.5" stroke-linecap="round" stroke-linejoin="round" class="h-5 w-5 text-emerald-800">
                      <path d="M12 12m-9 0a9 9 0 1 0 18 0a9 9 0 1 0 -18 0"></path>
                      <path d="M9 12l2 2l4 -4"></path>
                    </svg>
                    <span class="text-sm font-normal leading-[100%] tracking-[0%] text-emerald-800">{{ library.status === 'active' ? 'Active' : library.status }}</span>
                  </div>
                  <div class="flex items-center gap-1 rounded-md bg-stone-100 px-3 py-1.5">
                    <span class="text-sm font-normal leading-[100%] tracking-[0%] text-stone-500">Tokens:</span>
                    <span class="text-sm font-normal leading-[100%] tracking-[0%] text-stone-800">{{ formatNumber(library.token_count || 0) }}</span>
                  </div>
                  <div class="flex items-center gap-1 rounded-md bg-stone-100 px-3 py-1.5">
                    <span class="text-sm font-normal leading-[100%] tracking-[0%] text-stone-500">Documents:</span>
                    <span class="text-sm font-normal leading-[100%] tracking-[0%] text-stone-800">{{ formatNumber(library.document_count || 0) }}</span>
                  </div>
                  <div class="flex items-center gap-1 rounded-md bg-stone-100 px-3 py-1.5">
                    <span class="text-sm font-normal leading-[100%] tracking-[0%] text-stone-500">Update:</span>
                    <span class="text-sm font-normal leading-[100%] tracking-[0%] text-stone-800">{{ formatDate(library.updated_at) }}</span>
                  </div>
                </div>
              </div>
            </div>
          </div>

          <!-- Tabs 区域 - 在绿色卡片外面 -->
          <div class="mt-6">
            <div class="flex flex-col-reverse gap-2 sm:flex-row sm:items-start sm:justify-between">
              <div class="overflow-x-auto overflow-y-hidden sm:overflow-visible">
                <div class="relative flex flex-nowrap items-end gap-1">
                  <button 
                    :class="[
                      '-mb-px flex flex-shrink-0 items-center gap-2 whitespace-nowrap rounded-t-lg px-4 py-2 text-base font-medium',
                      activeTab === 'context' 
                        ? 'relative z-10 border border-stone-300 border-b-stone-50 bg-stone-50 text-stone-800' 
                        : 'border border-stone-300 border-b-transparent text-stone-500 hover:border-stone-400 hover:text-stone-600'
                    ]"
                    @click="activeTab = 'context'"
                  >
                    <svg xmlns="http://www.w3.org/2000/svg" width="20" height="20" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="1.5" stroke-linecap="round" stroke-linejoin="round">
                      <path d="M12 3l8 4.5l0 9l-8 4.5l-8 -4.5l0 -9l8 -4.5"></path>
                      <path d="M12 12l8 -4.5"></path>
                      <path d="M12 12l0 9"></path>
                      <path d="M12 12l-8 -4.5"></path>
                    </svg>
                    Context
                  </button>
                </div>
              </div>
              <!-- 工具栏 -->
              <div class="flex flex-wrap gap-2.5 sm:gap-1.5">
                <button 
                  class="flex h-8 items-center justify-center gap-1.5 rounded-lg border border-stone-300 text-base text-stone-500 transition hover:border-stone-400 w-8"
                  @click="fetchDocument"
                >
                  <svg xmlns="http://www.w3.org/2000/svg" width="24" height="24" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="1.5" stroke-linecap="round" stroke-linejoin="round" class="h-5 w-5 text-stone-500">
                    <path d="M20 11a8.1 8.1 0 0 0 -15.5 -2m-.5 -4v4h4"></path>
                    <path d="M4 13a8.1 8.1 0 0 0 15.5 2m.5 4v-4h-4"></path>
                  </svg>
                </button>
              </div>
            </div>
            <div class="border-t border-stone-300"></div>
          </div>

          <!-- Context Tab 内容 -->
          <div v-if="activeTab === 'context'" class="mt-8">
            <div class="flex flex-col gap-8">
              <!-- 搜索卡片 -->
              <div class="w-full rounded-3xl border border-stone-300 bg-white p-6 shadow-sm sm:p-8">
                <div class="flex w-full flex-col gap-1">
                  <label class="text-sm font-medium leading-[100%] tracking-[0%] text-stone-800 md:text-[16px]">Show doc for...</label>
                  <div class="flex flex-col gap-2 sm:flex-row sm:items-center">
                    <input 
                      v-model="searchTopic"
                      placeholder="e.g. data fetching, routing, middleware" 
                      class="h-[40px] w-full flex-1 rounded-lg border border-stone-300 bg-white px-3 py-2 text-sm text-stone-800 hover:border-emerald-600 focus:border-emerald-600 focus:outline-none focus:ring-1 focus:ring-emerald-600 md:text-[16px]"
                      @keyup.enter="handleSearch"
                    />
                    <div class="inline-flex items-center justify-center rounded-lg h-[40px] gap-1 border border-stone-300 bg-white p-1">
                      <button 
                        :class="[
                          'inline-flex items-center justify-center whitespace-nowrap rounded-md px-3 py-1 h-[32px] text-sm font-normal md:text-[16px]',
                          searchMode === 'code' ? 'bg-stone-200 text-stone-800 shadow-sm' : 'text-stone-600 hover:bg-stone-100'
                        ]"
                        @click="searchMode = 'code'"
                      >Code</button>
                      <button 
                        :class="[
                          'inline-flex items-center justify-center whitespace-nowrap rounded-md px-3 py-1 h-[32px] text-sm font-normal md:text-[16px]',
                          searchMode === 'info' ? 'bg-stone-200 text-stone-800 shadow-sm' : 'text-stone-600 hover:bg-stone-100'
                        ]"
                        @click="searchMode = 'info'"
                      >Info</button>
                    </div>
                    <button 
                      class="flex h-[40px] min-w-[130px] items-center justify-center gap-1 whitespace-nowrap rounded-lg bg-stone-200 px-3 text-sm font-normal leading-[100%] tracking-[0%] text-stone-600 hover:bg-stone-300 disabled:cursor-not-allowed disabled:opacity-50 md:text-[16px]"
                      :disabled="searching"
                      @click="handleSearch"
                    >
                      {{ searching ? 'Searching...' : 'Show Results' }}
                    </button>
                  </div>
                </div>
              </div>

              <!-- 结果卡片 -->
              <div class="w-full rounded-3xl border border-stone-300 bg-white p-6 shadow-sm sm:p-8">
                <div class="mb-4 flex flex-col flex-wrap items-start justify-between gap-3 sm:flex-row sm:items-center">
                  <div class="flex items-center gap-2">
                    <span v-if="hasSearched && searchResults.length > 0" class="text-sm text-stone-500">
                      {{ searchResults.length }} results
                    </span>
                  </div>
                  <div class="flex h-8 w-full flex-wrap gap-[1px] overflow-hidden rounded-lg sm:w-auto">
                    <button 
                      class="flex h-8 flex-1 items-center justify-center gap-1 bg-stone-200 px-3 text-sm font-normal leading-[100%] tracking-[0%] text-stone-600 hover:bg-stone-300 sm:flex-initial md:text-base"
                      @click="copyContent"
                    >
                      <svg xmlns="http://www.w3.org/2000/svg" width="20" height="20" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="1.5" stroke-linecap="round" stroke-linejoin="round">
                        <path d="M7 7m0 2.667a2.667 2.667 0 0 1 2.667 -2.667h8.666a2.667 2.667 0 0 1 2.667 2.667v8.666a2.667 2.667 0 0 1 -2.667 2.667h-8.666a2.667 2.667 0 0 1 -2.667 -2.667z"></path>
                        <path d="M4.012 16.737a2.005 2.005 0 0 1 -1.012 -1.737v-10c0 -1.1 .9 -2 2 -2h10c.75 0 1.158 .385 1.5 1"></path>
                      </svg>
                      Copy
                    </button>
                  </div>
                </div>
                
                <!-- 搜索结果列表 -->
                <div v-if="searchResults.length > 0" class="space-y-6 max-h-[500px] overflow-y-auto">
                  <div 
                    v-for="result in searchResults" 
                    :key="result.chunk_id"
                    class="border-b border-stone-200 pb-6 last:border-b-0"
                  >
                    <!-- 标题 -->
                    <h3 v-if="result.title" class="text-lg font-semibold text-stone-800 mb-2">
                      {{ result.title }}
                    </h3>
                    
                    <!-- 来源和元信息 -->
                    <div class="flex flex-wrap items-center gap-3 text-sm text-stone-500 mb-3">
                      <span v-if="result.source" class="flex items-center gap-1">
                        <svg xmlns="http://www.w3.org/2000/svg" width="14" height="14" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="1.5" stroke-linecap="round" stroke-linejoin="round">
                          <path d="M14 3v4a1 1 0 0 0 1 1h4"></path>
                          <path d="M17 21h-10a2 2 0 0 1 -2 -2v-14a2 2 0 0 1 2 -2h7l5 5v11a2 2 0 0 1 -2 2z"></path>
                        </svg>
                        {{ result.source }}
                      </span>
                      <span class="flex items-center gap-1">
                        <svg xmlns="http://www.w3.org/2000/svg" width="14" height="14" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="1.5" stroke-linecap="round" stroke-linejoin="round">
                          <path d="M4 7h16"></path>
                          <path d="M4 11h16"></path>
                          <path d="M4 15h16"></path>
                          <path d="M4 19h16"></path>
                        </svg>
                        {{ result.tokens }} tokens
                      </span>
                      <span class="flex items-center gap-1 text-emerald-600">
                        {{ (result.relevance * 100).toFixed(0) }}% match
                      </span>
                    </div>
                    
                    <!-- 内容 -->
                    <div class="rounded-lg bg-stone-50 p-4 overflow-x-auto">
                      <pre class="whitespace-pre-wrap font-mono text-sm text-stone-700">{{ result.content }}</pre>
                    </div>
                  </div>
                </div>
                
                <!-- 无搜索结果或初始状态 -->
                <div v-else class="overflow-hidden rounded-xl">
                  <textarea 
                    readonly 
                    class="h-[250px] w-full overflow-auto bg-stone-100 p-3 align-top font-mono text-xs text-stone-800 focus:outline-none sm:h-[350px] md:h-[434px] md:p-5 md:text-sm" 
                    spellcheck="false"
                    :value="searchResult"
                  ></textarea>
                </div>
              </div>
            </div>
          </div>
        </div>
      </div>
    </main>

    <!-- Footer -->
    <AppFooter />
  </div>
</template>

<script setup lang="ts">
import { ref, computed, onMounted, watch } from 'vue'
import { useRoute } from 'vue-router'
import AppHeader from '@/components/AppHeader.vue'
import AppFooter from '@/components/AppFooter.vue'
import { useUser } from '@/stores/user'
import { getLibrary } from '@/api/library'
import { getLatestCode } from '@/api/document'
import { searchDocuments } from '@/api/search'
import type { Library } from '@/api/library'
import type { SearchResultItem } from '@/api/search'

const route = useRoute()
const { isLoggedIn, userEmail, userPlan, initUserState, redirectToSSO } = useUser()

const libraryId = computed(() => Number(route.params.id))
const documentTitle = computed(() => route.params.title as string | undefined)
const library = ref<Library>({
  id: 0,
  name: '',
  version: '',
  description: '',
  status: '',
  document_count: 0,
  chunk_count: 0,
  token_count: 0,
  created_at: '',
  updated_at: ''
})

const activeTab = ref('context')
const searchTopic = ref('')
const searchMode = ref<'code' | 'info'>('code')
const searching = ref(false)
const searchResult = ref('Loading document...')
const searchResults = ref<SearchResultItem[]>([])
const hasSearched = ref(false)
const currentDocTitle = ref('')
const loadingDoc = ref(false)

const handleSignIn = () => {
  redirectToSSO()
}

const fetchLibrary = async () => {
  const res = await getLibrary(libraryId.value)
  if (res.code === 0) {
    library.value = res.data
  }
}

// 加载文档内容
const fetchDocument = async () => {
  loadingDoc.value = true
  searchResult.value = 'Loading document...'
  
  try {
    // 直接调用新接口获取最新文档内容
    const res = await getLatestCode(libraryId.value)
    if (res.code === 0) {
      currentDocTitle.value = res.data.title
      searchResult.value = res.data.content || 'No content available.'
    } else {
      searchResult.value = 'No documents available. Upload a document to get started.'
    }
  } catch (error) {
    searchResult.value = 'Failed to load document.'
  } finally {
    loadingDoc.value = false
  }
}

const handleSearch = async () => {
  if (!searchTopic.value.trim()) {
    searchResult.value = 'Please enter a topic to search.'
    searchResults.value = []
    hasSearched.value = false
    return
  }
  
  searching.value = true
  hasSearched.value = true
  searchResults.value = []
  searchResult.value = 'Searching...'
  
  try {
    const res = await searchDocuments({
      library_id: libraryId.value,
      query: searchTopic.value,
      mode: searchMode.value,
      limit: 10
    })
    
    if (res.code === 0 && res.data.results.length > 0) {
      searchResults.value = res.data.results
      searchResult.value = '' // 清空旧的文本结果
    } else {
      searchResults.value = []
      searchResult.value = `No results found for "${searchTopic.value}".`
    }
  } catch (error) {
    searchResults.value = []
    searchResult.value = 'Search failed. Please try again.'
  } finally {
    searching.value = false
  }
}

const copyContent = () => {
  // 如果有搜索结果，复制格式化的结果
  if (searchResults.value.length > 0) {
    const formatted = searchResults.value.map(r => {
      let text = ''
      if (r.title) text += `### ${r.title}\n\n`
      if (r.source) text += `Source: ${r.source}\n\n`
      text += r.content
      return text
    }).join('\n\n--------------------------------\n\n')
    navigator.clipboard.writeText(formatted)
  } else {
    navigator.clipboard.writeText(searchResult.value)
  }
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
  if (days === 1) return '1 day ago'
  if (days < 7) return `${days} days ago`
  if (days < 30) return `${Math.floor(days / 7)} weeks ago`
  if (days < 365) return `${Math.floor(days / 30)} months ago`
  return `${Math.floor(days / 365)} years ago`
}

onMounted(() => {
  initUserState()
  fetchLibrary()
  fetchDocument()
})

// 监听路由参数变化
watch(() => route.params.title, () => {
  fetchDocument()
})
</script>
