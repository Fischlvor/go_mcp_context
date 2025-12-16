package config

// Qiniu 七牛云存储配置
type Qiniu struct {
	AccessKey     string `json:"access_key" yaml:"access_key"`
	SecretKey     string `json:"secret_key" yaml:"secret_key"`
	Bucket        string `json:"bucket" yaml:"bucket"`
	Domain        string `json:"domain" yaml:"domain"` // CDN 域名
	Zone          string `json:"zone" yaml:"zone"`     // z0, z1, z2, na0, as0
	UseHTTPS      bool   `json:"use_https" yaml:"use_https"`
	UseCdnDomains bool   `json:"use_cdn_domains" yaml:"use_cdn_domains"`
	PathPrefix    string `json:"path_prefix" yaml:"path_prefix"` // 存储路径前缀，如 mcp/docs
}
