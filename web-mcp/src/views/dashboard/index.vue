<template>
  <div class="flex min-h-screen flex-col overflow-x-hidden bg-stone-50 antialiased">
    <!-- 顶部渐变背景 -->
    <div class="absolute inset-x-0 top-0 h-[260px] bg-gradient-to-b from-emerald-500/[0.15] to-transparent"></div>
    <!-- 底层背景色 -->
    <div class="fixed inset-0 -z-10 bg-stone-50"></div>

    <!-- 顶部 Header -->
    <AppHeader 
      :is-logged-in="isLoggedIn" 
      :user-email="userEmail" 
      :user-plan="userPlan"
      @add-docs="showAddDocsModal = true"
    />

    <!-- Add Docs 弹窗 -->
    <AddDocsModal v-model:visible="showAddDocsModal" />

    <!-- 主内容区 -->
    <main class="flex-grow pt-0">
      <main class="-mt-10 flex flex-col px-4 pt-10 sm:px-6 md:-mt-20 md:pt-20">
        <!-- 顶部 Tabs -->
        <div class="mt-3">
          <div class="relative mx-auto flex w-full max-w-[880px] justify-center">
            <div class="relative flex">
              <div class="absolute bottom-0 left-4 right-4 h-px bg-stone-200"></div>
              <button 
                v-for="tab in tabs" 
                :key="tab.id"
                :class="['relative px-4 py-3 text-base font-normal transition-colors duration-200', activeTab === tab.id ? 'text-emerald-600' : 'text-stone-800 hover:text-stone-600']"
                @click="activeTab = tab.id"
              >
                {{ tab.label }}
                <div v-if="activeTab === tab.id" class="absolute bottom-0 left-4 right-4 h-0.5 bg-emerald-600"></div>
              </button>
            </div>
          </div>
        </div>

        <!-- 内容区域 -->
        <div class="mt-6 flex flex-col gap-6 sm:mt-8 sm:gap-8">
          <!-- Overview Tab -->
          <template v-if="activeTab === 'overview'">
            <!-- 统计卡片 -->
            <div class="mx-auto w-full max-w-[880px]">
              <div class="rounded-3xl border border-stone-200 bg-white px-6 py-6 shadow-sm sm:px-8">
                <div class="grid grid-cols-1 gap-4 sm:grid-cols-4 sm:gap-8">
                  <div v-for="(stat, index) in stats" :key="stat.label" 
                    :class="['flex items-center justify-between sm:flex-col sm:items-start sm:px-4 sm:pb-0', 
                      index < stats.length - 1 ? 'border-b border-stone-200 pb-4 sm:border-b-0 sm:border-r' : '']">
                    <div class="flex items-center gap-1.5 text-left text-sm font-normal uppercase text-stone-500">{{ stat.label }}</div>
                    <div class="text-left text-base font-medium text-stone-800 sm:text-lg">{{ stat.value }}</div>
                  </div>
                </div>
              </div>
            </div>

            <!-- Connect 卡片 -->
            <div class="mx-auto w-full max-w-[880px]">
              <div class="rounded-3xl border border-stone-200 bg-white p-6 shadow-sm sm:p-8 lg:p-10">
                <div class="mb-6">
                  <div class="min-w-0 flex-1">
                    <h2 class="text-lg font-semibold text-stone-900 sm:text-xl">Connect</h2>
                    <div class="text-sm text-stone-500 sm:text-base">
                      <a href="#" class="underline transition-colors hover:text-stone-700">Check the docs</a> for installation
                    </div>
                  </div>
                </div>
                <div class="space-y-6">
                  <!-- MCP URL / API URL -->
                  <div class="rounded-xl bg-stone-100 px-4 py-2 sm:px-5">
                    <div class="flex flex-col gap-2 py-3 sm:grid sm:grid-cols-[80px_auto_1fr] sm:items-center sm:py-2">
                      <span class="text-sm font-normal uppercase text-stone-500">MCP URL</span>
                      <span class="hidden text-sm text-stone-500 sm:block">:</span>
                      <div class="flex items-center gap-2">
                        <span class="text-base font-medium text-stone-800">mcp.hsk423.cn/mcp</span>
                        <button class="text-stone-400 transition-colors hover:text-stone-600" @click="copyToClipboard('https://mcp.hsk423.cn/mcp')">
                          <svg xmlns="http://www.w3.org/2000/svg" width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round">
                            <path d="M7 7m0 2.667a2.667 2.667 0 0 1 2.667 -2.667h8.666a2.667 2.667 0 0 1 2.667 2.667v8.666a2.667 2.667 0 0 1 -2.667 2.667h-8.666a2.667 2.667 0 0 1 -2.667 -2.667z"></path>
                            <path d="M4.012 16.737a2.005 2.005 0 0 1 -1.012 -1.737v-10c0 -1.1 .9 -2 2 -2h10c.75 0 1.158 .385 1.5 1"></path>
                          </svg>
                        </button>
                      </div>
                    </div>
                    <!-- API URL 暂时注释
                    <div class="border-t border-stone-200"></div>
                    <div class="flex flex-col gap-2 py-3 sm:grid sm:grid-cols-[80px_auto_1fr] sm:items-center sm:py-2">
                      <span class="text-sm font-normal uppercase text-stone-500">API URL</span>
                      <span class="hidden text-sm text-stone-500 sm:block">:</span>
                      <div class="flex items-center gap-2">
                        <span class="text-base font-medium text-stone-800">mcp.hsk423.cn/api/v1</span>
                        <button class="text-stone-400 transition-colors hover:text-stone-600" @click="copyToClipboard('https://mcp.hsk423.cn/api/v1')">
                          <svg xmlns="http://www.w3.org/2000/svg" width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round">
                            <path d="M7 7m0 2.667a2.667 2.667 0 0 1 2.667 -2.667h8.666a2.667 2.667 0 0 1 2.667 2.667v8.666a2.667 2.667 0 0 1 -2.667 2.667h-8.666a2.667 2.667 0 0 1 -2.667 -2.667z"></path>
                            <path d="M4.012 16.737a2.005 2.005 0 0 1 -1.012 -1.737v-10c0 -1.1 .9 -2 2 -2h10c.75 0 1.158 .385 1.5 1"></path>
                          </svg>
                        </button>
                      </div>
                    </div>
                    -->
                  </div>

                  <!-- IDE Tabs -->
                  <div class="space-y-4">
                    <div class="-mx-4 overflow-x-auto px-4 sm:-mx-5 sm:px-5">
                      <div class="min-w-max">
                        <div class="relative flex w-full items-end gap-0">
                          <button 
                            v-for="ide in ides" 
                            :key="ide.id"
                            :class="['flex items-center font-medium gap-1.5 px-2.5 py-1.5 text-sm', 
                              activeIde === ide.id ? 'rounded-t-lg border border-stone-300 border-b-transparent text-stone-800' : 'border border-stone-300 border-l-transparent border-r-transparent border-t-transparent text-stone-500 hover:text-stone-600']"
                            @click="activeIde = ide.id"
                          >
                            <svg width="14" height="14" :viewBox="ide.viewBox" :fill="activeIde === ide.id ? ide.color : 'currentColor'" xmlns="http://www.w3.org/2000/svg" class="h-3.5 w-3.5">
                              <path :d="ide.path" />
                            </svg>
                            {{ ide.label }}
                          </button>
                          <div class="flex-grow border-b border-stone-300"></div>
                        </div>
                      </div>
                    </div>

                    <!-- Code Block -->
                    <div class="relative rounded-lg bg-stone-100 p-4">
                      <button class="absolute right-3 top-3 text-stone-400 transition-colors hover:text-stone-600" @click="copyCode">
                        <svg xmlns="http://www.w3.org/2000/svg" width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round">
                          <path d="M7 7m0 2.667a2.667 2.667 0 0 1 2.667 -2.667h8.666a2.667 2.667 0 0 1 2.667 2.667v8.666a2.667 2.667 0 0 1 -2.667 2.667h-8.666a2.667 2.667 0 0 1 -2.667 -2.667z"></path>
                          <path d="M4.012 16.737a2.005 2.005 0 0 1 -1.012 -1.737v-10c0 -1.1 .9 -2 2 -2h10c.75 0 1.158 .385 1.5 1"></path>
                        </svg>
                      </button>
                      <pre class="text-sm leading-relaxed"><code ref="codeBlock" class="rounded"></code></pre>
                    </div>
                  </div>
                </div>
              </div>
            </div>

            <!-- API Keys 卡片 -->
            <div class="mx-auto w-full max-w-[880px]">
              <div class="rounded-3xl border border-stone-200 bg-white p-6 shadow-sm sm:p-8 lg:p-10">
                <div class="mb-6 flex flex-col flex-nowrap items-start gap-4 sm:flex-row sm:justify-between">
                  <div class="min-w-0 flex-1">
                    <h2 class="text-lg font-semibold text-stone-900 sm:text-xl">API Keys</h2>
                    <div class="text-sm text-stone-500 sm:text-base">
                      <p>Manage your API keys to authenticate MCP requests</p>
                    </div>
                  </div>
                  <div class="flex-shrink-0">
                    <button 
                      class="flex items-center justify-center gap-2 whitespace-nowrap rounded-md border border-emerald-300 bg-emerald-50 px-3 py-2 text-sm font-normal leading-none text-emerald-800 transition-colors hover:bg-emerald-100 disabled:opacity-50"
                      :disabled="!isLoggedIn || apiKeys.length >= 5"
                      @click="showCreateDialog = true"
                    >
                      Create API Key...
                    </button>
                  </div>
                </div>
                <div class="space-y-4">
                  <!-- 未登录提示 -->
                  <div v-if="!isLoggedIn" class="rounded-xl border border-amber-300 bg-amber-50 p-5 text-amber-800">
                    <div class="flex items-start gap-3">
                      <div class="hidden sm:block">
                        <svg xmlns="http://www.w3.org/2000/svg" width="24" height="24" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="1.5" stroke-linecap="round" stroke-linejoin="round">
                          <path d="M3 12a9 9 0 1 0 18 0a9 9 0 0 0 -18 0"></path>
                          <path d="M12 8v4"></path>
                          <path d="M12 16h.01"></path>
                        </svg>
                      </div>
                      <div class="text-base font-normal">
                        Please login to manage your API keys.
                      </div>
                    </div>
                  </div>
                  <!-- 加载中 -->
                  <div v-else-if="apiKeysLoading" class="py-8 text-center text-stone-500">
                    Loading...
                  </div>
                  <!-- 空状态 -->
                  <div v-else-if="apiKeys.length === 0" class="rounded-xl border border-blue-300 bg-blue-50 p-5 text-blue-800">
                    <div class="flex items-start gap-3">
                      <div class="hidden sm:block">
                        <svg xmlns="http://www.w3.org/2000/svg" width="24" height="24" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="1.5" stroke-linecap="round" stroke-linejoin="round" class="text-blue-800">
                          <path d="M3 12a9 9 0 1 0 18 0a9 9 0 0 0 -18 0"></path>
                          <path d="M12 8v4"></path>
                          <path d="M12 16h.01"></path>
                        </svg>
                      </div>
                      <div class="flex flex-col">
                        <div class="text-base font-normal">
                          No API keys yet. <button class="font-semibold underline hover:text-blue-900" @click="showCreateDialog = true">Click here to generate your first API key.</button>
                        </div>
                      </div>
                    </div>
                  </div>
                  <!-- API Keys 列表（表格形式） -->
                  <div v-else class="w-full overflow-x-auto md:overflow-x-visible">
                    <table class="w-full min-w-[600px] table-fixed border-b border-stone-200">
                      <thead class="border-b border-stone-200">
                        <tr>
                          <th class="w-[170px] px-2 py-3 text-left text-sm font-normal uppercase leading-none text-stone-400 sm:px-4">NAME</th>
                          <th class="w-[160px] px-2 py-3 text-left text-sm font-normal uppercase leading-none text-stone-400 sm:px-4">KEY</th>
                          <th class="w-[140px] px-2 py-3 text-left text-sm font-normal uppercase leading-none text-stone-400 sm:px-4">CREATED</th>
                          <th class="w-[140px] px-2 py-3 text-left text-sm font-normal uppercase leading-none text-stone-400 sm:px-4">LAST USED</th>
                          <th class="w-[30px] px-1 py-3 text-center text-sm font-normal uppercase leading-none text-stone-400"></th>
                        </tr>
                      </thead>
                      <tbody class="divide-y divide-stone-200">
                        <tr v-for="key in apiKeys" :key="key.id" class="group transition-colors hover:bg-white">
                          <td class="h-11 truncate px-2 align-middle text-base font-normal leading-tight text-stone-800 sm:px-4">{{ key.name }}</td>
                          <td class="h-11 px-2 align-middle text-base font-normal slashed-zero tabular-nums leading-tight text-stone-800 sm:px-4">
                            <code class="rounded bg-stone-100 px-2 py-1 text-xs">mcpsk-****{{ key.token_suffix }}</code>
                          </td>
                          <td class="h-11 px-2 align-middle text-base font-normal slashed-zero tabular-nums leading-tight text-stone-800 sm:px-4">{{ formatDate(key.created_at) }}</td>
                          <td class="h-11 px-2 align-middle text-base font-normal slashed-zero tabular-nums leading-tight text-stone-800 sm:px-4">{{ formatLastUsed(key.last_used_at) }}</td>
                          <td class="h-11 px-1 text-center align-middle">
                            <button 
                              type="button" 
                              aria-label="Revoke" 
                              class="flex items-center justify-center text-stone-500 transition-colors hover:text-red-600"
                              @click="handleDeleteKey(key.id)"
                            >
                              <svg xmlns="http://www.w3.org/2000/svg" width="20" height="20" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="1.5" stroke-linecap="round" stroke-linejoin="round">
                                <path d="M4 7h16"></path>
                                <path d="M5 7l1 12a2 2 0 0 0 2 2h8a2 2 0 0 0 2 -2l1 -12"></path>
                                <path d="M9 7v-3a1 1 0 0 1 1 -1h4a1 1 0 0 1 1 1v3"></path>
                                <path d="M10 12l4 4m0 -4l-4 4"></path>
                              </svg>
                            </button>
                          </td>
                        </tr>
                      </tbody>
                    </table>
                    <div v-if="apiKeys.length >= 5" class="mt-4 text-center text-sm text-stone-500">
                      Maximum 5 API keys allowed
                    </div>
                  </div>
                </div>
              </div>
            </div>

            <!-- 创建 API Key 弹窗 -->
            <div v-if="showCreateDialog" class="fixed inset-0 z-50 flex items-center justify-center bg-black/50" @click.self="showCreateDialog = false">
              <div class="w-full max-w-md rounded-2xl bg-white p-6 shadow-xl">
                <h3 class="mb-4 text-lg font-semibold text-stone-900">Create API Key</h3>
                <div class="mb-4">
                  <label class="mb-2 block text-sm font-medium text-stone-700">Name</label>
                  <input 
                    v-model="newKeyName"
                    type="text"
                    placeholder="e.g., Development, Production"
                    class="w-full rounded-lg border border-stone-300 px-3 py-2 text-stone-800 focus:border-emerald-500 focus:outline-none focus:ring-1 focus:ring-emerald-500"
                    maxlength="100"
                    @keyup.enter="handleCreateKey"
                  />
                </div>
                <div class="flex justify-end gap-3">
                  <button 
                    class="rounded-md px-4 py-2 text-sm text-stone-600 transition-colors hover:bg-stone-100"
                    @click="showCreateDialog = false"
                  >
                    Cancel
                  </button>
                  <button 
                    class="rounded-md bg-emerald-600 px-4 py-2 text-sm text-white transition-colors hover:bg-emerald-700 disabled:opacity-50"
                    :disabled="!newKeyName.trim() || creatingKey"
                    @click="handleCreateKey"
                  >
                    {{ creatingKey ? 'Creating...' : 'Create' }}
                  </button>
                </div>
              </div>
            </div>

            <!-- 新创建的 Key 显示弹窗 -->
            <div v-if="newlyCreatedKey" class="fixed inset-0 z-50 flex items-center justify-center bg-black/50">
              <div class="w-full max-w-lg rounded-2xl bg-white p-6 shadow-xl">
                <div class="mb-4 flex items-center gap-2 text-emerald-600">
                  <svg xmlns="http://www.w3.org/2000/svg" width="24" height="24" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round">
                    <path d="M12 12m-9 0a9 9 0 1 0 18 0a9 9 0 1 0 -18 0"></path>
                    <path d="M9 12l2 2l4 -4"></path>
                  </svg>
                  <h3 class="text-lg font-semibold">API Key Created</h3>
                </div>
                <div class="mb-4 rounded-lg border border-amber-300 bg-amber-50 p-4 text-sm text-amber-800">
                  <strong>Important:</strong> Copy your API key now. You won't be able to see it again!
                </div>
                <div class="mb-4">
                  <label class="mb-2 block text-sm font-medium text-stone-700">Your API Key</label>
                  <div class="flex items-center gap-2">
                    <code class="flex-1 rounded-lg bg-stone-100 px-3 py-2 font-mono text-sm text-stone-800 break-all">
                      {{ newlyCreatedKey.api_key }}
                    </code>
                    <button 
                      class="rounded-md bg-stone-200 px-3 py-2 text-sm text-stone-700 transition-colors hover:bg-stone-300"
                      @click="copyNewKey"
                    >
                      Copy
                    </button>
                  </div>
                </div>
                <div class="flex justify-end">
                  <button 
                    class="rounded-md bg-emerald-600 px-4 py-2 text-sm text-white transition-colors hover:bg-emerald-700"
                    @click="closeNewKeyDialog"
                  >
                    Done
                  </button>
                </div>
              </div>
            </div>

            <!-- API 卡片 -->
            <div class="mx-auto w-full max-w-[880px]">
              <div class="rounded-3xl border border-stone-200 bg-white p-6 shadow-sm sm:p-8 lg:p-10">
                <div class="mb-6">
                  <div class="min-w-0 flex-1">
                    <h2 class="text-lg font-semibold text-stone-900 sm:text-xl">MCP</h2>
                    <div class="text-sm text-stone-500 sm:text-base">
                      <p>Use MCP protocol to search libraries and fetch code snippets for AI IDE integration</p>
                    </div>
                  </div>
                </div>
                <div class="space-y-6">
                  <!-- Search/Docs Toggle -->
                  <div class="flex" role="group">
                    <button 
                      :class="['border px-3 py-1 text-sm font-medium shadow-sm transition-colors rounded-l-md', 
                        apiTab === 'search' ? 'border-stone-500 bg-stone-700 text-white' : 'border-stone-300 bg-white text-stone-800 hover:bg-stone-50']"
                      @click="apiTab = 'search'"
                    >Search</button>
                    <button 
                      :class="['border px-3 py-1 text-sm font-medium shadow-sm transition-colors rounded-r-md border-l-0', 
                        apiTab === 'docs' ? 'border-stone-500 bg-stone-700 text-white' : 'border-stone-300 bg-white text-stone-800 hover:bg-stone-50']"
                      @click="apiTab = 'docs'"
                    >Docs</button>
                  </div>

                  <!-- Docs Type Toggle (仅在 Docs tab 显示) -->
                  <div v-if="apiTab === 'docs'">
                    <h4 class="mb-2 text-sm font-medium text-stone-700">Docs Type</h4>
                    <div class="flex" role="group">
                      <button 
                        :class="['border px-3 py-1 text-sm font-medium shadow-sm transition-colors rounded-l-md', 
                          docsType === 'code' ? 'border-stone-500 bg-stone-700 text-white' : 'border-stone-300 bg-white text-stone-800 hover:bg-stone-50']"
                        @click="docsType = 'code'"
                      >Code</button>
                      <button 
                        :class="['border px-3 py-1 text-sm font-medium shadow-sm transition-colors rounded-r-md border-l-0', 
                          docsType === 'info' ? 'border-stone-500 bg-stone-700 text-white' : 'border-stone-300 bg-white text-stone-800 hover:bg-stone-50']"
                        @click="docsType = 'info'"
                      >Info</button>
                    </div>
                  </div>

                  <!-- API Code Block -->
                  <div class="relative rounded-lg bg-stone-100 p-4">
                    <button class="absolute right-3 top-3 text-stone-400 transition-colors hover:text-stone-600" @click="copyToClipboard(apiCommand)" aria-label="Copy code">
                      <svg xmlns="http://www.w3.org/2000/svg" width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round">
                        <path d="M7 7m0 2.667a2.667 2.667 0 0 1 2.667 -2.667h8.666a2.667 2.667 0 0 1 2.667 2.667v8.666a2.667 2.667 0 0 1 -2.667 2.667h-8.666a2.667 2.667 0 0 1 -2.667 -2.667z"></path>
                        <path d="M4.012 16.737a2.005 2.005 0 0 1 -1.012 -1.737v-10c0 -1.1 .9 -2 2 -2h10c.75 0 1.158 .385 1.5 1"></path>
                      </svg>
                    </button>
                    <pre style="display: block; overflow-x: auto; padding: 0px; color: rgb(56, 58, 66); background: transparent; margin: 0px; font-size: 14px; line-height: 1.5;"><code ref="apiCommandBlock" class="language-bash" style="white-space: pre;"></code></pre>
                  </div>

                  <!-- Parameters -->
                  <div class="space-y-4">
                    <div v-if="apiTab === 'search'">
                      <h4 class="mb-2 text-sm font-medium text-stone-700">Parameters</h4>
                      <p class="text-sm text-stone-500">
                        <code class="rounded bg-stone-100 px-1 text-xs">libraryName</code> - Search term for finding libraries
                      </p>
                    </div>
                    <div v-else>
                      <h4 class="mb-2 text-sm font-medium text-stone-700">Parameters</h4>
                      <div class="grid grid-cols-2 gap-x-4 gap-y-2 text-sm text-stone-500">
                        <div><code class="rounded bg-stone-100 px-1 text-xs">libraryId</code> - Library database ID</div>
                        <div><code class="rounded bg-stone-100 px-1 text-xs">version</code> - Library version (optional)</div>
                        <div><code class="rounded bg-stone-100 px-1 text-xs">topic</code> - Search by topic</div>
                        <div><code class="rounded bg-stone-100 px-1 text-xs">mode</code> - Documentation mode (code/info)</div>
                        <div><code class="rounded bg-stone-100 px-1 text-xs">page</code> - Page number (1-10)</div>
                      </div>
                    </div>
                    <div>
                      <h4 class="mb-2 text-sm font-medium text-stone-700">Response</h4>
                      <div class="relative rounded-lg bg-stone-100 p-4">
                        <button class="absolute right-3 top-3 text-stone-400 transition-colors hover:text-stone-600" @click="copyToClipboard(apiResponse)" aria-label="Copy code">
                          <svg xmlns="http://www.w3.org/2000/svg" width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round">
                            <path d="M7 7m0 2.667a2.667 2.667 0 0 1 2.667 -2.667h8.666a2.667 2.667 0 0 1 2.667 2.667v8.666a2.667 2.667 0 0 1 -2.667 2.667h-8.666a2.667 2.667 0 0 1 -2.667 -2.667z"></path>
                            <path d="M4.012 16.737a2.005 2.005 0 0 1 -1.012 -1.737v-10c0 -1.1 .9 -2 2 -2h10c.75 0 1.158 .385 1.5 1"></path>
                          </svg>
                        </button>
                        <pre style="display: block; overflow-x: auto; padding: 0px; color: rgb(56, 58, 66); background: transparent; margin: 0px; font-size: 14px; line-height: 1.5;"><code ref="apiResponseBlock" class="language-json" style="white-space: pre;"></code></pre>
                      </div>
                    </div>
                  </div>
                </div>
              </div>
            </div>
          </template>

          <!-- Libraries Tab -->
          <template v-else-if="activeTab === 'libraries'">
            <div class="mx-auto w-full max-w-[880px]">
              <div class="rounded-3xl border border-stone-200 bg-white p-6 shadow-sm sm:p-8 lg:p-10">
                <h2 class="text-lg font-semibold text-stone-900 sm:text-xl">Libraries</h2>
                <p class="mt-2 text-sm text-stone-500">Manage your private libraries here.</p>
              </div>
            </div>
          </template>

          <!-- Members Tab -->
          <template v-else-if="activeTab === 'members'">
            <div class="mx-auto w-full max-w-[880px]">
              <div class="rounded-3xl border border-stone-200 bg-white p-6 shadow-sm sm:p-8 lg:p-10">
                <h2 class="text-lg font-semibold text-stone-900 sm:text-xl">Members</h2>
                <p class="mt-2 text-sm text-stone-500">Add new members to the team or change their permissions.</p>
              </div>
            </div>
          </template>

          <!-- Rules Tab -->
          <template v-else-if="activeTab === 'rules'">
            <div class="mx-auto w-full max-w-[880px]">
              <div class="rounded-3xl border border-stone-200 bg-white p-6 shadow-sm sm:p-8 lg:p-10">
                <h2 class="text-lg font-semibold text-stone-900 sm:text-xl">Team Rules</h2>
                <p class="mt-2 text-sm text-stone-500">Add rules that will be included in the context when you fetch library documentation.</p>
              </div>
            </div>
          </template>
        </div>
      </main>
    </main>

    <!-- Footer -->
    <AppFooter />

    <!-- Report Issue Button -->
    <div class="fixed bottom-6 right-6 z-50">
      <a target="_blank" class="flex min-h-[50px] min-w-[50px] items-center justify-center gap-2 rounded-[50px] bg-stone-800 px-4 py-2.5 shadow-xl transition-all hover:bg-stone-700" href="#">
        <svg xmlns="http://www.w3.org/2000/svg" width="24" height="24" viewBox="0 0 24 24" fill="currentColor" stroke="none" class="h-5 w-5 text-white">
          <path d="M12 4a4 4 0 0 1 3.995 3.8l.005 .2a1 1 0 0 1 .428 .096l3.033 -1.938a1 1 0 1 1 1.078 1.684l-3.015 1.931a7.17 7.17 0 0 1 .476 2.227h3a1 1 0 0 1 0 2h-3v1a6.01 6.01 0 0 1 -.195 1.525l2.708 1.616a1 1 0 1 1 -1.026 1.718l-2.514 -1.501a6.002 6.002 0 0 1 -3.973 2.56v-5.918a1 1 0 0 0 -2 0v5.917a6.002 6.002 0 0 1 -3.973 -2.56l-2.514 1.503a1 1 0 1 1 -1.026 -1.718l2.708 -1.616a6.01 6.01 0 0 1 -.195 -1.526v-1h-3a1 1 0 0 1 0 -2h3.001v-.055a7 7 0 0 1 .474 -2.173l-3.014 -1.93a1 1 0 1 1 1.078 -1.684l3.032 1.939l.024 -.012l.068 -.027l.019 -.005l.016 -.006l.032 -.008l.04 -.013l.034 -.007l.034 -.004l.045 -.008l.015 -.001l.015 -.002l.087 -.004a4 4 0 0 1 4 -4zm0 2a2 2 0 0 0 -2 2h4a2 2 0 0 0 -2 -2z"></path>
        </svg>
        <span class="hidden text-base text-white sm:inline">Report Issue</span>
      </a>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, computed, h, onMounted, watch, nextTick } from 'vue'
