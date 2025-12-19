package cache

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"
)

// =============================================================================
// RedisCache
// =============================================================================

// RedisCache Redis 缓存实现，同时实现 Cache 和 TagAwareCache 接口
type RedisCache struct {
	client *redis.Client
	prefix string
	ctx    context.Context
}

// 编译时检查接口实现
var _ TagAwareCache = (*RedisCache)(nil)

// NewRedisCache 创建 Redis 缓存（新建连接）
func NewRedisCache(host string, port int, password string, db int, prefix string) (*RedisCache, error) {
	client := redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%d", host, port),
		Password: password,
		DB:       db,
	})

	ctx := context.Background()
	if err := client.Ping(ctx).Err(); err != nil {
		return nil, fmt.Errorf("failed to connect to Redis: %w", err)
	}

	return &RedisCache{client: client, prefix: prefix, ctx: ctx}, nil
}

// NewRedisCacheWithClient 创建 Redis 缓存（复用已有连接）
func NewRedisCacheWithClient(client *redis.Client, prefix string) *RedisCache {
	return &RedisCache{client: client, prefix: prefix, ctx: context.Background()}
}

// -----------------------------------------------------------------------------
// Cache 接口实现
// -----------------------------------------------------------------------------

// Get 从缓存获取值，反序列化到 dest
func (c *RedisCache) Get(key string, dest interface{}) error {
	data, err := c.client.Get(c.ctx, c.prefix+key).Bytes()
	if err == redis.Nil {
		return ErrCacheMiss
	}
	if err != nil {
		return err
	}
	return json.Unmarshal(data, dest)
}

// Set 设置缓存值（带 TTL）
func (c *RedisCache) Set(key string, value interface{}, ttl time.Duration) error {
	data, err := json.Marshal(value)
	if err != nil {
		return err
	}
	return c.client.Set(c.ctx, c.prefix+key, data, ttl).Err()
}

// Delete 删除指定 key
func (c *RedisCache) Delete(key string) error {
	return c.client.Del(c.ctx, c.prefix+key).Err()
}

// Exists 检查 key 是否存在
func (c *RedisCache) Exists(key string) (bool, error) {
	n, err := c.client.Exists(c.ctx, c.prefix+key).Result()
	return n > 0, err
}

// Clear 清除指定前缀的所有缓存（使用 SCAN 避免阻塞）
func (c *RedisCache) Clear(prefix string) error {
	iter := c.client.Scan(c.ctx, 0, c.prefix+prefix+"*", 100).Iterator()
	for iter.Next(c.ctx) {
		if err := c.client.Del(c.ctx, iter.Val()).Err(); err != nil {
			return err
		}
	}
	return iter.Err()
}

// Close 关闭 Redis 连接
func (c *RedisCache) Close() error {
	return c.client.Close()
}

// -----------------------------------------------------------------------------
// TagAwareCache 接口实现
// -----------------------------------------------------------------------------

const tagVersionKeyPrefix = "tag:version:"

// GetTagVersion 获取 tag 版本号（不存在返回 0）
func (c *RedisCache) GetTagVersion(tag string) (int64, error) {
	version, err := c.client.Get(c.ctx, c.prefix+tagVersionKeyPrefix+tag).Int64()
	if err == redis.Nil {
		return 0, nil
	}
	return version, err
}

// InvalidateTags 递增 tag 版本号，使关联缓存失效
func (c *RedisCache) InvalidateTags(tags []string) error {
	for _, tag := range tags {
		if err := c.client.Incr(c.ctx, c.prefix+tagVersionKeyPrefix+tag).Err(); err != nil {
			return err
		}
	}
	return nil
}

// -----------------------------------------------------------------------------
// 业务扩展方法
// -----------------------------------------------------------------------------

// AddToBlacklist 添加 token 到黑名单（用于 JWT 撤销）
func (c *RedisCache) AddToBlacklist(tokenID string, ttl time.Duration) error {
	return c.Set("blacklist:"+tokenID, true, ttl)
}

// IsBlacklisted 检查 token 是否在黑名单
func (c *RedisCache) IsBlacklisted(tokenID string) (bool, error) {
	return c.Exists("blacklist:" + tokenID)
}
