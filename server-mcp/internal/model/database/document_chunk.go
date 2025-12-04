package database

import (
	"go-mcp-context/pkg/global"

	"github.com/pgvector/pgvector-go"
)

// DocumentChunk 文档块
type DocumentChunk struct {
	global.MODEL
	DocumentID  uint            `json:"document_id" gorm:"not null;index"`
	LibraryID   uint            `json:"library_id" gorm:"not null;index"`
	ChunkIndex  int             `json:"chunk_index" gorm:"not null"`
	ChunkText   string          `json:"chunk_text" gorm:"type:text;not null"`
	Tokens      int             `json:"tokens"`
	Embedding   pgvector.Vector `json:"-" gorm:"type:vector(1536)"`
	ChunkType   string          `json:"chunk_type" gorm:"size:10;default:'mixed'"`
	AccessCount int             `json:"access_count" gorm:"default:0"`
	Metadata    JSON            `json:"metadata" gorm:"type:jsonb"`
	Status      string          `json:"status" gorm:"size:20;default:'active'"`
}

func (DocumentChunk) TableName() string {
	return "document_chunks"
}
