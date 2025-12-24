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
                  <!-- 源链接或 Local -->
                  <div class="w-fit max-w-full">
                    <a 
                      v-if="library.source_type === 'github' && library.source_url"
                      :href="`https://github.com/${library.source_url}`"
                      target="_blank"
                      rel="noopener noreferrer"
                      class="block overflow-hidden text-ellipsis whitespace-nowrap text-base font-normal leading-normal text-stone-500 underline decoration-solid decoration-from-font hover:text-stone-700"
                      :title="`https://github.com/${library.source_url}`"
                    >
                      https://github.com/{{ library.source_url }}
                    </a>
                    <span 
                      v-else
                      class="block overflow-hidden text-ellipsis whitespace-nowrap text-base font-normal leading-normal text-stone-500"
                    >
                      Local
                    </span>
                  </div>
                  <!-- 可展开的描述 -->
                  <h3 class="text-base font-normal leading-[140%] text-stone-500">
                    <span v-if="!expandDescription && isTruncated(library.description)" class="inline-flex items-center gap-0">
                      <span class="overflow-hidden text-ellipsis whitespace-nowrap">{{ getTruncatedText(library.description) }}</span><span 
                        class="cursor-pointer text-emerald-600 hover:text-emerald-700 hover:underline flex-shrink-0"
                        @click="expandDescription = true"
                      >...</span>
                    </span>
                    <span v-else-if="expandDescription">
                      {{ library.description || 'No description' }}<span 
                        class="cursor-pointer text-emerald-600 hover:text-emerald-700 hover:underline"
                        @click="expandDescription = false"
                      > collapse</span>
                    </span>
                    <span v-else>
                      {{ library.description || 'No description' }}
                    </span>
                  </h3>
                </div>
                <!-- Manage 按钮 -->
                <div class="relative inline-flex">
                  <router-link 
                    :to="`/libraries/${libraryId}/admin`"
                    class="flex h-8 items-center justify-center gap-1.5 rounded-lg border border-stone-300 text-base text-stone-500 transition hover:border-stone-400 px-3 py-2 !border-emerald-300 bg-emerald-50 hover:bg-emerald-100"
                  >
                    <svg xmlns="http://www.w3.org/2000/svg" width="24" height="24" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="1.5" stroke-linecap="round" stroke-linejoin="round" class="h-5 w-5 text-emerald-600">
                      <path d="M10.325 4.317c.426 -1.756 2.924 -1.756 3.35 0a1.724 1.724 0 0 0 2.573 1.066c1.543 -.94 3.31 .826 2.37 2.37a1.724 1.724 0 0 0 1.065 2.572c1.756 .426 1.756 2.924 0 3.35a1.724 1.724 0 0 0 -1.066 2.573c.94 1.543 -.826 3.31 -2.37 2.37a1.724 1.724 0 0 0 -2.572 1.065c-.426 1.756 -2.924 1.756 -3.35 0a1.724 1.724 0 0 0 -2.573 -1.066c-1.543 .94 -3.31 -.826 -2.37 -2.37a1.724 1.724 0 0 0 -1.065 -2.572c-1.756 -.426 -1.756 -2.924 0 -3.35a1.724 1.724 0 0 0 1.066 -2.573c-.94 -1.543 .826 -3.31 2.37 -2.37c1 .608 2.296 .07 2.572 -1.065z"></path>
                      <path d="M9 12a3 3 0 1 0 6 0a3 3 0 0 0 -6 0"></path>
                    </svg>
                    <span class="text-emerald-600">Manage</span>
                  </router-link>
                </div>
              </div>

              <!-- 状态标签 -->
              <div class="flex flex-col-reverse gap-4 sm:flex-row sm:flex-wrap sm:items-start sm:justify-between">
                <div class="flex flex-wrap gap-2 text-sm sm:gap-1">
                  <div class="flex items-center gap-1 rounded-md bg-emerald-50 px-2 py-1">
                    <div class="flex h-5 w-5 items-center justify-center text-emerald-800">
                      <svg xmlns="http://www.w3.org/2000/svg" width="24" height="24" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="1.5" stroke-linecap="round" stroke-linejoin="round" class="tabler-icon tabler-icon-circle-check">
                        <path d="M12 12m-9 0a9 9 0 1 0 18 0a9 9 0 1 0 -18 0"></path>
                        <path d="M9 12l2 2l4 -4"></path>
                      </svg>
                    </div>
                    <span class="text-sm font-normal leading-[100%] tracking-[0%] text-emerald-800">{{ library.status === 'active' ? 'Completed' : library.status }}</span>
                  </div>
                  <div class="flex items-center gap-1 rounded-md bg-stone-100 px-3 py-1.5">
                    <span class="text-sm font-normal leading-[100%] tracking-[0%] text-stone-500">Tokens:</span>
                    <span class="font-variant-numeric-zero:slashed-zero text-sm font-normal leading-[100%] tracking-[0%] text-stone-800">{{ formatNumber(library.token_count || 0) }}</span>
                  </div>
                  <div class="flex items-center gap-1 rounded-md bg-stone-100 px-3 py-1.5">
                    <span class="text-sm font-normal leading-[100%] tracking-[0%] text-stone-500">Documents:</span>
                    <span class="font-variant-numeric-zero:slashed-zero text-sm font-normal leading-[100%] tracking-[0%] text-stone-800">{{ formatNumber(library.document_count || 0) }}</span>
                  </div>
                  <div class="flex items-center gap-1 rounded-md bg-stone-100 px-3 py-1.5">
                    <span class="text-sm font-normal leading-[100%] tracking-[0%] text-stone-500">Update:</span>
                    <span class="font-variant-numeric-zero:slashed-zero text-sm font-normal leading-[100%] tracking-[0%] text-stone-800">{{ formatDate(library.updated_at) }}</span>
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
              <div class="flex flex-wrap items-center gap-2.5 sm:gap-1.5">
                <!-- Terminal 按钮 -->
                <button 
                  :class="[
                    'flex h-8 items-center justify-center gap-1.5 rounded-lg border text-base transition w-8',
                    activeTab === 'logs' 
                      ? 'border-stone-700 bg-stone-700 hover:border-stone-700' 
                      : 'border-stone-300 text-stone-500 hover:border-stone-400'
                  ]"
                  @click="activeTab = activeTab === 'logs' ? 'context' : 'logs'"
                  title="Activity Logs"
                >
                  <svg xmlns="http://www.w3.org/2000/svg" width="24" height="24" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="1.5" stroke-linecap="round" stroke-linejoin="round" :class="['h-5 w-5', activeTab === 'logs' ? 'text-stone-50' : 'text-stone-500']">
                    <path d="M8 9l3 3l-3 3"></path>
                    <path d="M13 15l3 0"></path>
                    <path d="M3 4m0 2a2 2 0 0 1 2 -2h14a2 2 0 0 1 2 2v12a2 2 0 0 1 -2 2h-14a2 2 0 0 1 -2 -2z"></path>
                  </svg>
                </button>
                
                <!-- 刷新版本按钮（重新处理文档） -->
                <button 
                  class="flex h-8 items-center justify-center gap-1.5 rounded-lg border border-stone-300 text-base text-stone-500 transition hover:border-stone-400 w-8"
                  :class="{ 'opacity-50 cursor-not-allowed': refreshingVersion }"
                  :disabled="refreshingVersion"
                  @click="handleRefreshVersion"
                  title="Refresh version (reprocess all documents)"
                >
                  <svg xmlns="http://www.w3.org/2000/svg" width="24" height="24" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="1.5" stroke-linecap="round" stroke-linejoin="round" class="h-5 w-5 text-stone-500" :class="{ 'animate-spin': refreshingVersion }">
                    <path d="M20 11a8.1 8.1 0 0 0 -15.5 -2m-.5 -4v4h4"></path>
                    <path d="M4 13a8.1 8.1 0 0 0 15.5 2m.5 4v-4h-4"></path>
                  </svg>
                </button>
                
                <!-- 版本选择下拉框 -->
                <div class="relative">
                  <button 
                    class="flex h-8 items-center gap-1.5 rounded-lg border border-stone-300 px-3 text-sm text-stone-600 transition hover:border-stone-400 hover:bg-stone-50"
                    @click="showVersionDropdown = !showVersionDropdown"
                  >
                    <svg xmlns="http://www.w3.org/2000/svg" width="24" height="24" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="1.5" stroke-linecap="round" stroke-linejoin="round" class="h-4 w-4 text-stone-500">
                      <path d="M12 12m-3 0a3 3 0 1 0 6 0a3 3 0 1 0 -6 0"></path>
                      <path d="M12 3l0 6"></path>
                      <path d="M12 15l0 6"></path>
                    </svg>
                    <span>{{ currentVersionDisplay }}</span>
                    <svg xmlns="http://www.w3.org/2000/svg" width="24" height="24" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="1.5" stroke-linecap="round" stroke-linejoin="round" class="h-4 w-4 text-stone-400">
                      <path d="M6 9l6 6l6 -6"></path>
                    </svg>
                  </button>
                  
                  <!-- 下拉菜单 -->
                  <div 
                    v-if="showVersionDropdown" 
                    class="absolute right-0 top-full z-50 mt-1 w-48 rounded-md border border-stone-200 bg-white py-1 shadow-lg"
                  >
                    <!-- Latest (默认版本) 选项 -->
                    <button 
                      class="mx-1 my-0.5 flex w-[calc(100%-8px)] cursor-pointer items-center gap-1.5 rounded-md px-4 py-2 text-sm text-stone-700 transition-all duration-150 hover:bg-stone-50"
                      :class="{ 'bg-stone-100': isLatestVersion }"
                      @click="selectVersion('')"
                    >
                      <svg xmlns="http://www.w3.org/2000/svg" width="24" height="24" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="1.5" stroke-linecap="round" stroke-linejoin="round" class="h-5 w-5 text-stone-500">
                        <path d="M12 12m-3 0a3 3 0 1 0 6 0a3 3 0 1 0 -6 0"></path>
                        <path d="M12 3l0 6"></path>
                        <path d="M12 15l0 6"></path>
                      </svg>
                      <span class="truncate" :title="library.default_version">{{ library.default_version }}</span>
                    </button>
                    
                    <!-- 版本列表 -->
                    <button 
                      v-for="ver in sortedVersions" 
                      :key="ver"
                      class="mx-1 my-0.5 flex w-[calc(100%-8px)] cursor-pointer items-center gap-1.5 rounded-md px-4 py-2 text-sm text-stone-700 transition-all duration-150 hover:bg-stone-50"
                      :class="{ 'bg-stone-100': isCurrentVersion(ver) }"
                      @click="selectVersion(ver)"
                    >
                      <svg xmlns="http://www.w3.org/2000/svg" width="24" height="24" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="1.5" stroke-linecap="round" stroke-linejoin="round" class="h-5 w-5 text-stone-500">
                        <path d="M7 18m-2 0a2 2 0 1 0 4 0a2 2 0 1 0 -4 0"></path>
                        <path d="M7 6m-2 0a2 2 0 1 0 4 0a2 2 0 1 0 -4 0"></path>
                        <path d="M17 12m-2 0a2 2 0 1 0 4 0a2 2 0 1 0 -4 0"></path>
                        <path d="M7 8l0 8"></path>
                        <path d="M7 8a4 4 0 0 0 4 4h4"></path>
                      </svg>
                      <span class="truncate" :title="ver">{{ ver }}</span>
                    </button>
                    
                    <!-- 分隔线 -->
                    <div class="mx-1 my-1 h-px border-t border-stone-200"></div>
                    
                    <!-- New Version 按钮 -->
                    <button 
                      class="mx-1 my-0.5 flex w-[calc(100%-8px)] cursor-pointer items-center gap-1.5 rounded-md px-4 py-2 text-sm text-stone-700 transition-all duration-150 hover:bg-stone-50"
                      @click="openAddVersionModal"
                    >
                      <svg xmlns="http://www.w3.org/2000/svg" width="24" height="24" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="1.5" stroke-linecap="round" stroke-linejoin="round" class="h-5 w-5 text-stone-500">
                        <path d="M12 5l0 14"></path>
                        <path d="M5 12l14 0"></path>
                      </svg>
                      <span>New Version</span>
                    </button>
                  </div>
                </div>
              </div>
            </div>
            <div class="border-t border-stone-300"></div>
          </div>

          <!-- Terminal 日志面板 -->
          <div v-if="activeTab === 'logs'" class="mt-8">
            <div class="rounded-xl bg-stone-800 p-6 shadow-sm">
              <div class="h-[360px] overflow-y-auto font-mono text-sm text-stone-50 [scrollbar-color:theme(colors.stone.400)_theme(colors.stone.800)] [&::-webkit-scrollbar-thumb]:rounded-full [&::-webkit-scrollbar-thumb]:bg-stone-400 [&::-webkit-scrollbar-track]:bg-stone-800 [&::-webkit-scrollbar]:w-2">
                <div v-if="activityLogs.length === 0" class="text-stone-400">
                  No activity logs yet. Logs will appear here when you perform actions like importing documents or refreshing content.
                </div>
                <div 
                  v-for="(log, index) in activityLogs" 
                  :key="index" 
                  class="mb-1 flex"
                >
                  <span class="shrink-0 pr-3 tracking-tighter text-stone-400">{{ formatLogTime(log.time) }}</span>
                  <span :class="getLogClass(log)">{{ log.message }}</span>
                </div>
              </div>
            </div>
            <div class="mt-4 text-center">
              <button 
                class="group inline-flex items-center gap-1.5 text-stone-600 decoration-current underline-offset-2 transition-colors hover:text-emerald-600 hover:underline"
                @click="fetchLogs"
              >
                <svg xmlns="http://www.w3.org/2000/svg" width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="1.5" stroke-linecap="round" stroke-linejoin="round" class="group-hover:animate-spin">
                  <path d="M20 11a8.1 8.1 0 0 0 -15.5 -2m-.5 -4v4h4"></path>
                  <path d="M4 13a8.1 8.1 0 0 0 15.5 2m.5 4v-4h-4"></path>
                </svg>
                <span class="text-sm font-medium">Refresh Logs</span>
              </button>
            </div>
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
                
                <!-- 文档内容展示 -->
                <div class="overflow-hidden rounded-xl">
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

          <!-- Documents Tab 内容 -->
          <div v-if="activeTab === 'documents'" class="mt-8">
            <div class="rounded-3xl border border-stone-200 bg-white p-6 shadow-sm sm:p-8">
              <div class="space-y-6">
                <!-- 标题和上传按钮 -->
                <div class="flex items-center justify-between">
                  <div>
                    <h3 class="text-base font-semibold text-stone-800">Documents</h3>
                    <p class="mt-1 text-sm text-stone-500">
                      {{ version ? `Documents in version ${version}` : `Documents in version ${library.default_version || 'default'}` }}
                    </p>
                  </div>
                  <label 
                    v-if="isLoggedIn"
                    :class="[
                      'flex h-10 items-center justify-center gap-2 rounded-lg px-4 text-sm font-medium text-white transition-colors whitespace-nowrap cursor-pointer',
                      uploading ? 'bg-stone-400 cursor-not-allowed' : 'bg-emerald-600 hover:bg-emerald-700'
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
                    <span>{{ uploading ? 'Processing...' : 'Upload' }}</span>
                    <input 
                      type="file" 
                      class="hidden" 
                      accept=".md,.pdf,.docx"
                      :disabled="uploading"
                      @change="handleFileUpload"
                    />
                  </label>
                </div>

                <!-- 上传中提示（普通接口无进度，SSE 版本有进度条） -->
                <div v-if="uploading" class="rounded-lg border border-emerald-200 bg-emerald-50 p-4">
                  <div class="flex items-center gap-2">
                    <svg class="animate-spin h-5 w-5 text-emerald-600" xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 24 24">
                      <circle class="opacity-25" cx="12" cy="12" r="10" stroke="currentColor" stroke-width="4"></circle>
                      <path class="opacity-75" fill="currentColor" d="M4 12a8 8 0 018-8V0C5.373 0 0 5.373 0 12h4zm2 5.291A7.962 7.962 0 014 12H0c0 3.042 1.135 5.824 3 7.938l3-2.647z"></path>
                    </svg>
                    <span class="text-sm font-medium text-emerald-800">上传中...</span>
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
                      </tr>
                    </thead>
                    <tbody class="divide-y divide-stone-200">
                      <!-- 空状态 -->
                      <tr v-if="documents.length === 0 && !loadingDocs">
                        <td colspan="4" class="py-12 text-center">
                          <div class="flex flex-col items-center gap-2">
                            <svg xmlns="http://www.w3.org/2000/svg" width="40" height="40" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="1" stroke-linecap="round" stroke-linejoin="round" class="text-stone-300">
                              <path d="M14 3v4a1 1 0 0 0 1 1h4"></path>
                              <path d="M17 21h-10a2 2 0 0 1 -2 -2v-14a2 2 0 0 1 2 -2h7l5 5v11a2 2 0 0 1 -2 2z"></path>
                            </svg>
                            <p class="text-sm font-medium text-stone-500">No documents</p>
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
                            <span class="truncate">{{ getDisplayPath(doc.file_path) }}</span>
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
                      @click="page--; fetchDocumentsList()"
                    >
                      Previous
                    </button>
                    <button 
                      class="h-8 px-3 rounded-lg border border-stone-300 text-sm text-stone-600 hover:bg-stone-50 disabled:opacity-50"
                      :disabled="page * pageSize >= totalDocs"
                      @click="page++; fetchDocumentsList()"
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

    <!-- Add Version Modal -->
    <AddVersionModal
      v-model:visible="showAddVersionModal"
      :library-id="libraryId"
      :library-name="library.name"
      :source-type="library.source_type"
      :source-url="library.source_url"
      @success="handleVersionCreated"
    />

    <!-- 点击外部关闭下拉框 -->
    <div 
      v-if="showVersionDropdown" 
      class="fixed inset-0 z-40" 
      @click="showVersionDropdown = false"
    ></div>
  </div>
