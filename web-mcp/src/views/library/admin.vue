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
      <div class="mx-auto flex w-full max-w-[880px] flex-col items-center justify-center px-4 pt-10 lg:px-0">
        <div class="w-full space-y-6">
          <!-- 返回链接 -->
          <router-link 
            :to="`/libraries/${libraryId}`"
            class="inline-flex items-center gap-2 text-sm text-stone-600 hover:text-stone-900"
          >
            <svg xmlns="http://www.w3.org/2000/svg" width="24" height="24" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="1.5" stroke-linecap="round" stroke-linejoin="round" class="h-4 w-4">
              <path d="M5 12l14 0"></path>
              <path d="M5 12l6 6"></path>
              <path d="M5 12l6 -6"></path>
            </svg>
            {{ library.name }}
          </router-link>

          <!-- 标题行 -->
          <div class="flex flex-col gap-4 sm:flex-row sm:items-center sm:justify-between">
            <h1 class="text-2xl font-semibold tracking-tight text-stone-800">Admin Configuration</h1>
          </div>

          <!-- Tabs 区域 -->
          <div class="mt-10">
            <div class="flex flex-col-reverse gap-2 sm:flex-row sm:items-start sm:justify-between">
              <div class="overflow-x-auto overflow-y-hidden sm:overflow-visible">
                <div class="relative flex flex-nowrap items-end gap-1 border-b border-stone-300">
                  <button 
                    :class="[
                      '-mb-px flex flex-shrink-0 items-center gap-2 whitespace-nowrap rounded-t-lg px-4 py-2 text-base font-medium',
                      activeTab === 'configuration' 
                        ? 'relative z-10 border border-stone-300 border-b-stone-50 bg-stone-50 text-stone-800' 
                        : 'border border-stone-300 border-b-transparent text-stone-500 hover:border-stone-400 hover:text-stone-600'
                    ]"
                    @click="activeTab = 'configuration'"
                  >
                    <svg xmlns="http://www.w3.org/2000/svg" width="20" height="20" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="1.5" stroke-linecap="round" stroke-linejoin="round">
                      <path d="M7 8l-4 4l4 4"></path>
                      <path d="M17 8l4 4l-4 4"></path>
                      <path d="M14 4l-4 16"></path>
                    </svg>
                    Configuration
                  </button>
                  <button 
                    :class="[
                      '-mb-px flex flex-shrink-0 items-center gap-2 whitespace-nowrap rounded-t-lg px-4 py-2 text-base font-medium',
                      activeTab === 'versions' 
                        ? 'relative z-10 border border-stone-300 border-b-stone-50 bg-stone-50 text-stone-800' 
                        : 'border border-stone-300 border-b-transparent text-stone-500 hover:border-stone-400 hover:text-stone-600'
                    ]"
                    @click="activeTab = 'versions'"
                  >
                    <svg xmlns="http://www.w3.org/2000/svg" width="20" height="20" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="1.5" stroke-linecap="round" stroke-linejoin="round">
                      <path d="M7.5 7.5m-1 0a1 1 0 1 0 2 0a1 1 0 1 0 -2 0"></path>
                      <path d="M3 6v5.172a2 2 0 0 0 .586 1.414l7.71 7.71a2.41 2.41 0 0 0 3.408 0l5.592 -5.592a2.41 2.41 0 0 0 0 -3.408l-7.71 -7.71a2 2 0 0 0 -1.414 -.586h-5.172a3 3 0 0 0 -3 3z"></path>
                    </svg>
                    Versions
                  </button>
                </div>
              </div>
              <!-- 工具栏 -->
              <div class="flex flex-wrap gap-2.5 sm:gap-1.5">
                <button 
                  class="flex h-8 items-center justify-center gap-1.5 rounded-lg border border-stone-300 text-base text-stone-500 transition hover:border-stone-400 px-3 py-2"
                  @click="refreshData"
                >
                  <svg xmlns="http://www.w3.org/2000/svg" width="24" height="24" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="1.5" stroke-linecap="round" stroke-linejoin="round" class="h-5 w-5 text-stone-500">
                    <path d="M20 11a8.1 8.1 0 0 0 -15.5 -2m-.5 -4v4h4"></path>
                    <path d="M4 13a8.1 8.1 0 0 0 15.5 2m.5 4v-4h-4"></path>
                  </svg>
                  <span>Refresh</span>
                </button>
                <button 
                  class="flex h-8 items-center justify-center gap-1.5 rounded-lg border border-red-300 text-base text-red-600 transition hover:border-red-400 hover:bg-red-50 px-3 py-2"
                  @click="handleDeleteLibrary"
                >
                  <svg xmlns="http://www.w3.org/2000/svg" width="24" height="24" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="1.5" stroke-linecap="round" stroke-linejoin="round" class="h-5 w-5 text-red-600">
                    <path d="M4 7l16 0"></path>
                    <path d="M10 11l0 6"></path>
                    <path d="M14 11l0 6"></path>
                    <path d="M5 7l1 12a2 2 0 0 0 2 2h8a2 2 0 0 0 2 -2l1 -12"></path>
                    <path d="M9 7v-3a1 1 0 0 1 1 -1h4a1 1 0 0 1 1 1v3"></path>
                  </svg>
                  <span>Delete Library</span>
                </button>
              </div>
            </div>
          </div>

          <!-- Configuration Tab -->
          <div v-if="activeTab === 'configuration'" class="mt-8">
            <div class="rounded-3xl border border-stone-200 bg-white p-6 shadow-sm sm:p-8">
              <form class="space-y-6">
                <div class="flex items-center justify-between gap-4">
                  <h3 class="text-base font-semibold text-stone-800">Basic Information</h3>
                </div>
                
                <div class="-mt-2">
                  <div class="space-y-4">
                    <div>
                      <label class="block text-sm font-medium text-stone-700 mb-2">Name</label>
                      <input 
                        v-model="editForm.name"
                        type="text" 
                        class="w-full h-10 px-3 rounded-lg border border-stone-300 text-sm focus:outline-none focus:border-emerald-600 focus:ring-1 focus:ring-emerald-600 bg-white hover:border-emerald-600"
                        placeholder="Library name"
                      />
                    </div>
                    <div>
                      <label class="block text-sm font-medium text-stone-700 mb-2">Description</label>
                      <textarea 
                        v-model="editForm.description"
                        rows="3" 
                        class="w-full px-3 py-2 rounded-lg border border-stone-300 text-sm resize-none focus:outline-none focus:border-emerald-600 focus:ring-1 focus:ring-emerald-600 bg-white hover:border-emerald-600"
                        placeholder="Brief description of the library"
                      ></textarea>
                    </div>
                  </div>
                </div>

                <div class="border-t border-stone-200 pt-6">
                  <div class="flex items-center gap-4">
                    <button 
                      type="button"
                      class="inline-flex items-center gap-2 rounded-lg bg-emerald-600 px-4 py-2.5 text-sm font-medium text-white hover:bg-emerald-700 disabled:cursor-not-allowed disabled:opacity-50"
                      :disabled="saving"
                      @click="saveConfiguration"
                    >
                      {{ saving ? 'Saving...' : 'Save Configuration' }}
                    </button>
                  </div>
                </div>
              </form>
            </div>
          </div>

          <!-- Versions Tab -->
          <div v-if="activeTab === 'versions'" class="mt-8">
            <div class="rounded-3xl border border-stone-200 bg-white p-6 shadow-sm sm:p-8">
              <div class="space-y-6">
                <!-- 标题和添加版本按钮 -->
                <div class="flex items-center justify-between">
                  <div>
                    <h3 class="text-base font-semibold text-stone-800">Versions</h3>
                    <p class="mt-1 text-sm text-stone-500">Manage different versions and tags of this library</p>
                  </div>
                  <button 
                    class="inline-flex items-center gap-2 rounded-lg px-4 py-2 text-sm font-medium shadow-sm transition-all bg-emerald-600 text-white hover:bg-emerald-700"
                    title="Add a new version"
                    @click="openAddVersionModal"
                  >
                    <svg xmlns="http://www.w3.org/2000/svg" width="24" height="24" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="1.5" stroke-linecap="round" stroke-linejoin="round" class="h-4 w-4">
                      <path d="M12 5l0 14"></path>
                      <path d="M5 12l14 0"></path>
                    </svg>
                    Add Version
                  </button>
                </div>

                <!-- 版本列表表格 -->
                <div class="w-full overflow-x-auto md:overflow-x-visible">
                  <table class="w-full min-w-[600px] table-fixed border-b border-stone-200">
                    <thead class="border-b border-stone-200">
                      <tr>
                        <th class="w-[200px] px-2 py-3 text-left text-sm font-normal uppercase leading-none text-stone-400 sm:px-4">Version</th>
                        <th class="w-[120px] px-2 py-3 text-right text-sm font-normal uppercase leading-none text-stone-400 sm:px-4">Tokens</th>
                        <th class="w-[120px] px-2 py-3 text-right text-sm font-normal uppercase leading-none text-stone-400 sm:px-4">Snippets</th>
                        <th class="w-[160px] px-2 py-3 text-right text-sm font-normal uppercase leading-none text-stone-400 sm:px-4">Last Updated</th>
                        <th class="w-[100px] px-1 py-3 text-center text-sm font-normal uppercase leading-none text-stone-400">Actions</th>
                      </tr>
                    </thead>
                    <tbody class="divide-y divide-stone-200">
                      <!-- 空状态 -->
                      <tr v-if="versions.length === 0 && !loadingVersions">
                        <td colspan="5" class="py-12 text-center">
                          <div class="flex flex-col items-center gap-2">
                            <svg xmlns="http://www.w3.org/2000/svg" width="40" height="40" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="1" stroke-linecap="round" stroke-linejoin="round" class="text-stone-300">
                              <path d="M7.5 7.5m-1 0a1 1 0 1 0 2 0a1 1 0 1 0 -2 0"></path>
                              <path d="M3 6v5.172a2 2 0 0 0 .586 1.414l7.71 7.71a2.41 2.41 0 0 0 3.408 0l5.592 -5.592a2.41 2.41 0 0 0 0 -3.408l-7.71 -7.71a2 2 0 0 0 -1.414 -.586h-5.172a3 3 0 0 0 -3 3z"></path>
                            </svg>
                            <p class="text-sm font-medium text-stone-500">No versions yet</p>
                            <p class="text-sm text-stone-400">Create your first version to get started</p>
                          </div>
                        </td>
                      </tr>
                      <!-- 版本行 -->
                      <tr v-for="version in versions" :key="version.version" class="group transition-colors hover:bg-white">
                        <td class="h-11 px-2 align-middle sm:px-4">
                          <div class="flex items-center gap-2 text-base font-normal leading-tight text-stone-800">
                            <svg xmlns="http://www.w3.org/2000/svg" width="24" height="24" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="1.5" stroke-linecap="round" stroke-linejoin="round" class="h-4 w-4 flex-shrink-0">
                              <path d="M7.5 7.5m-1 0a1 1 0 1 0 2 0a1 1 0 1 0 -2 0"></path>
                              <path d="M3 6v5.172a2 2 0 0 0 .586 1.414l7.71 7.71a2.41 2.41 0 0 0 3.408 0l5.592 -5.592a2.41 2.41 0 0 0 0 -3.408l-7.71 -7.71a2 2 0 0 0 -1.414 -.586h-5.172a3 3 0 0 0 -3 3z"></path>
                            </svg>
                            <router-link 
                              :to="version.version === library.default_version 
                                ? `/libraries/${libraryId}` 
                                : `/libraries/${libraryId}/${version.version}`"
                              class="transition-colors hover:text-emerald-600 hover:underline"
                            >
                              {{ version.version }}
                            </router-link>
                            <span v-if="version.version === library.default_version" class="ml-1 rounded bg-emerald-600 px-1.5 py-0.5 text-xs font-semibold text-white">Default</span>
                          </div>
                        </td>
                        <td class="h-11 whitespace-nowrap px-2 text-right align-middle text-base font-normal slashed-zero tabular-nums leading-tight text-stone-800 sm:px-4">
                          {{ formatNumber(version.token_count || 0) }}
                        </td>
                        <td class="h-11 whitespace-nowrap px-2 text-right align-middle text-base font-normal slashed-zero tabular-nums leading-tight text-stone-800 sm:px-4">
                          {{ formatNumber(version.chunk_count || 0) }}
                        </td>
                        <td class="h-11 px-2 text-right align-middle text-base font-normal slashed-zero tabular-nums leading-tight text-stone-800 sm:px-4">
                          {{ formatDateShort(version.last_updated) }}
                        </td>
                        <td class="h-11 px-1 text-center align-middle">
                          <div class="flex items-center justify-center gap-2">
                            <!-- 刷新按钮 -->
                            <button 
                              class="flex items-center justify-center text-stone-500 transition-colors hover:text-emerald-600 disabled:opacity-50"
                              title="Refresh version"
                              :disabled="version.version === library.default_version"
                              @click="handleRefreshVersion(version.version)"
                            >
                              <svg xmlns="http://www.w3.org/2000/svg" width="20" height="20" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="1.5" stroke-linecap="round" stroke-linejoin="round">
                                <path d="M20 11a8.1 8.1 0 0 0 -15.5 -2m-.5 -4v4h4"></path>
                                <path d="M4 13a8.1 8.1 0 0 0 15.5 2m.5 4v-4h-4"></path>
                              </svg>
                            </button>
                            <!-- 删除按钮 -->
                            <button 
                              class="flex items-center justify-center text-stone-300 transition-colors hover:text-red-600 disabled:opacity-50 disabled:cursor-not-allowed"
                              title="Delete version"
                              :disabled="version.version === library.default_version"
                              @click="handleDeleteVersion(version.version)"
                            >
                              <svg xmlns="http://www.w3.org/2000/svg" width="20" height="20" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="1.5" stroke-linecap="round" stroke-linejoin="round">
                                <path d="M4 7l16 0"></path>
                                <path d="M10 11l0 6"></path>
                                <path d="M14 11l0 6"></path>
                                <path d="M5 7l1 12a2 2 0 0 0 2 2h8a2 2 0 0 0 2 -2l1 -12"></path>
                                <path d="M9 7v-3a1 1 0 0 1 1 -1h4a1 1 0 0 1 1 1v3"></path>
                              </svg>
                            </button>
                          </div>
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

    <!-- Add Version Modal -->
    <AddVersionModal
      v-model:visible="showAddVersionModal"
      :library-id="libraryId"
      :library-name="library.name"
      :source-type="library.source_type"
      :source-url="library.source_url"
      @success="handleVersionCreated"
    />

    <!-- Footer -->
    <AppFooter />
  </div>
