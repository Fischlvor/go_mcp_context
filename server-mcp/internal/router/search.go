package router

import (
	// "go-mcp-context/internal/api"

	"github.com/gin-gonic/gin"
)

type SearchRouter struct{}

// InitSearchPublicRouter 初始化搜索公开路由（无需认证）
func (s *SearchRouter) InitSearchPublicRouter(Router *gin.RouterGroup) {
	// searchRouter := Router.Group("search")
	// searchApi := api.ApiGroupApp.SearchApi
	{
		// 暂时禁用搜索API，前端使用 /documents/chunks API 实现搜索功能
		// searchRouter.POST("", searchApi.Search) // 搜索文档
	}
}