</template>

<script setup lang="ts">
import { ref, computed, onMounted, onUnmounted, watch } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { ElMessage, ElMessageBox } from 'element-plus'
import AppHeader from '@/components/AppHeader.vue'
import AppFooter from '@/components/AppFooter.vue'
import AddVersionModal from '@/components/AddVersionModal.vue'
import { useUser } from '@/stores/user'
import { getLibrary, getActivityLogs, refreshVersion } from '@/api/library'
import type { ActivityLog } from '@/api/library'
import { getDocuments, getChunks, uploadDocument } from '@/api/document'
import type { Library } from '@/api/library'
const route = useRoute()
const router = useRouter()
const { isLoggedIn, userEmail, userPlan, initUserState, redirectToSSO } = useUser()

const libraryId = computed(() => Number(route.params.id))
const version = computed(() => {
  // 从路由获取版本，如果没有则返回 undefined（让后端使用默认版本）
  const routeVersion = route.params.version as string | undefined
  return routeVersion || undefined
})
const library = ref<Library>({
  id: 0,
  name: '',
  default_version: '',
  versions: [],
  source_type: '',
  source_url: '',
  description: '',
  status: '',
  document_count: 0,
  chunk_count: 0,
  token_count: 0,
  created_at: '',
  updated_at: ''
})

