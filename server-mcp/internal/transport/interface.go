package transport

import (
	"go-mcp-context/internal/model/response"
)

// ResponseWriter 响应写入器接口
// 所有传输协议的响应写入器都必须实现此接口
type ResponseWriter interface {
	// WriteResponse 写入成功响应
	WriteResponse(resp *response.MCPResponse) error

	// WriteError 写入错误响应
	WriteError(err *response.MCPError, id interface{}) error

	// Close 关闭写入器，释放资源
	Close() error
}

// ConnectionManager 连接管理器接口 (预留给SSE协议使用)
// 用于管理多个客户端的SSE连接
type ConnectionManager interface {
	// Register 注册新连接
	Register(conn Connection) error

	// Unregister 注销连接
	Unregister(sessionID string) error

	// GetConnection 根据SessionID获取连接
	GetConnection(sessionID string) (Connection, error)

	// SendToSession 向指定Session发送消息
	SendToSession(sessionID string, data interface{}) error

	// BroadcastAll 广播消息给所有连接
	BroadcastAll(data interface{}) error

	// CleanupExpired 清理过期连接
	CleanupExpired() error
}

// Connection 连接接口 (预留给SSE协议使用)
// 表示单个客户端的SSE连接
type Connection interface {
	// SessionID 获取连接的SessionID
	SessionID() string

	// Send 发送消息到此连接
	Send(data interface{}) error

	// Close 关闭连接
	Close() error

	// IsAlive 检查连接是否存活
	IsAlive() bool
}
