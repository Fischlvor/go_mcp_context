package config

// Zap 日志配置
type Zap struct {
	Level          string `json:"level" yaml:"level"`                       // 日志级别：debug, info, warn, error
	Filename       string `json:"filename" yaml:"filename"`                 // 日志文件路径
	MaxSize        int    `json:"max_size" yaml:"max_size"`                 // 单个日志文件最大大小（MB）
	MaxBackups     int    `json:"max_backups" yaml:"max_backups"`           // 保留的旧日志文件数量
	MaxAge         int    `json:"max_age" yaml:"max_age"`                   // 日志文件保留天数
	IsConsolePrint bool   `json:"is_console_print" yaml:"is_console_print"` // 是否同时输出到控制台
}
