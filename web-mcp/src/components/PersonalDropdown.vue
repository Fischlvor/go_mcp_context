<template>
  <div class="hidden md:block">
    <!-- 背景遮罩 -->
    <div 
      v-if="isOpen" 
      class="fixed inset-0 z-40 transition-all duration-300 visible opacity-100 backdrop-blur-xs"
      @click="isOpen = false"
    ></div>
    
    <div class="relative">
      <!-- 触发按钮 -->
      <div 
        class="flex h-10 items-center gap-2 rounded-lg border border-stone-300 bg-transparent px-3 text-base font-medium hover:border-stone-400 cursor-pointer"
        @click="isOpen = !isOpen"
      >
        <button class="flex flex-1 items-center justify-between py-2 text-stone-800 hover:text-stone-900 md:justify-start md:gap-2">
          <div class="flex items-center gap-2">
            <span class="max-w-[120px] truncate leading-normal">{{ currentTeam }}</span>
            <svg xmlns="http://www.w3.org/2000/svg" width="18" height="18" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="1.5" stroke-linecap="round" stroke-linejoin="round">
              <path d="M8 9l4 -4l4 4"></path>
              <path d="M16 15l-4 4l-4 -4"></path>
            </svg>
          </div>
        </button>
        <router-link 
          to="/dashboard" 
          class="flex items-center justify-end py-2 font-normal text-stone-500 underline underline-offset-2 hover:text-stone-700"
          @click.stop
        >
          Dashboard
        </router-link>
      </div>

      <!-- 下拉菜单 -->
      <div 
        :class="[
          'absolute -left-2 -top-[9px] z-50 w-[calc(100vw-1rem)] space-y-3 transition-all duration-300 ease-in-out sm:w-[280px]',
          isOpen ? 'visible opacity-100' : 'pointer-events-none invisible opacity-0'
        ]"
      >
        <!-- 团队选择区域 -->
        <div class="rounded-xl border border-stone-300 bg-white shadow-xl">
          <div class="rounded-xl bg-white px-2 py-2">
            <!-- 当前选中的团队 -->
            <div>
              <button 
                class="flex w-full items-center rounded-lg bg-emerald-50 px-3 py-2 text-left text-base font-medium text-stone-800 hover:bg-emerald-100"
                @click="selectTeam('Personal')"
              >
                <div class="flex w-full items-center justify-between">
                  <span class="truncate leading-normal">Personal</span>
                  <span class="inline-flex items-center justify-center">
                    <svg xmlns="http://www.w3.org/2000/svg" width="20" height="20" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="1.5" stroke-linecap="round" stroke-linejoin="round" class="text-stone-500">
                      <path d="M9 10a3 3 0 1 0 6 0a3 3 0 0 0 -6 0"></path>
                      <path d="M6 21v-1a4 4 0 0 1 4 -4h4a4 4 0 0 1 4 4v1"></path>
                      <path d="M3 5a2 2 0 0 1 2 -2h14a2 2 0 0 1 2 2v14a2 2 0 0 1 -2 2h-14a2 2 0 0 1 -2 -2v-14z"></path>
                    </svg>
                  </span>
                </div>
              </button>
            </div>
            
            <!-- 分隔线 -->
            <div class="py-2">
              <div class="border-b border-stone-200"></div>
            </div>
            
            <!-- 创建团队按钮 -->
            <button class="flex w-full items-center gap-2 rounded-lg px-3 py-2 text-base font-normal text-emerald-600 hover:bg-emerald-50 hover:text-emerald-700">
              <svg xmlns="http://www.w3.org/2000/svg" width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round">
                <path d="M12 5l0 14"></path>
                <path d="M5 12l14 0"></path>
              </svg>
              <span>Create Team</span>
            </button>
          </div>
        </div>

        <!-- 用户信息区域 -->
        <div class="gap-2 rounded-xl bg-stone-800 px-6 py-4 shadow-xl">
          <div class="flex items-center gap-2 truncate text-base text-stone-50">
            <span class="truncate">{{ userEmail }}</span>
            <span class="inline-flex items-center rounded-md px-1.5 py-1 font-mono text-[11px] font-medium leading-none tracking-wide text-white bg-stone-600">
              {{ userPlan }}
            </span>
          </div>
          <div>
            <button class="text-sm font-normal text-stone-50 underline hover:text-stone-300">
              Change Plan
            </button>
          </div>
          <div>
            <button class="text-sm font-normal text-stone-50 underline hover:text-stone-300">
              Billing
            </button>
          </div>
          <div>
            <button 
              class="text-sm font-normal text-stone-50 underline hover:text-stone-300"
              @click="handleSignOut"
            >
              Sign Out
            </button>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref } from 'vue'
import { useRouter } from 'vue-router'
import { useUser } from '@/stores/user'

const router = useRouter()
const { logout } = useUser()

// Props
interface Props {
  userEmail?: string
  userPlan?: string
}

const props = withDefaults(defineProps<Props>(), {
  userEmail: 'user@example.com',
  userPlan: 'FREE'
})

// State
const isOpen = ref(false)
const currentTeam = ref('Personal')

// Methods
const selectTeam = (team: string) => {
  currentTeam.value = team
  isOpen.value = false
}

const handleSignOut = async () => {
  isOpen.value = false
  await logout()
  router.push('/')
}

// 点击外部关闭
const closeDropdown = () => {
  isOpen.value = false
}

defineExpose({
  closeDropdown
})
</script>
