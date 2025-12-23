<template>
  <Teleport to="body">
    <div 
      v-if="visible" 
      class="fixed inset-0 z-50 flex items-center justify-center bg-black/50"
      @click.self="handleClose"
    >
      <div class="w-full max-w-md rounded-2xl bg-white p-6 shadow-xl">
        <!-- 标题 -->
        <div class="mb-6">
          <h3 class="text-xl font-semibold text-stone-900">Add Documentation</h3>
          <p class="mt-1 text-sm text-stone-500">{{ sourceType === 'github' ? 'Import from GitHub repository' : 'Create a local documentation library' }}</p>
        </div>

        <!-- 来源选择 -->
        <div class="mb-4">
          <label class="mb-2 block text-sm font-medium text-stone-700">Source</label>
          <div class="flex gap-2">
            <button
              :class="[
                'flex-1 rounded-lg border px-4 py-2 text-sm font-medium transition-colors',
                sourceType === 'github'
                  ? 'border-emerald-500 bg-emerald-50 text-emerald-700'
                  : 'border-stone-300 bg-white text-stone-700 hover:bg-stone-50'
              ]"
              @click="sourceType = 'github'"
            >
              <svg class="mr-2 inline-block h-4 w-4" viewBox="0 0 24 24" fill="currentColor">
                <path d="M12 0c-6.626 0-12 5.373-12 12 0 5.302 3.438 9.8 8.207 11.387.599.111.793-.261.793-.577v-2.234c-3.338.726-4.033-1.416-4.033-1.416-.546-1.387-1.333-1.756-1.333-1.756-1.089-.745.083-.729.083-.729 1.205.084 1.839 1.237 1.839 1.237 1.07 1.834 2.807 1.304 3.492.997.107-.775.418-1.305.762-1.604-2.665-.305-5.467-1.334-5.467-5.931 0-1.311.469-2.381 1.236-3.221-.124-.303-.535-1.524.117-3.176 0 0 1.008-.322 3.301 1.23.957-.266 1.983-.399 3.003-.404 1.02.005 2.047.138 3.006.404 2.291-1.552 3.297-1.23 3.297-1.23.653 1.653.242 2.874.118 3.176.77.84 1.235 1.911 1.235 3.221 0 4.609-2.807 5.624-5.479 5.921.43.372.823 1.102.823 2.222v3.293c0 .319.192.694.801.576 4.765-1.589 8.199-6.086 8.199-11.386 0-6.627-5.373-12-12-12z"/>
              </svg>
              GitHub
            </button>
            <button
              :class="[
                'flex-1 rounded-lg border px-4 py-2 text-sm font-medium transition-colors',
                sourceType === 'local'
                  ? 'border-emerald-500 bg-emerald-50 text-emerald-700'
                  : 'border-stone-300 bg-white text-stone-700 hover:bg-stone-50'
              ]"
              @click="sourceType = 'local'"
            >
              <svg class="mr-2 inline-block h-4 w-4" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
                <path d="M3 7v10a2 2 0 002 2h14a2 2 0 002-2V9a2 2 0 00-2-2h-6l-2-2H5a2 2 0 00-2 2z"/>
              </svg>
              Local
            </button>
          </div>
        </div>

        <!-- GitHub URL 输入 -->
        <div v-if="sourceType === 'github'" class="mb-6">
          <label class="mb-2 block text-sm font-medium text-stone-700">GitHub Repository URL</label>
          <input
            v-model="githubUrl"
            type="text"
            placeholder="https://github.com/owner/repo"
            class="w-full rounded-lg border border-stone-300 px-3 py-2 text-stone-800 placeholder-stone-400 focus:border-emerald-500 focus:outline-none focus:ring-1 focus:ring-emerald-500"
            :disabled="loading"
            @keyup.enter="handleSubmit"
          />
          <p class="mt-1 text-xs text-stone-500">
            Example: https://github.com/gin-gonic/gin
          </p>
        </div>

        <!-- Local 输入 -->
        <div v-if="sourceType === 'local'" class="mb-6 space-y-4">
          <div>
            <label class="mb-2 block text-sm font-medium text-stone-700">Library Name</label>
            <input
              v-model="localName"
              type="text"
              placeholder="My Documentation"
              class="w-full rounded-lg border border-stone-300 px-3 py-2 text-stone-800 placeholder-stone-400 focus:border-emerald-500 focus:outline-none focus:ring-1 focus:ring-emerald-500"
              :disabled="loading"
              @keyup.enter="handleSubmit"
            />
          </div>
          <div>
            <label class="mb-2 block text-sm font-medium text-stone-700">Description (optional)</label>
            <input
              v-model="localDescription"
              type="text"
              placeholder="A brief description"
              class="w-full rounded-lg border border-stone-300 px-3 py-2 text-stone-800 placeholder-stone-400 focus:border-emerald-500 focus:outline-none focus:ring-1 focus:ring-emerald-500"
              :disabled="loading"
            />
          </div>
        </div>

        <!-- 错误提示 -->
        <div v-if="errorMessage" class="mb-4 rounded-lg border border-red-300 bg-red-50 p-3 text-sm text-red-700">
          {{ errorMessage }}
        </div>

        <!-- 按钮 -->
        <div class="flex justify-end gap-3">
          <button
            class="rounded-lg px-4 py-2 text-sm font-medium text-stone-600 transition-colors hover:bg-stone-100"
            :disabled="loading"
            @click="handleClose"
          >
            Cancel
          </button>
          <button
            class="rounded-lg bg-emerald-600 px-4 py-2 text-sm font-medium text-white transition-colors hover:bg-emerald-700 disabled:cursor-not-allowed disabled:opacity-50"
            :disabled="loading || !canSubmit"
            @click="handleSubmit"
          >
            <span v-if="loading" class="flex items-center gap-2">
              <svg class="h-4 w-4 animate-spin" viewBox="0 0 24 24" fill="none">
                <circle class="opacity-25" cx="12" cy="12" r="10" stroke="currentColor" stroke-width="4"></circle>
                <path class="opacity-75" fill="currentColor" d="M4 12a8 8 0 018-8V0C5.373 0 0 5.373 0 12h4zm2 5.291A7.962 7.962 0 014 12H0c0 3.042 1.135 5.824 3 7.938l3-2.647z"></path>
              </svg>
              {{ sourceType === 'github' ? 'Importing...' : 'Creating...' }}
            </span>
            <span v-else>{{ sourceType === 'github' ? 'Import' : 'Create' }}</span>
          </button>
        </div>
      </div>
    </div>
  </Teleport>
