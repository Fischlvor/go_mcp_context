package cache

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"
)

// RedisCache implements Cache using Redis
type RedisCache struct {
	client *redis.Client
	prefix string
	ctx    context.Context
}

// NewRedisCache creates a new Redis cache
func NewRedisCache(host string, port int, password string, db int, prefix string) (*RedisCache, error) {
	client := redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%d", host, port),
		Password: password,
		DB:       db,
	})

	ctx := context.Background()

	// Test connection
	if err := client.Ping(ctx).Err(); err != nil {
		return nil, fmt.Errorf("failed to connect to Redis: %w", err)
	}

	return &RedisCache{
		client: client,
		prefix: prefix,
		ctx:    ctx,
	}, nil
}

// Get retrieves a value from cache
func (c *RedisCache) Get(key string, dest interface{}) error {
	fullKey := c.prefix + key
	data, err := c.client.Get(c.ctx, fullKey).Bytes()
	if err == redis.Nil {
		return ErrCacheMiss
	}
	if err != nil {
		return err
	}

	return json.Unmarshal(data, dest)
}

// Set stores a value in cache with TTL
func (c *RedisCache) Set(key string, value interface{}, ttl time.Duration) error {
	fullKey := c.prefix + key
	data, err := json.Marshal(value)
	if err != nil {
		return err
	}

	return c.client.Set(c.ctx, fullKey, data, ttl).Err()
}

// Delete removes a key from cache
func (c *RedisCache) Delete(key string) error {
	fullKey := c.prefix + key
	return c.client.Del(c.ctx, fullKey).Err()
}

// Exists checks if a key exists in cache
func (c *RedisCache) Exists(key string) (bool, error) {
	fullKey := c.prefix + key
	n, err := c.client.Exists(c.ctx, fullKey).Result()
	if err != nil {
		return false, err
	}
	return n > 0, nil
}

// Clear removes all keys with the given prefix
func (c *RedisCache) Clear(prefix string) error {
	fullPrefix := c.prefix + prefix + "*"
	iter := c.client.Scan(c.ctx, 0, fullPrefix, 0).Iterator()

	for iter.Next(c.ctx) {
		if err := c.client.Del(c.ctx, iter.Val()).Err(); err != nil {
			return err
		}
	}

	return iter.Err()
}

// Close closes the Redis connection
func (c *RedisCache) Close() error {
	return c.client.Close()
}

// AddToBlacklist adds a token to the blacklist (for JWT revocation)
func (c *RedisCache) AddToBlacklist(tokenID string, ttl time.Duration) error {
	key := "blacklist:" + tokenID
	return c.Set(key, true, ttl)
}

// IsBlacklisted checks if a token is blacklisted
func (c *RedisCache) IsBlacklisted(tokenID string) (bool, error) {
	key := "blacklist:" + tokenID
	return c.Exists(key)
}