</template>

<script setup lang="ts">
import { ref, reactive, computed, onMounted, watch } from 'vue'
import { useRoute } from 'vue-router'
import { ElMessage, ElMessageBox } from 'element-plus'
import AppHeader from '@/components/AppHeader.vue'
import AppFooter from '@/components/AppFooter.vue'
import AddVersionModal from '@/components/AddVersionModal.vue'
import { useUser } from '@/stores/user'
import { getLibrary, updateLibrary, deleteVersion, refreshVersion, deleteLibrary, getVersions } from '@/api/library'
import { useRouter } from 'vue-router'
import { getDocuments, deleteDocument, uploadDocumentWithSSE } from '@/api/document'
import type { Library } from '@/api/library'
import type { Document, ProcessStatus } from '@/api/document'

const route = useRoute()
const router = useRouter()
const { isLoggedIn, userEmail, userPlan, initUserState, redirectToSSO } = useUser()

const libraryId = computed(() => Number(route.params.id))
const library = ref<Library>({
  id: 0,
  name: '',
  default_version: '',
  versions: [],
  source_type: '',
  source_url: '',
  description: '',
  document_count: 0,
  chunk_count: 0,
  token_count: 0,
  status: '',
  created_at: '',
  updated_at: ''
})

const activeTab = ref('configuration')
const saving = ref(false)

