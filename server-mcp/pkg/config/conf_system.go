package config

// System 系统配置
type System struct {
	Host         string `json:"host" yaml:"host"`                   // 服务器地址
	Port         int    `json:"port" yaml:"port"`                   // 服务器端口
	Env          string `json:"env" yaml:"env"`                     // 环境：debug, release, test
	RouterPrefix string `json:"router_prefix" yaml:"router_prefix"` // 路由前缀
}
