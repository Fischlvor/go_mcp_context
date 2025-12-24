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
	LibraryID      uint     `json:"libraryId"`      // 库的数据库 ID
	Name           string   `json:"name"`           // 库名
	Versions       []string `json:"versions"`       // 所有版本
	DefaultVersion string   `json:"defaultVersion"` // 默认版本
	Description    string   `json:"description"`    // 描述
	Snippets       int      `json:"snippets"`       // 文档片段数
	Score          float64  `json:"score"`          // 匹配分数
}

// MCPGetLibraryDocsResult get-library-docs 结果
type MCPGetLibraryDocsResult struct {
	LibraryID uint               `json:"libraryId"` // 库的数据库 ID
	Documents []MCPDocumentChunk `json:"documents"`
	Page      int                `json:"page"`
	HasMore   bool               `json:"hasMore"`
}

// MCPDocumentChunk 文档片段
type MCPDocumentChunk struct {
	Title     string  `json:"title"`     // 标题（从 Metadata 提取）
	Source    string  `json:"source"`    // 来源文档标题
	Content   string  `json:"content"`   // 内容
	Tokens    int     `json:"tokens"`    // token 数
	Relevance float64 `json:"relevance"` // 相关性分数 0-1
}
