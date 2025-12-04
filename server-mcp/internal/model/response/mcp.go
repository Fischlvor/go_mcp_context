package response

// MCPResponse JSON-RPC 2.0 响应
type MCPResponse struct {
	JSONRPC string      `json:"jsonrpc"`
	ID      interface{} `json:"id"`
	Result  interface{} `json:"result,omitempty"`
	Error   *MCPError   `json:"error,omitempty"`
}

// MCPError JSON-RPC 2.0 错误
type MCPError struct {
	Code    int         `json:"code"`
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
}

// MCPToolDefinition MCP 工具定义
type MCPToolDefinition struct {
	Name        string                 `json:"name"`
	Description string                 `json:"description"`
	InputSchema map[string]interface{} `json:"inputSchema"`
}

// MCPSearchLibrariesResult search-libraries 结果
type MCPSearchLibrariesResult struct {
	Libraries []MCPLibraryInfo `json:"libraries"`
}

// MCPLibraryInfo 库信息
type MCPLibraryInfo struct {
	ID          string  `json:"id"` // 格式：name/version
	Name        string  `json:"name"`
	Version     string  `json:"version"`
	Description string  `json:"description"`
	Snippets    int     `json:"snippets"` // 文档片段数
	Score       float64 `json:"score"`    // 匹配分数
}

// MCPGetLibraryDocsResult get-library-docs 结果
type MCPGetLibraryDocsResult struct {
	Documents []MCPDocumentChunk `json:"documents"`
	Page      int                `json:"page"`
	HasMore   bool               `json:"hasMore"`
}

// MCPDocumentChunk 文档片段
type MCPDocumentChunk struct {
	Content   string  `json:"content"`
	ChunkType string  `json:"chunkType"` // code, info, mixed
	Score     float64 `json:"score"`     // 相关性分数 0-1
}