// 从 URL query 读取初始 tab，默认 context
const activeTab = ref((route.query.tab as string) || 'context')
const searchTopic = ref('')
const searchMode = ref<'code' | 'info'>('code')
const searching = ref(false)
const searchResult = ref('Loading document...')
const hasSearched = ref(false)
const expandDescription = ref(false)

// Documents tab
const documents = ref<any[]>([])
const loadingDocs = ref(false)
const page = ref(1)
const pageSize = ref(10)
const totalDocs = ref(0)

// 上传状态
const uploading = ref(false)
const uploadProgress = ref(0)
const uploadMessage = ref('')

// 版本选择
const showVersionDropdown = ref(false)
const showAddVersionModal = ref(false)

// 刷新版本状态
const refreshingVersion = ref(false)

// Terminal 日志
const loadingLogs = ref(false)
const activityLogs = ref<ActivityLog[]>([])
const logPollingTimer = ref<ReturnType<typeof setInterval> | null>(null)

// 从 API 获取日志
const fetchLogs = async () => {
  // 防止重复请求
  if (loadingLogs.value) return
  
  loadingLogs.value = true
  try {
    const res = await getActivityLogs(libraryId.value, 50)
    activityLogs.value = res.logs || []
    
    // 如果状态是 complete，停止轮询
    if (res.status === 'complete') {
      stopLogPolling()
    }
  } catch (error) {
    console.error('Failed to fetch logs:', error)
  } finally {
    loadingLogs.value = false
  }
}

