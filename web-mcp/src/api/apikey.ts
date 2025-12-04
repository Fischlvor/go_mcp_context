import service from '@/utils/request'
import type { ApiResponse } from '@/utils/request'

export interface APIKey {
  id: number
  name: string
  token_suffix: string
  last_used_at: string | null
  created_at: string
}

export interface APIKeyCreateRequest {
  name: string
}

export interface APIKeyCreateResponse {
  id: number
  name: string
  api_key: string  // 完整 key，仅创建时返回
  token_suffix: string
  created_at: string
}

// 获取 API Key 列表
export const getAPIKeys = (): Promise<ApiResponse<APIKey[]>> => {
  return service({
    url: '/api-keys/list',
    method: 'get'
  })
}

// 创建 API Key
export const createAPIKey = (data: APIKeyCreateRequest): Promise<ApiResponse<APIKeyCreateResponse>> => {
  return service({
    url: '/api-keys/create',
    method: 'post',
    data
  })
}

// 删除 API Key
export const deleteAPIKey = (id: number): Promise<ApiResponse<null>> => {
  return service({
    url: `/api-keys/${id}`,
    method: 'delete'
  })
}
