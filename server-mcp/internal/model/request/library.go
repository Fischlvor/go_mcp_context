package request

// LibraryCreate 创建库请求
type LibraryCreate struct {
	Name        string                 `json:"name" binding:"required"`
	Version     string                 `json:"version" binding:"required"`
	Description string                 `json:"description"`
	Metadata    map[string]interface{} `json:"metadata"`
}

// LibraryUpdate 更新库请求
type LibraryUpdate struct {
	ID          uint                   `json:"id" binding:"required"`
	Name        string                 `json:"name" binding:"required"`
	Version     string                 `json:"version" binding:"required"`
	Description string                 `json:"description"`
	Metadata    map[string]interface{} `json:"metadata"`
}

// LibraryList 库列表请求
type LibraryList struct {
	Name   *string `json:"name" form:"name"`
	Status *string `json:"status" form:"status"`
	PageInfo
}

// LibraryDelete 删除库请求
type LibraryDelete struct {
	IDs []uint `json:"ids" binding:"required"`
}
