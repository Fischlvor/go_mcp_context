package database

import (
	"time"
)

// ========== 事件常量 ==========

// 文档事件
const (
	EventDocUpload   = "document.upload"
	EventDocParse    = "document.parse"
	EventDocChunk    = "document.chunk"
	EventDocEmbed    = "document.embed"
	EventDocComplete = "document.complete"
	EventDocFailed   = "document.failed"
	EventDocDelete   = "document.delete"
)

// 版本事件
const (
	EventVerCreate  = "version.create"
	EventVerDelete  = "version.delete"
	EventVerRefresh = "version.refresh"
)

// GitHub 导入事件
const (
	EventGHImportStart    = "github.import.start"
	EventGHImportDownload = "github.import.download"
	EventGHImportComplete = "github.import.complete"
	EventGHImportFailed   = "github.import.failed"
)

// 库事件
const (
	EventLibCreate = "library.create"
	EventLibUpdate = "library.update"
	EventLibDelete = "library.delete"
)

// 状态常量
const (
	LogStatusInfo    = "info"
	LogStatusSuccess = "success"
	LogStatusWarning = "warning"
	LogStatusError   = "error"
)

// ========== 数据模型 ==========

// ActivityLog 活动日志
type ActivityLog struct {
	ID         uint      `json:"id" gorm:"primaryKey"`
	LibraryID  uint      `json:"library_id" gorm:"not null;index"`
	ActorID    string    `json:"actor_id,omitempty" gorm:"size:36"`
	Event      string    `json:"event" gorm:"size:64;not null;index"`
	Status     string    `json:"status" gorm:"size:16;not null;default:info"`
	Message    string    `json:"message" gorm:"type:text;not null"`
	TargetType string    `json:"target_type,omitempty" gorm:"size:32"`
	TargetID   string    `json:"target_id,omitempty" gorm:"size:64"`
	TaskID     string    `json:"task_id,omitempty" gorm:"size:26;index"`
	Version    string    `json:"version,omitempty" gorm:"size:32"`
	Metadata   JSON      `json:"metadata,omitempty" gorm:"type:jsonb;default:'{}'"`
	CreatedAt  time.Time `json:"time" gorm:"autoCreateTime"`
}

func (ActivityLog) TableName() string {
	return "activity_logs"
}
