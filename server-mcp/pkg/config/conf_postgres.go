package config

import (
	"fmt"
	"strings"

	"gorm.io/gorm/logger"
)

// Postgres 数据库配置
type Postgres struct {
	Host         string `json:"host" yaml:"host"`                     // 数据库服务器地址
	Port         int    `json:"port" yaml:"port"`                     // 数据库服务器端口
	Config       string `json:"config" yaml:"config"`                 // 连接配置参数
	DBName       string `json:"db_name" yaml:"db_name"`               // 数据库名称
	Username     string `json:"username" yaml:"username"`             // 用户名
	Password     string `json:"password" yaml:"password"`             // 密码
	MaxIdleConns int    `json:"max_idle_conns" yaml:"max_idle_conns"` // 最大空闲连接数
	MaxOpenConns int    `json:"max_open_conns" yaml:"max_open_conns"` // 最大打开连接数
	LogMode      string `json:"log_mode" yaml:"log_mode"`             // 日志模式
	GormLogFile  string `json:"gorm_log_file" yaml:"gorm_log_file"`   // GORM 日志输出文件
}

// Dsn 返回数据库连接字符串
func (p Postgres) Dsn() string {
	return fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s %s",
		p.Host, p.Port, p.Username, p.Password, p.DBName, p.Config)
}

// LogLevel 返回 GORM 日志级别
func (p Postgres) LogLevel() logger.LogLevel {
	switch strings.ToLower(p.LogMode) {
	case "silent":
		return logger.Silent
	case "error":
		return logger.Error
	case "warn":
		return logger.Warn
	case "info":
		return logger.Info
	default:
		return logger.Info
	}
}
