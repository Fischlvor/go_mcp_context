package router

import (
	"go-mcp-context/internal/api"

	"github.com/gin-gonic/gin"
)

type ActivityLogRouter struct{}

// InitActivityLogPublicRouter 初始化活动日志公开路由（无需认证）
func (r *ActivityLogRouter) InitActivityLogPublicRouter(Router *gin.RouterGroup) {
	activityLogApi := api.ApiGroupApp.ActivityLogApi
	{
		Router.GET("logs", activityLogApi.List) // GET /api/v1/logs?libraryId=123
	}
}
