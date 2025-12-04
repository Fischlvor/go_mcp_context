package initialize

import (
	"go-mcp-context/pkg/embedding"
	"go-mcp-context/pkg/global"
)

// InitEmbedding 初始化 Embedding 服务
func InitEmbedding() embedding.EmbeddingService {
	cfg := global.Config.Embedding

	if cfg.BaseURL != "" {
		// 使用第三方代理
		return embedding.NewOpenAIProxyEmbedding(cfg.APIKey, cfg.BaseURL, cfg.Model, cfg.Dimension)
	}
	// 使用官方 API
	return embedding.NewOpenAIEmbedding(cfg.APIKey, cfg.Model, cfg.Dimension)
}
