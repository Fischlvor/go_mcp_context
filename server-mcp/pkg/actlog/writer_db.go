package actlog

import (
	"go-mcp-context/internal/model/database"

	"gorm.io/gorm"
)

// DBWriter 数据库写入器
type DBWriter struct {
	db *gorm.DB
}

// NewDBWriter 创建数据库写入器
func NewDBWriter(db *gorm.DB) *DBWriter {
	return &DBWriter{db: db}
}

// Write 写入单条日志
func (w *DBWriter) Write(entry *LogEntry) error {
	log := w.toDBModel(entry)
	return w.db.Create(log).Error
}

// WriteBatch 批量写入日志
func (w *DBWriter) WriteBatch(entries []*LogEntry) error {
	if len(entries) == 0 {
		return nil
	}

	logs := make([]*database.ActivityLog, len(entries))
	for i, entry := range entries {
		logs[i] = w.toDBModel(entry)
	}

	return w.db.CreateInBatches(logs, len(logs)).Error
}

// Close 关闭写入器
func (w *DBWriter) Close() error {
	return nil
}

// toDBModel 转换为数据库模型
func (w *DBWriter) toDBModel(entry *LogEntry) *database.ActivityLog {
	var metadata database.JSON
	if entry.Metadata != nil {
		metadata = database.JSON(entry.Metadata)
	}

	return &database.ActivityLog{
		LibraryID:  entry.LibraryID,
		ActorID:    entry.ActorID,
		Event:      entry.Event,
		Status:     entry.Status,
		Message:    entry.Message,
		TargetType: entry.TargetType,
		TargetID:   entry.TargetID,
		TaskID:     entry.TaskID,
		Version:    entry.Version,
		Metadata:   metadata,
		CreatedAt:  entry.CreatedAt,
	}
}
