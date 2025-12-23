import service from '@/utils/request'
import type { ApiResponse } from '@/utils/request'
// SSE 相关导入（已改用普通接口，通过日志查看进度）
// import { createSSEPost, type SSEEvent } from '@/utils/sse'

// 库列表项（精简字段，用于主页表格）
export interface LibraryListItem {
  id: number
  name: string
  source_type: string      // github, website, local
  source_url: string       // vuejs/docs
  default_version: string
  token_count: number      // 对应 TOKENS
  chunk_count: number      // 对应 SNIPPETS
  updated_at: string       // 对应 UPDATE
}

// 库详情（完整字段）
export interface Library {
  id: number
  name: string
  default_version: string
  versions: string[]
  source_type: string
  source_url: string
  description: string
  document_count: number
  chunk_count: number
  token_count: number
  status: string
  created_at: string
  updated_at: string
}

export interface LibraryListResponse {
  list: LibraryListItem[]
  total: number
  page: number
  page_size: number
}

export interface LibraryCreateRequest {
  name: string
  description?: string
  source_type?: string  // github, website, local
  source_url?: string   // vuejs/docs
}

// 活动日志
export interface ActivityLog {
  id: number
  library_id: number
  actor_id?: number
  event: string           // document.upload, version.create 等
  status: string          // info, success, warning, error
  message: string
  target_type?: string    // document, version
  target_id?: string
  task_id?: string
  version?: string
  metadata?: Record<string, unknown>
  time: string            // ISO 时间戳
}

// 获取库列表
export const getLibraries = (params?: { name?: string; page?: number; page_size?: number; sort?: string }): Promise<LibraryListResponse> => {
  return service({
    url: '/libraries',
    method: 'get',
    params
  })
}

// 创建库
export const createLibrary = (data: LibraryCreateRequest): Promise<ApiResponse<Library>> => {
  return service({
    url: '/libraries',
    method: 'post',
    data
  })
}

// 获取库详情
export const getLibrary = (id: number): Promise<Library> => {
  return service({
    url: `/libraries/${id}`,
    method: 'get'
  })
}

// 更新库
export const updateLibrary = (id: number, data: LibraryCreateRequest): Promise<Library> => {
  return service({
    url: `/libraries/${id}`,
    method: 'put',
    data
  })
}

// 删除库
export const deleteLibrary = (id: number): Promise<null> => {
  return service({
    url: `/libraries/${id}`,
    method: 'delete'
  })
}

// 获取版本列表
export const getVersions = (libraryId: number): Promise<any> => {
  return service({
    url: `/libraries/${libraryId}/versions`,
    method: 'get'
  })
}

// 创建版本
export const createVersion = (libraryId: number, version: string): Promise<null> => {
  return service({
    url: `/libraries/${libraryId}/versions`,
    method: 'post',
    data: { version }
  })
}

// 删除版本
export const deleteVersion = (libraryId: number, version: string): Promise<null> => {
  return service({
    url: `/libraries/${libraryId}/versions/${version}`,
    method: 'delete'
  })
}

// 刷新版本（异步）
export const refreshVersion = (libraryId: number, version: string): Promise<null> => {
  return service({
    url: `/libraries/${libraryId}/versions/${version}/refresh`,
    method: 'post'
  })
}

// 刷新状态类型
export interface RefreshStatus {
  doc_id?: number
  doc_title?: string
  stage: 'started' | 'doc_processing' | 'doc_completed' | 'doc_failed' | 'all_completed' | 'error'
  current: number
  total: number
  message: string
}

// ====== 以下是 SSE 版本的函数，保留备用 ======
// 刷新版本（SSE 实时推送）
// export const refreshVersionWithSSE = (
//   libraryId: number,
//   version: string,
//   callbacks: {
//     onProgress?: (status: RefreshStatus) => void
//     onComplete?: (status: RefreshStatus) => void
//     onError?: (error: Error) => void
//   }
// ): Promise<void> => {
//   return createSSEPost(
//     `/libraries/${libraryId}/versions/${version}/refresh-sse`,
//     {},
//     {
//       onMessage: (event: SSEEvent<RefreshStatus>) => {
//         console.log('[SSE Refresh] Received event:', event)
//         const stage = event.data.stage
//         
//         if (stage === 'all_completed') {
//           callbacks.onComplete?.(event.data)
//         } else if (stage === 'error' || stage === 'doc_failed') {
//           if (stage === 'error') {
//             callbacks.onError?.(new Error(event.data.message || 'Refresh failed'))
//           } else {
//             // doc_failed 只是单个文档失败，继续处理
//             callbacks.onProgress?.(event.data)
//           }
//         } else {
//           callbacks.onProgress?.(event.data)
//         }
//       },
//       onError: callbacks.onError,
//     }
//   )
// }

