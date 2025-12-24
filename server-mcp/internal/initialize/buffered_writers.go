package initialize

import (
	"go-mcp-context/pkg/bufferedwriter/actlog"
	"go-mcp-context/pkg/bufferedwriter/mcplog"
	"go-mcp-context/pkg/bufferedwriter/stats"
	"go-mcp-context/pkg/global"
)

// InitBufferedWriters 初始化所有缓冲写入器
func InitBufferedWriters() {
	actlog.Init(global.DB) // 活动日志
	stats.Init()           // 统计系统
	mcplog.Init()          // MCP 调用日志
}

// CloseBufferedWriters 关闭所有缓冲写入器（刷新缓冲区）
func CloseBufferedWriters() {
	mcplog.Close()
	stats.Shutdown()
	actlog.Close()
}
