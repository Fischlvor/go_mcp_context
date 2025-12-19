package response

// SearchResult 搜索结果
type SearchResult struct {
	Results []SearchResultItem `json:"results"`
	Total   int64              `json:"total"`
	Page    int                `json:"page"`
	Limit   int                `json:"limit"`
	HasMore bool               `json:"hasMore"`
}

// SearchResultItem 搜索结果项
type SearchResultItem struct {
	ChunkID     uint    `json:"chunk_id"`
	UploadID    uint    `json:"upload_id"`
	LibraryID   uint    `json:"library_id"`
	Version     string  `json:"version"`     // 文档版本
	Title       string  `json:"title"`       // LLM 生成的标题（code mode）或 headers 层级（info mode）
	Description string  `json:"description"` // LLM 生成的描述（code mode），info mode 为空
	Source      string  `json:"source"`      // 文件来源路径
	Language    string  `json:"language"`    // 代码语言（code mode），info mode 为空
	Code        string  `json:"code"`        // 代码内容（code mode），info mode 为空
	Content     string  `json:"content"`     // ChunkText 原文
	Tokens      int     `json:"tokens"`      // token 数
	Relevance   float64 `json:"relevance"`   // 最终相关性分数 0-1
}
