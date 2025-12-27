package transport

import (
	"go-mcp-context/internal/transport/http"
	"go-mcp-context/internal/transport/streamable"

	"github.com/gin-gonic/gin"
)

// CreateResponseWriter 创建响应写入器
// 根据传输协议类型创建对应的ResponseWriter实现
func CreateResponseWriter(c *gin.Context, transportType TransportType) ResponseWriter {
	switch transportType {
	case TransportHTTP:
		// HTTP协议: 直接返回JSON
		return http.NewHTTPResponseWriter(c)

	case TransportStreamable:
		// Streamable HTTP协议: 可以返回JSON或SSE流
		return streamable.NewStreamableResponseWriter(c)

	case TransportSSE:
		// SSE协议: 通过预先建立的连接推送响应 (预留)
		// 未来实现时的代码示例:
		// sessionID := c.GetHeader("MCP-Session-Id")
		// return sse.NewSSEResponseWriter(c, sessionID)

		// 当前暂不支持，降级到HTTP
		return http.NewHTTPResponseWriter(c)

	default:
		// 未知协议，使用HTTP作为默认
		return http.NewHTTPResponseWriter(c)
	}
}

// CreateResponseWriterWithConfig 创建带配置的响应写入器
// 允许自定义响应配置
func CreateResponseWriterWithConfig(c *gin.Context, transportType TransportType, config ResponseConfig) ResponseWriter {
	// 当前简化实现，未来可以根据config调整写入器行为
	_ = config // 避免未使用变量警告
	return CreateResponseWriter(c, transportType)
}
