package global

import (
	"io"

	"go-mcp-context/pkg/config"
	"go-mcp-context/pkg/embedding"

	"github.com/go-redis/redis"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

var (
	Config    *config.Config
	Log       *zap.Logger
	LogWriter io.Writer // 全局日志写入器，供 GORM 等使用
	DB        *gorm.DB
	Redis     *redis.Client
	Embedding embedding.EmbeddingService
)
