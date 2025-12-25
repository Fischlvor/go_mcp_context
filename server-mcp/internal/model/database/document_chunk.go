package database

import (
	"go-mcp-context/pkg/global"

	"github.com/pgvector/pgvector-go"
)

// DocumentChunk 文档块（核心检索单元）
// 两层架构中的核心表，直接关联 Library，包含版本信息
type DocumentChunk struct {
	global.MODEL
	LibraryID  uint   `json:"library_id" gorm:"not null;index:idx_chunk_library_version"`      // 关联库
	UploadID   uint   `json:"upload_id" gorm:"index"`                                          // 关联上传记录（可选）
	Version    string `json:"version" gorm:"size:50;not null;index:idx_chunk_library_version"` // 版本号
	ChunkIndex int    `json:"chunk_index" gorm:"not null"`                                     // 块序号

	// Context7 风格的 Enrich 字段
	Title       string `json:"title" gorm:"size:500"`        // LLM 生成的标题
	Description string `json:"description" gorm:"type:text"` // LLM 生成的描述
	Source      string `json:"source" gorm:"type:text"`      // 来源路径/URL
	Language    string `json:"language" gorm:"size:20"`      // 代码语言 (js, go, python, markdown)
	Code        string `json:"code" gorm:"type:text"`        // 代码块内容（如果有）

	// 原始内容、预计算 tsvector 与向量
	ChunkText           string          `json:"chunk_text" gorm:"type:text;not null"`                // 原始文本内容
	ChunkTSVectorSimple string          `json:"-" gorm:"column:chunk_tsvector_simple;type:tsvector"` // simple 配置预计算 tsvector
	Tokens              int             `json:"tokens"`                                              // token 数量
	Embedding           pgvector.Vector `json:"-" gorm:"type:vector(1536)"`                          // 向量

	// 分类和统计
	ChunkType   string `json:"chunk_type" gorm:"size:10;default:'mixed'"` // code, info, mixed
	AccessCount int    `json:"access_count" gorm:"default:0"`             // 访问次数（用于热度计算）
	Metadata    JSON   `json:"metadata" gorm:"type:jsonb"`                // 扩展元数据
	Status      string `json:"status" gorm:"size:20;default:'active'"`    // active, pending, deleted

	// 版本控制（用于无感知更新）
	BatchVersion int64 `json:"batch_version" gorm:"default:0;index"` // 批次版本号，支持原子切换

	// 关联
	Library Library `json:"-" gorm:"foreignKey:LibraryID"`
}

func (DocumentChunk) TableName() string {
	return "document_chunks"
}
