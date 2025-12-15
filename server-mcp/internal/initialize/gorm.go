package initialize

import (
	"fmt"
	"os"

	dbmodel "go-mcp-context/internal/model/database"
	"go-mcp-context/pkg/global"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

// InitGorm 初始化数据库连接
func InitGorm() *gorm.DB {
	pgCfg := global.Config.Postgres

	db, err := gorm.Open(postgres.Open(pgCfg.Dsn()), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Silent),
	})
	if err != nil {
		fmt.Printf("Failed to connect to PostgreSQL: %v\n", err)
		os.Exit(1)
	}

	// 获取底层 SQL 连接
	sqlDB, _ := db.DB()
	sqlDB.SetMaxIdleConns(pgCfg.MaxIdleConns)
	sqlDB.SetMaxOpenConns(pgCfg.MaxOpenConns)

	// 启用 pgvector 扩展
	if err := db.Exec("CREATE EXTENSION IF NOT EXISTS vector").Error; err != nil {
		fmt.Printf("Warning: Failed to enable pgvector extension: %v\n", err)
	}

	return db
}

// InitTables 初始化数据库表
func InitTables() {
	if err := global.DB.AutoMigrate(
		&dbmodel.Library{},
		&dbmodel.DocumentUpload{}, // 原 Document 改为 DocumentUpload
		&dbmodel.DocumentChunk{},
		&dbmodel.SearchCache{},
		&dbmodel.APIKey{},
		&dbmodel.Statistics{},
	); err != nil {
		fmt.Printf("Failed to migrate database: %v\n", err)
		os.Exit(1)
	}

	// 创建索引
	createIndexes()
}

// createIndexes 创建数据库索引
func createIndexes() {
	// 向量索引 (HNSW)
	indexSQL := `
		CREATE INDEX IF NOT EXISTS idx_chunks_embedding 
		ON document_chunks 
		USING hnsw (embedding vector_cosine_ops)
		WITH (m = 16, ef_construction = 64)
	`
	if err := global.DB.Exec(indexSQL).Error; err != nil {
		fmt.Printf("Warning: Could not create vector index: %v\n", err)
	}

	// 全文搜索索引
	ftsSQL := `
		CREATE INDEX IF NOT EXISTS idx_chunks_text 
		ON document_chunks 
		USING gin(to_tsvector('english', chunk_text))
	`
	if err := global.DB.Exec(ftsSQL).Error; err != nil {
		fmt.Printf("Warning: Could not create full-text index: %v\n", err)
	}
}
