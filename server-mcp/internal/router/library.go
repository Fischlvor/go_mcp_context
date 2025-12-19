package router

import (
	"go-mcp-context/internal/api"

	"github.com/gin-gonic/gin"
)

type LibraryRouter struct{}

// InitLibraryPublicRouter 初始化库公开路由（无需认证）
func (l *LibraryRouter) InitLibraryPublicRouter(Router *gin.RouterGroup) {
	libraryRouter := Router.Group("libraries")
	libraryApi := api.ApiGroupApp.LibraryApi
	{
		libraryRouter.GET("", libraryApi.List)                    // 列表查询
		libraryRouter.GET(":id", libraryApi.Get)                  // 详情查询
		libraryRouter.GET(":id/versions", libraryApi.GetVersions) // 获取版本列表
	}
}

// InitLibraryRouter 初始化库私有路由（需要认证）
func (l *LibraryRouter) InitLibraryRouter(Router *gin.RouterGroup) {
	libraryRouter := Router.Group("libraries")
	libraryApi := api.ApiGroupApp.LibraryApi
	{
		libraryRouter.POST("", libraryApi.Create)                                             // 创建
		libraryRouter.PUT(":id", libraryApi.Update)                                           // 更新
		libraryRouter.DELETE(":id", libraryApi.Delete)                                        // 删除
		libraryRouter.POST(":id/versions", libraryApi.CreateVersion)                          // 创建版本
		libraryRouter.DELETE(":id/versions/:version", libraryApi.DeleteVersion)               // 删除版本
		libraryRouter.POST(":id/versions/:version/refresh", libraryApi.RefreshVersion)        // 刷新版本（异步）
		libraryRouter.POST(":id/versions/:version/refresh-sse", libraryApi.RefreshVersionSSE) // 刷新版本（SSE 实时推送）
	}
}
