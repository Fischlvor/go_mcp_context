package config

// Embedding 向量嵌入配置
type Embedding struct {
	Provider  string `json:"provider" yaml:"provider"`     // 提供商：openai, local
	BaseURL   string `json:"base_url" yaml:"base_url"`     // API Base URL（可选，用于第三方代理）
	APIKey    string `json:"api_key" yaml:"api_key"`       // OpenAI API Key
	Model     string `json:"model" yaml:"model"`           // 模型名称
	Dimension int    `json:"dimension" yaml:"dimension"`   // 向量维度
	ModelPath string `json:"model_path" yaml:"model_path"` // 本地模型路径（仅 local 使用）
}
