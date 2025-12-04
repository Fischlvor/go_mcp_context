package core

import (
	"log"

	"go-mcp-context/pkg/config"
	"go-mcp-context/pkg/utils"

	"gopkg.in/yaml.v3"
)

// InitConf 初始化配置
func InitConf() *config.Config {
	yamlConf, err := utils.LoadYAML()
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	var cfg config.Config
	if err := yaml.Unmarshal(yamlConf, &cfg); err != nil {
		log.Fatalf("Failed to unmarshal config: %v", err)
	}

	return &cfg
}
