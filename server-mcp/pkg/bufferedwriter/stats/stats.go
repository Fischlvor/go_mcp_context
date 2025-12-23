// Package stats 提供异步批量统计记录功能
package stats

import (
	"fmt"
	"log"
	"sync"
	"time"

	"go-mcp-context/pkg/bufferedwriter"
	"go-mcp-context/pkg/global"

	"gorm.io/gorm"
)

// MetricEntry 统计条目
type MetricEntry struct {
	LibraryID  *uint
	MetricName string
	Delta      int64
}

// DefaultConfig 默认配置（5条或5秒触发）
var DefaultConfig = bufferedwriter.Config{
	Size:     100,
	Batch:    5,
	Interval: 5 * time.Second,
}

// DBWriter 数据库写入器
type DBWriter struct {
	db *gorm.DB
}

// NewDBWriter 创建数据库写入器
func NewDBWriter(db *gorm.DB) *DBWriter {
	return &DBWriter{db: db}
}

// WriteBatch 批量写入统计
func (w *DBWriter) WriteBatch(batch []*MetricEntry) error {
	if len(batch) == 0 {
		return nil
	}

	// 聚合相同 metric
	aggregated := make(map[string]*MetricEntry)
	for _, entry := range batch {
		key := entry.MetricName
		if entry.LibraryID != nil {
			key += fmt.Sprintf("_%d", *entry.LibraryID)
		}
		if existing, ok := aggregated[key]; ok {
			existing.Delta += entry.Delta
		} else {
			aggregated[key] = &MetricEntry{
				LibraryID:  entry.LibraryID,
				MetricName: entry.MetricName,
				Delta:      entry.Delta,
			}
		}
	}

	// 批量 upsert
	for _, entry := range aggregated {
		if err := w.upsertMetric(entry); err != nil {
			log.Printf("[stats] Failed to upsert metric %s: %v", entry.MetricName, err)
		}
	}
	return nil
}

// Close 关闭写入器
func (w *DBWriter) Close() error {
	return nil
}

// upsertMetric 更新或插入统计
func (w *DBWriter) upsertMetric(entry *MetricEntry) error {
	// PostgreSQL upsert
	// 使用 COALESCE 处理 NULL，将 NULL 转为 0 用于唯一约束匹配
	sql := `
		INSERT INTO statistics (library_id, metric_name, metric_value, recorded_at)
		VALUES ($1, $2, $3, NOW())
		ON CONFLICT (COALESCE(library_id, 0), metric_name) 
		DO UPDATE SET metric_value = statistics.metric_value + EXCLUDED.metric_value, recorded_at = NOW()
	`
	return w.db.Exec(sql, entry.LibraryID, entry.MetricName, entry.Delta).Error
}

var (
	defaultBuffer *bufferedwriter.Buffer[*MetricEntry]
	once          sync.Once
)

// Init 初始化统计系统
func Init() {
	once.Do(func() {
		writer := NewDBWriter(global.DB)
		defaultBuffer = bufferedwriter.New("stats", writer, DefaultConfig)
	})
}

// Increment 增加统计（使用默认缓冲区）
func Increment(metricName string, delta int64) {
	if defaultBuffer == nil {
		return
	}
	defaultBuffer.Write(&MetricEntry{
		MetricName: metricName,
		Delta:      delta,
	})
}

// IncrementWithLibrary 增加库相关统计
func IncrementWithLibrary(libraryID uint, metricName string, delta int64) {
	if defaultBuffer == nil {
		return
	}
	defaultBuffer.Write(&MetricEntry{
		LibraryID:  &libraryID,
		MetricName: metricName,
		Delta:      delta,
	})
}

// Shutdown 关闭统计系统
func Shutdown() error {
	if defaultBuffer != nil {
		return defaultBuffer.Close()
	}
	return nil
}
