package database

import (
	"time"
)

// Statistics 系统统计
type Statistics struct {
	ID          uint      `json:"id" gorm:"primaryKey"`
	LibraryID   *uint     `json:"library_id" gorm:"index"`
	MetricName  string    `json:"metric_name" gorm:"size:100;index"`
	MetricValue int64     `json:"metric_value"`
	RecordedAt  time.Time `json:"recorded_at" gorm:"index;default:now()"`
}

func (Statistics) TableName() string {
	return "statistics"
}
