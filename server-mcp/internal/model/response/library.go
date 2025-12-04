package response

// LibraryInfo 库信息响应
type LibraryInfo struct {
	ID            uint   `json:"id"`
	Name          string `json:"name"`
	Version       string `json:"version"`
	Description   string `json:"description"`
	DocumentCount int    `json:"document_count"`
	ChunkCount    int    `json:"chunk_count"`
	Status        string `json:"status"`
}
