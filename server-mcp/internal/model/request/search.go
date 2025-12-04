package request

// Search 搜索请求
type Search struct {
	LibraryID uint   `json:"library_id" binding:"required"`
	Query     string `json:"query" binding:"required"`
	Mode      string `json:"mode"`  // code, info, 或空（全部）
	Page      int    `json:"page"`  // 页码，默认 1
	Limit     int    `json:"limit"` // 每页数量，默认 10，最大 50
}