</template>

<script setup lang="ts">
import { ref, watch, computed } from 'vue'
import { useRouter } from 'vue-router'
import { initImportFromGitHub, createLibrary } from '@/api/library'
import { ElMessage } from 'element-plus'

interface Props {
  visible: boolean
}

const props = defineProps<Props>()

const emit = defineEmits<{
  (e: 'update:visible', value: boolean): void
  (e: 'success', libraryId: number): void
}>()

const router = useRouter()

const sourceType = ref<'github' | 'local'>('github')
const githubUrl = ref('')
const localName = ref('')
const localDescription = ref('')
const loading = ref(false)
const errorMessage = ref('')

// 是否可提交
const canSubmit = computed(() => {
  if (sourceType.value === 'github') {
    return githubUrl.value.trim().length > 0
  } else {
    return localName.value.trim().length > 0
  }
})

// 重置状态
watch(() => props.visible, (val) => {
  if (val) {
    githubUrl.value = ''
    localName.value = ''
    localDescription.value = ''
    errorMessage.value = ''
    loading.value = false
  }
})

const handleClose = () => {
  if (!loading.value) {
    emit('update:visible', false)
  }
}

const handleSubmit = async () => {
  if (!canSubmit.value || loading.value) return

  loading.value = true
  errorMessage.value = ''

  try {
    if (sourceType.value === 'github') {
      // GitHub 导入
      const result = await initImportFromGitHub(githubUrl.value.trim())
      
      ElMessage.success('导入已启动，正在跳转...')
      emit('update:visible', false)
      emit('success', result.library_id)
      
      // 跳转到库详情页的 logs tab
      router.push({
        name: 'library-detail',
        params: { id: result.library_id },
        query: { tab: 'logs' }
      })
    } else {
      // Local 创建
      const result = await createLibrary({
        name: localName.value.trim(),
        description: localDescription.value.trim(),
        source_type: 'local'
      })
      
      ElMessage.success('库创建成功')
      emit('update:visible', false)
      emit('success', result.data.id)
      
      // 跳转到库详情页
      router.push({
        name: 'library-detail',
        params: { id: result.data.id }
      })
    }
  } catch (error: any) {
    console.error('Operation failed:', error)
    errorMessage.value = error.message || '操作失败'
  } finally {
    loading.value = false
  }
}
</script>
