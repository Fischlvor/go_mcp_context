package initialize

import (
	"go-mcp-context/pkg/global"
	"go-mcp-context/pkg/llm"
)

// InitLLM 初始化 LLM 服务
func InitLLM() {
	cfg := global.Config.LLM

	// 如果 LLM 配置为空，尝试复用 Embedding 配置
	apiKey := cfg.APIKey
	baseURL := cfg.BaseURL
	if apiKey == "" {
		apiKey = global.Config.Embedding.APIKey
	}
	if baseURL == "" {
		baseURL = global.Config.Embedding.BaseURL
	}

	if apiKey == "" {
		global.Log.Warn("LLM service not configured, skipping initialization")
		return
	}

	global.LLM = llm.NewOpenAILLM(llm.OpenAIConfig{
		APIKey:      apiKey,
		BaseURL:     baseURL,
		Model:       cfg.Model,
		MaxTokens:   cfg.MaxTokens,
		Temperature: cfg.Temperature,
	})

	global.Log.Info("LLM service initialized successfully")
}
