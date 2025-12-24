package request

// MCPRequest JSON-RPC 2.0 请求
type MCPRequest struct {
	JSONRPC string                 `json:"jsonrpc"`
	ID      interface{}            `json:"id"`
	Method  string                 `json:"method"`
	Params  map[string]interface{} `json:"params"`
}

// MCPSearchLibraries search-libraries 工具参数
type MCPSearchLibraries struct {
	LibraryName string `json:"libraryName"`
}

// MCPGetLibraryDocs get-library-docs 工具参数
type MCPGetLibraryDocs struct {
	LibraryID uint   `json:"libraryId"` // 库的数据库 ID
	Version   string `json:"version"`   // 版本（可选，默认使用 defaultVersion）
	Topic     string `json:"topic"`
	Mode      string `json:"mode"` // code, info
	Page      int    `json:"page"` // 1-10
}
