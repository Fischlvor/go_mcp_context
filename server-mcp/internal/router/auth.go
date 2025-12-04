package router

import (
	"go-mcp-context/internal/api"

	"github.com/gin-gonic/gin"
)

type AuthRouter struct{}

func (a *AuthRouter) InitAuthRouter(PublicRouter *gin.RouterGroup) {
	authPublicRouter := PublicRouter.Group("auth")
	authApi := api.ApiGroupApp.AuthApi
	{
		// 获取 SSO 登录 URL
		authPublicRouter.GET("sso_login_url", authApi.GetSSOLoginURL)
		// SSO 回调接口（后端用 code 换 token，refresh_token 存 session）
		authPublicRouter.GET("callback", authApi.SSOCallback)
		// 登出
		authPublicRouter.POST("logout", authApi.Logout)
	}
}
