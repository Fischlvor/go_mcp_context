package database

import (
	"go-mcp-context/pkg/global"

	"github.com/lib/pq"
)

// Library 文档库（两层架构：Library -> DocumentChunk）
// Library 不再包含 version 字段，版本信息存储在 DocumentUpload 和 DocumentChunk 中
type Library struct {
	global.MODEL
	Name           string         `json:"name" gorm:"size:255;not null;uniqueIndex:idx_libraries_name_deleted"` // 库名唯一（支持软删除）
	Description    string         `json:"description" gorm:"type:text"`
	DefaultVersion string         `json:"default_version" gorm:"size:50"`             // 默认版本（最新版本）
	Versions       pq.StringArray `json:"versions" gorm:"type:text[]"`                // 所有可用版本列表
	SourceType     string         `json:"source_type" gorm:"size:20;default:'local'"` // github, website, local
	SourceURL      string         `json:"source_url" gorm:"size:500"`                 // vuejs/docs 或 vuejs.org/guide
	EmbeddingModel string         `json:"embedding_model" gorm:"size:100;default:'text-embedding-3-small'"`
	Status         string         `json:"status" gorm:"size:20;default:'active'"` // active, archived, deleted
	CreatedBy      string         `json:"created_by" gorm:"size:36;index"`        // 创建者 UUID，空值表示公共库
	// 关联
	Uploads []DocumentUpload `json:"uploads,omitempty" gorm:"foreignKey:LibraryID"`
	Chunks  []DocumentChunk  `json:"chunks,omitempty" gorm:"foreignKey:LibraryID"`
}

func (Library) TableName() string {
	return "libraries"
}
