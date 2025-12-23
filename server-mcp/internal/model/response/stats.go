package response

// UserStats 用户统计数据
type UserStats struct {
	Libraries int64 `json:"libraries"` // 我的库数量
	Documents int64 `json:"documents"` // 我的文档数量
	Tokens    int64 `json:"tokens"`    // 我的 Token 总数
	MCPCalls  int64 `json:"mcp_calls"` // 我的 MCP 调用次数
}
