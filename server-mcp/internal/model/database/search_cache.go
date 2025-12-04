package database

import (
	"time"

	"go-mcp-context/pkg/global"
)

// SearchCache 搜索缓存
type SearchCache struct {
	global.MODEL
	QueryHash string    `json:"query_hash" gorm:"size:64;uniqueIndex"`
	LibraryID uint      `json:"library_id" gorm:"not null;index"`
	Results   JSON      `json:"results" gorm:"type:jsonb;not null"`
	HitCount  int       `json:"hit_count" gorm:"default:0"`
	TTL       time.Time `json:"ttl" gorm:"index"`
}

func (SearchCache) TableName() string {
	return "search_cache"
}
