import service from '@/utils/request'

export interface SSOLoginUrlResponse {
  sso_login_url: string
}

export interface SSOCallbackResponse {
  access_token: string
}

// 获取 SSO 登录地址
export const getSSOLoginUrl = (redirectUri: string, returnUrl: string): Promise<SSOLoginUrlResponse> => {
  return service({
    url: '/auth/sso_login_url',
    method: 'get',
    params: { redirect_uri: redirectUri, return_url: returnUrl }
  })
}

// SSO 回调处理
export const handleSSOCallback = (code: string, redirectUri: string): Promise<SSOCallbackResponse> => {
  return service({
    url: '/auth/callback',
    method: 'get',
    params: { code, redirect_uri: redirectUri }
  })
}
