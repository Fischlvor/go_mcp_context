<template>
  <Teleport to="body">
    <div v-if="visible" class="fixed inset-0 z-50 flex items-center justify-center bg-black/50" @click.self="handleClose">
      <div class="w-full max-w-md rounded-xl bg-white p-6 shadow-lg md:p-10" @click="handleClickOutside">
        <!-- Header -->
        <div class="mb-6 flex items-center justify-between">
          <h2 class="text-xl font-semibold text-stone-900">Add a Specific Version</h2>
          <button 
            type="button" 
            class="flex-shrink-0 text-stone-400 transition-colors hover:text-stone-600"
            :disabled="importing"
            @click="handleClose"
          >
            <svg xmlns="http://www.w3.org/2000/svg" width="24" height="24" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="1.5" stroke-linecap="round" stroke-linejoin="round">
              <path d="M18 6l-12 12"></path>
              <path d="M6 6l12 12"></path>
            </svg>
          </button>
        </div>

        <div class="space-y-6">
          <!-- GitHub 模式 -->
          <template v-if="isGitHubSource">
            <!-- 仓库信息卡片 -->
            <div class="rounded-xl bg-stone-100 p-6">
              <div class="flex flex-col gap-4">
                <div class="flex flex-col">
                  <p class="text-base font-semibold text-stone-800">{{ libraryName }}</p>
                  <a 
                    :href="`https://github.com/${sourceUrl}`" 
                    target="_blank" 
                    rel="noopener noreferrer" 
                    class="w-fit text-base text-emerald-600 underline hover:text-emerald-700"
                  >
                    https://github.com/{{ sourceUrl }}
                  </a>
                </div>
                
                <!-- Tag/Branch 选择器 -->
                <div class="relative">
                  <button 
                    type="button"
                    class="inline-flex h-[40px] w-full items-center justify-between gap-2 rounded-[6px] border border-stone-400 bg-white px-4 text-sm font-medium text-stone-800 shadow-md transition-colors hover:border-emerald-500 focus-visible:outline-none focus-visible:ring-1 focus-visible:ring-emerald-600"
                    :disabled="importing"
                    @click="showRefDropdown = !showRefDropdown"
                  >
                    <span v-if="selectedGitRef" class="truncate text-base font-normal text-stone-800">{{ selectedGitRef }}</span>
                    <span v-else class="truncate text-base font-normal text-stone-400">Search for a tag or branch</span>
                    <svg xmlns="http://www.w3.org/2000/svg" width="24" height="24" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round" class="h-4 w-4 flex-shrink-0 text-stone-500">
                      <path d="M8 9l4 -4l4 4"></path>
                      <path d="M16 15l-4 4l-4 -4"></path>
                    </svg>
                  </button>
                  
                  <!-- 下拉列表 -->
                  <div 
                    v-if="showRefDropdown" 
                    class="absolute z-10 mt-1 w-full rounded-lg border border-stone-200 bg-white shadow-lg"
                  >
                    <!-- 搜索框 -->
                    <div class="p-2 border-b border-stone-100">
                      <input
                        v-model="refSearchQuery"
                        type="text"
                        placeholder="Search tags..."
                        class="w-full h-9 px-3 rounded-md border border-stone-300 text-sm focus:outline-none focus:border-emerald-500"
                        @click.stop
                      />
                    </div>
                    
                    <!-- 加载状态 -->
                    <div v-if="loadingReleases" class="flex items-center justify-center py-6">
                      <svg class="h-5 w-5 animate-spin text-emerald-500" xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 24 24">
                        <circle class="opacity-25" cx="12" cy="12" r="10" stroke="currentColor" stroke-width="4"></circle>
                        <path class="opacity-75" fill="currentColor" d="M4 12a8 8 0 018-8V0C5.373 0 0 5.373 0 12h4zm2 5.291A7.962 7.962 0 014 12H0c0 3.042 1.135 5.824 3 7.938l3-2.647z"></path>
                      </svg>
                      <span class="ml-2 text-sm text-stone-500">Loading versions...</span>
                    </div>
                    
                    <!-- 版本列表 -->
                    <ul v-else class="max-h-48 overflow-y-auto py-1">
                      <li v-if="filteredReleases.length === 0" class="px-4 py-3 text-sm text-stone-500 text-center">
                        No matching versions found
                      </li>
                      <li 
                        v-for="ref in filteredReleases" 
                        :key="ref"
                        class="flex items-center gap-2 px-4 py-2 text-sm text-stone-700 cursor-pointer hover:bg-stone-50"
                        @click="selectGitRef(ref)"
                      >
                        <svg xmlns="http://www.w3.org/2000/svg" width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="1.5" stroke-linecap="round" stroke-linejoin="round" class="text-stone-400">
                          <path d="M7.5 7.5m-1 0a1 1 0 1 0 2 0a1 1 0 1 0 -2 0"></path>
                          <path d="M3 6v5.172a2 2 0 0 0 .586 1.414l7.71 7.71a2.41 2.41 0 0 0 3.408 0l5.592 -5.592a2.41 2.41 0 0 0 0 -3.408l-7.71 -7.71a2 2 0 0 0 -1.414 -.586h-5.172a3 3 0 0 0 -3 3z"></path>
                        </svg>
                        {{ ref }}
                      </li>
                    </ul>
                  </div>
                </div>
              </div>
            </div>

            <!-- 导入中提示 -->
            <div v-if="importing" class="rounded-xl border border-stone-200 p-4 text-center">
              <div class="flex items-center justify-center gap-2 text-stone-600">
                <svg class="h-5 w-5 animate-spin" xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 24 24">
                  <circle class="opacity-25" cx="12" cy="12" r="10" stroke="currentColor" stroke-width="4"></circle>
                  <path class="opacity-75" fill="currentColor" d="M4 12a8 8 0 018-8V0C5.373 0 0 5.373 0 12h4zm2 5.291A7.962 7.962 0 014 12H0c0 3.042 1.135 5.824 3 7.938l3-2.647z"></path>
                </svg>
                <span>正在启动导入...</span>
              </div>
            </div>
          </template>

          <!-- 本地模式 -->
          <template v-else>
            <div>
              <label class="block text-sm font-medium text-stone-700 mb-2">Version Name</label>
              <input 
                v-model="newVersionName"
                type="text" 
                placeholder="e.g., 1.0.0, 1.2.3-beta, 2.0.0-rc.1"
                class="w-full h-10 px-3 rounded-lg border border-stone-300 text-sm focus:outline-none focus:border-emerald-600 focus:ring-1 focus:ring-emerald-600"
                @keyup.enter="handleLocalAddVersion"
              />
              <p class="mt-1 text-xs text-stone-500">Semantic Versioning (e.g., 1.0.0, 1.2.3-beta). The 'v' prefix will be added automatically.</p>
            </div>
          </template>
        </div>

        <!-- 底部按钮 -->
        <div class="mt-6 flex justify-end gap-2">
          <button 
            type="button" 
            class="rounded-lg px-3 py-2 text-base transition-colors bg-stone-200 text-stone-800 hover:bg-stone-300 disabled:opacity-50"
            :disabled="importing"
            @click="handleClose"
          >
            Cancel
          </button>
          
          <!-- GitHub 模式的 Add Version 按钮 -->
          <button 
            v-if="isGitHubSource"
            type="button" 
            class="rounded-lg px-3 py-2 text-base transition-colors bg-emerald-600 text-white hover:bg-emerald-700 disabled:bg-stone-300 disabled:cursor-not-allowed flex items-center gap-2"
            :disabled="!selectedGitRef || importing"
            @click="handleGitHubImport"
          >
            <svg v-if="importing" class="h-5 w-5 animate-spin" xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 24 24">
              <circle class="opacity-25" cx="12" cy="12" r="10" stroke="currentColor" stroke-width="4"></circle>
              <path class="opacity-75" fill="currentColor" d="M4 12a8 8 0 018-8V0C5.373 0 0 5.373 0 12h4zm2 5.291A7.962 7.962 0 014 12H0c0 3.042 1.135 5.824 3 7.938l3-2.647z"></path>
            </svg>
            <svg v-else xmlns="http://www.w3.org/2000/svg" width="20" height="20" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="1.5" stroke-linecap="round" stroke-linejoin="round">
              <path d="M12 5l0 14"></path>
              <path d="M5 12l14 0"></path>
            </svg>
            {{ importing ? 'Importing...' : 'Add Version' }}
          </button>
          
          <!-- 本地模式的 Create Version 按钮 -->
          <button 
            v-else
            type="button" 
            class="rounded-lg px-3 py-2 text-base transition-colors bg-emerald-600 text-white hover:bg-emerald-700 disabled:bg-stone-300 disabled:cursor-not-allowed flex items-center gap-2"
            :disabled="!newVersionName.trim()"
            @click="handleLocalAddVersion"
          >
            <svg xmlns="http://www.w3.org/2000/svg" width="20" height="20" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="1.5" stroke-linecap="round" stroke-linejoin="round">
              <path d="M12 5l0 14"></path>
              <path d="M5 12l14 0"></path>
            </svg>
            Create Version
          </button>
        </div>
      </div>
    </div>
  </Teleport>
