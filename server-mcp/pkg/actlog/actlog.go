// Package actlog 提供异步批量活动日志记录功能
//
// 使用示例:
//
//	// 初始化（在 main.go 中）
//	actlog.Init(global.DB)
//	defer actlog.Close()
//
//	// 记录日志（业务代码中）
//	actlog.Log(libraryID, actlog.EventDocUpload, actlog.StatusInfo, "开始上传: README.md",
//	    actlog.WithTaskID(taskID),
//	    actlog.WithTarget("document", "123"),
//	)
package actlog

import (
	"log"
	"sync"
	"time"

	"gorm.io/gorm"
)

var (
	defaultBuffer *Buffer
	initOnce      sync.Once
	mu            sync.RWMutex
)

// Init 初始化活动日志系统
func Init(db *gorm.DB) {
	initOnce.Do(func() {
		writer := NewDBWriter(db)
		defaultBuffer = NewBuffer(writer, DefaultBufferConfig)
		log.Println("[actlog] Activity logger initialized")
	})
}

// InitWithWriter 使用自定义 Writer 初始化
func InitWithWriter(writer Writer, config BufferConfig) {
	initOnce.Do(func() {
		defaultBuffer = NewBuffer(writer, config)
		log.Println("[actlog] Activity logger initialized with custom writer")
	})
}

// Close 关闭活动日志系统
func Close() {
	mu.Lock()
	defer mu.Unlock()

	if defaultBuffer != nil {
		if err := defaultBuffer.Close(); err != nil {
			log.Printf("[actlog] Failed to close: %v", err)
		}
		defaultBuffer = nil
		log.Println("[actlog] Activity logger closed")
	}
}

// Log 记录活动日志
func Log(libraryID uint, event, status, message string, opts ...Option) {
	mu.RLock()
	buf := defaultBuffer
	mu.RUnlock()

	if buf == nil {
		log.Printf("[actlog] Not initialized, dropping log: %s", message)
		return
	}

	// 应用选项
	options := &Options{}
	for _, opt := range opts {
		opt(options)
	}

	entry := &LogEntry{
		LibraryID:  libraryID,
		ActorID:    options.ActorID,
		Event:      event,
		Status:     status,
		Message:    message,
		TargetType: options.TargetType,
		TargetID:   options.TargetID,
		TaskID:     options.TaskID,
		Version:    options.Version,
		Metadata:   options.Metadata,
		CreatedAt:  time.Now(),
	}

	buf.Write(entry)
}

// Info 记录 info 级别日志
func Info(libraryID uint, event, message string, opts ...Option) {
	Log(libraryID, event, StatusInfo, message, opts...)
}

// Success 记录 success 级别日志
func Success(libraryID uint, event, message string, opts ...Option) {
	Log(libraryID, event, StatusSuccess, message, opts...)
}

// Warning 记录 warning 级别日志
func Warning(libraryID uint, event, message string, opts ...Option) {
	Log(libraryID, event, StatusWarning, message, opts...)
}

// Error 记录 error 级别日志
func Error(libraryID uint, event, message string, opts ...Option) {
	Log(libraryID, event, StatusError, message, opts...)
}

// ============================================================================
// TaskLogger: 任务级别日志器，预填充公共字段
// ============================================================================

// TaskLogger 任务日志器，预填充 libraryID, taskID, version 等公共字段
type TaskLogger struct {
	libraryID  uint
	taskID     string
	version    string
	targetType string
	targetID   string
	actorID    string
}

// NewTaskLogger 创建任务日志器
func NewTaskLogger(libraryID uint, taskID, version string) *TaskLogger {
	return &TaskLogger{
		libraryID: libraryID,
		taskID:    taskID,
		version:   version,
	}
}

// WithTarget 返回带目标信息的新日志器（不修改原实例）
func (l *TaskLogger) WithTarget(targetType, targetID string) *TaskLogger {
	return &TaskLogger{
		libraryID:  l.libraryID,
		taskID:     l.taskID,
		version:    l.version,
		targetType: targetType,
		targetID:   targetID,
		actorID:    l.actorID,
	}
}

// WithActor 返回带操作者的新日志器（不修改原实例）
func (l *TaskLogger) WithActor(actorID string) *TaskLogger {
	return &TaskLogger{
		libraryID:  l.libraryID,
		taskID:     l.taskID,
		version:    l.version,
		targetType: l.targetType,
		targetID:   l.targetID,
		actorID:    actorID,
	}
}

// buildOpts 构建选项列表
func (l *TaskLogger) buildOpts(extraOpts ...Option) []Option {
	opts := make([]Option, 0, 5+len(extraOpts))
	if l.taskID != "" {
		opts = append(opts, WithTaskID(l.taskID))
	}
	if l.version != "" {
		opts = append(opts, WithVersion(l.version))
	}
	if l.targetType != "" {
		opts = append(opts, WithTarget(l.targetType, l.targetID))
	}
	if l.actorID != "" {
		opts = append(opts, WithActor(l.actorID))
	}
	opts = append(opts, extraOpts...)
	return opts
}

// Log 记录日志
func (l *TaskLogger) Log(event, status, message string, opts ...Option) {
	Log(l.libraryID, event, status, message, l.buildOpts(opts...)...)
}

// Info 记录 info 级别日志
func (l *TaskLogger) Info(event, message string, opts ...Option) {
	Info(l.libraryID, event, message, l.buildOpts(opts...)...)
}

// Success 记录 success 级别日志
func (l *TaskLogger) Success(event, message string, opts ...Option) {
	Success(l.libraryID, event, message, l.buildOpts(opts...)...)
}

// Warning 记录 warning 级别日志
func (l *TaskLogger) Warning(event, message string, opts ...Option) {
	Warning(l.libraryID, event, message, l.buildOpts(opts...)...)
}

// Error 记录 error 级别日志
func (l *TaskLogger) Error(event, message string, opts ...Option) {
	Error(l.libraryID, event, message, l.buildOpts(opts...)...)
}
