package config

// Config 总配置结构
type Config struct {
	System    System    `json:"system" yaml:"system"`
	Postgres  Postgres  `json:"postgres" yaml:"postgres"`
	Redis     Redis     `json:"redis" yaml:"redis"`
	Embedding Embedding `json:"embedding" yaml:"embedding"`
	Chunker   Chunker   `json:"chunker" yaml:"chunker"`
	Cache     Cache     `json:"cache" yaml:"cache"`
	JWT       JWT       `json:"jwt" yaml:"jwt"`
	SSO       SSO       `json:"sso" yaml:"sso"`
	Zap       Zap       `json:"zap" yaml:"zap"`
}
