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
                    @click="showAddVersionModal = true"
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
    <div v-if="showAddVersionModal" class="fixed inset-0 z-50 flex items-center justify-center bg-black/50">
      <div class="w-full max-w-md rounded-lg bg-white p-6 shadow-lg">
        <div class="mb-4">
          <h2 class="text-lg font-semibold text-stone-800">Add New Version</h2>
          <p class="mt-1 text-sm text-stone-500">Create a new version for this library</p>
        </div>

        <div class="space-y-4">
          <div>
            <label class="block text-sm font-medium text-stone-700 mb-2">Version Name</label>
            <input 
              v-model="newVersionName"
              type="text" 
              placeholder="e.g., 1.0.0, 1.2.3-beta, 2.0.0-rc.1"
              class="w-full h-10 px-3 rounded-lg border border-stone-300 text-sm focus:outline-none focus:border-emerald-600 focus:ring-1 focus:ring-emerald-600"
            />
            <p class="mt-1 text-xs text-stone-500">Semantic Versioning (e.g., 1.0.0, 1.2.3-beta). The 'v' prefix will be added automatically.</p>
          </div>
        </div>

        <div class="mt-6 flex gap-3 justify-end">
          <button 
            class="px-4 py-2 rounded-lg border border-stone-300 text-sm font-medium text-stone-700 hover:bg-stone-50"
            @click="showAddVersionModal = false"
          >
            Cancel
          </button>
          <button 
            class="px-4 py-2 rounded-lg bg-emerald-600 text-sm font-medium text-white hover:bg-emerald-700 disabled:opacity-50"
            :disabled="!newVersionName.trim()"
            @click="handleAddVersion"
          >
            Create Version
          </button>
        </div>
      </div>
    </div>

    <!-- Refresh Progress Modal -->
    <div v-if="refreshing" class="fixed inset-0 z-50 flex items-center justify-center bg-black/50">
      <div class="w-full max-w-lg rounded-lg bg-white p-6 shadow-lg">
        <div class="mb-4">
          <h2 class="text-lg font-semibold text-stone-800">Refreshing Version: {{ refreshingVersion }}</h2>
          <p class="mt-1 text-sm text-stone-500">Reprocessing all documents in this version...</p>
        </div>

        <!-- Progress Bar -->
        <div v-if="refreshProgress" class="mb-4">
          <div class="flex justify-between text-sm text-stone-600 mb-1">
            <span>{{ refreshProgress.message }}</span>
            <span>{{ refreshProgress.current }} / {{ refreshProgress.total }}</span>
          </div>
          <div class="w-full h-2 bg-stone-200 rounded-full overflow-hidden">
            <div 
              class="h-full bg-emerald-500 transition-all duration-300"
              :style="{ width: `${refreshProgress.total > 0 ? (refreshProgress.current / refreshProgress.total) * 100 : 0}%` }"
            ></div>
          </div>
        </div>

        <!-- Document Status List -->
        <div v-if="refreshDocStatuses.length > 0" class="max-h-48 overflow-y-auto border border-stone-200 rounded-lg">
          <div 
            v-for="doc in refreshDocStatuses" 
            :key="doc.id"
            class="flex items-center gap-2 px-3 py-2 border-b border-stone-100 last:border-b-0"
          >
            <!-- Status Icon -->
            <div class="flex-shrink-0">
              <svg v-if="doc.status === 'processing'" class="h-4 w-4 text-blue-500 animate-spin" xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 24 24">
                <circle class="opacity-25" cx="12" cy="12" r="10" stroke="currentColor" stroke-width="4"></circle>
                <path class="opacity-75" fill="currentColor" d="M4 12a8 8 0 018-8V0C5.373 0 0 5.373 0 12h4zm2 5.291A7.962 7.962 0 014 12H0c0 3.042 1.135 5.824 3 7.938l3-2.647z"></path>
              </svg>
              <svg v-else-if="doc.status === 'completed'" class="h-4 w-4 text-emerald-500" xmlns="http://www.w3.org/2000/svg" viewBox="0 0 24 24" fill="currentColor">
                <path fill-rule="evenodd" d="M2.25 12c0-5.385 4.365-9.75 9.75-9.75s9.75 4.365 9.75 9.75-4.365 9.75-9.75 9.75S2.25 17.385 2.25 12zm13.36-1.814a.75.75 0 10-1.22-.872l-3.236 4.53L9.53 12.22a.75.75 0 00-1.06 1.06l2.25 2.25a.75.75 0 001.14-.094l3.75-5.25z" clip-rule="evenodd" />
              </svg>
              <svg v-else class="h-4 w-4 text-red-500" xmlns="http://www.w3.org/2000/svg" viewBox="0 0 24 24" fill="currentColor">
                <path fill-rule="evenodd" d="M12 2.25c-5.385 0-9.75 4.365-9.75 9.75s4.365 9.75 9.75 9.75 9.75-4.365 9.75-9.75S17.385 2.25 12 2.25zm-1.72 6.97a.75.75 0 10-1.06 1.06L10.94 12l-1.72 1.72a.75.75 0 101.06 1.06L12 13.06l1.72 1.72a.75.75 0 101.06-1.06L13.06 12l1.72-1.72a.75.75 0 10-1.06-1.06L12 10.94l-1.72-1.72z" clip-rule="evenodd" />
              </svg>
            </div>
            <!-- Document Title -->
            <span class="text-sm text-stone-700 truncate">{{ doc.title }}</span>
          </div>
        </div>

        <!-- Loading State (before documents start) -->
        <div v-else class="flex items-center justify-center py-8">
          <svg class="h-8 w-8 text-emerald-500 animate-spin" xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 24 24">
            <circle class="opacity-25" cx="12" cy="12" r="10" stroke="currentColor" stroke-width="4"></circle>
            <path class="opacity-75" fill="currentColor" d="M4 12a8 8 0 018-8V0C5.373 0 0 5.373 0 12h4zm2 5.291A7.962 7.962 0 014 12H0c0 3.042 1.135 5.824 3 7.938l3-2.647z"></path>
          </svg>
        </div>
      </div>
    </div>

    <!-- Footer -->
    <AppFooter />
  </div>
