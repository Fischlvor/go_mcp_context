// Package cache 提供缓存抽象和实现
//
// 包含两层接口：
//   - Cache: 基础缓存接口（Get/Set/Delete）
//   - TagAwareCache: 带 Tag 版本的缓存接口，支持 O(1) 批量失效
//
// 使用示例：
//
//	// 基础缓存
//	cache.GetOrSet(c, "key", time.Hour, fetchFunc)
//
//	// 带 Tag 的缓存（更新 tag 后自动失效）
//	cache.GetOrSetWithTags(c, "key", []string{"library:1"}, time.Hour, fetchFunc)
package cache

import (
	"errors"
	"time"
)

// ErrCacheMiss 缓存未命中错误
var ErrCacheMiss = errors.New("cache miss")

// =============================================================================
// Cache 接口
// =============================================================================

// Cache 基础缓存接口
type Cache interface {
	// Get 从缓存获取值，反序列化到 dest
	Get(key string, dest interface{}) error
	// Set 设置缓存值，ttl 为过期时间
	Set(key string, value interface{}, ttl time.Duration) error
	// Delete 删除指定 key
	Delete(key string) error
	// Exists 检查 key 是否存在
	Exists(key string) (bool, error)
	// Clear 清除指定前缀的所有缓存
	Clear(prefix string) error
	// Close 关闭缓存连接
	Close() error
}

// =============================================================================
// 工具函数
// =============================================================================

// GetOrSet 缓存旁路模式（Cache-Aside Pattern）
//
// 工作流程：
//  1. 尝试从缓存获取
//  2. 命中则直接返回
//  3. 未命中则调用 fetchFunc 获取数据
//  4. 将数据写入缓存并返回
func GetOrSet[T any](c Cache, key string, ttl time.Duration, fetchFunc func() (T, error)) (T, error) {
	var result T

	// 尝试从缓存获取
	if c != nil {
		if err := c.Get(key, &result); err == nil {
			return result, nil
		}
	}

	// 缓存未命中，调用 fetchFunc
	data, err := fetchFunc()
	if err != nil {
		return result, err
	}

	// 写入缓存
	if c != nil {
		_ = c.Set(key, data, ttl)
	}

	return data, nil
}
