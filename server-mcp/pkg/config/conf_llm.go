package config

// LLM 配置
type LLM struct {
	Provider    string  `json:"provider" yaml:"provider"`       // openai
	BaseURL     string  `json:"base_url" yaml:"base_url"`       // 为空则复用 embedding.base_url
	APIKey      string  `json:"api_key" yaml:"api_key"`         // 为空则复用 embedding.api_key
	Model       string  `json:"model" yaml:"model"`             // gpt-4o-mini
	MaxTokens   int     `json:"max_tokens" yaml:"max_tokens"`   // 输出限制
	Temperature float32 `json:"temperature" yaml:"temperature"` // 温度
}