// 上传状态
const uploading = ref(false)
const uploadProgress = ref(0)
const uploadMessage = ref('')

// 刷新版本状态
const refreshingVersion = ref('')

// Configuration form
const editForm = reactive({
  name: '',
  description: ''
})

// Versions
interface VersionInfo {
  version: string
  token_count: number
  chunk_count: number
  last_updated: string
}

const versions = ref<VersionInfo[]>([])
const loadingVersions = ref(false)
const showAddVersionModal = ref(false)
const selectedVersion = ref('')

// Documents
const documents = ref<Document[]>([])
const loadingDocs = ref(false)
const page = ref(1)
const pageSize = ref(10)
const totalDocs = ref(0)

const handleSignIn = () => {
  redirectToSSO()
}

const fetchLibrary = async () => {
  const res = await getLibrary(libraryId.value)
  library.value = res
  editForm.name = res.name
  editForm.description = res.description
}

const fetchVersions = async () => {
  loadingVersions.value = true
  try {
    const res = await getVersions(libraryId.value)
    versions.value = res || []
    // 自动选择第一个版本
    if (versions.value.length > 0 && !selectedVersion.value) {
      selectedVersion.value = versions.value[0].version
    }
  } finally {
    loadingVersions.value = false
  }
}

const fetchDocuments = async () => {
  loadingDocs.value = true
  try {
    const res = await getDocuments({
      library_id: libraryId.value,
      page: page.value,
      page_size: pageSize.value
    })
    documents.value = res.list || []
    totalDocs.value = res.total
  } finally {
    loadingDocs.value = false
  }
}