import hljs from 'highlight.js/lib/core'
import json from 'highlight.js/lib/languages/json'
import bash from 'highlight.js/lib/languages/bash'
import ini from 'highlight.js/lib/languages/ini'
import 'highlight.js/styles/github.css'

// 注册语言
hljs.registerLanguage('json', json)
hljs.registerLanguage('bash', bash)
hljs.registerLanguage('toml', ini) // TOML 使用 ini 高亮
import { ElMessage } from 'element-plus'
import AppHeader from '@/components/AppHeader.vue'
import AppFooter from '@/components/AppFooter.vue'
import AddDocsModal from '@/components/AddDocsModal.vue'
import { useUser } from '@/stores/user'
import { getAPIKeys, createAPIKey, deleteAPIKey, type APIKey, type APIKeyCreateResponse } from '@/api/apikey'
import { getMyStats } from '@/api/library'

// 用户状态
const { isLoggedIn, userEmail, userPlan, initUserState } = useUser()

// Add Docs 弹窗状态
const showAddDocsModal = ref(false)

// API Keys 状态
const apiKeys = ref<APIKey[]>([])
const apiKeysLoading = ref(false)
const showCreateDialog = ref(false)
const newKeyName = ref('')
const creatingKey = ref(false)
const newlyCreatedKey = ref<APIKeyCreateResponse | null>(null)

