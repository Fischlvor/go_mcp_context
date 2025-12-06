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
                <div class="relative flex flex-nowrap items-end gap-1">
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
                      activeTab === 'documents' 
                        ? 'relative z-10 border border-stone-300 border-b-stone-50 bg-stone-50 text-stone-800' 
                        : 'border border-stone-300 border-b-transparent text-stone-500 hover:border-stone-400 hover:text-stone-600'
                    ]"
                    @click="activeTab = 'documents'"
                  >
                    <svg xmlns="http://www.w3.org/2000/svg" width="20" height="20" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="1.5" stroke-linecap="round" stroke-linejoin="round">
                      <path d="M14 3v4a1 1 0 0 0 1 1h4"></path>
                      <path d="M17 21h-10a2 2 0 0 1 -2 -2v-14a2 2 0 0 1 2 -2h7l5 5v11a2 2 0 0 1 -2 2z"></path>
                      <path d="M9 17h6"></path>
                      <path d="M9 13h6"></path>
                    </svg>
                    Documents
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
              </div>
            </div>
            <div class="border-t border-stone-300"></div>
          </div>

          <!-- Configuration Tab -->
          <div v-if="activeTab === 'configuration'" class="mt-8">
            <div class="rounded-3xl border border-stone-200 bg-white p-6 shadow-sm sm:p-8">
              <div class="space-y-6">
                <h3 class="text-lg font-semibold text-stone-800">Library Information</h3>
                
                <div class="space-y-4">
                  <div>
                    <label class="block text-sm font-medium text-stone-700 mb-2">Name</label>
                    <input 
                      v-model="editForm.name"
                      type="text" 
                      class="w-full h-10 px-3 rounded-lg border border-stone-300 text-sm focus:outline-none focus:border-emerald-600 focus:ring-1 focus:ring-emerald-600"
                      placeholder="Library name"
                    />
                  </div>
                  <div>
                    <label class="block text-sm font-medium text-stone-700 mb-2">Version</label>
                    <input 
                      v-model="editForm.version"
                      type="text" 
                      class="w-full h-10 px-3 rounded-lg border border-stone-300 text-sm focus:outline-none focus:border-emerald-600 focus:ring-1 focus:ring-emerald-600"
                      placeholder="e.g. 3.4.0"
                    />
                  </div>
                  <div>
                    <label class="block text-sm font-medium text-stone-700 mb-2">Description</label>
                    <textarea 
                      v-model="editForm.description"
                      rows="3" 
                      class="w-full px-3 py-2 rounded-lg border border-stone-300 text-sm resize-none focus:outline-none focus:border-emerald-600 focus:ring-1 focus:ring-emerald-600"
                      placeholder="Brief description of the library"
                    ></textarea>
                  </div>
                </div>

                <div class="flex justify-end gap-3 pt-4 border-t border-stone-200">
                  <button 
                    class="h-10 px-4 rounded-lg border border-stone-300 text-sm font-medium text-stone-700 hover:bg-stone-50"
                    @click="resetForm"
                  >
                    Reset
                  </button>
                  <button 
                    class="h-10 px-4 rounded-lg bg-emerald-600 text-sm font-medium text-white hover:bg-emerald-700 disabled:opacity-50"
                    :disabled="saving"
                    @click="saveConfiguration"
                  >
                    {{ saving ? 'Saving...' : 'Save Changes' }}
                  </button>
                </div>
              </div>
            </div>
          </div>

          <!-- Documents Tab -->
          <div v-if="activeTab === 'documents'" class="mt-8">
            <div class="rounded-3xl border border-stone-200 bg-white p-6 shadow-sm sm:p-8">
              <div class="space-y-6">
                <!-- 标题和上传按钮 -->
                <div class="flex items-center justify-between">
                  <div>
                    <h3 class="text-base font-semibold text-stone-800">Documents</h3>
                    <p class="mt-1 text-sm text-stone-500">Manage documents in this library</p>
                  </div>
                  <label 
                    :class="[
                      'flex h-10 items-center justify-center gap-2 rounded-lg px-4 text-sm font-medium text-white transition-colors',
                      uploading ? 'bg-stone-400 cursor-not-allowed' : 'bg-emerald-600 hover:bg-emerald-700 cursor-pointer'
                    ]"
                  >
                    <svg v-if="!uploading" xmlns="http://www.w3.org/2000/svg" width="20" height="20" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round">
                      <path d="M4 17v2a2 2 0 0 0 2 2h12a2 2 0 0 0 2 -2v-2"></path>
                      <path d="M7 9l5 -5l5 5"></path>
                      <path d="M12 4l0 12"></path>
                    </svg>
                    <svg v-else class="h-5 w-5 animate-spin" xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 24 24">
                      <circle class="opacity-25" cx="12" cy="12" r="10" stroke="currentColor" stroke-width="4"></circle>
                      <path class="opacity-75" fill="currentColor" d="M4 12a8 8 0 018-8V0C5.373 0 0 5.373 0 12h4zm2 5.291A7.962 7.962 0 014 12H0c0 3.042 1.135 5.824 3 7.938l3-2.647z"></path>
                    </svg>
                    {{ uploading ? 'Processing...' : 'Upload' }}
                    <input 
                      type="file" 
                      class="hidden" 
                      accept=".md,.pdf,.docx"
                      :disabled="uploading"
                      @change="handleFileUpload"
                    />
                  </label>
                </div>

                <!-- 上传进度条 -->
                <div v-if="uploading" class="rounded-lg border border-emerald-200 bg-emerald-50 p-4">
                  <div class="flex items-center justify-between mb-2">
                    <span class="text-sm font-medium text-emerald-800">{{ uploadMessage }}</span>
                    <span class="text-sm font-medium text-emerald-800">{{ uploadProgress }}%</span>
                  </div>
                  <div class="h-2 w-full overflow-hidden rounded-full bg-emerald-200">
                    <div 
                      class="h-full rounded-full bg-emerald-600 transition-all duration-300"
                      :style="{ width: uploadProgress + '%' }"
                    ></div>
                  </div>
                </div>

                <!-- 文档列表表格 -->
                <div class="w-full overflow-x-auto md:overflow-x-visible">
                  <table class="w-full min-w-[600px] table-fixed border-b border-stone-200">
                    <thead class="border-b border-stone-200">
                      <tr>
                        <th class="w-[240px] px-2 py-3 text-left text-sm font-normal uppercase leading-none text-stone-400 sm:px-4">Title</th>
                        <th class="w-[120px] px-2 py-3 text-right text-sm font-normal uppercase leading-none text-stone-400 sm:px-4">Tokens</th>
                        <th class="w-[120px] px-2 py-3 text-right text-sm font-normal uppercase leading-none text-stone-400 sm:px-4">Snippets</th>
                        <th class="w-[160px] px-2 py-3 text-right text-sm font-normal uppercase leading-none text-stone-400 sm:px-4">Last Updated</th>
                        <th class="w-[100px] px-1 py-3 text-center text-sm font-normal uppercase leading-none text-stone-400">Actions</th>
                      </tr>
                    </thead>
                    <tbody class="divide-y divide-stone-200">
                      <!-- 空状态 -->
                      <tr v-if="documents.length === 0 && !loadingDocs">
                        <td colspan="5" class="py-12 text-center">
                          <div class="flex flex-col items-center gap-2">
                            <svg xmlns="http://www.w3.org/2000/svg" width="40" height="40" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="1" stroke-linecap="round" stroke-linejoin="round" class="text-stone-300">
                              <path d="M14 3v4a1 1 0 0 0 1 1h4"></path>
                              <path d="M17 21h-10a2 2 0 0 1 -2 -2v-14a2 2 0 0 1 2 -2h7l5 5v11a2 2 0 0 1 -2 2z"></path>
                            </svg>
                            <p class="text-sm font-medium text-stone-500">No documents yet</p>
                            <p class="text-sm text-stone-400">Upload your first document to get started</p>
                          </div>
                        </td>
                      </tr>
                      <!-- 文档行 -->
                      <tr v-for="doc in documents" :key="doc.id" class="group transition-colors hover:bg-white">
                        <td class="h-11 px-2 align-middle sm:px-4">
                          <div class="flex items-center gap-2 text-base font-normal leading-tight text-stone-800">
                            <svg xmlns="http://www.w3.org/2000/svg" width="24" height="24" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="1.5" stroke-linecap="round" stroke-linejoin="round" class="h-4 w-4 flex-shrink-0">
                              <path d="M14 3v4a1 1 0 0 0 1 1h4"></path>
                              <path d="M17 21h-10a2 2 0 0 1 -2 -2v-14a2 2 0 0 1 2 -2h7l5 5v11a2 2 0 0 1 -2 2z"></path>
                            </svg>
                            <router-link 
                              :to="`/libraries/${libraryId}/${encodeURIComponent(doc.title)}`"
                              class="transition-colors hover:text-emerald-600 hover:underline truncate"
                            >
                              {{ doc.title }}
                            </router-link>
                          </div>
                        </td>
                        <td class="h-11 whitespace-nowrap px-2 text-right align-middle text-base font-normal slashed-zero tabular-nums leading-tight text-stone-800 sm:px-4">
                          {{ formatNumber(doc.token_count || 0) }}
                        </td>
                        <td class="h-11 whitespace-nowrap px-2 text-right align-middle text-base font-normal slashed-zero tabular-nums leading-tight text-stone-800 sm:px-4">
                          {{ formatNumber(doc.chunk_count || 0) }}
                        </td>
                        <td class="h-11 px-2 text-right align-middle text-base font-normal slashed-zero tabular-nums leading-tight text-stone-800 sm:px-4">
                          {{ formatDateShort(doc.updated_at) }}
                        </td>
                        <td class="h-11 px-1 text-center align-middle">
                          <div class="flex items-center justify-center gap-2">
                            <!-- 刷新按钮 -->
                            <button 
                              class="flex items-center justify-center text-stone-500 transition-colors hover:text-emerald-600"
                              title="Reprocess"
                              @click="handleReprocess(doc.id)"
                            >
                              <svg xmlns="http://www.w3.org/2000/svg" width="20" height="20" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="1.5" stroke-linecap="round" stroke-linejoin="round">
                                <path d="M20 11a8.1 8.1 0 0 0 -15.5 -2m-.5 -4v4h4"></path>
                                <path d="M4 13a8.1 8.1 0 0 0 15.5 2m.5 4v-4h-4"></path>
                              </svg>
                            </button>
                            <!-- 删除按钮 -->
                            <button 
                              class="flex items-center justify-center text-stone-500 transition-colors hover:text-red-600"
                              title="Delete"
                              @click="handleDelete(doc.id)"
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

                <!-- 分页 -->
                <div v-if="totalDocs > pageSize" class="flex items-center justify-between border-t border-stone-200 pt-4">
                  <span class="text-sm text-stone-500">{{ totalDocs }} documents</span>
                  <div class="flex gap-2">
                    <button 
                      class="h-8 px-3 rounded-lg border border-stone-300 text-sm text-stone-600 hover:bg-stone-50 disabled:opacity-50"
                      :disabled="page === 1"
                      @click="page--; fetchDocuments()"
                    >
                      Previous
                    </button>
                    <button 
                      class="h-8 px-3 rounded-lg border border-stone-300 text-sm text-stone-600 hover:bg-stone-50 disabled:opacity-50"
                      :disabled="page * pageSize >= totalDocs"
                      @click="page++; fetchDocuments()"
                    >
                      Next
                    </button>
                  </div>
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
import { ref, reactive, computed, onMounted, watch } from 'vue'
import { useRoute } from 'vue-router'
import AppHeader from '@/components/AppHeader.vue'
import AppFooter from '@/components/AppFooter.vue'
import { useUser } from '@/stores/user'
import { getLibrary, updateLibrary } from '@/api/library'
import { getDocuments, deleteDocument, uploadDocumentWithSSE } from '@/api/document'
import type { Library } from '@/api/library'
import type { Document, ProcessStatus } from '@/api/document'

