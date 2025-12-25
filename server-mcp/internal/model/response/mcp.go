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
	Title       string  `json:"title"`                 // 标题（code mode: LLM 生成, info mode: headers 层级）
	Description string  `json:"description,omitempty"` // LLM 生成的描述（仅 code mode）
	Source      string  `json:"source"`                // 来源文件路径
	Version     string  `json:"version"`               // 版本号
	Mode        string  `json:"mode"`                  // 类型：code 或 info
	Language    string  `json:"language,omitempty"`    // 代码语言（仅 code mode）
	Code        string  `json:"code,omitempty"`        // 代码内容（仅 code mode）
	Content     string  `json:"content,omitempty"`     // ChunkText 原文（仅 info mode）
	Tokens      int     `json:"tokens"`                // token 数
	Relevance   float64 `json:"relevance"`             // 相关性分数 0-1
}
