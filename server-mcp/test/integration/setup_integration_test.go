package integration_test

import (
	"fmt"
	"os"
	"testing"

	"go-mcp-context/internal/initialize"
	"go-mcp-context/pkg/config"
	"go-mcp-context/pkg/core"
	"go-mcp-context/pkg/global"

	"gopkg.in/yaml.v3"
)

// TestMain 集成测试的入口点
// 集成测试使用真实的外部服务（GitHub API、OpenAI API 等）
func TestMain(m *testing.M) {
	fmt.Println("🚀 Starting integration tests...")
	fmt.Println("⚠️  Integration tests will use real external services")

	// 1. 加载测试配置
	setupIntegrationConfig()

	// 2. 初始化日志
	global.Log = core.InitLogger()

	// 3. 初始化数据库
	global.DB = initialize.InitGorm()
	initialize.InitTables()

	// 4. 初始化 Redis
	global.Redis = initialize.ConnectRedis()

	// 5. 初始化缓存
	global.Cache = initialize.InitCache()

	// 6. 初始化 Embedding（使用真实的 OpenAI API）
	global.Embedding = initialize.InitEmbedding()

	// 7. 初始化存储服务（使用真实的 Qiniu）
	initialize.InitStorage()

	// 8. 初始化 LLM 服务（使用真实的 OpenAI API）
	initialize.InitLLM()

	fmt.Println("✅ Integration test environment initialized")

	// 运行测试
	code := m.Run()

	// 清理
	fmt.Println("🧹 Cleaning up integration test environment...")
	if global.DB != nil {
		sqlDB, _ := global.DB.DB()
		if sqlDB != nil {
			sqlDB.Close()
		}
	}
	if global.Redis != nil {
		global.Redis.Close()
	}

	fmt.Println("✅ Integration test environment cleaned up")
	os.Exit(code)
}

// setupIntegrationConfig 加载集成测试配置
func setupIntegrationConfig() {
	configPath := os.Getenv("CONFIG_PATH")
	if configPath == "" {
		configPath = "configs/config.test.yaml"
	}

	data, err := os.ReadFile(configPath)
	if err != nil {
		// 尝试从上一级目录读取
		data, err = os.ReadFile("../../" + configPath)
		if err != nil {
			panic(fmt.Sprintf("❌ ERROR: Failed to read config file %s: %v", configPath, err))
		}
	}

	global.Config = &config.Config{}
	if err := yaml.Unmarshal(data, global.Config); err != nil {
		panic(fmt.Sprintf("❌ ERROR: Failed to parse config file: %v", err))
	}

	// 验证配置
	if global.Config.Postgres.DBName != "mcp_context_test" {
		panic("❌ ERROR: Integration test must use 'mcp_context_test' database!")
	}

	fmt.Printf("✅ Integration Test Database: %s\n", global.Config.Postgres.DBName)
	fmt.Printf("✅ Integration Test Redis DB: %d\n", global.Config.Redis.DB)
}
