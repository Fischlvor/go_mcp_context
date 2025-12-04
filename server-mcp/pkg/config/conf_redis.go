package config

// Redis 缓存数据库配置
type Redis struct {
	Address  string `json:"address" yaml:"address"`   // Redis 服务器地址，如 "localhost:6379"
	Password string `json:"password" yaml:"password"` // 连接密码，如果没有设置则留空
	DB       int    `json:"db" yaml:"db"`             // 数据库索引，默认为 0
}