// ============ GitHub 导入相关 ============

// GitHub 仓库信息响应
export interface GitHubRepoInfo {
  repo: string
  default_branch: string
  description: string
  versions: string[]  // 主版本的最新 tag 列表
}

// GitHub 导入进度
export interface GitHubImportProgress {
  stage: 'fetching_tree' | 'downloading' | 'processing' | 'completed' | 'failed' | 'info' | 'warning'
  current: number
  total: number
  message: string
  filename?: string
}

// GitHub 导入请求
export interface GitHubImportRequest {
  repo: string         // owner/repo
  branch?: string      // 分支名
  tag?: string         // 特定 tag（与 branch 二选一）
  version?: string     // 存储为的版本名（默认使用 tag 名）
  path_filter?: string // 只导入指定路径
  excludes?: string[]  // 排除模式
}

// 获取 GitHub 仓库的版本列表
export const getGitHubReleases = (repo: string): Promise<GitHubRepoInfo> => {
  return service({
    url: '/libraries/github/releases',
    method: 'get',
    params: { repo }
  })
}

// ============ 活动日志相关 ============

// 活动日志响应
export interface ActivityLogsResponse {
  logs: ActivityLog[]
  task_id?: string
  status: 'complete' | 'processing'
}

// 获取活动日志
export const getActivityLogs = (libraryId: number, limit?: number): Promise<ActivityLogsResponse> => {
  return service({
    url: '/logs',
    method: 'get',
    params: { libraryId, limit }
  })
}

// 从 GitHub 导入文档（异步，立即返回，通过日志查看进度）
export const importFromGitHub = (libraryId: number, data: GitHubImportRequest): Promise<null> => {
  return service({
    url: '/libraries/github/import',
    method: 'post',
    params: { id: libraryId },
    data
  })
}

// GitHub 初始化导入响应
export interface GitHubInitImportResponse {
  library_id: number
  version: string
}

// 从 GitHub URL 初始化导入（创建库 + 导入默认分支）
export const initImportFromGitHub = (githubUrl: string): Promise<GitHubInitImportResponse> => {
  return service({
    url: '/libraries/github/init-import',
    method: 'post',
    data: { github_url: githubUrl }
  })
}

// 从 GitHub 导入文档（SSE 实时推送）- 已改用普通接口，通过日志查看进度
// export const importFromGitHubSSE = (
//   libraryId: number,
//   data: GitHubImportRequest,
//   callbacks: {
//     onProgress?: (status: GitHubImportProgress) => void
//     onComplete?: (status: GitHubImportProgress) => void
//     onError?: (error: Error) => void
//   }
// ): Promise<void> => {
//   return createSSEPost(
//     `/libraries/${libraryId}/import-github-sse`,
//     data,
//     {
//       onMessage: (event: SSEEvent<GitHubImportProgress>) => {
//         console.log('[SSE GitHub Import] Received event:', event)
//         const stage = event.data.stage
//         
//         if (stage === 'completed') {
//           callbacks.onComplete?.(event.data)
//         } else if (stage === 'failed') {
//           callbacks.onError?.(new Error(event.data.message || 'Import failed'))
//         } else {
//           callbacks.onProgress?.(event.data)
//         }
//       },
//       onError: callbacks.onError,
//     }
//   )
// }
// ====== SSE 版本代码结束 ======

// ====== 统计接口 ======

// 用户统计数据
export interface UserStats {
  libraries: number
  documents: number
  tokens: number
  mcp_calls: number
}

// 获取当前用户统计
export const getMyStats = (): Promise<UserStats> => {
  return service({
    url: '/stats/my',
    method: 'get'
  })
}