// 开始轮询日志
const startLogPolling = () => {
  // 先停止已有的轮询
  stopLogPolling()
  // 立即请求一次
  fetchLogs()
  // 每 2 秒轮询一次
  logPollingTimer.value = setInterval(fetchLogs, 2000)
}

// 停止轮询日志
const stopLogPolling = () => {
  if (logPollingTimer.value) {
    clearInterval(logPollingTimer.value)
    logPollingTimer.value = null
  }
}

// 刷新版本（重新处理所有文档）
const handleRefreshVersion = async () => {
  if (refreshingVersion.value) return
  
  const currentVersion = currentVersionDisplay.value
  if (!currentVersion) {
    ElMessage.warning('请先选择版本')
    return
  }
  
  try {
    await ElMessageBox.confirm(
      '这将重新处理该版本下的所有文档',
      `刷新版本 "${currentVersion}"？`,
      {
        confirmButtonText: '确定',
        cancelButtonText: '取消',
        type: 'warning'
      }
    )
  } catch {
    return // 用户取消
  }
  
  refreshingVersion.value = true
  try {
    await refreshVersion(libraryId.value, currentVersion)
    ElMessage.success('版本刷新已启动，请查看控制台日志')
    // 清空其他 tab 数据，刷新完成后需要重新加载
    documents.value = []
    searchResult.value = 'Loading document...'
    hasSearched.value = false
    // 切换到控制台 tab 并开始轮询
    activeTab.value = 'logs'
    startLogPolling()
  } catch (error: any) {
    ElMessage.error('刷新版本失败: ' + (error.message || '未知错误'))
  } finally {
    refreshingVersion.value = false
  }
}

