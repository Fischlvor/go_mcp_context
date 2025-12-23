package service

import (
	"fmt"
	"time"

	dbmodel "go-mcp-context/internal/model/database"
	"go-mcp-context/internal/model/response"
	"go-mcp-context/pkg/cache"
	"go-mcp-context/pkg/global"
)

const (
	// 用户统计缓存 TTL
	userStatsCacheTTL = 5 * time.Minute
)

type StatsService struct{}

// GetUserStats 获取用户统计数据（带缓存）
func (s *StatsService) GetUserStats(userUUID string) (*response.UserStats, error) {
	cacheKey := fmt.Sprintf("user_stats:%s", userUUID)

	return cache.GetOrSet(global.Cache, cacheKey, userStatsCacheTTL, func() (*response.UserStats, error) {
		return s.fetchUserStats(userUUID)
	})
}

// fetchUserStats 从数据库获取用户统计数据
func (s *StatsService) fetchUserStats(userUUID string) (*response.UserStats, error) {
	var result response.UserStats

	// 1. 我的库数量
	if err := global.DB.Table("libraries").
		Where("created_by = ? AND deleted_at IS NULL", userUUID).
		Count(&result.Libraries).Error; err != nil {
		return nil, err
	}

	// 2. 我的文档数量（通过 library 关联）
	if err := global.DB.Table("document_uploads").
		Joins("JOIN libraries ON libraries.id = document_uploads.library_id").
		Where("libraries.created_by = ? AND libraries.deleted_at IS NULL AND document_uploads.deleted_at IS NULL", userUUID).
		Count(&result.Documents).Error; err != nil {
		return nil, err
	}

	// 3. 我的 Token 总数（通过 library 关联）
	var tokenSum struct {
		Total int64
	}
	if err := global.DB.Table("document_chunks").
		Select("COALESCE(SUM(tokens), 0) as total").
		Joins("JOIN libraries ON libraries.id = document_chunks.library_id").
		Where("libraries.created_by = ? AND libraries.deleted_at IS NULL AND document_chunks.status = 'active'", userUUID).
		Scan(&tokenSum).Error; err != nil {
		return nil, err
	}
	result.Tokens = tokenSum.Total

	// 4. 我的 MCP 调用次数（通过 statistics 表）
	var callSum struct {
		Total int64
	}
	if err := global.DB.Table("statistics").
		Select("COALESCE(SUM(metric_value), 0) as total").
		Joins("JOIN libraries ON libraries.id = statistics.library_id").
		Where("libraries.created_by = ? AND statistics.metric_name = ?", userUUID, dbmodel.MetricMCPGetLibraryDocs).
		Scan(&callSum).Error; err != nil {
		return nil, err
	}
	result.MCPCalls = callSum.Total

	return &result, nil
}