const resetForm = () => {
  editForm.name = library.value.name
  editForm.description = library.value.description
}

const saveConfiguration = async () => {
  saving.value = true
  try {
    const res = await updateLibrary(libraryId.value, editForm)
    library.value = res
    console.log('✓ Configuration saved')
  } finally {
    saving.value = false
  }
}

const handleFileUpload = async (event: Event) => {
  const input = event.target as HTMLInputElement
  const file = input.files?.[0]
  if (!file) return

  if (!selectedVersion.value) {
    ElMessage.warning('Please select a version first')
    return
  }

  const allowedTypes = ['.md', '.pdf', '.docx']
  const ext = file.name.substring(file.name.lastIndexOf('.')).toLowerCase()
  if (!allowedTypes.includes(ext)) {
    ElMessage.warning('Only .md, .pdf, .docx formats are supported')
    return
  }

  // 重置状态
  uploading.value = true
  uploadProgress.value = 0
  uploadMessage.value = 'Uploading...'

  try {
    // 创建 FormData 并添加版本信息
    const formData = new FormData()
    formData.append('library_id', libraryId.value.toString())
    formData.append('version', selectedVersion.value)
    formData.append('file', file)

    // 使用普通上传接口（后台异步处理，通过日志查看进度）
    const response = await fetch(`/api/documents/upload?library_id=${libraryId.value}&version=${selectedVersion.value}`, {
      method: 'POST',
      body: formData
    })
    
    if (response.ok) {
      ElMessage.success('上传已启动，跳转到控制台查看进度')
      uploading.value = false
      uploadProgress.value = 0
      uploadMessage.value = ''
      // 跳转到详情页的控制台 tab
      router.push({ 
        name: 'library-detail', 
        params: { id: libraryId.value }, 
        query: { tab: 'logs' } 
      })
    } else {
      throw new Error('Upload failed')
    }

    // ====== 以下是 SSE 版本的代码，保留备用 ======
    // const eventSource = new EventSource(`/api/documents/upload-sse?library_id=${libraryId.value}&version=${selectedVersion.value}`)
    // 
    // eventSource.addEventListener('parsing', (event) => {
    //   const data = JSON.parse(event.data)
    //   uploadProgress.value = 20
    //   uploadMessage.value = 'Parsing document...'
    // })
    //
    // eventSource.addEventListener('chunking', (event) => {
    //   const data = JSON.parse(event.data)
    //   uploadProgress.value = 50
    //   uploadMessage.value = 'Chunking document...'
    // })
    //
    // eventSource.addEventListener('embedding', (event) => {
    //   const data = JSON.parse(event.data)
    //   uploadProgress.value = 80
    //   uploadMessage.value = 'Generating embeddings...'
    // })
    //
    // eventSource.addEventListener('completed', (event) => {
    //   const data = JSON.parse(event.data)
    //   uploadProgress.value = 100
    //   uploadMessage.value = 'Upload successful!'
    //   console.log('✓ Upload successful')
    //   eventSource.close()
    //   
    //   setTimeout(() => {
    //     uploading.value = false
    //     uploadProgress.value = 0
    //     uploadMessage.value = ''
    //     fetchDocuments()
    //     fetchVersions()
    //   }, 500)
    // })
    //
    // eventSource.addEventListener('error', (event) => {
    //   console.error('Upload error:', event)
    //   ElMessage.error('Upload failed')
    //   eventSource.close()
    //   uploading.value = false
    //   uploadProgress.value = 0
    //   uploadMessage.value = ''
    // })
    // ====== SSE 版本代码结束 ======
  } catch (error) {
    ElMessage.error('Upload failed: ' + (error instanceof Error ? error.message : 'Unknown error'))
    uploading.value = false
  }
  
  input.value = ''
}

