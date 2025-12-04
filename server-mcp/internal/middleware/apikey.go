package middleware

import (
	"go-mcp-context/internal/model/response"
	"go-mcp-context/internal/service"
	"go-mcp-context/pkg/global"

	"github.com/gin-gonic/gin"
	"github.com/gofrs/uuid"
	"go.uber.org/zap"
)

var apiKeyService = service.ServiceGroupApp.ApiKeyService

// APIKeyAuth API Key 认证中间件（用于 MCP 调用）
func APIKeyAuth() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 从 MCP_API_KEY header 获取 API Key
		apiKey := c.GetHeader("MCP_API_KEY")
		if apiKey == "" {
			response.NoAuth("未提供 API Key", c)
			c.Abort()
			return
		}

		// 验证 API Key
		userUUID, err := apiKeyService.ValidateAPIKey(apiKey)
		if err != nil {
			global.Log.Warn("API Key 验证失败",
				zap.String("error", err.Error()),
				zap.String("key_prefix", safeKeyPrefix(apiKey)),
			)
			response.NoAuth("无效的 API Key", c)
			c.Abort()
			return
		}

		// 将用户 UUID 存入上下文
		parsedUUID, _ := uuid.FromString(userUUID)
		c.Set("user_uuid", parsedUUID)

		c.Next()
	}
}

// safeKeyPrefix 安全地获取 API Key 前缀用于日志
func safeKeyPrefix(apiKey string) string {
	if len(apiKey) > 10 {
		return apiKey[:10] + "..."
	}
	return "***"
}