</template>

<script setup lang="ts">
import { ref, computed, watch, nextTick } from 'vue'
import { ElMessage } from 'element-plus'
import { getGitHubReleases, importFromGitHub, createVersion } from '@/api/library'

// Props
const props = defineProps<{
  visible: boolean
  libraryId: number
  libraryName: string
  sourceType: string  // 'github' | 'local' | 'website'
  sourceUrl: string   // e.g., 'owner/repo'
}>()

// Emits
const emit = defineEmits<{
  (e: 'update:visible', value: boolean): void
  (e: 'success', version: string): void
}>()

// 状态
const newVersionName = ref('')
const githubReleases = ref<string[]>([])
const loadingReleases = ref(false)
const selectedGitRef = ref('')
const showRefDropdown = ref(false)
const refSearchQuery = ref('')
const importing = ref(false)

// 计算属性
const isGitHubSource = computed(() => props.sourceType === 'github')

const filteredReleases = computed(() => {
  if (!refSearchQuery.value) return githubReleases.value
  const query = refSearchQuery.value.toLowerCase()
  return githubReleases.value.filter(v => v.toLowerCase().includes(query))
})

// 监听 visible 变化
watch(() => props.visible, async (newVal) => {
  if (newVal && isGitHubSource.value && props.sourceUrl) {
    await fetchGitHubReleases()
  }
})

