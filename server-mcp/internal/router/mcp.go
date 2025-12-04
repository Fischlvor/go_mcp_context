package router

import (
	"go-mcp-context/internal/api"

	"github.com/gin-gonic/gin"
)

type MCPRouter struct{}

func (m *MCPRouter) InitMCPRouter(Router *gin.RouterGroup) {
	mcpRouter := Router.Group("mcp")
	mcpApi := api.ApiGroupApp.MCPApi
	{
		mcpRouter.GET("health", mcpApi.Health)
		mcpRouter.GET("tools", mcpApi.ListTools)
		mcpRouter.POST("", mcpApi.HandleRequest)
	}
}