// 加载 API Keys
const loadAPIKeys = async () => {
  if (!isLoggedIn.value) return
  apiKeysLoading.value = true
  try {
    const res = await getAPIKeys()
    apiKeys.value = res || []
  } finally {
    apiKeysLoading.value = false
  }
}

// 创建 API Key
const handleCreateKey = async () => {
  if (!newKeyName.value.trim()) return
  creatingKey.value = true
  try {
    const res = await createAPIKey({ name: newKeyName.value.trim() })
    newlyCreatedKey.value = res
    newKeyName.value = ''
    showCreateDialog.value = false
    await loadAPIKeys()
  } catch (e) {
    console.error('Failed to create API key:', e)
  } finally {
    creatingKey.value = false
  }
}

// 删除 API Key
const handleDeleteKey = async (id: number) => {
  if (!confirm('确定要删除这个 API Key 吗？')) return
  try {
    await deleteAPIKey(id)
    await loadAPIKeys()
  } catch (e) {
    console.error('Failed to delete API key:', e)
  }
}

// 复制新创建的 Key
const copyNewKey = async () => {
  if (!newlyCreatedKey.value) return
  
  const text = newlyCreatedKey.value.api_key
  try {
    // 优先使用 Clipboard API
    if (navigator.clipboard && window.isSecureContext) {
      await navigator.clipboard.writeText(text)
    } else {
      // Fallback: 使用 execCommand（兼容 HTTP 环境）
      const textarea = document.createElement('textarea')
      textarea.value = text
      textarea.style.position = 'fixed'
      textarea.style.opacity = '0'
      document.body.appendChild(textarea)
      textarea.select()
      document.execCommand('copy')
      document.body.removeChild(textarea)
    }
    ElMessage.success('已复制到剪贴板')
  } catch (err) {
    console.error('复制失败:', err)
    ElMessage.error('复制失败，请手动复制')
  }
}

