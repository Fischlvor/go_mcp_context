package transport

import (
	"github.com/gin-gonic/gin"
)

// TransportType 传输协议类型
type TransportType string

const (
	// TransportHTTP 传统HTTP协议 - 直接JSON响应
	TransportHTTP TransportType = "http"

	// TransportStreamable Streamable HTTP协议 - 可以是JSON或SSE流
	TransportStreamable TransportType = "streamable"

	// TransportSSE SSE协议 - 需要预先建立连接 (预留)
	TransportSSE TransportType = "sse"
)

// RequestContext MCP请求上下文
// 包含协议无关的请求信息，用于在各层之间传递
type RequestContext struct {
	// Transport 传输协议类型
	Transport TransportType

	// SessionID 会话ID (仅SSE协议使用，当前为空)
	SessionID string

	// Method MCP方法名
	Method string

	// Params 请求参数
	Params map[string]interface{}

	// ID 请求ID
	ID interface{}

	// GinCtx Gin上下文
	GinCtx *gin.Context
}

// ResponseConfig 响应配置
// 用于控制响应的行为
type ResponseConfig struct {
	// EnableStreaming 是否启用流式响应
	EnableStreaming bool

	// BufferSize 流式响应的缓冲区大小
	BufferSize int

	// Timeout 响应超时时间(秒)
	Timeout int
}

// DefaultResponseConfig 默认响应配置
var DefaultResponseConfig = ResponseConfig{
	EnableStreaming: true,
	BufferSize:      1024,
	Timeout:         60,
}