const route = useRoute()
const { isLoggedIn, userEmail, userPlan, initUserState, redirectToSSO } = useUser()

const libraryId = computed(() => Number(route.params.id))
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

const activeTab = ref('configuration')
const saving = ref(false)

// 上传状态
const uploading = ref(false)
const uploadProgress = ref(0)
const uploadMessage = ref('')

// Configuration form
const editForm = reactive({
  name: '',
  version: '',
  description: ''
})

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
  if (res.code === 0) {
    library.value = res.data
    editForm.name = res.data.name
    editForm.version = res.data.version
    editForm.description = res.data.description
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
    if (res.code === 0) {
      documents.value = res.data.list || []
      totalDocs.value = res.data.total
    }
  } finally {
    loadingDocs.value = false
  }
}

const resetForm = () => {
  editForm.name = library.value.name
  editForm.version = library.value.version
  editForm.description = library.value.description
}

const saveConfiguration = async () => {
  saving.value = true
  try {
    const res = await updateLibrary(libraryId.value, editForm)
    if (res.code === 0) {
      library.value = res.data
      console.log('✓ Configuration saved')
    }
  } finally {
    saving.value = false
  }
}

const handleFileUpload = async (event: Event) => {
  const input = event.target as HTMLInputElement
  const file = input.files?.[0]
  if (!file) return

  const allowedTypes = ['.md', '.pdf', '.docx']
  const ext = file.name.substring(file.name.lastIndexOf('.')).toLowerCase()
  if (!allowedTypes.includes(ext)) {
    alert('Only .md, .pdf, .docx formats are supported')
    return
  }

  // 重置状态
  uploading.value = true
  uploadProgress.value = 0
  uploadMessage.value = 'Uploading...'

  try {
    await uploadDocumentWithSSE(libraryId.value, file, {
      onProgress: (status: ProcessStatus) => {
        uploadProgress.value = status.progress
        uploadMessage.value = status.message
      },
      onComplete: () => {
        console.log('✓ Upload successful')
        uploading.value = false
        uploadProgress.value = 0
        uploadMessage.value = ''
        fetchDocuments()
      },
      onError: (error: Error) => {
        alert('Upload failed: ' + error.message)
        uploading.value = false
        uploadProgress.value = 0
        uploadMessage.value = ''
      }
    })
  } catch (error) {
    alert('Upload failed')
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
  alert('Reprocess feature coming soon')
}

// 切换到 documents tab 时加载文档
watch(activeTab, (newTab) => {
  if (newTab === 'documents' && documents.value.length === 0) {
    fetchDocuments()
  }
})

onMounted(() => {
  initUserState()
  fetchLibrary()
  
  // 检查 URL 参数
  if (route.query.tab === 'documents') {
    activeTab.value = 'documents'
  }
})
</script>
