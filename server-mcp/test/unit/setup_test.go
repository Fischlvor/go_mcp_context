package test_test

import (
	"context"
	"fmt"
	"os"
	"testing"

	"go-mcp-context/internal/initialize"
	"go-mcp-context/pkg/config"
	"go-mcp-context/pkg/core"
	"go-mcp-context/pkg/global"

	"gopkg.in/yaml.v3"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

// TestMain 全局测试初始化
// 在所有测试运行前执行一次
func TestMain(m *testing.M) {
	// 1. 设置测试配置路径
	os.Setenv("CONFIG_PATH", "./configs/config.test.yaml")

	// 2. 初始化测试环境
	setupTestEnvironment()

	// 3. 初始化数据库表结构
	initTestTables()

	// 4. 测试前清理数据（保留表结构）
	CleanupTestData()

	// 5. 运行测试
	code := m.Run()

	// 6. 测试完成后保留数据（方便查看）
	fmt.Println("✅ Test data preserved for inspection")

	// 7. 清理资源
	cleanupTestEnvironment()

	os.Exit(code)
}

// setupTestEnvironment 初始化测试环境
func setupTestEnvironment() {
	fmt.Println("🔧 Initializing test environment...")

	// 1. 加载测试配置
	configPath := os.Getenv("CONFIG_PATH")
	if configPath == "" {
		configPath = "configs/config.test.yaml"
	}

	// 手动加载配置文件
	data, err := os.ReadFile(configPath)
	if err != nil {
		// 尝试从上一级目录读取
		data, err = os.ReadFile("../" + configPath)
		if err != nil {
			// 尝试从两级上级目录读取
			data, err = os.ReadFile("../../" + configPath)
			if err != nil {
				panic(fmt.Sprintf("❌ ERROR: Failed to read config file %s: %v", configPath, err))
			}
		}
	}

	global.Config = &config.Config{}
	if err := yaml.Unmarshal(data, global.Config); err != nil {
		panic(fmt.Sprintf("❌ ERROR: Failed to parse config file: %v", err))
	}

	// 验证配置是否正确
	if global.Config.Postgres.DBName != "mcp_context_test" {
		panic("❌ ERROR: Test config must use 'mcp_context_test' database!")
	}
	if global.Config.Redis.DB != 15 {
		panic("❌ ERROR: Test config must use Redis DB 15!")
	}

	fmt.Printf("✅ Test Database: %s\n", global.Config.Postgres.DBName)
	fmt.Printf("✅ Test Redis DB: %d\n", global.Config.Redis.DB)

	// 2. 初始化日志
	global.Log = core.InitLogger()

	// 3. 初始化数据库连接
	global.DB = initTestDatabase()

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

	fmt.Println("✅ Test environment initialized")
}

// initTestDatabase 初始化测试数据库连接
func initTestDatabase() *gorm.DB {
	pgCfg := global.Config.Postgres

	// 使用静默日志（测试时不输出 SQL）
	db, err := gorm.Open(postgres.Open(pgCfg.Dsn()), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Silent),
	})
	if err != nil {
		panic(fmt.Sprintf("Failed to connect to test database: %v", err))
	}

	// 设置连接池
	sqlDB, _ := db.DB()
	sqlDB.SetMaxIdleConns(pgCfg.MaxIdleConns)
	sqlDB.SetMaxOpenConns(pgCfg.MaxOpenConns)

	// 启用 pgvector 扩展
	if err := db.Exec("CREATE EXTENSION IF NOT EXISTS vector").Error; err != nil {
		fmt.Printf("⚠️  Warning: Failed to enable pgvector extension: %v\n", err)
	}

	return db
}

// initTestTables 初始化测试数据库表结构
func initTestTables() {
	fmt.Println("🔧 Initializing test database tables...")

	// 检查表是否已存在（尝试查询表，如果成功说明表存在）
	var count int64
	err := global.DB.Raw("SELECT COUNT(*) FROM libraries").Scan(&count).Error
	if err == nil {
		fmt.Printf("✅ Test database tables already exist (found %d libraries)\n", count)
		return
	}

	fmt.Printf("📋 Creating tables (error: %v)...\n", err)
	// 调用生产环境的 InitTables() 函数
	// ✅ 复用逻辑，确保测试环境和生产环境表结构一致
	initialize.InitTables()

	fmt.Println("✅ Test database tables initialized")
}

// cleanupTestEnvironment 清理测试环境
func cleanupTestEnvironment() {
	fmt.Println("🧹 Cleaning up test environment...")

	// 关闭 Redis 连接
	if global.Redis != nil {
		global.Redis.Close()
	}

	// 关闭数据库连接
	if global.DB != nil {
		sqlDB, _ := global.DB.DB()
		sqlDB.Close()
	}

	fmt.Println("✅ Test environment cleaned up")
}

// CleanupTestData 清理测试数据（测试结束后调用一次）
// 导出函数，供其他测试包使用
func CleanupTestData() {
	fmt.Println("🧹 Cleaning up test data before running tests...")

	if global.DB == nil {
		return
	}

	// 清空所有表（保持表结构）
	// 使用 TRUNCATE CASCADE 确保外键约束不会阻止清理
	tables := []string{
		"activity_logs",
		"mcp_call_logs",
		"statistics",
		"api_keys",
		"search_cache",
		"document_chunks",
		"document_uploads",
		"libraries",
	}

	for _, table := range tables {
		// 使用 DELETE 而不是 TRUNCATE，因为 postgres_test 不是 sequence 的 owner
		global.DB.Exec(fmt.Sprintf("DELETE FROM %s", table))
	}

	// 重置序列（如果有权限）
	global.DB.Exec("SELECT setval('libraries_id_seq', 1, false)")
	global.DB.Exec("SELECT setval('document_uploads_id_seq', 1, false)")
	global.DB.Exec("SELECT setval('document_chunks_id_seq', 1, false)")

	// 清空 Redis 测试数据库
	if global.Redis != nil {
		ctx := context.Background()
		global.Redis.FlushDB(ctx)
	}
}
