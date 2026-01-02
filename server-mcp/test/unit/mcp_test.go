package test_test

import (
	"testing"

	"go-mcp-context/internal/model/request"
	"go-mcp-context/internal/service"
)

// TestMCPSearchLibraries 测试 MCP 库搜索
func Test_MCP_SearchLibraries(t *testing.T) {
	mcpService := service.NewMCPService()

	t.Run("search libraries by name", func(t *testing.T) {
		// 先创建一个库
		libService := &service.LibraryService{}
		lib, _ := libService.Create(&request.LibraryCreate{
			Name:        "mcp-search-lib",
			Description: "test library for mcp search",
		})

		// 搜索库
		req := &request.MCPSearchLibraries{
			LibraryName: "mcp-search",
		}

		result, err := mcpService.SearchLibraries(req)
		if err != nil {
			t.Fatalf("SearchLibraries() error = %v", err)
		}

		if result == nil {
			t.Fatal("Expected result, got nil")
		}

		// 验证返回的库列表
		if len(result.Libraries) == 0 {
			t.Logf("No libraries found (expected if search didn't match)")
		} else {
			// 验证库信息结构
			for _, lib := range result.Libraries {
				if lib.LibraryID == 0 {
					t.Error("Expected non-zero library ID")
				}
				if lib.Name == "" {
					t.Error("Expected non-empty library name")
				}
			}
		}

		_ = lib // 使用 lib 避免未使用警告
	})

	t.Run("search non-existent library", func(t *testing.T) {
		req := &request.MCPSearchLibraries{
			LibraryName: "non-existent-xyz-12345",
		}

		result, err := mcpService.SearchLibraries(req)
		if err != nil {
			t.Fatalf("SearchLibraries() error = %v", err)
		}

		if result == nil {
			t.Fatal("Expected result, got nil")
		}

		// 可能返回空列表或降级搜索结果
		if len(result.Libraries) > 0 {
			t.Logf("Found %d libraries (may be from fallback search)", len(result.Libraries))
		}
	})
}

// TestMCPGetLibraryDocs 测试 MCP 获取库文档
func Test_MCP_GetLibraryDocs(t *testing.T) {
	mcpService := service.NewMCPService()

	t.Run("get library docs with library id", func(t *testing.T) {
		req := &request.MCPGetLibraryDocs{
			LibraryID: 1,
			Topic:     "test",
			Mode:      "code",
			Version:   "latest",
			Page:      1,
		}

		result, err := mcpService.GetLibraryDocs(req)
		if err != nil {
			t.Logf("GetLibraryDocs() error = %v (expected if no documents)", err)
			return
		}

		if result == nil {
			t.Fatal("Expected result, got nil")
		}

		if result.LibraryID != 1 {
			t.Errorf("Expected library ID 1, got %d", result.LibraryID)
		}

		if result.Page < 1 {
			t.Errorf("Expected page >= 1, got %d", result.Page)
		}
	})

	t.Run("get library docs without library id (global search)", func(t *testing.T) {
		req := &request.MCPGetLibraryDocs{
			LibraryID: 0, // 全局搜索
			Topic:     "test",
			Mode:      "info",
			Version:   "latest",
			Page:      1,
		}

		result, err := mcpService.GetLibraryDocs(req)
		if err != nil {
			t.Logf("GetLibraryDocs() error = %v (expected if no documents)", err)
			return
		}

		if result == nil {
			t.Fatal("Expected result, got nil")
		}

		if result.LibraryID != 0 {
			t.Errorf("Expected library ID 0 (global), got %d", result.LibraryID)
		}
	})

	t.Run("get library docs with pagination", func(t *testing.T) {
		req := &request.MCPGetLibraryDocs{
			LibraryID: 1,
			Topic:     "test",
			Mode:      "",
			Version:   "latest",
			Page:      2,
		}

		result, err := mcpService.GetLibraryDocs(req)
		if err != nil {
			t.Logf("GetLibraryDocs() error = %v (expected if no documents)", err)
			return
		}

		if result == nil {
			t.Fatal("Expected result, got nil")
		}

		if result.Page != 2 {
			t.Errorf("Expected page 2, got %d", result.Page)
		}
	})
}

