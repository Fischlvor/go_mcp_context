package initialize

import (
	"fmt"
	"io"
	"log"
	"os"
	"time"

	dbmodel "go-mcp-context/internal/model/database"
	"go-mcp-context/pkg/global"

	"github.com/natefinch/lumberjack"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

// InitGorm 初始化数据库连接
func InitGorm() *gorm.DB {
	pgCfg := global.Config.Postgres

	gormLogSink := io.Discard
	if pgCfg.GormLogFile != "" {
		gormLogSink = &lumberjack.Logger{
			Filename:   pgCfg.GormLogFile,
			MaxSize:    100,
			MaxBackups: 7,
			MaxAge:     30,
			Compress:   true,
		}
	}

	gormLogger := logger.New(
		log.New(gormLogSink, "", log.LstdFlags|log.Lshortfile),
		logger.Config{
			SlowThreshold:             time.Second,
			LogLevel:                  logger.Info,
			IgnoreRecordNotFoundError: true,
			Colorful:                  false,
		},
	)

	db, err := gorm.Open(postgres.Open(pgCfg.Dsn()), &gorm.Config{
		Logger: gormLogger,
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
		&dbmodel.ActivityLog{},
		&dbmodel.MCPCallLog{},
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
		USING gin(chunk_tsvector_simple)
	`
	if err := global.DB.Exec(ftsSQL).Error; err != nil {
		fmt.Printf("Warning: Could not create full-text index: %v\n", err)
	}

	// 简单配置全文索引（simple 可选 chunk_type）
	simpleSQL := `
		CREATE INDEX IF NOT EXISTS idx_chunks_text_simple_active
		ON document_chunks
		USING gin(chunk_tsvector_simple)
		WHERE status = 'active' AND deleted_at IS NULL;
	`
	if err := global.DB.Exec(simpleSQL).Error; err != nil {
		fmt.Printf("Warning: Could not create simple full-text index: %v\n", err)
	}

	// simple 配置下的 chunk_type 定向全文索引
	simpleCodeSQL := `
		CREATE INDEX IF NOT EXISTS idx_chunks_text_simple_active_code
		ON document_chunks
		USING gin(chunk_tsvector_simple)
		WHERE status = 'active' AND deleted_at IS NULL AND chunk_type = 'code';
	`
	if err := global.DB.Exec(simpleCodeSQL).Error; err != nil {
		fmt.Printf("Warning: Could not create simple full-text code index: %v\n", err)
	}

	simpleInfoSQL := `
		CREATE INDEX IF NOT EXISTS idx_chunks_text_simple_active_info
		ON document_chunks
		USING gin(chunk_tsvector_simple)
		WHERE status = 'active' AND deleted_at IS NULL AND chunk_type = 'info';
	`
	if err := global.DB.Exec(simpleInfoSQL).Error; err != nil {
		fmt.Printf("Warning: Could not create simple full-text info index: %v\n", err)
	}

	// library/version/type 过滤索引，兼顾 chunk_index 顺序
	chunkFilterSQL := `
		CREATE INDEX IF NOT EXISTS idx_chunks_library_version_type
		ON document_chunks (library_id, version, chunk_type, chunk_index)
		WHERE status = 'active' AND deleted_at IS NULL;
	`
	if err := global.DB.Exec(chunkFilterSQL).Error; err != nil {
		fmt.Printf("Warning: Could not create library/version filter index: %v\n", err)
	}

	// 活动日志 BRIN 索引（时序数据优化）
	brinSQL := `
		CREATE INDEX IF NOT EXISTS idx_logs_library_time 
		ON activity_logs 
		USING BRIN(library_id, created_at)
	`
	if err := global.DB.Exec(brinSQL).Error; err != nil {
		fmt.Printf("Warning: Could not create BRIN index for activity_logs: %v\n", err)
	}
}