// 关闭新 Key 弹窗
const closeNewKeyDialog = () => {
  newlyCreatedKey.value = null
}

// 格式化日期（如 Dec 4, 2025）
const formatDate = (time: string) => {
  return new Date(time).toLocaleDateString('en-US', { 
    month: 'short', 
    day: 'numeric', 
    year: 'numeric' 
  })
}

// 格式化最后使用时间
const formatLastUsed = (time: string | null) => {
  if (!time) return 'Never'
  return formatDate(time)
}

// 格式化时间（保留兼容）
const formatTime = (time: string | null) => {
  if (!time) return '从未使用'
  return new Date(time).toLocaleString('zh-CN')
}

// 格式化数字（如 1.5M）
const formatNumber = (num: number): string => {
  if (num >= 1000000) return (num / 1000000).toFixed(1) + 'M'
  if (num >= 1000) return (num / 1000).toFixed(1) + 'K'
  return num.toString()
}

// Tabs
const tabs = [
  { id: 'overview', label: 'Overview' },
  { id: 'libraries', label: 'Libraries' },
  { id: 'members', label: 'Members' },
  { id: 'rules', label: 'Rules' }
]
const activeTab = ref('overview')

// Stats
const stats = ref([
  { label: 'Libraries', value: '0' },
  { label: 'Documents', value: '0' },
  { label: 'Tokens', value: '0' },
  { label: 'MCP Calls', value: '0' }
])

