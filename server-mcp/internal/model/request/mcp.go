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
	LibraryID string `json:"libraryID"`
	Topic     string `json:"topic"`
	Mode      string `json:"mode"` // code, info
	Page      int    `json:"page"` // 1-10
}
