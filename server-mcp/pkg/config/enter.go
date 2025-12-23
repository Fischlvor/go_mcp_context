package config

// Config 总配置结构
type Config struct {
	System    System    `json:"system" yaml:"system"`
	Postgres  Postgres  `json:"postgres" yaml:"postgres"`
	Redis     Redis     `json:"redis" yaml:"redis"`
	Embedding Embedding `json:"embedding" yaml:"embedding"`
	LLM       LLM       `json:"llm" yaml:"llm"`
	Qiniu     Qiniu     `json:"qiniu" yaml:"qiniu"`
	Chunker   Chunker   `json:"chunker" yaml:"chunker"`
	Cache     Cache     `json:"cache" yaml:"cache"`
	JWT       JWT       `json:"jwt" yaml:"jwt"`
	SSO       SSO       `json:"sso" yaml:"sso"`
	Zap       Zap       `json:"zap" yaml:"zap"`
	GitHub    GitHub    `json:"github" yaml:"github"`
}

// GitHub 配置
type GitHub struct {
	Token string `json:"token" yaml:"token"` // GitHub Personal Access Token
	Proxy string `json:"proxy" yaml:"proxy"` // 代理地址（可选，如 http://10.21.71.52:7890）
}
