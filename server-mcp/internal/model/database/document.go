package database

import (
	"go-mcp-context/pkg/global"
)

// Document 文档
type Document struct {
	global.MODEL
	LibraryID    uint            `json:"library_id" gorm:"not null;index"`
	Title        string          `json:"title" gorm:"size:500"`
	FilePath     string          `json:"file_path" gorm:"type:text;not null"`
	FileType     string          `json:"file_type" gorm:"size:50"`
	FileSize     int64           `json:"file_size"`
	ContentHash  string          `json:"content_hash" gorm:"size:64;index"`
	ChunkCount   int             `json:"chunk_count" gorm:"default:0"`
	TokenCount   int             `json:"token_count" gorm:"default:0"`
	ErrorMessage string          `json:"error_message,omitempty" gorm:"type:text"`
	Status       string          `json:"status" gorm:"size:20;default:'active'"`
	Library      Library         `json:"-" gorm:"foreignKey:LibraryID"`
	Chunks       []DocumentChunk `json:"chunks,omitempty" gorm:"foreignKey:DocumentID"`
}

func (Document) TableName() string {
	return "documents"
}
