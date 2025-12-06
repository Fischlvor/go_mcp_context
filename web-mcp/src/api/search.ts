import service from '@/utils/request'
import type { ApiResponse } from '@/utils/request'

export interface SearchResultItem {
  chunk_id: number
  document_id: number
  library_id: number
  title: string      // 标题（从 Metadata 提取）
  source: string     // 来源文档标题
  content: string    // 内容
  tokens: number     // token 数
  relevance: number  // 相关性分数 0-1
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
