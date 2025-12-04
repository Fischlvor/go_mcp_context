package config

// JWT 认证配置
type JWT struct {
	AccessTokenSecret      string `json:"access_token_secret" yaml:"access_token_secret"`             // Access Token 密钥
	RefreshTokenSecret     string `json:"refresh_token_secret" yaml:"refresh_token_secret"`           // Refresh Token 密钥
	AccessTokenExpiryTime  string `json:"access_token_expiry_time" yaml:"access_token_expiry_time"`   // Access Token 过期时间
	RefreshTokenExpiryTime string `json:"refresh_token_expiry_time" yaml:"refresh_token_expiry_time"` // Refresh Token 过期时间
	APITokenExpiryTime     string `json:"api_token_expiry_time" yaml:"api_token_expiry_time"`         // API Token 过期时间
	Issuer                 string `json:"issuer" yaml:"issuer"`                                       // 签发者
}
