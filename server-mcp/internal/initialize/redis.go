package initialize

import (
	"context"
	"fmt"
	"os"

	"go-mcp-context/pkg/cache"
	"go-mcp-context/pkg/global"

	"github.com/redis/go-redis/v9"
)

// ConnectRedis 初始化 Redis 连接
func ConnectRedis() *redis.Client {
	redisCfg := global.Config.Redis

	client := redis.NewClient(&redis.Options{
		Addr:     redisCfg.Address,
		Password: redisCfg.Password,
		DB:       redisCfg.DB,
	})

	// 测试连接
	ctx := context.Background()
	if _, err := client.Ping(ctx).Result(); err != nil {
		fmt.Printf("Failed to connect to Redis: %v\n", err)
		os.Exit(1)
	}

	return client
}

// InitCache 初始化带 Tag 版本的缓存服务（复用全局 Redis 客户端）
func InitCache() cache.TagAwareCache {
	if global.Redis == nil {
		fmt.Println("Redis client not initialized, please call ConnectRedis first")
		os.Exit(1)
	}

	// 使用已有的 Redis 客户端创建缓存
	return cache.NewRedisCacheWithClient(global.Redis, "mcp:")
}