</template>

<script setup lang="ts">
import { ref, reactive, computed, onMounted, watch } from 'vue'
import { useRoute } from 'vue-router'
import { ElMessage } from 'element-plus'
import AppHeader from '@/components/AppHeader.vue'
import AppFooter from '@/components/AppFooter.vue'
import { useUser } from '@/stores/user'
import { getLibrary, updateLibrary, createVersion, deleteVersion, refreshVersionWithSSE, deleteLibrary, getVersions } from '@/api/library'
import type { RefreshStatus } from '@/api/library'
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
const refreshing = ref(false)
const refreshingVersion = ref('')
const refreshProgress = ref<RefreshStatus | null>(null)
const refreshDocStatuses = ref<Array<{ id: number; title: string; status: 'processing' | 'completed' | 'failed' }>>([])

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
const newVersionName = ref('')
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

    // 使用 SSE 上传
    const eventSource = new EventSource(`/api/documents/upload-sse?library_id=${libraryId.value}&version=${selectedVersion.value}`)
    
    eventSource.addEventListener('parsing', (event) => {
      const data = JSON.parse(event.data)
      uploadProgress.value = 20
      uploadMessage.value = 'Parsing document...'
    })

    eventSource.addEventListener('chunking', (event) => {
      const data = JSON.parse(event.data)
      uploadProgress.value = 50
      uploadMessage.value = 'Chunking document...'
    })

    eventSource.addEventListener('embedding', (event) => {
      const data = JSON.parse(event.data)
      uploadProgress.value = 80
      uploadMessage.value = 'Generating embeddings...'
    })

    eventSource.addEventListener('completed', (event) => {
      const data = JSON.parse(event.data)
      uploadProgress.value = 100
      uploadMessage.value = 'Upload successful!'
      console.log('✓ Upload successful')
      eventSource.close()
      
      setTimeout(() => {
        uploading.value = false
        uploadProgress.value = 0
        uploadMessage.value = ''
        fetchDocuments()
        fetchVersions()
      }, 500)
    })

    eventSource.addEventListener('error', (event) => {
      console.error('Upload error:', event)
      ElMessage.error('Upload failed')
      eventSource.close()
      uploading.value = false
      uploadProgress.value = 0
      uploadMessage.value = ''
    })

    // 发送文件
    const response = await fetch(`/api/documents/upload?library_id=${libraryId.value}&version=${selectedVersion.value}`, {
      method: 'POST',
      body: formData
    })
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
  if (!confirm(`Refresh version "${version}"? This will reprocess all documents in this version.`)) return
  
  // 初始化刷新状态
  refreshing.value = true
  refreshingVersion.value = version
  refreshProgress.value = null
  refreshDocStatuses.value = []
  
  try {
    await refreshVersionWithSSE(libraryId.value, version, {
      onProgress: (status) => {
        refreshProgress.value = status
        
        // 更新文档状态列表
        if (status.doc_id && status.doc_title) {
          const existing = refreshDocStatuses.value.find(d => d.id === status.doc_id)
          if (existing) {
            existing.status = status.stage === 'doc_completed' ? 'completed' 
                            : status.stage === 'doc_failed' ? 'failed' 
                            : 'processing'
          } else {
            refreshDocStatuses.value.push({
              id: status.doc_id,
              title: status.doc_title,
              status: 'processing'
            })
          }
        }
      },
      onComplete: (status) => {
        refreshProgress.value = status
        ElMessage.success(`Version refresh completed. ${status.total} documents processed.`)
        refreshing.value = false
        fetchVersions()
      },
      onError: (error) => {
        console.error('Failed to refresh version:', error)
        ElMessage.error('Refresh failed: ' + error.message)
        refreshing.value = false
      }
    })
  } catch (error) {
    console.error('Failed to refresh version:', error)
    refreshing.value = false
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

const handleAddVersion = async () => {
  if (!newVersionName.value.trim()) {
    ElMessage.warning('Please enter a version name')
    return
  }

  try {
    // 自动添加 v 前缀
    const versionWithPrefix = newVersionName.value.startsWith('v') 
      ? newVersionName.value 
      : `v${newVersionName.value}`
    
    const res = await createVersion(libraryId.value, versionWithPrefix)
    // 拦截器已处理错误，这里只处理成功情况
    console.log('✓ Version created')
    showAddVersionModal.value = false
    newVersionName.value = ''
    await fetchVersions()
  } catch (error) {
    // 错误已由拦截器显示，这里只记录日志
    console.error('Failed to add version:', error)
  }
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
