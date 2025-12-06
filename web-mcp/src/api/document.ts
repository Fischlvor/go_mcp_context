import service from '@/utils/request'
import type { ApiResponse } from '@/utils/request'
import { uploadWithSSE, type SSEOptions, type SSEEvent } from '@/utils/sse'

export interface Document {
  id: number
  library_id: number
  title: string
  file_path: string
  file_type: string
  file_size: number
  content_hash: string
  chunk_count: number
  token_count: number
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

// 处理状态（SSE 推送）
// 后端 stage: uploaded → parsing(10%) → chunking(30%) → embedding(50%) → saving(80%) → completed(100%) / failed
export interface ProcessStatus {
  stage: 'uploaded' | 'parsing' | 'chunking' | 'embedding' | 'saving' | 'completed' | 'failed' | 'error'
  progress: number
  message: string
  status: string
  document_id?: number
  title?: string
}

// 获取文档列表
export const getDocuments = (params: { library_id: number; page?: number; page_size?: number }): Promise<ApiResponse<DocumentListResponse>> => {
  return service({
    url: '/documents',
    method: 'get',
    params
  })
}

// 上传文档（普通方式）
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

// 上传文档（SSE 实时状态）
export const uploadDocumentWithSSE = (
  libraryId: number,
  file: File,
  callbacks: {
    onProgress?: (status: ProcessStatus) => void
    onComplete?: (status: ProcessStatus) => void
    onError?: (error: Error) => void
  }
): Promise<void> => {
  return uploadWithSSE(
    '/documents/upload-sse',
    file,
    { library_id: String(libraryId) },
    {
      onMessage: (event: SSEEvent<ProcessStatus>) => {
        console.log('[SSE] Received event:', event)
        const stage = event.data.stage || event.type
        
        if (stage === 'completed') {
          callbacks.onComplete?.(event.data)
        } else if (stage === 'failed' || stage === 'error') {
          callbacks.onError?.(new Error(event.data.message || 'Processing failed'))
        } else {
          // 处理所有进度事件（uploaded, parsing, chunking, embedding, saving）
          callbacks.onProgress?.({
            ...event.data,
            stage: stage as ProcessStatus['stage'],
            progress: event.data.progress || 0,
            message: event.data.message || stage,
            status: event.data.status || 'processing'
          })
        }
      },
      onError: callbacks.onError,
    }
  )
}

// 获取文档详情
export const getDocument = (id: number): Promise<ApiResponse<Document>> => {
  return service({
    url: `/documents/${id}`,
    method: 'get'
  })
}

export interface DocumentContent {
  title: string
  content: string
}

// 获取库的最新文档内容
export const getLatestCode = (libraryId: number): Promise<ApiResponse<DocumentContent>> => {
  return service({
    url: `/documents/code/${libraryId}`,
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
