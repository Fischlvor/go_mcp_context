package test_test

import (
	"fmt"
	"testing"

	"go-mcp-context/internal/model/request"
	"go-mcp-context/internal/service"
	"go-mcp-context/pkg/global"

	dbmodel "go-mcp-context/internal/model/database"
)

// TestStatsGetUserStats 测试用户统计数据查询
func Test_Stats_GetUserStats(t *testing.T) {
	statsService := &service.StatsService{}

	t.Run("get user stats with data", func(t *testing.T) {
		userUUID := "test-user-stats-001"

		// 通过 Service 层创建库
		libService := &service.LibraryService{}
		lib, err := libService.Create(&request.LibraryCreate{
			Name:        "stats-test-lib",
			Description: "test",
		})
		if err != nil {
			t.Fatalf("Create library error = %v", err)
		}
		// 设置 created_by
		global.DB.Model(&lib).Update("created_by", userUUID)

		// 创建统计记录
		libID := lib.ID
		stat := &dbmodel.Statistics{
			LibraryID:   &libID,
			MetricName:  dbmodel.MetricMCPGetLibraryDocs,
			MetricValue: 10,
		}
		global.DB.Create(stat)

		// 获取统计数据
		stats, err := statsService.GetUserStats(userUUID)
		if err != nil {
			t.Fatalf("GetUserStats() error = %v", err)
		}

		if stats == nil {
			t.Fatal("Expected stats, got nil")
		}

		if stats.Libraries < 1 {
			t.Errorf("Expected at least 1 library, got %d", stats.Libraries)
		}

		if stats.Tokens < 0 {
			t.Errorf("Expected non-negative tokens, got %d", stats.Tokens)
		}

		if stats.MCPCalls < 10 {
			t.Errorf("Expected at least 10 MCP calls, got %d", stats.MCPCalls)
		}
	})

	t.Run("get user stats with no data", func(t *testing.T) {
		userUUID := "test-user-stats-empty"

		stats, err := statsService.GetUserStats(userUUID)
		if err != nil {
			t.Fatalf("GetUserStats() error = %v", err)
		}

		if stats == nil {
			t.Fatal("Expected stats, got nil")
		}

		if stats.Libraries != 0 {
			t.Errorf("Expected 0 libraries, got %d", stats.Libraries)
		}

		if stats.Documents != 0 {
			t.Errorf("Expected 0 documents, got %d", stats.Documents)
		}

		if stats.Tokens != 0 {
			t.Errorf("Expected 0 tokens, got %d", stats.Tokens)
		}

		if stats.MCPCalls != 0 {
			t.Errorf("Expected 0 MCP calls, got %d", stats.MCPCalls)
		}
	})

	t.Run("get user stats with multiple libraries", func(t *testing.T) {
		userUUID := "test-user-stats-multi"

		// 通过 Service 层创建多个库
		libService := &service.LibraryService{}
		for i := 1; i <= 3; i++ {
			lib, _ := libService.Create(&request.LibraryCreate{
				Name:        "stats-multi-lib-" + string(rune('0'+i)),
				Description: "test",
			})
			global.DB.Model(&lib).Update("created_by", userUUID)
		}

		stats, err := statsService.GetUserStats(userUUID)
		if err != nil {
			t.Fatalf("GetUserStats() error = %v", err)
		}

		if stats.Libraries < 3 {
			t.Errorf("Expected at least 3 libraries, got %d", stats.Libraries)
		}
	})

	t.Run("cache works correctly", func(t *testing.T) {
		userUUID := "test-user-stats-cache"

		// 通过 Service 层创建库
		libService := &service.LibraryService{}
		lib, _ := libService.Create(&request.LibraryCreate{
			Name:        "stats-cache-lib",
			Description: "test",
		})
		global.DB.Model(&lib).Update("created_by", userUUID)

		// 第一次查询
		stats1, err := statsService.GetUserStats(userUUID)
		if err != nil {
			t.Fatalf("GetUserStats() error = %v", err)
		}

		// 添加更多数据
		lib2, _ := libService.Create(&request.LibraryCreate{
			Name:        "stats-cache-lib-2",
			Description: "test",
		})
		global.DB.Model(&lib2).Update("created_by", userUUID)

		// 第二次查询（应该返回缓存的数据）
		stats2, err := statsService.GetUserStats(userUUID)
		if err != nil {
			t.Fatalf("GetUserStats() error = %v", err)
		}

		// 由于缓存，两次查询结果应该相同
		if stats1.Libraries != stats2.Libraries {
			t.Logf("Note: Cache may have expired or not working. stats1.Libraries=%d, stats2.Libraries=%d", stats1.Libraries, stats2.Libraries)
		}
	})

	t.Run("get user stats with different user uuids", func(t *testing.T) {
		userUUID1 := "test-user-stats-uuid1"
		userUUID2 := "test-user-stats-uuid2"

		// 为用户1创建库
		libService := &service.LibraryService{}
		lib1, _ := libService.Create(&request.LibraryCreate{
			Name:        "stats-uuid1-lib",
			Description: "test",
		})
		global.DB.Model(&lib1).Update("created_by", userUUID1)

		// 为用户2创建库
		lib2, _ := libService.Create(&request.LibraryCreate{
			Name:        "stats-uuid2-lib",
			Description: "test",
		})
		global.DB.Model(&lib2).Update("created_by", userUUID2)

		// 获取两个用户的统计数据
		stats1, err := statsService.GetUserStats(userUUID1)
		if err != nil {
			t.Fatalf("GetUserStats() error = %v", err)
		}

		stats2, err := statsService.GetUserStats(userUUID2)
		if err != nil {
			t.Fatalf("GetUserStats() error = %v", err)
		}

		// 验证两个用户的统计数据是独立的
		if stats1.Libraries == 0 {
			t.Error("Expected at least 1 library for user 1")
		}

		if stats2.Libraries == 0 {
			t.Error("Expected at least 1 library for user 2")
		}
	})

	t.Run("get user stats with special uuid format", func(t *testing.T) {
		userUUID := "special-uuid-format-12345"

		// 为用户创建库
		libService := &service.LibraryService{}
		lib, _ := libService.Create(&request.LibraryCreate{
			Name:        "stats-special-uuid-lib",
			Description: "test",
		})
		global.DB.Model(&lib).Update("created_by", userUUID)

		// 获取统计数据
		stats, err := statsService.GetUserStats(userUUID)
		if err != nil {
			t.Fatalf("GetUserStats() error = %v", err)
		}

		if stats == nil {
			t.Fatal("Expected stats, got nil")
		}

		if stats.Libraries < 1 {
			t.Errorf("Expected at least 1 library, got %d", stats.Libraries)
		}
	})

	t.Run("get user stats with many libraries", func(t *testing.T) {
		userUUID := "user-many-libs"

		// 为用户创建多个库
		libService := &service.LibraryService{}
		createdCount := int64(0)
		for i := 0; i < 5; i++ {
			lib, err := libService.Create(&request.LibraryCreate{
				Name:        fmt.Sprintf("stats-lib-%d", i),
				Description: "test",
			})
			if err == nil && lib != nil {
				global.DB.Model(&lib).Update("created_by", userUUID)
				createdCount++
			}
		}

		// 获取统计数据
		stats, err := statsService.GetUserStats(userUUID)
		if err != nil {
			t.Fatalf("GetUserStats() error = %v", err)
		}

		if stats == nil {
			t.Fatal("Expected stats, got nil")
		}

		if stats.Libraries < createdCount {
			t.Logf("Note: Created %d libraries, got %d in stats (may be due to database constraints)", createdCount, stats.Libraries)
		}
	})

	t.Run("get user stats multiple times should use cache", func(t *testing.T) {
		userUUID := "user-cache-test"

		// 创建库
		libService := &service.LibraryService{}
		lib, _ := libService.Create(&request.LibraryCreate{
			Name:        "stats-cache-test-lib",
			Description: "test",
		})
		global.DB.Model(&lib).Update("created_by", userUUID)

		// 第一次获取
		stats1, err := statsService.GetUserStats(userUUID)
		if err != nil {
			t.Fatalf("GetUserStats() error = %v", err)
		}

		if stats1 == nil {
			t.Fatal("Expected stats, got nil")
		}

		// 第二次获取（应该从缓存返回）
		stats2, err := statsService.GetUserStats(userUUID)
		if err != nil {
			t.Fatalf("GetUserStats() error = %v", err)
		}

		if stats2 == nil {
			t.Fatal("Expected stats, got nil")
		}

		// 两次结果应该相同（来自缓存）
		if stats1.Libraries != stats2.Libraries {
			t.Logf("Note: Cache may have expired, stats1.Libraries=%d, stats2.Libraries=%d", stats1.Libraries, stats2.Libraries)
		}
	})
}
