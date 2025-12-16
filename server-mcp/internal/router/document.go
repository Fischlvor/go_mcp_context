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
		documentRouter.GET("list", documentApi.List)                              // 查询文档列表
		documentRouter.GET("detail/:id", documentApi.Get)                         // 查询文档详情
		documentRouter.GET("chunks/:mode/:libid/*version", documentApi.GetChunks) // 获取库的文档块 (mode: code/info, version 可选)
	}
}

// InitDocumentRouter 初始化文档私有路由（需要认证）
func (d *DocumentRouter) InitDocumentRouter(Router *gin.RouterGroup) {
	documentRouter := Router.Group("documents")
	documentApi := api.ApiGroupApp.DocumentApi
	{
		documentRouter.POST("upload", documentApi.Upload)            // 上传（普通）
		documentRouter.POST("upload-sse", documentApi.UploadWithSSE) // 上传（SSE 实时状态）
		documentRouter.DELETE(":id", documentApi.Delete)             // 删除
	}
}
