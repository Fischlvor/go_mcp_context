package initialize

import (
	"go-mcp-context/pkg/embedding"
	"go-mcp-context/pkg/global"
)

// InitEmbedding 初始化 Embedding 服务（带 Redis 缓存）
func InitEmbedding() embedding.EmbeddingService {
	cfg := global.Config.Embedding

	var inner embedding.EmbeddingService
	if cfg.BaseURL != "" {
		// 使用第三方代理
		inner = embedding.NewOpenAIProxyEmbedding(cfg.APIKey, cfg.BaseURL, cfg.Model, cfg.Dimension)
	} else {
		// 使用官方 API
		inner = embedding.NewOpenAIEmbedding(cfg.APIKey, cfg.Model, cfg.Dimension)
	}

	// 用 CachedEmbeddingService 包装，自动缓存查询 embedding
	// 注意：此时 global.Redis 和 global.Log 可能还未初始化，需要在 main.go 中调整顺序
	return embedding.NewCachedEmbeddingService(inner, global.Redis, global.Log)
}
