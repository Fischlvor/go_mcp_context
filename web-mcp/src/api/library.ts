import service from '@/utils/request'
import type { ApiResponse } from '@/utils/request'

export interface Library {
  id: number
  name: string
  version: string
  description: string
  status: string
  document_count?: number
  chunk_count?: number
  created_at: string
  updated_at: string
}

export interface LibraryListResponse {
  list: Library[]
  total: number
  page: number
  page_size: number
}

export interface LibraryCreateRequest {
  name: string
  version: string
  description: string
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