// 获取日志样式（优先根据 status 渲染，info 再根据 event 渲染）
const getLogClass = (log: ActivityLog) => {
  // 优先根据 status 判断
  if (log.status === 'start') return 'text-purple-400'     // 开始 - 紫色
  if (log.status === 'success') return 'text-emerald-400'  // 成功 - 绿色
  if (log.status === 'error') return 'text-red-400'        // 错误 - 红色
  if (log.status === 'warning') return 'text-yellow-400'   // 警告 - 黄色
  
  // status === 'info' 时，根据 event 类型设置颜色
  const event = log.event || ''
  
  if (event === 'document.enrich') return 'text-amber-400'         // AI增强 - 琥珀色
  if (event === 'document.embed') return 'text-indigo-400'         // Embedding - 靛蓝色
  if (event === 'version.refresh') return 'text-sky-400'           // 版本刷新 - 蓝色
  
  // 其他 info 默认白色
  return 'text-stone-300'
}

// 获取文档显示路径（去掉 lib 和 version 前缀）
// 例如：mcp/docs/gin/v1.6.3/README.md -> README.md
// 例如：mcp/docs/gin/v1.6.3/examples/README.md -> examples/README.md
const getDisplayPath = (filePath: string) => {
  if (!filePath) return ''
  const parts = filePath.split('/')
  // 路径格式: mcp/docs/{lib}/{version}/{relative_path...}
  // 跳过前 4 段，返回剩余部分
  if (parts.length > 4) {
    return parts.slice(4).join('/')
  }
  // 如果路径格式不符合预期，返回文件名
  return parts[parts.length - 1]
}

