package cache

import (
	"fmt"
	"time"
)

// =============================================================================
// TagAwareCache 接口
// =============================================================================

// TagAwareCache 带 Tag 版本的缓存接口
//
// 通过 Version Tag 模式实现 O(1) 缓存失效：
//   - 每个 tag 维护一个版本号
//   - 缓存 key 包含 tag 版本号
//   - 失效时只需递增 tag 版本号，旧缓存自动失效
//
// 示例：
//
//	tag = "library:6:v1.0"
//	key = "search:xxx:tv:3"  // tv:3 表示 tag version = 3
//
//	// 失效时
//	cache.InvalidateTags(["library:6:v1.0"])  // tag version 变成 4
//	// 下次构建 key 时变成 "search:xxx:tv:4"，旧缓存自动失效
type TagAwareCache interface {
	Cache

	// GetTagVersion 获取 tag 的当前版本号（不存在返回 0）
	GetTagVersion(tag string) (int64, error)

	// InvalidateTags 递增 tag 版本号，使关联的缓存失效
	InvalidateTags(tags []string) error
}

// =============================================================================
// 工具函数
// =============================================================================

// BuildTaggedKey 构建带 tag 版本的缓存 key
//
// 格式: {baseKey}:tv:{version1}:{version2}:...
//
// 示例：
//
//	baseKey = "search:topic:6:v1.0:code:abc123"
//	tags = ["library:6:v1.0"]
//	tag version = 3
//	result = "search:topic:6:v1.0:code:abc123:tv:3"
func BuildTaggedKey(c TagAwareCache, baseKey string, tags []string) (string, error) {
	if len(tags) == 0 {
		return baseKey, nil
	}

	key := baseKey + ":tv"
	for _, tag := range tags {
		version, err := c.GetTagVersion(tag)
		if err != nil {
			return "", err
		}
		key += fmt.Sprintf(":%d", version)
	}
	return key, nil
}

// GetOrSetWithTags 带 Tag 版本的缓存旁路模式
//
// 与 GetOrSet 的区别：
//   - 缓存 key 包含 tag 版本号
//   - 当 tag 失效（版本号递增）后，旧缓存自动失效
//
// 参数：
//   - baseKey: 基础缓存 key
//   - tags: 关联的 tag 列表（用于构建带版本的 key）
//   - ttl: 缓存过期时间
//   - fetchFunc: 缓存未命中时的数据获取函数
func GetOrSetWithTags[T any](c TagAwareCache, baseKey string, tags []string, ttl time.Duration, fetchFunc func() (T, error)) (T, error) {
	var result T

	if c == nil {
		return fetchFunc()
	}

	// 构建带 tag version 的 key
	key, err := BuildTaggedKey(c, baseKey, tags)
	if err != nil {
		return fetchFunc()
	}

	// 尝试从缓存获取
	if err := c.Get(key, &result); err == nil {
		return result, nil
	}

	// 缓存未命中，调用 fetchFunc
	data, err := fetchFunc()
	if err != nil {
		return result, err
	}

	// 写入缓存
	_ = c.Set(key, data, ttl)

	return data, nil
}
