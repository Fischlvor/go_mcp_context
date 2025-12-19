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
                <!-- 刷新按钮 -->
                <button 
                  class="flex h-8 items-center justify-center gap-1.5 rounded-lg border border-stone-300 text-base text-stone-500 transition hover:border-stone-400 w-8"
                  @click="handleSearch"
                >
                  <svg xmlns="http://www.w3.org/2000/svg" width="24" height="24" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="1.5" stroke-linecap="round" stroke-linejoin="round" class="h-5 w-5 text-stone-500">
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
                            <span class="truncate">{{ doc.title }}</span>
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
    <div v-if="showAddVersionModal" class="fixed inset-0 z-50 flex items-center justify-center bg-black/50" @click.self="showAddVersionModal = false">
      <div class="w-full max-w-md rounded-lg bg-white p-6 shadow-lg">
        <div class="mb-4">
          <h2 class="text-lg font-semibold text-stone-800">Add New Version</h2>
          <p class="mt-1 text-sm text-stone-500">Create a new version for this library</p>
        </div>
        <div>
          <label class="block text-sm font-medium text-stone-700">Version Name</label>
          <input 
            v-model="newVersionName"
            type="text"
            placeholder="e.g. 1.0.0 or v1.0.0"
            class="mt-1 w-full rounded-lg border border-stone-300 px-3 py-2 text-sm focus:border-emerald-500 focus:outline-none focus:ring-1 focus:ring-emerald-500"
            @keyup.enter="handleAddVersion"
          />
          <p class="mt-1 text-xs text-stone-400">Version will be prefixed with 'v' if not already</p>
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

    <!-- 点击外部关闭下拉框 -->
    <div 
      v-if="showVersionDropdown" 
      class="fixed inset-0 z-40" 
      @click="showVersionDropdown = false"
    ></div>
  </div>
</template>

<script setup lang="ts">
import { ref, computed, onMounted, watch } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { ElMessage } from 'element-plus'
import AppHeader from '@/components/AppHeader.vue'
import AppFooter from '@/components/AppFooter.vue'
import { useUser } from '@/stores/user'
import { getLibrary, createVersion } from '@/api/library'
import { getDocuments, uploadDocumentWithSSE, getChunks } from '@/api/document'
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

const activeTab = ref('context')
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
const newVersionName = ref('')

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

// 创建新版本
const handleAddVersion = async () => {
  if (!newVersionName.value.trim()) {
    ElMessage.warning('Please enter a version name')
    return
  }

  // 自动添加 v 前缀
  let versionWithPrefix = newVersionName.value.trim()
  if (!versionWithPrefix.startsWith('v')) {
    versionWithPrefix = 'v' + versionWithPrefix
  }

  try {
    await createVersion(libraryId.value, versionWithPrefix)
    showAddVersionModal.value = false
    newVersionName.value = ''
    // 刷新库信息以获取新版本列表
    await fetchLibrary()
    // 跳转到新版本
    router.push(`/libraries/${libraryId.value}/${versionWithPrefix}`)
  } catch (error) {
    console.error('Failed to create version:', error)
  }
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

// 上传文档
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

  // 重置状态
  uploading.value = true
  uploadProgress.value = 0
  uploadMessage.value = 'Uploading...'

  try {
    await uploadDocumentWithSSE(
      libraryId.value,
      file,
      {
        onProgress: (status) => {
          const progressMap: Record<string, number> = {
            uploaded: 10,
            parsing: 30,
            chunking: 50,
            embedding: 70,
            saving: 90
          }
          uploadProgress.value = progressMap[status.stage] || status.progress || 0
          uploadMessage.value = status.message || status.stage
        },
        onComplete: () => {
          uploadProgress.value = 100
          uploadMessage.value = 'Upload successful!'
          setTimeout(() => {
            uploading.value = false
            uploadProgress.value = 0
            uploadMessage.value = ''
            fetchDocumentsList()
            fetchLibrary()
          }, 500)
        },
        onError: (error) => {
          const errorMsg = error.message || 'Unknown error'
          const status = (error as any).status
          const code = (error as any).code
          
          // 区分错误类型
          let displayMsg = errorMsg
          if (status) {
            // HTTP 错误
            displayMsg = `HTTP Error: ${displayMsg}`
          } else if (code !== undefined) {
            // 业务错误
            displayMsg = `Error (${code}): ${displayMsg}`
          }
          
          ElMessage.error('Upload failed: ' + displayMsg)
          uploading.value = false
          uploadProgress.value = 0
          uploadMessage.value = ''
        }
      },
      uploadVersion
    )
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
  
  try {
    // 调用统一的 getChunks API，通过 topic 参数控制是否搜索
    const res = await getChunks(searchMode.value, libraryId.value, {
      version: version.value, // 传入版本参数
      topic: searchTopic.value.trim() || undefined // 空字符串转为 undefined，返回全部文档
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
      if (searchTopic.value.trim()) {
        searchResult.value = `No results found for "${searchTopic.value}".`
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
  handleSearch()
})

// 监听 activeTab 变化，切换到 documents 时加载文档列表
watch(activeTab, (newTab) => {
  if (newTab === 'documents' && documents.value.length === 0) {
    fetchDocumentsList()
  }
})

// 监听版本和 searchMode 变化，统一调用 handleSearch
// handleSearch 会根据 searchTopic 是否有值决定是搜索还是获取全部文档
watch(
  [() => route.params.version, searchMode],
  () => {
    handleSearch()
    // 如果在 documents tab，也重新加载文档列表
    if (activeTab.value === 'documents') {
      fetchDocumentsList()
    }
  }
)
</script>
