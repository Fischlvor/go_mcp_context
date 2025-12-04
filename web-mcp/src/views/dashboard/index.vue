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
    />

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
                        <span class="text-base font-medium text-stone-800">mcp.context7.com/mcp</span>
                        <button class="text-stone-400 transition-colors hover:text-stone-600" @click="copyToClipboard('mcp.context7.com/mcp')">
                          <svg xmlns="http://www.w3.org/2000/svg" width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round">
                            <path d="M7 7m0 2.667a2.667 2.667 0 0 1 2.667 -2.667h8.666a2.667 2.667 0 0 1 2.667 2.667v8.666a2.667 2.667 0 0 1 -2.667 2.667h-8.666a2.667 2.667 0 0 1 -2.667 -2.667z"></path>
                            <path d="M4.012 16.737a2.005 2.005 0 0 1 -1.012 -1.737v-10c0 -1.1 .9 -2 2 -2h10c.75 0 1.158 .385 1.5 1"></path>
                          </svg>
                        </button>
                      </div>
                    </div>
                    <div class="border-t border-stone-200"></div>
                    <div class="flex flex-col gap-2 py-3 sm:grid sm:grid-cols-[80px_auto_1fr] sm:items-center sm:py-2">
                      <span class="text-sm font-normal uppercase text-stone-500">API URL</span>
                      <span class="hidden text-sm text-stone-500 sm:block">:</span>
                      <div class="flex items-center gap-2">
                        <span class="text-base font-medium text-stone-800">context7.com/api/v2</span>
                        <button class="text-stone-400 transition-colors hover:text-stone-600" @click="copyToClipboard('context7.com/api/v2')">
                          <svg xmlns="http://www.w3.org/2000/svg" width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round">
                            <path d="M7 7m0 2.667a2.667 2.667 0 0 1 2.667 -2.667h8.666a2.667 2.667 0 0 1 2.667 2.667v8.666a2.667 2.667 0 0 1 -2.667 2.667h-8.666a2.667 2.667 0 0 1 -2.667 -2.667z"></path>
                            <path d="M4.012 16.737a2.005 2.005 0 0 1 -1.012 -1.737v-10c0 -1.1 .9 -2 2 -2h10c.75 0 1.158 .385 1.5 1"></path>
                          </svg>
                        </button>
                      </div>
                    </div>
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
                            <component :is="ide.icon" class="h-3.5 w-3.5" />
                            {{ ide.label }}
                          </button>
                          <div class="flex-grow border-b border-stone-300"></div>
                        </div>
                      </div>
                    </div>

                    <!-- Remote/Local Toggle -->
                    <div class="flex" role="group">
                      <button 
                        :class="['border px-3 py-1 text-sm font-medium shadow-sm transition-colors rounded-l-md', 
                          connectionType === 'remote' ? 'border-stone-500 bg-stone-700 text-white' : 'border-stone-300 bg-white text-stone-800 hover:bg-stone-50']"
                        @click="connectionType = 'remote'"
                      >Remote</button>
                      <button 
                        :class="['border px-3 py-1 text-sm font-medium shadow-sm transition-colors rounded-r-md border-l-0', 
                          connectionType === 'local' ? 'border-stone-500 bg-stone-700 text-white' : 'border-stone-300 bg-white text-stone-800 hover:bg-stone-50']"
                        @click="connectionType = 'local'"
                      >Local</button>
                    </div>

                    <!-- Code Block -->
                    <div class="relative rounded-lg bg-stone-100 p-4">
                      <button class="absolute right-3 top-3 text-stone-400 transition-colors hover:text-stone-600" @click="copyCode">
                        <svg xmlns="http://www.w3.org/2000/svg" width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round">
                          <path d="M7 7m0 2.667a2.667 2.667 0 0 1 2.667 -2.667h8.666a2.667 2.667 0 0 1 2.667 2.667v8.666a2.667 2.667 0 0 1 -2.667 2.667h-8.666a2.667 2.667 0 0 1 -2.667 -2.667z"></path>
                          <path d="M4.012 16.737a2.005 2.005 0 0 1 -1.012 -1.737v-10c0 -1.1 .9 -2 2 -2h10c.75 0 1.158 .385 1.5 1"></path>
                        </svg>
                      </button>
                      <pre class="text-sm leading-relaxed text-stone-700"><code>{{ mcpConfig }}</code></pre>
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
                    <h2 class="text-lg font-semibold text-stone-900 sm:text-xl">API</h2>
                    <div class="text-sm text-stone-500 sm:text-base">
                      <p>Use the Context7 API to search libraries and fetch documentation programmatically</p>
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

                  <!-- API Code Block -->
                  <div class="relative rounded-lg bg-stone-100 p-4">
                    <button class="absolute right-3 top-3 text-stone-400 transition-colors hover:text-stone-600">
                      <svg xmlns="http://www.w3.org/2000/svg" width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round">
                        <path d="M7 7m0 2.667a2.667 2.667 0 0 1 2.667 -2.667h8.666a2.667 2.667 0 0 1 2.667 2.667v8.666a2.667 2.667 0 0 1 -2.667 2.667h-8.666a2.667 2.667 0 0 1 -2.667 -2.667z"></path>
                        <path d="M4.012 16.737a2.005 2.005 0 0 1 -1.012 -1.737v-10c0 -1.1 .9 -2 2 -2h10c.75 0 1.158 .385 1.5 1"></path>
                      </svg>
                    </button>
                    <pre class="text-sm leading-relaxed text-stone-700"><code>{{ apiCommand }}</code></pre>
                  </div>

                  <!-- Parameters -->
                  <div class="space-y-4">
                    <div>
                      <h4 class="mb-2 text-sm font-medium text-stone-700">Parameters</h4>
                      <p class="text-sm text-stone-500">
                        <code class="rounded bg-stone-100 px-1 text-xs">query</code> - Search term for finding libraries
                      </p>
                    </div>
                    <div>
                      <h4 class="mb-2 text-sm font-medium text-stone-700">Response</h4>
                      <div class="relative rounded-lg bg-stone-100 p-4">
                        <button class="absolute right-3 top-3 text-stone-400 transition-colors hover:text-stone-600">
                          <svg xmlns="http://www.w3.org/2000/svg" width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round">
                            <path d="M7 7m0 2.667a2.667 2.667 0 0 1 2.667 -2.667h8.666a2.667 2.667 0 0 1 2.667 2.667v8.666a2.667 2.667 0 0 1 -2.667 2.667h-8.666a2.667 2.667 0 0 1 -2.667 -2.667z"></path>
                            <path d="M4.012 16.737a2.005 2.005 0 0 1 -1.012 -1.737v-10c0 -1.1 .9 -2 2 -2h10c.75 0 1.158 .385 1.5 1"></path>
                          </svg>
                        </button>
                        <pre class="max-h-[300px] overflow-auto text-sm leading-relaxed text-stone-700"><code>{{ apiResponse }}</code></pre>
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
import { ref, computed, h, onMounted } from 'vue'
import AppHeader from '@/components/AppHeader.vue'
import AppFooter from '@/components/AppFooter.vue'
import { useUser } from '@/stores/user'
import { getAPIKeys, createAPIKey, deleteAPIKey, type APIKey, type APIKeyCreateResponse } from '@/api/apikey'