// 格式化日志时间
const formatLogTime = (isoTime: string) => {
  const date = new Date(isoTime)
  return date.toLocaleString('zh-CN', {
    month: 'short',
    day: 'numeric',
    hour: '2-digit',
    minute: '2-digit',
    second: '2-digit'
  })
}

// 当前版本显示
const currentVersionDisplay = computed(() => {
  const routeVersion = route.params.version as string | undefined
  // 如果路由没有指定版本，显示默认版本
  if (!routeVersion) {
    return library.value.default_version || 'Latest'
  }
  return routeVersion
})

// 排序后的版本列表（不包含默认版本）
const sortedVersions = computed(() => {
  const versions = library.value.versions || []
  return versions.filter(v => v !== library.value.default_version)
})

// 判断是否是最新版本（默认版本）
const isLatestVersion = computed(() => {
  // 如果路由没有指定版本，或者版本等于默认版本，则是最新版本
  const routeVersion = route.params.version as string | undefined
  return !routeVersion || routeVersion === library.value.default_version
})

// 判断是否是当前版本
const isCurrentVersion = (ver: string) => {
  const routeVersion = route.params.version as string | undefined
  return ver === routeVersion
}

// 选择版本
const selectVersion = (ver: string) => {
  showVersionDropdown.value = false
  if (ver === '' || ver === library.value.default_version) {
    // 选择默认版本，不带 version 参数
    router.push(`/libraries/${libraryId.value}`)
  } else {
    router.push(`/libraries/${libraryId.value}/${ver}`)
  }
}

// 打开添加版本弹窗
const openAddVersionModal = () => {
  showVersionDropdown.value = false
  showAddVersionModal.value = true
}

