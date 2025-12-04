package config

// Cache 缓存配置
type Cache struct {
	TTL    string `json:"ttl" yaml:"ttl"`       // 缓存过期时间，如 "24h"
	Prefix string `json:"prefix" yaml:"prefix"` // 缓存键前缀
}