// 加载统计数据
const loadStats = async () => {
  try {
    const res = await getMyStats()
    if (res) {
      stats.value = [
        { label: 'Libraries', value: res.libraries.toString() },
        { label: 'Documents', value: res.documents.toString() },
        { label: 'Tokens', value: formatNumber(res.tokens) },
        { label: 'MCP Calls', value: res.mcp_calls.toString() }
      ]
    }
  } catch (e) {
    console.error('Failed to load stats:', e)
  }
}

// IDE 配置
const ides = [
  { 
    id: 'cursor', 
    label: 'Cursor', 
    color: '#000000',
    viewBox: '0 0 24 24',
    path: 'M11.503.131 1.891 5.678a.84.84 0 0 0-.42.726v11.188c0 .3.162.575.42.724l9.609 5.55a1 1 0 0 0 .998 0l9.61-5.55a.84.84 0 0 0 .42-.724V6.404a.84.84 0 0 0-.42-.726L12.497.131a1.01 1.01 0 0 0-.996 0M2.657 6.338h18.55c.263 0 .43.287.297.515L12.23 22.918c-.062.107-.229.064-.229-.06V12.335a.59.59 0 0 0-.295-.51l-9.11-5.257c-.109-.063-.064-.23.061-.23'
  },
  { 
    id: 'claude', 
    label: 'Claude Code', 
    color: '#D97757',
    viewBox: '0 0 16 16',
    path: 'm3.127 10.604 3.135-1.76.053-.153-.053-.085H6.11l-.525-.032-1.791-.048-1.554-.065-1.505-.08-.38-.081L0 7.832l.036-.234.32-.214.455.04 1.009.069 1.513.105 1.097.064 1.626.17h.259l.036-.105-.089-.065-.068-.064-1.566-1.062-1.695-1.121-.887-.646-.48-.327-.243-.306-.104-.67.435-.48.585.04.15.04.593.456 1.267.981 1.654 1.218.242.202.097-.068.012-.049-.109-.181-.9-1.626-.96-1.655-.428-.686-.113-.411a2 2 0 0 1-.068-.484l.496-.674L4.446 0l.662.089.279.242.411.94.666 1.48 1.033 2.014.302.597.162.553.06.17h.105v-.097l.085-1.134.157-1.392.154-1.792.052-.504.25-.605.497-.327.387.186.319.456-.045.294-.19 1.23-.37 1.93-.243 1.29h.142l.161-.16.654-.868 1.097-1.372.484-.545.565-.601.363-.287h.686l.505.751-.226.775-.707.895-.585.759-.839 1.13-.524.904.048.072.125-.012 1.897-.403 1.024-.186 1.223-.21.553.258.06.263-.218.536-1.307.323-1.533.307-2.284.54-.028.02.032.04 1.029.098.44.024h1.077l2.005.15.525.346.315.424-.053.323-.807.411-3.631-.863-.872-.218h-.12v.073l.726.71 1.331 1.202 1.667 1.55.084.383-.214.302-.226-.032-1.464-1.101-.565-.497-1.28-1.077h-.084v.113l.295.432 1.557 2.34.08.718-.112.234-.404.141-.444-.08-.911-1.28-.94-1.44-.759-1.291-.093.053-.448 4.821-.21.246-.484.186-.403-.307-.214-.496.214-.98.258-1.28.21-1.016.19-1.263.112-.42-.008-.028-.092.012-.953 1.307-1.448 1.957-1.146 1.227-.274.109-.477-.247.045-.44.266-.39 1.586-2.018.956-1.25.617-.723-.004-.105h-.036l-4.212 2.736-.75.096-.324-.302.04-.496.154-.162 1.267-.871z'
  },
  { 
    id: 'vscode', 
    label: 'VS Code', 
    color: '#0078D4',
    viewBox: '0 0 24 24',
    path: 'M23.15 2.587L18.21.21a1.494 1.494 0 0 0-1.705.29l-9.46 8.63-4.12-3.128a.999.999 0 0 0-1.276.057L.327 7.261A1 1 0 0 0 .326 8.74L3.899 12 .326 15.26a1 1 0 0 0 .001 1.479L1.65 17.94a.999.999 0 0 0 1.276.057l4.12-3.128 9.46 8.63a1.492 1.492 0 0 0 1.704.29l4.942-2.377A1.5 1.5 0 0 0 24 20.06V3.939a1.5 1.5 0 0 0-.85-1.352zm-5.146 14.861L10.826 12l7.178-5.448v10.896z'
  },
  { 
    id: 'codex', 
    label: 'Codex', 
    color: '#000000',
    viewBox: '0 0 16 16',
    path: 'M14.949 6.547a3.94 3.94 0 0 0-.348-3.273 4.11 4.11 0 0 0-4.4-1.934A4.1 4.1 0 0 0 8.423.2 4.15 4.15 0 0 0 6.305.086a4.1 4.1 0 0 0-1.891.948 4.04 4.04 0 0 0-1.158 1.753 4.1 4.1 0 0 0-1.563.679A4 4 0 0 0 .554 4.72a3.99 3.99 0 0 0 .502 4.731 3.94 3.94 0 0 0 .346 3.274 4.11 4.11 0 0 0 4.402 1.933c.382.425.852.764 1.377.995.526.231 1.095.35 1.67.346 1.78.002 3.358-1.132 3.901-2.804a4.1 4.1 0 0 0 1.563-.68 4 4 0 0 0 1.14-1.253 3.99 3.99 0 0 0-.506-4.716m-6.097 8.406a3.05 3.05 0 0 1-1.945-.694l.096-.054 3.23-1.838a.53.53 0 0 0 .265-.455v-4.49l1.366.778q.02.011.025.035v3.722c-.003 1.653-1.361 2.992-3.037 2.996m-6.53-2.75a2.95 2.95 0 0 1-.36-2.01l.095.057L5.29 12.09a.53.53 0 0 0 .527 0l3.949-2.246v1.555a.05.05 0 0 1-.022.041L6.473 13.3c-1.454.826-3.311.335-4.15-1.098m-.85-6.94A3.02 3.02 0 0 1 3.07 3.949v3.785a.51.51 0 0 0 .262.451l3.93 2.237-1.366.779a.05.05 0 0 1-.048 0L2.585 9.342a2.98 2.98 0 0 1-1.113-4.094zm11.216 2.571L8.747 5.576l1.362-.776a.05.05 0 0 1 .048 0l3.265 1.86a3 3 0 0 1 1.173 1.207 2.96 2.96 0 0 1-.27 3.2 3.05 3.05 0 0 1-1.36.997V8.279a.52.52 0 0 0-.276-.445m1.36-2.015-.097-.057-3.226-1.855a.53.53 0 0 0-.53 0L6.249 6.153V4.598a.04.04 0 0 1 .019-.04L9.533 2.7a3.07 3.07 0 0 1 3.257.139c.474.325.843.778 1.066 1.303.223.526.289 1.103.191 1.664zM5.503 8.575 4.139 7.8a.05.05 0 0 1-.026-.037V4.049c0-.57.166-1.127.476-1.607s.752-.864 1.275-1.105a3.08 3.08 0 0 1 3.234.41l-.096.054-3.23 1.838a.53.53 0 0 0-.265.455zm.742-1.577 1.758-1 1.762 1v2l-1.755 1-1.762-1z'
  },
  { 
    id: 'windsurf', 
    label: 'Windsurf', 
    color: '#000000',
    viewBox: '0 0 24 24',
    path: 'M23.55 5.067c-1.2038-.002-2.1806.973-2.1806 2.1765v4.8676c0 .972-.8035 1.7594-1.7597 1.7594-.568 0-1.1352-.286-1.4718-.7659l-4.9713-7.1003c-.4125-.5896-1.0837-.941-1.8103-.941-1.1334 0-2.1533.9635-2.1533 2.153v4.8957c0 .972-.7969 1.7594-1.7596 1.7594-.57 0-1.1363-.286-1.4728-.7658L.4076 5.1598C.2822 4.9798 0 5.0688 0 5.2882v4.2452c0 .2147.0656.4228.1884.599l5.4748 7.8183c.3234.462.8006.8052 1.3509.9298 1.3771.313 2.6446-.747 2.6446-2.0977v-4.893c0-.972.7875-1.7593 1.7596-1.7593h.003a1.798 1.798 0 0 1 1.4718.7658l4.9723 7.0994c.4135.5905 1.05.941 1.8093.941 1.1587 0 2.1515-.9645 2.1515-2.153v-4.8948c0-.972.7875-1.7594 1.7596-1.7594h.194a.22.22 0 0 0 .2204-.2202v-4.622a.22.22 0 0 0-.2203-.2203Z'
  },
  { 
    id: 'gemini', 
    label: 'Gemini CLI', 
    color: '#4285F4',
    viewBox: '0 0 24 24',
    path: 'M11.04 19.32Q12 21.51 12 24q0-2.49.93-4.68.96-2.19 2.58-3.81t3.81-2.55Q21.51 12 24 12q-2.49 0-4.68-.93a12.3 12.3 0 0 1-3.81-2.58 12.3 12.3 0 0 1-2.58-3.81Q12 2.49 12 0q0 2.49-.96 4.68-.93 2.19-2.55 3.81a12.3 12.3 0 0 1-3.81 2.58Q2.49 12 0 12q2.49 0 4.68.96 2.19.93 3.81 2.55t2.55 3.81'
  }
]
const activeIde = ref('cursor')