// 处理版本创建成功
const handleVersionCreated = async (newVersion: string) => {
  console.log('✓ Version created:', newVersion)
  // 清空旧版本数据，确保新版本重新加载
  documents.value = []
  searchResult.value = ''
  hasSearched.value = false
  // 刷新库信息以获取新版本列表
  await fetchLibrary()
  // 跳转到新版本的 logs tab
  router.push(`/libraries/${libraryId.value}/${newVersion}?tab=logs`)
}

const handleSignIn = () => {
  redirectToSSO()
}

const fetchLibrary = async () => {
  const library_data = await getLibrary(libraryId.value)
  library.value = library_data
}

// 加载文档列表
const fetchDocumentsList = async () => {
  loadingDocs.value = true
  try {
    const res = await getDocuments({
      library_id: libraryId.value,
      version: version.value,
      page: page.value,
      page_size: pageSize.value
    })
    documents.value = res.list || []
    totalDocs.value = res.total
  } finally {
    loadingDocs.value = false
  }
}

// 上传文档（普通接口，后台异步处理，通过日志查看进度）
const handleFileUpload = async (event: Event) => {
  const input = event.target as HTMLInputElement
  const file = input.files?.[0]
  if (!file) return

  const allowedTypes = ['.md', '.pdf', '.docx']
  const ext = file.name.substring(file.name.lastIndexOf('.')).toLowerCase()
  if (!allowedTypes.includes(ext)) {
    ElMessage.warning('Only .md, .pdf, .docx formats are supported')
    return
  }

  // 使用当前版本或默认版本
  const uploadVersion = version.value || library.value.default_version || 'default'

  // 显示上传中状态（不显示进度条）
  uploading.value = true

  try {
    // 使用统一的 API 接口上传（后台异步处理，通过日志查看进度）
    await uploadDocument(libraryId.value, file, uploadVersion)
    
    uploading.value = false
    ElMessage.success('上传已启动，跳转到控制台查看进度')
    // 跳转到 logs tab
    activeTab.value = 'logs'
    startLogPolling()

    // ====== 以下是 SSE 版本的代码，保留备用 ======
    // await uploadDocumentWithSSE(
    //   libraryId.value,
    //   file,
    //   {
    //     onProgress: (status) => {
    //       const progressMap: Record<string, number> = {
    //         uploaded: 10,
    //         parsing: 30,
    //         chunking: 50,
    //         embedding: 70,
    //         saving: 90
    //       }
    //       uploadProgress.value = progressMap[status.stage] || status.progress || 0
    //       uploadMessage.value = status.message || status.stage
    //     },
    //     onComplete: () => {
    //       uploadProgress.value = 100
    //       uploadMessage.value = 'Upload successful!'
    //       setTimeout(() => {
    //         uploading.value = false
    //         uploadProgress.value = 0
    //         uploadMessage.value = ''
    //         fetchDocumentsList()
    //         fetchLibrary()
    //       }, 500)
    //     },
    //     onError: (error) => {
    //       const errorMsg = error.message || 'Unknown error'
    //       const status = (error as any).status
    //       const code = (error as any).code
    //       
    //       let displayMsg = errorMsg
    //       if (status) {
    //         displayMsg = `HTTP Error: ${displayMsg}`
    //       } else if (code !== undefined) {
    //         displayMsg = `Error (${code}): ${displayMsg}`
    //       }
    //       
    //       ElMessage.error('Upload failed: ' + displayMsg)
    //       uploading.value = false
    //       uploadProgress.value = 0
    //       uploadMessage.value = ''
    //     }
    //   },
    //   uploadVersion
    // )
    // ====== SSE 版本代码结束 ======
  } catch (error) {
    ElMessage.error('Upload failed: ' + (error instanceof Error ? error.message : 'Unknown error'))
    uploading.value = false
    uploadProgress.value = 0
    uploadMessage.value = ''
  }
  
  input.value = ''
}

