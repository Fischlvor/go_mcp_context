import service from '@/utils/request'
import type { ApiResponse } from '@/utils/request'

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

// 获取库列表
export const getLibraries = (params?: { name?: string; page?: number; page_size?: number }): Promise<ApiResponse<LibraryListResponse>> => {
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
export const getLibrary = (id: number): Promise<ApiResponse<Library>> => {
  return service({
    url: `/libraries/${id}`,
    method: 'get'
  })
}

// 更新库
export const updateLibrary = (id: number, data: LibraryCreateRequest): Promise<ApiResponse<Library>> => {
  return service({
    url: `/libraries/${id}`,
    method: 'put',
    data
  })
}

// 删除库
export const deleteLibrary = (id: number): Promise<ApiResponse<null>> => {
  return service({
    url: `/libraries/${id}`,
    method: 'delete'
  })
}

// 获取版本列表
export const getVersions = (libraryId: number): Promise<ApiResponse<any>> => {
  return service({
    url: `/libraries/${libraryId}/versions`,
    method: 'get'
  })
}

// 创建版本
export const createVersion = (libraryId: number, version: string): Promise<ApiResponse<null>> => {
  return service({
    url: `/libraries/${libraryId}/versions`,
    method: 'post',
    data: { version }
  })
}

// 删除版本
export const deleteVersion = (libraryId: number, version: string): Promise<ApiResponse<null>> => {
  return service({
    url: `/libraries/${libraryId}/versions/${version}`,
    method: 'delete'
  })
}

// 刷新版本
export const refreshVersion = (libraryId: number, version: string): Promise<ApiResponse<null>> => {
  return service({
    url: `/libraries/${libraryId}/versions/${version}/refresh`,
    method: 'post'
  })
}
