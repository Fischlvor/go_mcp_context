package database

import (
	"time"

	"go-mcp-context/pkg/global"
)

// APIKey API 密钥（用于 MCP 调用认证）
type APIKey struct {
	global.MODEL
	UserUUID    string     `json:"user_uuid" gorm:"type:uuid;not null;index"` // 用户 UUID
	TokenHash   string     `json:"-" gorm:"size:64;uniqueIndex"`              // SHA256 哈希，不返回前端
	TokenSuffix string     `json:"token_suffix" gorm:"size:4;not null"`       // 后 4 位，用于显示
	Name        string     `json:"name" gorm:"size:100;not null"`             // 用户自定义名称
	UsageCount  int64      `json:"usage_count" gorm:"default:0"`              // 使用次数
	LastUsedAt  *time.Time `json:"last_used_at"`                              // 上次使用时间
}

func (APIKey) TableName() string {
	return "api_keys"
}