// MCP Config - 根据不同 IDE 返回不同配置
const mcpConfigs: Record<string, { language: string; code: string }> = {
  cursor: {
    language: 'json',
    code: `{
  "mcpServers": {
    "go-mcp-context": {
      "url": "https://mcp.hsk423.cn/mcp",
      "headers": {
        "Authorization": "Bearer YOUR_API_KEY"
      }
    }
  }
}`
  },
  claude: {
    language: 'bash',
    code: `claude mcp add --transport http go-mcp-context https://mcp.hsk423.cn/mcp \\
  --header "Authorization: Bearer YOUR_API_KEY"`
  },
  vscode: {
    language: 'json',
    code: `"mcp": {
  "servers": {
    "go-mcp-context": {
      "type": "http",
      "url": "https://mcp.hsk423.cn/mcp",
      "headers": {
        "Authorization": "Bearer YOUR_API_KEY"
      }
    }
  }
}`
  },
  codex: {
    language: 'toml',
    code: `[mcp_servers.go-mcp-context]
url = "https://mcp.hsk423.cn/mcp"
http_headers = { "Authorization" = "Bearer YOUR_API_KEY" }`
  },
  windsurf: {
    language: 'json',
    code: `{
  "mcpServers": {
    "go-mcp-context": {
      "serverUrl": "https://mcp.hsk423.cn/mcp",
      "headers": {
        "Authorization": "Bearer YOUR_API_KEY"
      }
    }
  }
}`
  },
  gemini: {
    language: 'json',
    code: `{
  "mcpServers": {
    "go-mcp-context": {
      "httpUrl": "https://mcp.hsk423.cn/mcp",
      "headers": {
        "Authorization": "Bearer YOUR_API_KEY",
        "Accept": "application/json, text/event-stream"
      }
    }
  }
}`
  }
}

