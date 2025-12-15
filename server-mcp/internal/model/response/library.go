package response

import "time"

// LibraryListItem 库列表项（前端主页表格，精简字段）
type LibraryListItem struct {
	ID             uint      `json:"id"`
	Name           string    `json:"name"`
	SourceType     string    `json:"source_type"`     // github, website, local
	SourceURL      string    `json:"source_url"`      // vuejs/docs
	DefaultVersion string    `json:"default_version"` // 当前版本
	TokenCount     int       `json:"token_count"`     // 对应 Context7 的 TOKENS
	ChunkCount     int       `json:"chunk_count"`     // 对应 Context7 的 SNIPPETS
	UpdatedAt      time.Time `json:"updated_at"`      // 对应 Context7 的 UPDATE
}

// LibraryInfo 库详情响应（完整信息）
type LibraryInfo struct {
	ID             uint      `json:"id"`
	Name           string    `json:"name"`
	DefaultVersion string    `json:"default_version"`
	Versions       []string  `json:"versions"`
	SourceType     string    `json:"source_type"`
	SourceURL      string    `json:"source_url"`
	Description    string    `json:"description"`
	DocumentCount  int       `json:"document_count"`
	ChunkCount     int       `json:"chunk_count"`
	TokenCount     int       `json:"token_count"`
	Status         string    `json:"status"`
	CreatedAt      time.Time `json:"created_at"`
	UpdatedAt      time.Time `json:"updated_at"`
}

// VersionInfo 版本信息（用于上传时选择）
type VersionInfo struct {
	Version     string    `json:"version" gorm:"column:version"`
	TokenCount  int       `json:"token_count" gorm:"column:token_count"`
	ChunkCount  int       `json:"chunk_count" gorm:"column:chunk_count"`
	LastUpdated time.Time `json:"last_updated" gorm:"column:last_updated"`
}