const handleDelete = async (id: number) => {
  if (!confirm('Are you sure you want to delete this document?')) return
  
  await deleteDocument(id)
  console.log('✓ Document deleted')
  fetchDocuments()
}

const refreshData = () => {
  fetchLibrary()
  if (activeTab.value === 'documents') {
    fetchDocuments()
  }
}

const handleDeleteLibrary = async () => {
  if (!confirm(`Are you sure you want to delete the library "${library.value.name}"? This action cannot be undone.`)) return
  
  try {
    await deleteLibrary(libraryId.value)
    console.log('✓ Library deleted')
    ElMessage.success('Library deleted successfully')
    router.push('/')
  } catch (error) {
    console.error('Failed to delete library:', error)
  }
}

const formatSize = (bytes: number) => {
  if (bytes < 1024) return bytes + ' B'
  if (bytes < 1024 * 1024) return (bytes / 1024).toFixed(1) + ' KB'
  return (bytes / 1024 / 1024).toFixed(1) + ' MB'
}

const getStatusText = (status: string) => {
  const map: Record<string, string> = {
    active: 'Completed',
    processing: 'Processing',
    failed: 'Failed'
  }
  return map[status] || status
}

const formatDateShort = (dateStr: string) => {
  if (!dateStr) return '-'
  const date = new Date(dateStr)
  return date.toLocaleDateString('en-US', { month: 'short', day: 'numeric', year: 'numeric' })
}

