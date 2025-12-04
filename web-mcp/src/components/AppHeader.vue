<template>
  <header class="left-0 right-0 top-0 z-40">
    <div class="mx-auto flex w-full max-w-[880px] flex-col items-start justify-between border-b border-stone-200 px-4 md:h-[88px] md:flex-row md:items-center lg:px-0">
      <!-- 左侧 -->
      <div class="flex w-full items-center justify-between py-4 md:w-auto md:py-0">
        <div class="flex items-center gap-3">
          <!-- Logo 框 -->
          <router-link class="inline-flex items-center" to="/">
            <div class="flex h-10 items-center justify-center gap-2 rounded-lg border border-stone-300 bg-white px-3 hover:bg-stone-50">
              <svg class="h-5 w-5" viewBox="0 0 24 24" fill="none" xmlns="http://www.w3.org/2000/svg">
                <rect x="3" y="3" width="18" height="18" rx="3" fill="#1c1917"/>
                <rect x="6" y="6" width="5" height="5" rx="1" fill="#10b981"/>
                <rect x="13" y="6" width="5" height="5" rx="1" fill="#10b981"/>
                <rect x="6" y="13" width="5" height="5" rx="1" fill="#10b981"/>
                <rect x="13" y="13" width="5" height="5" rx="1" fill="#10b981"/>
              </svg>
              <span class="text-base font-semibold text-stone-900">Context7</span>
            </div>
          </router-link>
          <!-- 已登录：Personal 下拉菜单 -->
          <PersonalDropdown v-if="isLoggedIn" :user-email="userEmail" :user-plan="userPlan" />
          
          <!-- 未登录：Sign in 按钮（在 Logo 旁边） -->
          <button 
            v-else
            class="hidden h-10 items-center justify-center gap-2 rounded-lg border border-stone-300 bg-white px-3 text-base font-normal leading-none text-stone-800 transition-colors hover:bg-stone-50 md:inline-flex md:whitespace-nowrap"
            @click="$emit('sign-in')"
          >
            <svg xmlns="http://www.w3.org/2000/svg" width="18" height="18" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="1.5" stroke-linecap="round" stroke-linejoin="round">
              <path d="M8 7a4 4 0 1 0 8 0a4 4 0 0 0 -8 0"></path>
              <path d="M6 21v-2a4 4 0 0 1 4 -4h4a4 4 0 0 1 4 4v2"></path>
            </svg>
            Sign in
          </button>
        </div>
      </div>
      <!-- 右侧导航 -->
      <div class="hidden md:flex md:w-auto md:flex-row md:items-center md:gap-3">
        <span class="hidden text-base text-stone-500 underline underline-offset-2 md:block">Plans</span>
        <span class="hidden h-4 w-px bg-stone-400 md:block"></span>
        <span class="hidden text-base text-stone-500 underline underline-offset-2 md:block">Learn</span>
        <span class="hidden h-4 w-px bg-stone-400 md:block"></span>
        <span class="hidden text-base text-stone-500 underline underline-offset-2 md:block">Try Live</span>
        <span class="hidden h-4 w-px bg-stone-400 md:block"></span>
        <span class="hidden text-base text-stone-500 underline underline-offset-2 md:block">Install</span>
        
        <!-- Add Docs 按钮（始终显示在右侧） -->
        <button 
          class="hidden h-10 items-center justify-center gap-2 rounded-lg bg-emerald-600 px-3 text-base font-normal leading-none text-white transition-colors hover:bg-emerald-700 md:inline-flex md:whitespace-nowrap"
          @click="$emit('add-docs')"
        >
          <svg xmlns="http://www.w3.org/2000/svg" width="18" height="18" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round">
            <path d="M12 5l0 14"></path>
            <path d="M5 12l14 0"></path>
          </svg>
          Add Docs
        </button>
      </div>
    </div>
  </header>
</template>

<script setup lang="ts">
import PersonalDropdown from '@/components/PersonalDropdown.vue'

interface Props {
  isLoggedIn?: boolean
  userEmail?: string
  userPlan?: string
}

withDefaults(defineProps<Props>(), {
  isLoggedIn: false,
  userEmail: 'user@example.com',
  userPlan: 'FREE'
})

defineEmits<{
  (e: 'add-docs'): void
  (e: 'sign-in'): void
}>()
</script>