const mcpConfig = computed(() => {
  return mcpConfigs[activeIde.value]?.code || mcpConfigs.cursor.code
})

const mcpLanguage = computed(() => {
  return mcpConfigs[activeIde.value]?.language || 'json'
})

// 代码块 ref
const codeBlock = ref<HTMLElement | null>(null)

// 高亮代码
const highlightCode = () => {
  nextTick(() => {
    if (codeBlock.value) {
      // 清除之前的高亮状态
      codeBlock.value.removeAttribute('data-highlighted')
      codeBlock.value.innerHTML = ''
      codeBlock.value.textContent = mcpConfig.value
      codeBlock.value.className = `language-${mcpLanguage.value} rounded`
      hljs.highlightElement(codeBlock.value)
    }
  })
}

// 监听 IDE 切换，重新高亮
watch(activeIde, highlightCode)

// 初始化
onMounted(async () => {
  await initUserState()
  loadAPIKeys()
  loadStats()
  highlightCode()
  highlightApiCode()
})

// API Tab
const apiTab = ref('search')
const apiCommandBlock = ref<HTMLElement | null>(null)
const apiResponseBlock = ref<HTMLElement | null>(null)
const docsType = ref('code') // code 或 info

const apiCommand = computed(() => {
  if (apiTab.value === 'search') {
    return `curl -X POST "https://mcp.hsk423.cn/mcp" \\
  -H "Content-Type: application/json" \\
  -H "MCP_API_KEY: YOUR_API_KEY" \\
  -d '{
    "jsonrpc": "2.0",
    "id": 1,
    "method": "tools/call",
    "params": {
      "name": "search-libraries",
      "arguments": {
        "libraryName": "gin"
      }
    }
  }'`
  }
  
  // Docs tab - 根据 docsType 构建请求
  const mode = docsType.value
  return `curl -X POST "https://mcp.hsk423.cn/mcp" \\
  -H "Content-Type: application/json" \\
  -H "MCP_API_KEY: YOUR_API_KEY" \\
  -d '{
    "jsonrpc": "2.0",
    "id": 1,
    "method": "tools/call",
    "params": {
      "name": "get-library-docs",
      "arguments": {
        "libraryId": 6,
        "version": "v1.9.1",
        "topic": "routing",
        "mode": "${mode}",
        "page": 1
      }
    }
  }'`
})

