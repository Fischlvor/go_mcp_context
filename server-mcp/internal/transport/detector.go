package transport

import (
	"strings"

	"github.com/gin-gonic/gin"
)

// DetectTransport 检测传输协议类型
// 根据HTTP请求的特征判断使用哪种传输协议
func DetectTransport(c *gin.Context) TransportType {
	// 1. 检测SessionID (SSE协议的特征)
	// SSE协议需要先通过GET建立连接获取SessionID，后续POST请求会携带此ID
	sessionID := c.GetHeader("MCP-Session-Id")
	if sessionID != "" {
		// 预留: 未来实现SSE协议时启用
		// return TransportSSE
		_ = sessionID // 避免未使用变量警告
	}

	// 2. 检测Accept头
	// Streamable HTTP协议会同时接受JSON和SSE
	accept := c.GetHeader("Accept")
	if accept != "" {
		// 如果Accept头包含text/event-stream，说明客户端支持流式响应
		if strings.Contains(accept, "text/event-stream") {
			return TransportStreamable
		}
	}

	// 3. 默认使用HTTP协议
	return TransportHTTP
}

// DetectFromMethod 从HTTP方法检测协议
// GET请求通常用于建立SSE连接
func DetectFromMethod(c *gin.Context) TransportType {
	if c.Request.Method == "GET" {
		// 预留: GET请求用于建立SSE连接
		return TransportSSE
	}
	return DetectTransport(c)
}
