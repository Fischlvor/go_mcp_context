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
	ChunkID    uint    `json:"chunk_id"`
	DocumentID uint    `json:"document_id"`
	LibraryID  uint    `json:"library_id"`
	Title      string  `json:"title"`     // 从 Metadata 取最深层级标题
	Source     string  `json:"source"`    // Document.Title
	Content    string  `json:"content"`   // ChunkText 原文
	Tokens     int     `json:"tokens"`    // token 数
	Relevance  float64 `json:"relevance"` // 最终相关性分数 0-1
}