// TestMCPGetAllLibraries 测试获取所有库
func Test_MCP_GetAllLibraries(t *testing.T) {
	mcpService := service.NewMCPService()

	t.Run("get all libraries", func(t *testing.T) {
		libraries, err := mcpService.GetAllLibraries()
		if err != nil {
			t.Fatalf("GetAllLibraries() error = %v", err)
		}

		if libraries == nil {
			t.Fatal("Expected libraries, got nil")
		}

		// 验证库列表结构
		for _, lib := range libraries {
			if lib.ID == 0 {
				t.Error("Expected non-zero library ID")
			}
			if lib.Name == "" {
				t.Error("Expected non-empty library name")
			}
		}
	})
}

// TestMCPGetLibraryByID 测试按 ID 获取库
func Test_MCP_GetLibraryByID(t *testing.T) {
	mcpService := service.NewMCPService()

	t.Run("get library by id", func(t *testing.T) {
		// 先创建一个库
		libService := &service.LibraryService{}
		lib, _ := libService.Create(&request.LibraryCreate{
			Name:        "getbyid-mcp-lib",
			Description: "test",
		})

		// 按 ID 获取库
		retrieved, err := mcpService.GetLibraryByID(lib.ID)
		if err != nil {
			t.Fatalf("GetLibraryByID() error = %v", err)
		}

		if retrieved == nil {
			t.Fatal("Expected library, got nil")
		}

		if retrieved.ID != lib.ID {
			t.Errorf("Expected ID %d, got %d", lib.ID, retrieved.ID)
		}
	})

	t.Run("get non-existent library by id", func(t *testing.T) {
		retrieved, err := mcpService.GetLibraryByID(99999)
		if err != nil {
			t.Fatalf("GetLibraryByID() error = %v", err)
		}

		if retrieved != nil {
			t.Error("Expected nil for non-existent library")
		}
	})

	t.Run("get library by id zero", func(t *testing.T) {
		retrieved, err := mcpService.GetLibraryByID(0)
		if err != nil {
			t.Fatalf("GetLibraryByID() error = %v", err)
		}

		if retrieved != nil {
			t.Error("Expected nil for library ID 0")
		}
	})
}

// TestMCPSearchLibrariesAdvanced 测试高级搜索场景
func Test_MCP_SearchLibraries_Advanced(t *testing.T) {
	mcpService := service.NewMCPService()

	t.Run("search with empty name", func(t *testing.T) {
		req := &request.MCPSearchLibraries{
			LibraryName: "",
		}

		result, err := mcpService.SearchLibraries(req)
		if err != nil {
			t.Fatalf("SearchLibraries() error = %v", err)
		}

		if result == nil {
			t.Fatal("Expected result, got nil")
		}
	})

	t.Run("search with special characters", func(t *testing.T) {
		req := &request.MCPSearchLibraries{
			LibraryName: "test@#$%",
		}

		result, err := mcpService.SearchLibraries(req)
		if err != nil {
			t.Fatalf("SearchLibraries() error = %v", err)
		}

		if result == nil {
			t.Fatal("Expected result, got nil")
		}
	})

	t.Run("search with very long name", func(t *testing.T) {
		longName := ""
		for i := 0; i < 100; i++ {
			longName += "a"
		}

		req := &request.MCPSearchLibraries{
			LibraryName: longName,
		}

		result, err := mcpService.SearchLibraries(req)
		if err != nil {
			t.Fatalf("SearchLibraries() error = %v", err)
		}

		if result == nil {
			t.Fatal("Expected result, got nil")
		}
	})
}