// 用户状态
const { isLoggedIn, userEmail, userPlan, initUserState } = useUser()

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
    if (res.code === 0) {
      apiKeys.value = res.data || []
    }
  } catch (e) {
    console.error('Failed to load API keys:', e)
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
    if (res.code === 0) {
      newlyCreatedKey.value = res.data
      newKeyName.value = ''
      showCreateDialog.value = false
      await loadAPIKeys()
    } else {
      alert(res.msg || '创建失败')
    }
  } catch (e) {
    console.error('Failed to create API key:', e)
    alert('创建失败')
  } finally {
    creatingKey.value = false
  }
}

// 删除 API Key
const handleDeleteKey = async (id: number) => {
  if (!confirm('确定要删除这个 API Key 吗？')) return
  try {
    const res = await deleteAPIKey(id)
    if (res.code === 0) {
      await loadAPIKeys()
    } else {
      alert(res.msg || '删除失败')
    }
  } catch (e) {
    console.error('Failed to delete API key:', e)
    alert('删除失败')
  }
}

// 复制新创建的 Key
const copyNewKey = () => {
  if (newlyCreatedKey.value) {
    navigator.clipboard.writeText(newlyCreatedKey.value.api_key)
    alert('已复制到剪贴板')
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

onMounted(() => {
  initUserState()
  loadAPIKeys()
})

// Tabs
const tabs = [
  { id: 'overview', label: 'Overview' },
  { id: 'libraries', label: 'Libraries' },
  { id: 'members', label: 'Members' },
  { id: 'rules', label: 'Rules' }
]
const activeTab = ref('overview')

// Stats
const stats = [
  { label: 'Search Requests', value: '0' },
  { label: 'Docs Requests', value: '0' },
  { label: 'Parsing Tokens', value: '0' },
  { label: 'Monthly Cost', value: 'N/A' }
]

// IDE Icons (简化版)
const CursorIcon = {
  render() {
    return h('svg', { width: 14, height: 14, viewBox: '0 0 24 24', fill: 'currentColor' }, [
      h('path', { d: 'M11.503.131 1.891 5.678a.84.84 0 0 0-.42.726v11.188c0 .3.162.575.42.724l9.609 5.55a1 1 0 0 0 .998 0l9.61-5.55a.84.84 0 0 0 .42-.724V6.404a.84.84 0 0 0-.42-.726L12.497.131a1.01 1.01 0 0 0-.996 0' })
    ])
  }
}

const ides = [
  { id: 'cursor', label: 'Cursor', icon: CursorIcon },
  { id: 'claude', label: 'Claude Code', icon: CursorIcon },
  { id: 'vscode', label: 'VS Code', icon: CursorIcon },
  { id: 'codex', label: 'Codex', icon: CursorIcon },
  { id: 'windsurf', label: 'Windsurf', icon: CursorIcon },
  { id: 'gemini', label: 'Gemini CLI', icon: CursorIcon },
  { id: 'more', label: 'More', icon: CursorIcon }
]
const activeIde = ref('cursor')

// Connection type
const connectionType = ref('remote')

// MCP Config
const mcpConfig = computed(() => {
  return `{
  "mcpServers": {
    "context7": {
      "url": "https://mcp.context7.com/mcp",
      "headers": {
        "CONTEXT7_API_KEY": "YOUR_API_KEY"
      }
    }
  }
}`
})

// API Tab
const apiTab = ref('search')

const apiCommand = computed(() => {
  if (apiTab.value === 'search') {
    return `curl -X GET "https://context7.com/api/v2/search?query=next.js" \\
  -H "Authorization: Bearer CONTEXT7_API_KEY"`
  }
  return `curl -X GET "https://context7.com/api/v2/docs?id=/vercel/next.js" \\
  -H "Authorization: Bearer CONTEXT7_API_KEY"`
})

const apiResponse = computed(() => {
  return `{
  "results": [
    {
      "id": "/websites/nextjs",
      "title": "Next.js",
      "description": "Next.js is a React framework...",
      "totalTokens": 1526838,
      "totalSnippets": 7382
    }
  ]
}`
})

// Copy functions
const copyToClipboard = (text: string) => {
  navigator.clipboard.writeText(text)
}

const copyCode = () => {
  navigator.clipboard.writeText(mcpConfig.value)
}
</script>