const apiResponse = computed(() => {
  if (apiTab.value === 'search') {
    return `{
  "jsonrpc": "2.0",
  "id": 1,
  "result": {
    "libraries": [
      {
        "libraryId": 6,
        "name": "gin",
        "versions": ["v1.11.0", "v1.10.1", "v1.9.1", "v1.8.2", "v1.7.7"],
        "defaultVersion": "latest",
        "description": "Gin is a high-performance HTTP web framework written in Go. It provides a Martini-like API but with significantly better performance—up to 40 times faster—thanks to httprouter.",
        "snippets": 716,
        "score": 1
      }
    ]
  }
}`
  }
  
  // Docs tab - 根据 mode 返回不同的响应示例
  if (docsType.value === 'code') {
    // Code 模式：包含 title, description, source, language, code（不含 content）
    return `{
  "jsonrpc": "2.0",
  "id": 1,
  "result": {
    "libraryId": 6,
    "documents": [
      {
        "title": "Defining Routes with Different HTTP Methods in Gin",
        "description": "This code snippet demonstrates how to define routes for various HTTP methods using the Gin framework in Go. It includes examples for GET, POST, PUT, PATCH, DELETE, and OPTIONS requests.",
        "source": "mcp/docs/gin/v1.9.1/docs/doc.md",
        "language": "go",
        "code": "func main() {\\n  router := gin.Default()\\n\\n  router.GET(\"/someGet\", getting)\\n  router.POST(\"/somePost\", posting)\\n  router.PUT(\"/somePut\", putting)\\n  router.DELETE(\"/someDelete\", deleting)\\n  router.PATCH(\"/somePatch\", patching)\\n  router.HEAD(\"/someHead\", head)\\n  router.OPTIONS(\"/someOptions\", options)\\n\\n  router.Run()\\n}",
        "tokens": 319,
        "relevance": 0.134
      },
      {
        "title": "Grouping API Routes in Gin Framework",
        "description": "This code snippet demonstrates how to group API routes using the Gin framework in Go. It defines two separate route groups, v1 and v2, each containing three endpoints: login, submit, and read.",
        "source": "mcp/docs/gin/v1.9.1/docs/doc.md",
        "language": "go",
        "code": "func main() {\\n  router := gin.Default()\\n\\n  v1 := router.Group(\"/v1\")\\n  {\\n    v1.POST(\"/login\", loginEndpoint)\\n    v1.POST(\"/submit\", submitEndpoint)\\n    v1.POST(\"/read\", readEndpoint)\\n  }\\n\\n  v2 := router.Group(\"/v2\")\\n  {\\n    v2.POST(\"/login\", loginEndpoint)\\n    v2.POST(\"/submit\", submitEndpoint)\\n    v2.POST(\"/read\", readEndpoint)\\n  }\\n\\n  router.Run(\":8080\")\\n}",
        "tokens": 126,
        "relevance": 0.133
      }
    ],
    "page": 1,
    "hasMore": true
  }
}`
  } else {
    // Info 模式：只包含 title, source, content（无 description, language, code）
    return `{
  "jsonrpc": "2.0",
  "id": 1,
  "result": {
    "libraryId": 6,
    "documents": [
      {
        "title": "Gin Web Framework > Getting started > Installation",
        "source": "mcp/docs/gin/v1.9.1/README.md",
        "content": "To install Gin package, you need to install Go and set your Go workspace first.\\n\\n1. Download and install it:\\n\\n$ go get -u github.com/gin-gonic/gin\\n\\n2. Import it in your code:\\n\\nimport \\"github.com/gin-gonic/gin\\"\\n\\n3. (Optional) Import net/http. This is required for example if using constants such as http.StatusOK.",
        "tokens": 105,
        "relevance": 0.096
      },
      {
        "title": "Gin Web Framework > Benchmarks",
        "source": "mcp/docs/gin/v1.9.1/README.md",
        "content": "- (1): Total Repetitions achieved in constant time, higher means more confident result\\n- (2): Single Repetition Duration (ns/op), lower is better\\n- (3): Heap Memory (B/op), lower is better\\n- (4): Average Allocations per Repetition (allocs/op), lower is better",
        "tokens": 69,
        "relevance": 0.091
      }
    ],
    "page": 1,
    "hasMore": true
  }
}`
  }
})

// API 代码高亮
const highlightApiCode = () => {
  nextTick(() => {
    // 高亮 API 命令
    if (apiCommandBlock.value) {
      apiCommandBlock.value.removeAttribute('data-highlighted')
      apiCommandBlock.value.innerHTML = ''
      apiCommandBlock.value.textContent = apiCommand.value
      hljs.highlightElement(apiCommandBlock.value)
    }
    
    // 高亮 API 响应
    if (apiResponseBlock.value) {
      apiResponseBlock.value.removeAttribute('data-highlighted')
      apiResponseBlock.value.innerHTML = ''
      apiResponseBlock.value.textContent = apiResponse.value
      hljs.highlightElement(apiResponseBlock.value)
    }
  })
}

// 监听 API Tab 切换，重新高亮
watch(apiTab, highlightApiCode)
// 监听 Docs Type 切换，重新高亮
watch(docsType, highlightApiCode)

// Copy functions（兼容 HTTP 环境）
const copyToClipboard = async (text: string) => {
  try {
    if (navigator.clipboard && window.isSecureContext) {
      await navigator.clipboard.writeText(text)
    } else {
      const textarea = document.createElement('textarea')
      textarea.value = text
      textarea.style.position = 'fixed'
      textarea.style.opacity = '0'
      document.body.appendChild(textarea)
      textarea.select()
      document.execCommand('copy')
      document.body.removeChild(textarea)
    }
    ElMessage.success('已复制')
  } catch (err) {
    console.error('复制失败:', err)
  }
}

const copyCode = () => {
  copyToClipboard(mcpConfig.value)
}
</script>
