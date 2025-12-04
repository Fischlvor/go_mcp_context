package initialize

import (
	"fmt"
	"os"

	"go-mcp-context/pkg/global"

	"github.com/go-redis/redis"
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
	if _, err := client.Ping().Result(); err != nil {
		fmt.Printf("Failed to connect to Redis: %v\n", err)
		os.Exit(1)
	}

	return client
}
