package router

import (
	"go-mcp-context/internal/api"

	"github.com/gin-gonic/gin"
)

type ApiKeyRouter struct{}

// InitApiKeyRouter 初始化 API Key 路由（需要 SSO JWT 认证）
func (r *ApiKeyRouter) InitApiKeyRouter(Router *gin.RouterGroup) {
	apiKeyRouter := Router.Group("api-keys")
	apiKeyApi := api.ApiGroupApp.ApiKeyApi
	{
		apiKeyRouter.POST("create", apiKeyApi.Create) // 创建 API Key
		apiKeyRouter.GET("list", apiKeyApi.List)      // 获取列表
		apiKeyRouter.DELETE(":id", apiKeyApi.Delete)  // 删除
	}
}
