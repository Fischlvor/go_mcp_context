package router

import (
	"go-mcp-context/internal/api"

	"github.com/gin-gonic/gin"
)

type DocumentRouter struct{}

// InitDocumentPublicRouter 初始化文档公开路由（无需认证）
func (d *DocumentRouter) InitDocumentPublicRouter(Router *gin.RouterGroup) {
	documentRouter := Router.Group("documents")
	documentApi := api.ApiGroupApp.DocumentApi
	{
		documentRouter.GET(":id", documentApi.Get) // 查询文档详情
	}
}

// InitDocumentRouter 初始化文档私有路由（需要认证）
func (d *DocumentRouter) InitDocumentRouter(Router *gin.RouterGroup) {
	documentRouter := Router.Group("documents")
	documentApi := api.ApiGroupApp.DocumentApi
	{
		documentRouter.POST("upload", documentApi.Upload) // 上传
		documentRouter.DELETE(":id", documentApi.Delete)  // 删除
	}
}
