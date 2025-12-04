import service from '@/utils/request'
import type { ApiResponse } from '@/utils/request'

export interface SearchResultItem {
  chunk_id: number
  document_id: number
  library_id: number
  content: string
  chunk_type: string
  score: number
  vector_score: number
  bm25_score: number
}

export interface SearchResult {
  results: SearchResultItem[]
  total: number
  page: number
  limit: number
  hasMore: boolean
}

export interface SearchRequest {
  library_id: number
  query: string
  mode?: string
  page?: number
  limit?: number
}

// 搜索文档
export const searchDocuments = (data: SearchRequest): Promise<ApiResponse<SearchResult>> => {
  return service({
    url: '/search',
    method: 'post',
    data
  })
}