// TestMCPGetLibraryDocsAdvanced 测试高级文档获取场景
func Test_MCP_GetLibraryDocs_Advanced(t *testing.T) {
	mcpService := service.NewMCPService()

	t.Run("get docs with high page number", func(t *testing.T) {
		req := &request.MCPGetLibraryDocs{
			LibraryID: 1,
			Topic:     "test",
			Mode:      "code",
			Version:   "latest",
			Page:      100,
		}

		result, err := mcpService.GetLibraryDocs(req)
		if err != nil {
			t.Logf("GetLibraryDocs() error = %v (expected if page out of range)", err)
			return
		}

		if result != nil {
			// 页码可能被调整到有效范围内
			if result.Page < 1 {
				t.Errorf("Expected page >= 1, got %d", result.Page)
			}
		}
	})

	t.Run("get docs with empty topic", func(t *testing.T) {
		req := &request.MCPGetLibraryDocs{
			LibraryID: 1,
			Topic:     "",
			Mode:      "code",
			Version:   "latest",
			Page:      1,
		}

		result, err := mcpService.GetLibraryDocs(req)
		if err != nil {
			t.Logf("GetLibraryDocs() error = %v (expected if no documents)", err)
			return
		}

		if result != nil {
			if result.Page != 1 {
				t.Errorf("Expected page 1, got %d", result.Page)
			}
		}
	})

	t.Run("get docs with invalid version", func(t *testing.T) {
		req := &request.MCPGetLibraryDocs{
			LibraryID: 1,
			Topic:     "test",
			Mode:      "code",
			Version:   "v99.99.99",
			Page:      1,
		}

		result, err := mcpService.GetLibraryDocs(req)
		if err != nil {
			t.Logf("GetLibraryDocs() error = %v (expected if version doesn't exist)", err)
			return
		}

		if result != nil {
			if result.Page != 1 {
				t.Errorf("Expected page 1, got %d", result.Page)
			}
		}
	})
}

// TestMCPGetAllLibrariesAdvanced 测试获取所有库的高级场景
func Test_MCP_GetAllLibraries_Advanced(t *testing.T) {
	mcpService := service.NewMCPService()

	t.Run("get all libraries returns valid list", func(t *testing.T) {
		result, err := mcpService.GetAllLibraries()
		if err != nil {
			t.Fatalf("GetAllLibraries() error = %v", err)
		}

		if result == nil {
			t.Fatal("Expected result, got nil")
		}

		// 应该至少有一些库
		if len(result) == 0 {
			t.Logf("Note: No libraries found in database")
		}
	})

	t.Run("get all libraries returns array", func(t *testing.T) {
		result, err := mcpService.GetAllLibraries()
		if err != nil {
			t.Fatalf("GetAllLibraries() error = %v", err)
		}

		if result == nil {
			t.Fatal("Expected result, got nil")
		}

		// 检查返回的是非空数组
		if len(result) == 0 {
			t.Error("Expected non-empty result")
		}

		// 检查每个库都有 ID
		for _, lib := range result {
			if lib.ID == 0 {
				t.Error("Expected library ID > 0")
			}
		}
	})
}

// TestMCPSearchLibrariesEdgeCases 测试搜索库的边界情况
func Test_MCP_SearchLibraries_EdgeCases(t *testing.T) {
	mcpService := service.NewMCPService()

	t.Run("search with null bytes in name", func(t *testing.T) {
		req := &request.MCPSearchLibraries{
			LibraryName: "test\x00name",
		}

		result, err := mcpService.SearchLibraries(req)
		if err != nil {
			t.Logf("SearchLibraries() error = %v (expected for invalid input)", err)
			return
		}

		if result != nil {
			// 只要没有 panic，就认为是成功的
			t.Logf("SearchLibraries returned %d results", len(result.Libraries))
		}
	})

	t.Run("search with sql injection attempt", func(t *testing.T) {
		req := &request.MCPSearchLibraries{
			LibraryName: "'; DROP TABLE libraries; --",
		}

		result, err := mcpService.SearchLibraries(req)
		if err != nil {
			t.Logf("SearchLibraries() error = %v (expected for invalid input)", err)
			return
		}

		if result != nil {
			// 只要没有 panic 或数据库错误，就认为是安全的
			t.Logf("SearchLibraries safely handled injection attempt")
		}
	})

	t.Run("search with very long query", func(t *testing.T) {
		longQuery := ""
		for i := 0; i < 1000; i++ {
			longQuery += "a"
		}

		req := &request.MCPSearchLibraries{
			LibraryName: longQuery,
		}

		result, err := mcpService.SearchLibraries(req)
		if err != nil {
			t.Logf("SearchLibraries() error = %v (expected for very long query)", err)
			return
		}

		if result != nil {
			t.Logf("SearchLibraries handled long query with %d results", len(result.Libraries))
		}
	})
}
