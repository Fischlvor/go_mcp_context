<template>
  <footer class="relative z-10 mt-auto w-full">
    <div class="mx-auto mt-20 flex h-auto min-h-[69px] w-full max-w-[880px] flex-col items-start justify-start gap-3 border-t border-stone-200 px-4 py-6 lg:flex-row lg:items-center lg:justify-between lg:px-0">
      <div class="w-full text-sm font-normal leading-tight text-stone-500 md:text-base lg:w-auto">
        © {{ currentYear }}, Context7 is an 
        <a href="https://upstash.com" target="_blank" rel="noopener noreferrer" class="font-normal text-stone-500 underline transition-colors hover:text-stone-700">Upstash</a> 
        project
      </div>
      <div class="flex w-full items-center justify-start gap-3 overflow-x-auto whitespace-nowrap px-1 py-1 text-sm font-normal leading-tight text-stone-500 md:text-base lg:w-auto lg:justify-end">
        <router-link to="/about" class="font-normal text-stone-500 underline transition-colors hover:text-stone-700">About</router-link>
        <router-link to="/contact" class="font-normal text-stone-500 underline transition-colors hover:text-stone-700">Contact</router-link>
        <a href="https://upstash.com/docs/common/help/legal" target="_blank" rel="noopener noreferrer" class="font-normal text-stone-500 underline transition-colors hover:text-stone-700">Legal</a>
        <a href="https://x.com/Context7AI" target="_blank" rel="noopener noreferrer" class="font-normal text-stone-500 underline transition-colors hover:text-stone-700">Follow on X</a>
      </div>
    </div>
  </footer>
</template>

<script setup lang="ts">
import { computed } from 'vue'

const currentYear = computed(() => new Date().getFullYear())
</script>
