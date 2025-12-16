package request

// DocumentUpload 文档上传请求（form-data）
type DocumentUpload struct {
	LibraryID uint   `form:"library_id" binding:"required"`
	Version   string `form:"version"` // 文档版本，默认为 "latest"
	// file 通过 FormFile 获取
}

// DocumentList 文档列表请求
type DocumentList struct {
	LibraryID *uint   `json:"library_id" form:"library_id" binding:"required"`
	Version   *string `json:"version" form:"version"` // 版本过滤（可选）
	PageInfo
}

// DocumentDelete 删除文档请求
type DocumentDelete struct {
	IDs []uint `json:"ids" binding:"required"`
}
