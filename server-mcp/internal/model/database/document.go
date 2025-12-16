package database

import (
	"go-mcp-context/pkg/global"
)

// DocumentUpload 文档上传记录（记录上传历史，不参与检索）
// 两层架构中的辅助表，用于追踪文档来源
type DocumentUpload struct {
	global.MODEL
	LibraryID    uint    `json:"library_id" gorm:"not null;index"`
	Version      string  `json:"version" gorm:"size:50;not null;index"` // 文档版本
	Title        string  `json:"title" gorm:"size:500"`
	FilePath     string  `json:"file_path" gorm:"type:text;not null"` // 存储路径（Key）
	FileType     string  `json:"file_type" gorm:"size:50"`            // md, pdf, docx, swagger
	FileSize     int64   `json:"file_size"`
	ContentHash  string  `json:"content_hash" gorm:"size:64;index"` // 文件内容哈希，用于去重
	ChunkCount   int     `json:"chunk_count" gorm:"default:0"`      // 生成的 chunk 数量
	TokenCount   int     `json:"token_count" gorm:"default:0"`      // 总 token 数
	ErrorMessage string  `json:"error_message,omitempty" gorm:"type:text"`
	Status       string  `json:"status" gorm:"size:20;default:'pending'"` // pending, processing, completed, failed, deleted
	Library      Library `json:"-" gorm:"foreignKey:LibraryID"`
}

func (DocumentUpload) TableName() string {
	return "document_uploads"
}