const formatNumber = (num: number) => {
  if (num >= 1000000) return (num / 1000000).toFixed(1) + 'M'
  if (num >= 1000) return (num / 1000).toFixed(1) + 'K'
  return num.toLocaleString()
}

const handleReprocess = async (id: number) => {
  // TODO: 实现重新处理文档的 API
  console.log('Reprocess document:', id)
  ElMessage.info('Reprocess feature coming soon')
}

const handleRefreshVersion = async (version: string) => {
  try {
    await ElMessageBox.confirm(
      '这将重新处理该版本下的所有文档',
      `刷新版本 "${version}"？`,
      {
        confirmButtonText: '确定',
        cancelButtonText: '取消',
        type: 'warning'
      }
    )
  } catch {
    return // 用户取消
  }
  
  refreshingVersion.value = version
  
  try {
    await refreshVersion(libraryId.value, version)
    ElMessage.success('版本刷新已启动，跳转到控制台查看进度')
    // 跳转到详情页的控制台 tab（包含版本信息）
    router.push({ 
      name: 'library-version', 
      params: { id: libraryId.value, version: version }, 
      query: { tab: 'logs' } 
    })
  } catch (error: any) {
    console.error('Failed to refresh version:', error)
    ElMessage.error('刷新失败: ' + (error.message || '未知错误'))
  } finally {
    refreshingVersion.value = ''
  }
}

const handleDeleteVersion = async (version: string) => {
  if (!confirm(`Are you sure you want to delete version "${version}"? This will delete all documents in this version.`)) return
  
  try {
    await deleteVersion(libraryId.value, version)
    console.log('✓ Version deleted')
    ElMessage.success('Version deleted successfully')
    await fetchVersions()
  } catch (error) {
    console.error('Failed to delete version:', error)
  }
}

// 打开添加版本弹窗
const openAddVersionModal = () => {
  showAddVersionModal.value = true
}

// 处理版本创建成功
const handleVersionCreated = async (version: string) => {
  console.log('✓ Version created:', version)
  // 跳转到详情页的控制台 tab 查看进度（包含版本信息）
  router.push({ 
    name: 'library-version', 
    params: { id: libraryId.value, version: version }, 
    query: { tab: 'logs' } 
  })
}

// 切换到 documents tab 时加载文档和版本
watch(activeTab, (newTab) => {
  if (newTab === 'documents') {
    if (documents.value.length === 0) {
      fetchDocuments()
    }
    if (versions.value.length === 0) {
      fetchVersions()
    }
  }
  if (newTab === 'versions') {
    if (versions.value.length === 0) {
      fetchVersions()
    }
  }
})

onMounted(() => {
  initUserState()
  fetchLibrary()
  
  // 检查 URL 参数
  if (route.query.tab === 'documents') {
    activeTab.value = 'documents'
    fetchDocuments()
    fetchVersions()
  }
})
</script>