const handleSearch = async () => {
  searching.value = true
  hasSearched.value = true
  
  const topic = searchTopic.value.trim()
  
  try {
    // 调用统一的 getChunks API，通过 topic 参数控制是否搜索
    const res = await getChunks(searchMode.value, libraryId.value, {
      version: version.value, // 传入版本参数
      topic: topic || undefined // 空字符串转为 undefined，返回全部文档
    })
    
    const chunks = (res.chunks || []) as any[]
    if (chunks.length > 0) {
      // 格式化为 ---- 分割的文本格式
      // code mode: title → source → description → code
      // info mode: title → source → description (content)
      const formatted = chunks.map((chunk: any) => {
        let text = ''
        if (chunk.title) text += `### ${chunk.title}\n\n`
        if (chunk.source) text += `Source: ${chunk.source}\n\n`
        if (chunk.description) text += `${chunk.description}\n\n`
        // code mode: 显示代码块（带语言标记）
        // info mode: 显示 chunk_text 原文
        if (searchMode.value === 'code') {
          if (chunk.code) {
            const lang = chunk.language || ''
            text += `\`\`\`${lang}\n${chunk.code}\n\`\`\``
          }
        } else {
          if (chunk.chunk_text) text += chunk.chunk_text
        }
        return text.trim()
      }).join('\n\n--------------------------------\n\n')
      searchResult.value = formatted
    } else {
      if (topic) {
        searchResult.value = `No results found for "${topic}".`
      } else {
        searchResult.value = 'No documents found.'
      }
    }
  } catch (error) {
    searchResult.value = 'Search failed. Please try again.'
    console.error('Search error:', error)
  } finally {
    searching.value = false
  }
}

const copyContent = () => {
  navigator.clipboard.writeText(searchResult.value)
}

const formatNumber = (num: number) => {
  if (num >= 1000000) return (num / 1000000).toFixed(1) + 'M'
  if (num >= 1000) return (num / 1000).toFixed(1) + 'K'
  return num.toString()
}

const formatDateShort = (dateStr: string) => {
  if (!dateStr) return '-'
  const date = new Date(dateStr)
  return date.toLocaleDateString('en-US', { month: 'short', day: 'numeric', year: 'numeric' })
}

const formatDate = (dateStr: string) => {
  if (!dateStr) return '-'
  const date = new Date(dateStr)
  const now = new Date()
  
  // 如果时间戳无效或是未来时间，显示 'now'
  if (isNaN(date.getTime()) || date > now) {
    return 'now'
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

// 检查描述是否需要截断（超过150个字符）
const isTruncated = (text: string | undefined) => {
  if (!text) return false
  return text.length > 70
}

// 获取截断的文本
const getTruncatedText = (text: string | undefined) => {
  if (!text) return 'No description'
  if (text.length > 70) {
    return text.substring(0, 70)
  }
  return text
}

onMounted(async () => {
  initUserState()
  await fetchLibrary()
  
  // 根据当前 tab 加载对应数据
  if (activeTab.value === 'context') {
    handleSearch()
  } else if (activeTab.value === 'documents') {
    fetchDocumentsList()
  } else if (activeTab.value === 'logs') {
    startLogPolling()
  }
})

// 组件卸载时清理定时器
onUnmounted(() => {
  stopLogPolling()
})

// 监听 activeTab 变化
watch(activeTab, (newTab) => {
  // 更新 URL query 参数（不刷新页面）
  router.replace({ query: { ...route.query, tab: newTab } })
  
  // 切换到 context 时，如果没有搜索结果则加载
  if (newTab === 'context' && !hasSearched.value) {
    handleSearch()
  }
  // 切换到 documents 时加载文档列表
  if (newTab === 'documents' && documents.value.length === 0) {
    fetchDocumentsList()
  }
  // 切换到 terminal 时开始轮询日志
  if (newTab === 'logs') {
    startLogPolling()
  } else {
    // 离开 terminal 时停止轮询
    stopLogPolling()
  }
})

// 监听 URL query.tab 变化（路由跳转时更新 activeTab）
watch(
  () => route.query.tab,
  (newTab) => {
    if (newTab && newTab !== activeTab.value) {
      activeTab.value = newTab as string
    }
  }
)

// 监听版本变化，刷新所有 tab 的数据
watch(
  () => route.params.version,
  () => {
    // 版本切换时，清空并重新加载各 tab 数据
    // 清空 documents 数据，强制下次切换时重新加载
    documents.value = []
    // 清空 context 搜索结果
    searchResult.value = 'Loading document...'
    hasSearched.value = false
    // 根据当前 tab 加载对应数据
    if (activeTab.value === 'context') {
      handleSearch()
    } else if (activeTab.value === 'documents') {
      fetchDocumentsList()
    } else if (activeTab.value === 'logs') {
      fetchLogs()
    }
  }
)

// 监听 searchMode 变化，重新搜索
watch(searchMode, () => {
  handleSearch()
})
</script>
