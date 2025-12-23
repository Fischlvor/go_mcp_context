package database

import (
	"time"
)

// Metric 名称常量
const (
	MetricMCPGetLibraryDocs  = "mcp.func.get_library_docs"
	MetricMCPSearchLibraries = "mcp.func.search_libraries"
)

// Statistics 系统统计
type Statistics struct {
	ID          uint      `json:"id" gorm:"primaryKey"`
	LibraryID   *uint     `json:"library_id" gorm:"uniqueIndex:idx_stats_lib_metric"`
	MetricName  string    `json:"metric_name" gorm:"size:100;uniqueIndex:idx_stats_lib_metric"`
	MetricValue int64     `json:"metric_value"`
	RecordedAt  time.Time `json:"recorded_at" gorm:"index;default:now()"`
}

func (Statistics) TableName() string {
	return "statistics"
}
