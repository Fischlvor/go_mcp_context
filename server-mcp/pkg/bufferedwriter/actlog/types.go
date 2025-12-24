// Package actlog 提供异步批量活动日志记录功能
package actlog

import "time"

// LogEntry 日志条目
type LogEntry struct {
	LibraryID  uint                   `json:"library_id"`
	ActorID    string                 `json:"actor_id,omitempty"`
	Event      string                 `json:"event"`
	Status     string                 `json:"status"`
	Message    string                 `json:"message"`
	TargetType string                 `json:"target_type,omitempty"`
	TargetID   string                 `json:"target_id,omitempty"`
	TaskID     string                 `json:"task_id,omitempty"`
	Version    string                 `json:"version,omitempty"`
	Metadata   map[string]interface{} `json:"metadata,omitempty"`
	CreatedAt  time.Time              `json:"created_at"`
}

// Options 日志选项
type Options struct {
	ActorID    string
	TargetType string
	TargetID   string
	TaskID     string
	Version    string
	Metadata   map[string]interface{}
}

// Option 函数式选项
type Option func(*Options)

// WithActor 设置操作者
func WithActor(actorID string) Option {
	return func(o *Options) {
		o.ActorID = actorID
	}
}

// WithTarget 设置操作目标
func WithTarget(targetType, targetID string) Option {
	return func(o *Options) {
		o.TargetType = targetType
		o.TargetID = targetID
	}
}

// WithTaskID 设置任务ID
func WithTaskID(taskID string) Option {
	return func(o *Options) {
		o.TaskID = taskID
	}
}

// WithVersion 设置版本
func WithVersion(version string) Option {
	return func(o *Options) {
		o.Version = version
	}
}

// WithMetadata 设置元数据
func WithMetadata(metadata map[string]interface{}) Option {
	return func(o *Options) {
		o.Metadata = metadata
	}
}

// 日志状态常量
const (
	StatusStart   = "start" // 任务开始
	StatusInfo    = "info"
	StatusSuccess = "success"
	StatusWarning = "warning"
	StatusError   = "error"
)

// 事件类型常量
const (
	// 文档事件
	EventDocUpload   = "document.upload"
	EventDocParse    = "document.parse"
	EventDocChunk    = "document.chunk"
	EventDocEnrich   = "document.enrich"
	EventDocEmbed    = "document.embed"
	EventDocComplete = "document.complete"
	EventDocFailed   = "document.failed"
	EventDocDelete   = "document.delete"

	// 版本事件
	EventVerCreate  = "version.create"
	EventVerDelete  = "version.delete"
	EventVerRefresh = "version.refresh"

	// GitHub 导入事件
	EventGHImportStart    = "github.import.start"
	EventGHImportDownload = "github.import.download"
	EventGHImportComplete = "github.import.complete"
	EventGHImportFailed   = "github.import.failed"

	// 库事件
	EventLibCreate = "library.create"
	EventLibUpdate = "library.update"
	EventLibDelete = "library.delete"
)
