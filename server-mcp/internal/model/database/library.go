package database

import (
	"go-mcp-context/pkg/global"
)

// Library 文档库
type Library struct {
	global.MODEL
	Name           string     `json:"name" gorm:"size:255;not null"`
	Version        string     `json:"version" gorm:"size:50;not null"`
	Description    string     `json:"description" gorm:"type:text"`
	Metadata       JSON       `json:"metadata" gorm:"type:jsonb"`
	EmbeddingModel string     `json:"embedding_model" gorm:"size:100;default:'text-embedding-3-small'"`
	Status         string     `json:"status" gorm:"size:20;default:'active'"`
	Documents      []Document `json:"documents,omitempty" gorm:"foreignKey:LibraryID"`
}

func (Library) TableName() string {
	return "libraries"
}