// 获取 GitHub 版本列表
const fetchGitHubReleases = async () => {
  if (!props.sourceUrl) return
  
  loadingReleases.value = true
  try {
    const res = await getGitHubReleases(props.sourceUrl)
    githubReleases.value = res.versions || []
  } catch (error) {
    console.error('Failed to fetch GitHub releases:', error)
    ElMessage.error('Failed to fetch GitHub versions')
  } finally {
    loadingReleases.value = false
  }
}

// 选择 GitHub ref
const selectGitRef = (ref: string) => {
  selectedGitRef.value = ref
  refSearchQuery.value = ref
  showRefDropdown.value = false
}

// 处理 GitHub 导入
const handleGitHubImport = async () => {
  if (!selectedGitRef.value) {
    ElMessage.warning('Please select a tag or branch')
    return
  }

  importing.value = true

  try {
    await importFromGitHub(
      props.libraryId,
      {
        repo: props.sourceUrl,
        tag: selectedGitRef.value,
        version: selectedGitRef.value
      }
    )
    ElMessage.success('导入已启动，跳转到控制台查看进度')
    const version = selectedGitRef.value
    resetState()
    emit('update:visible', false)
    // 使用 nextTick 确保模态框关闭后再触发跳转
    await nextTick()
    emit('success', version)
  } catch (error: any) {
    console.error('GitHub import failed:', error)
    ElMessage.error('导入失败: ' + (error.message || '未知错误'))
  } finally {
    importing.value = false
  }
}

// 处理本地模式添加版本
const handleLocalAddVersion = async () => {
  if (!newVersionName.value.trim()) {
    ElMessage.warning('Please enter a version name')
    return
  }

  try {
    const versionWithPrefix = newVersionName.value.startsWith('v') 
      ? newVersionName.value 
      : `v${newVersionName.value}`
    
    await createVersion(props.libraryId, versionWithPrefix)
    ElMessage.success('版本创建成功，跳转到控制台')
    resetState()
    emit('update:visible', false)
    // 使用 nextTick 确保模态框关闭后再触发跳转
    await nextTick()
    emit('success', versionWithPrefix)
  } catch (error) {
    console.error('Failed to add version:', error)
  }
}

// 关闭弹窗
const handleClose = () => {
  if (importing.value) return
  resetState()
  emit('update:visible', false)
}

// 重置状态
const resetState = () => {
  newVersionName.value = ''
  selectedGitRef.value = ''
  refSearchQuery.value = ''
  showRefDropdown.value = false
  importing.value = false
}

// 点击外部关闭下拉列表
const handleClickOutside = (event: MouseEvent) => {
  const target = event.target as HTMLElement
  if (!target.closest('.relative')) {
    showRefDropdown.value = false
  }
}
</script>
