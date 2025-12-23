package router

import (
	"go-mcp-context/internal/api"

	"github.com/gin-gonic/gin"
)

type StatsRouter struct{}

// InitStatsRouter 初始化统计路由（需要 SSO JWT 认证）
func (r *StatsRouter) InitStatsRouter(Router *gin.RouterGroup) {
	statsRouter := Router.Group("stats")
	statsApi := api.ApiGroupApp.StatsApi
	{
		statsRouter.GET("my", statsApi.GetMyStats) // 获取当前用户统计
	}
}
