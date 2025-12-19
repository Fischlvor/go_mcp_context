import service from '@/utils/request'
import type { ApiResponse } from '@/utils/request'
import { uploadWithSSE, type SSEEvent } from '@/utils/sse'

// ============================================================
// 类型定义
// ============================================================

// 文档上传记录
export interface Document {
  id: number
  library_id: number
  version: string
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

// 文档列表响应
export interface DocumentListResponse {
  list: Document[]
  total: number
  page: number
  page_size: number
}

// 文档块
export interface DocumentChunk {
  id: number
  library_id: number
  upload_id: number
  version: string
  chunk_index: number
  title: string
  description: string
  source: string
  language: string
  code: string
  chunk_text: string
  tokens: number
  chunk_type: 'code' | 'info' | 'mixed'
  access_count: number
  status: string
  created_at: string
  updated_at: string
}

// 文档块响应
export interface ChunksResponse {
  chunks: DocumentChunk[]
}

// 文档内容（合并后）
export interface DocumentContent {
  title: string
  content: string
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

// ============================================================
// 常量
// ============================================================

// 块分隔符
const CHUNK_SEPARATOR = '\n\n--------------------------------\n\n'

// ============================================================
// API 请求
// ============================================================

// 获取文档列表（支持版本过滤）
export const getDocuments = (params: { 
  library_id: number
  version?: string
  page?: number
  page_size?: number 
}): Promise<DocumentListResponse> => {
  return service({
    url: '/documents/list',
    method: 'get',
    params
  })
}

// 获取文档详情
export const getDocument = (id: number): Promise<ApiResponse<Document>> => {
  return service({
    url: `/documents/detail/${id}`,
    method: 'get'
  })
}

// 获取库的文档块（统一入口，支持搜索和列表）
// mode: 'code' 返回代码块, 'info' 返回文档块
// version: 可选，不传则使用库的默认版本
// topic: 可选，传入则进行向量搜索，不传则返回全部文档块
export const getChunks = (
  mode: 'code' | 'info', 
  libraryId: number, 
  options?: { version?: string; topic?: string }
): Promise<ChunksResponse & { total?: number; topic?: string }> => {
  const params: Record<string, string> = {}
  if (options?.version) params.version = options.version
  if (options?.topic) params.topic = options.topic
  
  return service({
    url: `/documents/chunks/${mode}/${libraryId}`,
    method: 'get',
    params: Object.keys(params).length > 0 ? params : undefined
  })
}

// 格式化单个代码块（code 模式）
const formatCodeChunk = (chunk: DocumentChunk): string => {
  const parts: string[] = []
  if (chunk.title) parts.push(`### ${chunk.title}`)
  if (chunk.source) parts.push(`Source: ${chunk.source}`)
  if (chunk.description) parts.push(chunk.description)
  // 优先使用 code 字段，否则使用 chunk_text
  const content = chunk.code || chunk.chunk_text
  if (content) {
    const lang = chunk.language || ''
    parts.push(`\`\`\`${lang}\n${content}\n\`\`\``)
  }
  return parts.join('\n\n')
}

// 格式化单个信息块（info 模式）
const formatInfoChunk = (chunk: DocumentChunk): string => {
  const parts: string[] = []
  if (chunk.title) parts.push(`### ${chunk.title}`)
  if (chunk.source) parts.push(`Source: ${chunk.source}`)
  if (chunk.description) parts.push(chunk.description)
  if (chunk.chunk_text) parts.push(chunk.chunk_text)
  return parts.join('\n\n')
}

// 获取代码块内容（合并后）
export const getLatestCode = async (libraryId: number, version?: string): Promise<DocumentContent> => {
  const res = await getChunks('code', libraryId, { version })
  const chunks = res.chunks || []
  const content = chunks.map(formatCodeChunk).join(CHUNK_SEPARATOR)
  const title = chunks.length > 0 ? chunks[0].title : ''
  return { title, content }
}

// 获取文档信息块内容（合并后）
export const getLatestInfo = async (libraryId: number, version?: string): Promise<DocumentContent> => {
  const res = await getChunks('info', libraryId, { version })
  const chunks = res.chunks || []
  const content = chunks.map(formatInfoChunk).join(CHUNK_SEPARATOR)
  const title = chunks.length > 0 ? chunks[0].title : ''
  return { title, content }
}

// 上传文档（普通方式）
export const uploadDocument = (libraryId: number, file: File, version: string = 'latest'): Promise<ApiResponse<Document>> => {
  const formData = new FormData()
  formData.append('file', file)
  formData.append('library_id', String(libraryId))
  formData.append('version', version)
  
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
  },
  version: string = 'latest'
): Promise<void> => {
  return uploadWithSSE(
    '/documents/upload-sse',
    file,
    { library_id: String(libraryId), version },
    {
      onMessage: (event: SSEEvent<ProcessStatus>) => {
        console.log('[SSE] Received event:', event)
        const stage = event.data.stage || event.type
        
        if (stage === 'completed') {
          callbacks.onComplete?.(event.data)
        } else if (stage === 'failed' || stage === 'error') {
          callbacks.onError?.(new Error(event.data.message || 'Processing failed'))
        } else {
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

// 删除文档
export const deleteDocument = (id: number): Promise<ApiResponse<null>> => {
  return service({
    url: `/documents/${id}`,
    method: 'delete'
  })
}
