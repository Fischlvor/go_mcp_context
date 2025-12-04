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
	DocumentID  uint    `json:"document_id"`
	LibraryID   uint    `json:"library_id"`
	Content     string  `json:"content"`
	ChunkType   string  `json:"chunk_type"`
	Score       float64 `json:"score"`        // 综合分数 0-1
	VectorScore float64 `json:"vector_score"` // 向量相似度
	BM25Score   float64 `json:"bm25_score"`   // BM25 分数
}
