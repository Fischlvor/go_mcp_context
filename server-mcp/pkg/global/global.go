package global

import (
	"io"

	"go-mcp-context/pkg/cache"
	"go-mcp-context/pkg/config"
	"go-mcp-context/pkg/embedding"
	"go-mcp-context/pkg/llm"
	"go-mcp-context/pkg/storage"

	"github.com/redis/go-redis/v9"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

var (
	Config    *config.Config
	Log       *zap.Logger
	LogWriter io.Writer // 全局日志写入器，供 GORM 等使用
	DB        *gorm.DB
	Redis     *redis.Client       // Redis 客户端（用于 CachedEmbeddingService）
	Cache     cache.TagAwareCache // 带 Tag 版本的缓存接口（用于搜索结果缓存等）
	Embedding embedding.EmbeddingService
	Storage   storage.Storage // 文件存储服务
	LLM       llm.LLMService  // LLM 服务
)
