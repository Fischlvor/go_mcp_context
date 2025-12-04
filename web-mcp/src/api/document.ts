import service from '@/utils/request'
import type { ApiResponse } from '@/utils/request'

export interface Document {
  id: number
  library_id: number
  title: string
  file_path: string
  file_type: string
  file_size: number
  content_hash: string
  status: string
  created_at: string
  updated_at: string
}

export interface DocumentListResponse {
  list: Document[]
  total: number
  page: number
  page_size: number
}

// 获取文档列表
export const getDocuments = (params: { library_id: number; page?: number; page_size?: number }): Promise<ApiResponse<DocumentListResponse>> => {
  return service({
    url: '/documents',
    method: 'get',
    params
  })
}

// 上传文档
export const uploadDocument = (libraryId: number, file: File): Promise<ApiResponse<Document>> => {
  const formData = new FormData()
  formData.append('file', file)
  formData.append('library_id', String(libraryId))
  
  return service({
    url: '/documents/upload',
    method: 'post',
    data: formData,
    headers: {
      'Content-Type': 'multipart/form-data'
    }
  })
}

// 获取文档详情
export const getDocument = (id: number): Promise<ApiResponse<Document>> => {
  return service({
    url: `/documents/${id}`,
    method: 'get'
  })
}

// 删除文档
export const deleteDocument = (id: number): Promise<ApiResponse<null>> => {
  return service({
    url: `/documents/${id}`,
    method: 'delete'
  })
}
