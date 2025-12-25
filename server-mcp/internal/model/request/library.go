package request

// LibraryCreate 创建库请求（Local 类型）
type LibraryCreate struct {
	Name        string `json:"name" binding:"required"`
	Description string `json:"description"`
	CreatedBy   string `json:"-"` // 创建者 UUID（从 JWT 获取，不从请求体读取）
}

// LibraryUpdate 更新库请求（ID 从 URL 路径获取）
type LibraryUpdate struct {
	Name        string `json:"name" binding:"required"`
	Description string `json:"description"`
}

// LibraryList 库列表请求
type LibraryList struct {
	Name   *string `json:"name" form:"name"`
	Status *string `json:"status" form:"status"`
	Sort   *string `json:"sort" form:"sort"` // popular, recent（默认 recent）
	PageInfo
}

// LibraryDelete 删除库请求
type LibraryDelete struct {
	IDs []uint `json:"ids" binding:"required"`
}

// VersionCreate 创建版本请求
type VersionCreate struct {
	Version string `json:"version" binding:"required,min=1,max=50"`
}
