package config

// SSO 单点登录配置
type SSO struct {
	ServiceURL     string `json:"service_url" yaml:"service_url"`         // SSO 后端 API 地址
	WebURL         string `json:"web_url" yaml:"web_url"`                 // SSO 前端登录页面地址
	ClientID       string `json:"client_id" yaml:"client_id"`             // 应用 ID
	ClientSecret   string `json:"client_secret" yaml:"client_secret"`     // 应用密钥
	CallbackURL    string `json:"callback_url" yaml:"callback_url"`       // 回调地址
	PublicKeyPath  string `json:"public_key_path" yaml:"public_key_path"` // JWT 公钥路径
	SessionsSecret string `json:"sessions_secret" yaml:"sessions_secret"` // Session 加密密钥
}
