package request

// LibraryCreate 创建库请求
type LibraryCreate struct {
	Name           string `json:"name" binding:"required"`
	Description    string `json:"description"`
	SourceType     string `json:"source_type"`     // github, website, local（默认 local）
	SourceURL      string `json:"source_url"`      // vuejs/docs 或 vuejs.org/guide
	DefaultVersion string `json:"default_version"` // 默认版本（默认 default）
}

// LibraryUpdate 更新库请求
type LibraryUpdate struct {
	ID          uint   `json:"id" binding:"required"`
	Name        string `json:"name" binding:"required"`
	Description string `json:"description"`
	SourceType  string `json:"source_type"`
	SourceURL   string `json:"source_url"`
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

// VersionCreate 创建版本请求
type VersionCreate struct {
	Version string `json:"version" binding:"required,min=1,max=50"`
}
