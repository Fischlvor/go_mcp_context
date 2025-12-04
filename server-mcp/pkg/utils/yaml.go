package utils

import (
	"os"
)

// getConfigFile 根据环境变量获取配置文件路径
func getConfigFile() string {
	env := os.Getenv("APP_ENV")
	if env == "prod" || env == "production" {
		return "configs/config.prod.yaml"
	}
	return "configs/config.yaml"
}

// LoadYAML 从文件中读取 YAML 数据并返回字节数组
func LoadYAML() ([]byte, error) {
	configFile := getConfigFile()
	return os.ReadFile(configFile)
}
