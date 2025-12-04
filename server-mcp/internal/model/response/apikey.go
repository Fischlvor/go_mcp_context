package response

import "time"

// APIKeyCreateResponse 创建 API Key 响应（仅创建时返回完整 key）
type APIKeyCreateResponse struct {
	ID          uint      `json:"id"`
	Name        string    `json:"name"`
	APIKey      string    `json:"api_key"`      // 完整 key，仅此一次显示
	TokenSuffix string    `json:"token_suffix"` // 后 4 位
	CreatedAt   time.Time `json:"created_at"`
}

// APIKeyListItem 列表项（脱敏显示）
type APIKeyListItem struct {
	ID          uint       `json:"id"`
	Name        string     `json:"name"`
	TokenSuffix string     `json:"token_suffix"` // 后 4 位
	LastUsedAt  *time.Time `json:"last_used_at"`
	CreatedAt   time.Time  `json:"created_at"`
}
