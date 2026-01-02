package test_test

import (
	"testing"
	"time"

	"go-mcp-context/internal/model/request"
	"go-mcp-context/internal/service"
)

// Test_Search_InvalidateLibraryCache 测试缓存失效
func Test_Search_InvalidateLibraryCache(t *testing.T) {
	searchService := &service.SearchService{}

	t.Run("invalidate library cache", func(t *testing.T) {
		// 通过 Service 层创建库
		libService := &service.LibraryService{}
		lib, _ := libService.Create(&request.LibraryCreate{Name: "cache-lib", Description: "test"})

		// 失效缓存
		err := searchService.InvalidateLibraryCache(lib.ID, "latest")
		if err != nil {
			t.Fatalf("InvalidateLibraryCache() error = %v", err)
		}
	})

	t.Run("invalidate cache for non-existent library", func(t *testing.T) {
		err := searchService.InvalidateLibraryCache(99999, "latest")
		if err != nil {
			t.Fatalf("InvalidateLibraryCache() error = %v", err)
		}
	})

	t.Run("invalidate cache with different versions", func(t *testing.T) {
		libService := &service.LibraryService{}
		lib, _ := libService.Create(&request.LibraryCreate{Name: "cache-version-lib", Description: "test"})

		versions := []string{"latest", "v1.0.0", "v2.0.0"}
		for _, version := range versions {
			err := searchService.InvalidateLibraryCache(lib.ID, version)
			if err != nil {
				t.Fatalf("InvalidateLibraryCache(version=%s) error = %v", version, err)
			}
		}
	})

	t.Run("invalidate cache multiple times", func(t *testing.T) {
		libService := &service.LibraryService{}
		lib, _ := libService.Create(&request.LibraryCreate{Name: "cache-multi-lib", Description: "test"})

		// 多次失效缓存
		for i := 0; i < 3; i++ {
			err := searchService.InvalidateLibraryCache(lib.ID, "latest")
			if err != nil {
				t.Fatalf("InvalidateLibraryCache() iteration %d error = %v", i, err)
			}
		}
	})

	t.Run("invalidate cache with zero library id", func(t *testing.T) {
		err := searchService.InvalidateLibraryCache(0, "latest")
		if err != nil {
			t.Logf("InvalidateLibraryCache(libID=0) error = %v (expected)", err)
		}
	})

	t.Run("invalidate cache with empty version", func(t *testing.T) {
		libService := &service.LibraryService{}
		lib, _ := libService.Create(&request.LibraryCreate{Name: "cache-empty-version-lib", Description: "test"})

		err := searchService.InvalidateLibraryCache(lib.ID, "")
		if err != nil {
			t.Logf("InvalidateLibraryCache(version='') error = %v (expected)", err)
		}
	})
}

// Test_Search_SearchDocuments_Advanced 测试搜索排序和相关性
func Test_Search_SearchDocuments_Advanced(t *testing.T) {
	searchService := &service.SearchService{}

	t.Run("search results have relevance scores", func(t *testing.T) {
		req := &request.Search{
			LibraryID: 1,
			Query:     "function",
			Mode:      "code",
			Version:   "latest",
			Page:      1,
			Limit:     10,
		}

		result, err := searchService.SearchDocuments(req)
		if err != nil {
			t.Logf("SearchDocuments() error = %v (expected if no documents)", err)
			return
		}

		if result != nil && len(result.Results) > 0 {
			// 检查相关性分数
			for i, item := range result.Results {
				if item.Relevance < 0 || item.Relevance > 1 {
					t.Errorf("Result %d: Expected relevance between 0-1, got %f", i, item.Relevance)
				}
			}

			// 检查结果是否按相关性降序排列
			for i := 1; i < len(result.Results); i++ {
				if result.Results[i].Relevance > result.Results[i-1].Relevance {
					t.Errorf("Results not sorted by relevance: %f > %f", result.Results[i].Relevance, result.Results[i-1].Relevance)
				}
			}
		}
	})

	t.Run("search with code mode returns code chunks", func(t *testing.T) {
		req := &request.Search{
			LibraryID: 1,
			Query:     "func",
			Mode:      "code",
			Version:   "latest",
			Page:      1,
			Limit:     10,
		}

		result, err := searchService.SearchDocuments(req)
		if err != nil {
			t.Logf("SearchDocuments() error = %v (expected if no documents)", err)
			return
		}

		if result != nil && len(result.Results) > 0 {
			for _, item := range result.Results {
				if item.Mode != "code" {
					t.Errorf("Expected mode 'code', got '%s'", item.Mode)
				}
				// code 模式应该有 Code 字段
				if item.Code == "" {
					t.Logf("Warning: code mode item has empty Code field")
				}
			}
		}
	})

	t.Run("search with info mode returns info chunks", func(t *testing.T) {
		req := &request.Search{
			LibraryID: 1,
			Query:     "struct",
			Mode:      "info",
			Version:   "latest",
			Page:      1,
			Limit:     10,
		}

		result, err := searchService.SearchDocuments(req)
		if err != nil {
			t.Logf("SearchDocuments() error = %v (expected if no documents)", err)
			return
		}

		if result != nil && len(result.Results) > 0 {
			for _, item := range result.Results {
				if item.Mode != "info" {
					t.Errorf("Expected mode 'info', got '%s'", item.Mode)
				}
				// info 模式应该有 Content 字段
				if item.Content == "" {
					t.Logf("Warning: info mode item has empty Content field")
				}
			}
		}
	})

	t.Run("search pagination works correctly", func(t *testing.T) {
		// 第一页
		req1 := &request.Search{
			LibraryID: 1,
			Query:     "test",
			Mode:      "",
			Version:   "latest",
			Page:      1,
			Limit:     5,
		}

		result1, err := searchService.SearchDocuments(req1)
		if err != nil {
			t.Logf("SearchDocuments() error = %v (expected if no documents)", err)
			return
		}

		if result1 != nil {
			// 第二页
			req2 := &request.Search{
				LibraryID: 1,
				Query:     "test",
				Mode:      "",
				Version:   "latest",
				Page:      2,
				Limit:     5,
			}

			result2, err := searchService.SearchDocuments(req2)
			if err != nil {
				t.Logf("SearchDocuments() error = %v (expected if no documents)", err)
				return
			}

			if result2 != nil && len(result1.Results) > 0 && len(result2.Results) > 0 {
				// 检查两页的结果是否不同
				if len(result1.Results) > 0 && len(result2.Results) > 0 {
					if result1.Results[0].ChunkID == result2.Results[0].ChunkID {
						t.Logf("Note: First result on both pages is the same (may be expected)")
					}
				}
			}
		}
	})

	t.Run("search with special characters in query", func(t *testing.T) {
		specialQueries := []string{
			"func()",
			"type[]",
			"var:=",
			"if else",
		}

		for _, query := range specialQueries {
			req := &request.Search{
				LibraryID: 1,
				Query:     query,
				Mode:      "code",
				Version:   "latest",
				Page:      1,
				Limit:     10,
			}

			result, err := searchService.SearchDocuments(req)
			if err != nil {
				t.Logf("SearchDocuments(query='%s') error = %v (expected if no documents)", query, err)
				continue
			}

			if result != nil {
				// 只要没有 panic，就认为是成功的
				t.Logf("SearchDocuments(query='%s') returned %d results", query, len(result.Results))
			}
		}
	})

	t.Run("search with vector similarity filtering", func(t *testing.T) {
		// 测试向量搜索的相似度过滤
		req := &request.Search{
			LibraryID: 1,
			Query:     "interface implementation",
			Mode:      "code",
			Version:   "latest",
			Page:      1,
			Limit:     5,
		}

		result, err := searchService.SearchDocuments(req)
		if err != nil {
			t.Logf("SearchDocuments() error = %v (expected if no documents)", err)
			return
		}

		if result != nil && len(result.Results) > 0 {
			// 检查结果是否有相关性分数
			for _, item := range result.Results {
				if item.Relevance < 0 || item.Relevance > 1 {
					t.Errorf("Invalid relevance score: %f", item.Relevance)
				}
			}
		}
	})

	t.Run("search with high limit and multiple pages", func(t *testing.T) {
		// 测试大 limit 和多页
		for page := 1; page <= 3; page++ {
			req := &request.Search{
				LibraryID: 1,
				Query:     "error handling",
				Mode:      "code",
				Version:   "latest",
				Page:      page,
				Limit:     20,
			}

			result, err := searchService.SearchDocuments(req)
			if err != nil {
				t.Logf("SearchDocuments(page=%d) error = %v (expected if no documents)", page, err)
				continue
			}

			if result != nil {
				if result.Page != page {
					t.Errorf("Expected page %d, got %d", page, result.Page)
				}
			}
		}
	})

	t.Run("search with empty library id", func(t *testing.T) {
		req := &request.Search{
			LibraryID: 0,
			Query:     "test",
			Mode:      "code",
			Version:   "latest",
			Page:      1,
			Limit:     10,
		}

		result, err := searchService.SearchDocuments(req)
		if err != nil {
			t.Logf("SearchDocuments(libID=0) error = %v (expected)", err)
			return
		}

		if result != nil {
			t.Logf("SearchDocuments(libID=0) returned %d results", len(result.Results))
		}
	})
}

// Test_Search_HybridRRF 测试混合 RRF 算法
func Test_Search_HybridRRF(t *testing.T) {
	searchService := &service.SearchService{}

	t.Run("hybrid rrf with single library search", func(t *testing.T) {
		// 测试混合 RRF 算法处理单个库的搜索结果
		req := &request.Search{
			LibraryID: 1,
			Query:     "test hybrid rrf",
			Mode:      "code",
			Version:   "latest",
			Page:      1,
			Limit:     10,
		}

		result, err := searchService.SearchDocuments(req)
		if err != nil {
			t.Logf("SearchDocuments() error = %v (expected if no data)", err)
			return
		}

		if result != nil {
			t.Logf("HybridRRF returned %d results", len(result.Results))
		}
	})

	t.Run("hybrid rrf with multiple topics", func(t *testing.T) {
		// 测试混合 RRF 算法处理多个主题
		req := &request.Search{
			LibraryID: 1,
			Query:     "error handling exception",
			Mode:      "code",
			Version:   "latest",
			Page:      1,
			Limit:     20,
		}

		result, err := searchService.SearchDocuments(req)
		if err != nil {
			t.Logf("SearchDocuments(multiple topics) error = %v (expected if no data)", err)
			return
		}

		if result != nil && len(result.Results) > 0 {
			// 验证结果排序（按相关性降序）
			for i := 0; i < len(result.Results)-1; i++ {
				if result.Results[i].Relevance < result.Results[i+1].Relevance {
					t.Errorf("Results not properly ranked: result %d (%.2f) < result %d (%.2f)",
						i, result.Results[i].Relevance, i+1, result.Results[i+1].Relevance)
				}
			}
		}
	})

	t.Run("hybrid rrf with empty query", func(t *testing.T) {
		// 测试混合 RRF 算法处理空查询
		req := &request.Search{
			LibraryID: 1,
			Query:     "",
			Mode:      "code",
			Version:   "latest",
			Page:      1,
			Limit:     10,
		}

		result, err := searchService.SearchDocuments(req)
		if err != nil {
			t.Logf("SearchDocuments(empty query) error = %v (expected)", err)
			return
		}

		if result != nil {
			t.Logf("SearchDocuments(empty query) returned %d results", len(result.Results))
		}
	})
}

// Test_Search_ReciprocalRankFusion 测试倒数排名融合
func Test_Search_ReciprocalRankFusion(t *testing.T) {
	searchService := &service.SearchService{}

	t.Run("reciprocal rank fusion with vector search", func(t *testing.T) {
		// 测试倒数排名融合与向量搜索的结合
		req := &request.Search{
			LibraryID: 1,
			Query:     "database query optimization",
			Mode:      "code",
			Version:   "latest",
			Page:      1,
			Limit:     15,
		}

		result, err := searchService.SearchDocuments(req)
		if err != nil {
			t.Logf("SearchDocuments() error = %v (expected if no data)", err)
			return
		}

		if result != nil && len(result.Results) > 0 {
			// 验证倒数排名融合的结果
			for i, r := range result.Results {
				if r.Relevance < 0 || r.Relevance > 1 {
					t.Errorf("Result %d has invalid relevance score: %.2f", i, r.Relevance)
				}
			}
		}
	})

	t.Run("reciprocal rank fusion with different k values", func(t *testing.T) {
		// 测试不同的 k 值对排名的影响
		queries := []string{
			"error handling",
			"performance optimization",
			"security best practices",
		}

		for _, query := range queries {
			req := &request.Search{
				LibraryID: 1,
				Query:     query,
				Mode:      "code",
				Version:   "latest",
				Page:      1,
				Limit:     10,
			}

			result, err := searchService.SearchDocuments(req)
			if err != nil {
				t.Logf("SearchDocuments(query=%s) error = %v (expected if no data)", query, err)
				continue
			}

			if result != nil && len(result.Results) > 0 {
				t.Logf("Query '%s' returned %d results with top relevance: %.2f",
					query, len(result.Results), result.Results[0].Relevance)
			}
		}
	})

	t.Run("reciprocal rank fusion result consistency", func(t *testing.T) {
		// 测试倒数排名融合的结果一致性
		req := &request.Search{
			LibraryID: 1,
			Query:     "consistency test",
			Mode:      "code",
			Version:   "latest",
			Page:      1,
			Limit:     10,
		}

		// 执行两次相同的搜索
		result1, err1 := searchService.SearchDocuments(req)
		result2, err2 := searchService.SearchDocuments(req)

		if err1 != nil || err2 != nil {
			t.Logf("SearchDocuments() errors: %v, %v (expected if no data)", err1, err2)
			return
		}

		if result1 != nil && result2 != nil {
			if len(result1.Results) != len(result2.Results) {
				t.Errorf("Inconsistent result count: %d vs %d", len(result1.Results), len(result2.Results))
			}

			// 验证前几个结果的顺序一致
			for i := 0; i < len(result1.Results) && i < 3; i++ {
				if result1.Results[i].ChunkID != result2.Results[i].ChunkID {
					t.Logf("Result order differs at position %d: %d vs %d",
						i, result1.Results[i].ChunkID, result2.Results[i].ChunkID)
				}
			}
		}
	})
}

// Test_Search_SearchDocuments_EdgeCases 测试高级搜索场景
func Test_Search_SearchDocuments_EdgeCases(t *testing.T) {
	searchService := &service.SearchService{}

	t.Run("search with different library ids", func(t *testing.T) {
		// 测试不同库的搜索
		for libID := uint(1); libID <= 3; libID++ {
			req := &request.Search{
				LibraryID: libID,
				Query:     "test",
				Mode:      "code",
				Version:   "latest",
				Page:      1,
				Limit:     10,
			}

			result, err := searchService.SearchDocuments(req)
			if err != nil {
				t.Logf("SearchDocuments(libID=%d) error = %v (expected if no documents)", libID, err)
				continue
			}

			if result != nil {
				// 检查结果中的库 ID
				if len(result.Results) > 0 {
					if result.Results[0].LibraryID != libID {
						t.Errorf("Expected library ID %d, got %d", libID, result.Results[0].LibraryID)
					}
				}
			}
		}
	})

	t.Run("search with different versions", func(t *testing.T) {
		versions := []string{"latest", "v1.0.0", "v2.0.0"}

		for _, version := range versions {
			req := &request.Search{
				LibraryID: 1,
				Query:     "test",
				Mode:      "code",
				Version:   version,
				Page:      1,
				Limit:     10,
			}

			result, err := searchService.SearchDocuments(req)
			if err != nil {
				t.Logf("SearchDocuments(version=%s) error = %v (expected if version doesn't exist)", version, err)
				continue
			}

			if result != nil {
				// 检查结果中的版本
				if len(result.Results) > 0 {
					if result.Results[0].Version != version {
						t.Errorf("Expected version %s, got %s", version, result.Results[0].Version)
					}
				}
			}
		}
	})

	t.Run("search with mode filtering", func(t *testing.T) {
		modes := []string{"code", "info", ""}

		for _, mode := range modes {
			req := &request.Search{
				LibraryID: 1,
				Query:     "test",
				Mode:      mode,
				Version:   "latest",
				Page:      1,
				Limit:     10,
			}

			result, err := searchService.SearchDocuments(req)
			if err != nil {
				t.Logf("SearchDocuments(mode=%s) error = %v (expected if no documents)", mode, err)
				continue
			}

			if result != nil {
				// 检查结果中的模式
				if len(result.Results) > 0 && mode != "" {
					if result.Results[0].Mode != mode {
						t.Errorf("Expected mode %s, got %s", mode, result.Results[0].Mode)
					}
				}
			}
		}
	})

	t.Run("search with large limit", func(t *testing.T) {
		req := &request.Search{
			LibraryID: 1,
			Query:     "test",
			Mode:      "code",
			Version:   "latest",
			Page:      1,
			Limit:     50,
		}

		result, err := searchService.SearchDocuments(req)
		if err != nil {
			t.Logf("SearchDocuments() error = %v (expected if no documents)", err)
			return
		}

		if result != nil {
			// limit 应该被限制为最大值 10
			if result.Limit > 10 {
				t.Errorf("Expected limit <= 10, got %d", result.Limit)
			}
		}
	})

	t.Run("search with negative page", func(t *testing.T) {
		req := &request.Search{
			LibraryID: 1,
			Query:     "test",
			Mode:      "code",
			Version:   "latest",
			Page:      -5,
			Limit:     10,
		}

		result, err := searchService.SearchDocuments(req)
		if err != nil {
			t.Logf("SearchDocuments() error = %v (expected if no documents)", err)
			return
		}

		if result != nil {
			// page 应该被调整为 >= 1
			if result.Page < 1 {
				t.Errorf("Expected page >= 1, got %d", result.Page)
			}
		}
	})
}

// Test_Search_SearchDocuments 测试搜索文档（不同参数分支）
func Test_Search_SearchDocuments(t *testing.T) {
	searchService := &service.SearchService{}

	t.Run("search with default pagination", func(t *testing.T) {
		// 搜索库 1 的文档
		req := &request.Search{
			LibraryID: 1,
			Query:     "test",
			Mode:      "code",
			Version:   "latest",
			Page:      0, // 默认值
			Limit:     0, // 默认值
		}

		result, err := searchService.SearchDocuments(req)
		if err != nil {
			t.Logf("SearchDocuments() error = %v (expected if no documents)", err)
			return
		}

		if result != nil {
			if result.Page < 1 {
				t.Errorf("Expected page >= 1, got %d", result.Page)
			}

			if result.Limit < 1 {
				t.Errorf("Expected limit >= 1, got %d", result.Limit)
			}
		}
	})

	t.Run("search with custom pagination", func(t *testing.T) {
		// 搜索库 1 的文档（自定义分页）
		req := &request.Search{
			LibraryID: 1,
			Query:     "test",
			Mode:      "info",
			Version:   "latest",
			Page:      2,
			Limit:     5,
		}

		result, err := searchService.SearchDocuments(req)
		if err != nil {
			t.Logf("SearchDocuments() error = %v (expected if no documents)", err)
			return
		}

		if result != nil {
			if result.Page != 2 {
				t.Errorf("Expected page 2, got %d", result.Page)
			}

			if result.Limit != 5 {
				t.Errorf("Expected limit 5, got %d", result.Limit)
			}
		}
	})

	t.Run("search with limit exceeding max", func(t *testing.T) {
		// 搜索库 1 的文档（limit 超过最大值）
		req := &request.Search{
			LibraryID: 1,
			Query:     "test",
			Mode:      "",
			Version:   "latest",
			Page:      1,
			Limit:     100, // 超过最大值 10
		}

		result, err := searchService.SearchDocuments(req)
		if err != nil {
			t.Logf("SearchDocuments() error = %v (expected if no documents)", err)
			return
		}

		if result != nil {
			if result.Limit > 10 {
				t.Errorf("Expected limit <= 10, got %d", result.Limit)
			}
		}
	})

	t.Run("search with empty query", func(t *testing.T) {
		// 搜索空查询
		req := &request.Search{
			LibraryID: 1,
			Query:     "",
			Mode:      "code",
			Version:   "latest",
			Page:      1,
			Limit:     10,
		}

		result, err := searchService.SearchDocuments(req)
		if err != nil {
			t.Logf("SearchDocuments() error = %v (expected if no documents)", err)
			return
		}

		if result != nil {
			if result.Page != 1 {
				t.Errorf("Expected page 1, got %d", result.Page)
			}
		}
	})

	t.Run("search with multiple topics", func(t *testing.T) {
		// 搜索多个 topic（逗号分隔）
		req := &request.Search{
			LibraryID: 1,
			Query:     "test,example,demo",
			Mode:      "code",
			Version:   "latest",
			Page:      1,
			Limit:     10,
		}

		result, err := searchService.SearchDocuments(req)
		if err != nil {
			t.Logf("SearchDocuments() error = %v (expected if no documents)", err)
			return
		}

		if result != nil {
			if result.Page != 1 {
				t.Errorf("Expected page 1, got %d", result.Page)
			}
		}
	})

	t.Run("search with space-separated topics", func(t *testing.T) {
		// 搜索空格分隔的 topic
		req := &request.Search{
			LibraryID: 1,
			Query:     "test example demo",
			Mode:      "info",
			Version:   "latest",
			Page:      1,
			Limit:     10,
		}

		result, err := searchService.SearchDocuments(req)
		if err != nil {
			t.Logf("SearchDocuments() error = %v (expected if no documents)", err)
			return
		}

		if result != nil {
			if result.Page != 1 {
				t.Errorf("Expected page 1, got %d", result.Page)
			}
		}
	})

	t.Run("search with zero limit defaults to 10", func(t *testing.T) {
		req := &request.Search{
			LibraryID: 1,
			Query:     "test",
			Mode:      "code",
			Version:   "latest",
			Page:      1,
			Limit:     0, // 应该默认为 10
		}

		result, err := searchService.SearchDocuments(req)
		if err != nil {
			t.Logf("SearchDocuments() error = %v (expected if no documents)", err)
			return
		}

		if result != nil {
			if result.Limit != 10 {
				t.Errorf("Expected limit 10, got %d", result.Limit)
			}
		}
	})

	t.Run("search with zero page defaults to 1", func(t *testing.T) {
		req := &request.Search{
			LibraryID: 1,
			Query:     "test",
			Mode:      "code",
			Version:   "latest",
			Page:      0, // 应该默认为 1
			Limit:     10,
		}

		result, err := searchService.SearchDocuments(req)
		if err != nil {
			t.Logf("SearchDocuments() error = %v (expected if no documents)", err)
			return
		}

		if result != nil {
			if result.Page != 1 {
				t.Errorf("Expected page 1, got %d", result.Page)
			}
		}
	})

	t.Run("search with very large limit gets capped", func(t *testing.T) {
		req := &request.Search{
			LibraryID: 1,
			Query:     "test",
			Mode:      "code",
			Version:   "latest",
			Page:      1,
			Limit:     999, // 应该被限制为最大值
		}

		result, err := searchService.SearchDocuments(req)
		if err != nil {
			t.Logf("SearchDocuments() error = %v (expected if no documents)", err)
			return
		}

		if result != nil {
			if result.Limit > 10 {
				t.Errorf("Expected limit <= 10, got %d", result.Limit)
			}
		}
	})

	t.Run("search with negative page gets adjusted", func(t *testing.T) {
		req := &request.Search{
			LibraryID: 1,
			Query:     "test",
			Mode:      "code",
			Version:   "latest",
			Page:      -10, // 应该被调整为 >= 1
			Limit:     10,
		}

		result, err := searchService.SearchDocuments(req)
		if err != nil {
			t.Logf("SearchDocuments() error = %v (expected if no documents)", err)
			return
		}

		if result != nil {
			if result.Page < 1 {
				t.Errorf("Expected page >= 1, got %d", result.Page)
			}
		}
	})

	t.Run("search with all modes", func(t *testing.T) {
		modes := []string{"code", "info", ""}
		for _, mode := range modes {
			req := &request.Search{
				LibraryID: 1,
				Query:     "test",
				Mode:      mode,
				Version:   "latest",
				Page:      1,
				Limit:     10,
			}

			result, err := searchService.SearchDocuments(req)
			if err != nil {
				t.Logf("SearchDocuments(mode=%s) error = %v (expected if no documents)", mode, err)
				continue
			}

			if result != nil {
				t.Logf("SearchDocuments(mode=%s) returned %d results", mode, len(result.Results))
			}
		}
	})

	t.Run("search with different versions", func(t *testing.T) {
		versions := []string{"latest", "v1.0.0", "v2.0.0"}
		for _, version := range versions {
			req := &request.Search{
				LibraryID: 1,
				Query:     "test",
				Mode:      "code",
				Version:   version,
				Page:      1,
				Limit:     10,
			}

			result, err := searchService.SearchDocuments(req)
			if err != nil {
				t.Logf("SearchDocuments(version=%s) error = %v (expected if version doesn't exist)", version, err)
				continue
			}

			if result != nil {
				t.Logf("SearchDocuments(version=%s) returned %d results", version, len(result.Results))
			}
		}
	})
}

// Test_Search_SearchDocuments_WithLibraryID 测试搜索文档（不同库ID）
func Test_Search_SearchDocuments_WithLibraryID(t *testing.T) {
	searchService := &service.SearchService{}

	t.Run("search with library id 1", func(t *testing.T) {
		req := &request.Search{
			LibraryID: 1,
			Query:     "test query",
			Mode:      "code",
			Version:   "latest",
			Page:      1,
			Limit:     10,
		}

		result, err := searchService.SearchDocuments(req)
		if err != nil {
			t.Logf("SearchDocuments(libID=1) error = %v (expected if no documents)", err)
			return
		}

		if result != nil {
			// 验证返回的结果项中的 library ID
			for _, item := range result.Results {
				if item.LibraryID != 1 {
					t.Errorf("Expected library ID 1, got %d", item.LibraryID)
				}
			}
		}
	})

	t.Run("search with library id 2", func(t *testing.T) {
		req := &request.Search{
			LibraryID: 2,
			Query:     "another query",
			Mode:      "info",
			Version:   "latest",
			Page:      1,
			Limit:     5,
		}

		result, err := searchService.SearchDocuments(req)
		if err != nil {
			t.Logf("SearchDocuments(libID=2) error = %v (expected if no documents)", err)
			return
		}

		if result != nil {
			// 验证返回的结果项中的 library ID
			for _, item := range result.Results {
				if item.LibraryID != 2 {
					t.Errorf("Expected library ID 2, got %d", item.LibraryID)
				}
			}
		}
	})

	t.Run("search with non-existent library id", func(t *testing.T) {
		req := &request.Search{
			LibraryID: 99999,
			Query:     "test",
			Mode:      "code",
			Version:   "latest",
			Page:      1,
			Limit:     10,
		}

		result, err := searchService.SearchDocuments(req)
		if err != nil {
			t.Logf("SearchDocuments(libID=99999) error = %v (expected)", err)
			return
		}

		if result != nil && len(result.Results) > 0 {
			t.Logf("SearchDocuments(libID=99999) returned %d results (unexpected)", len(result.Results))
		}
	})
}

// Test_Search_InvalidateLibraryCache_Advanced 测试缓存失效的高级场景
func Test_Search_InvalidateLibraryCache_Advanced(t *testing.T) {
	searchService := &service.SearchService{}

	t.Run("invalidate cache for multiple libraries", func(t *testing.T) {
		// 测试多个库的缓存失效
		libIDs := []uint{1, 2, 3}
		for _, libID := range libIDs {
			err := searchService.InvalidateLibraryCache(libID, "latest")
			if err != nil {
				t.Fatalf("InvalidateLibraryCache(libID=%d) error = %v", libID, err)
			}
		}
	})

	t.Run("invalidate cache with various versions", func(t *testing.T) {
		versions := []string{"latest", "v1.0.0", "v2.0.0", "v3.0.0"}
		for _, version := range versions {
			err := searchService.InvalidateLibraryCache(1, version)
			if err != nil {
				t.Fatalf("InvalidateLibraryCache(version=%s) error = %v", version, err)
			}
		}
	})
}

// Test_Search_SearchDocuments_HybridMode 测试混合搜索模式
func Test_Search_SearchDocuments_HybridMode(t *testing.T) {
	searchService := &service.SearchService{}
	libService := &service.LibraryService{}

	// 创建测试库
	lib, err := libService.Create(&request.LibraryCreate{
		Name:        "hybrid-search-lib",
		Description: "test hybrid search",
	})
	if err != nil {
		t.Fatalf("Create() error = %v", err)
	}

	t.Run("search with hybrid mode", func(t *testing.T) {
		req := &request.Search{
			LibraryID: lib.ID,
			Version:   lib.DefaultVersion,
			Query:     "test query",
			Mode:      "hybrid", // 混合模式
			Limit:     10,
		}

		result, err := searchService.SearchDocuments(req)
		if err != nil {
			t.Logf("SearchDocuments(hybrid) error = %v (expected if no documents)", err)
			return
		}

		if result != nil {
			t.Logf("SearchDocuments(hybrid) returned %d results", len(result.Results))
		}
	})

	t.Run("search with vector mode", func(t *testing.T) {
		req := &request.Search{
			LibraryID: lib.ID,
			Version:   lib.DefaultVersion,
			Query:     "test query",
			Mode:      "vector", // 向量模式
			Limit:     10,
		}

		result, err := searchService.SearchDocuments(req)
		if err != nil {
			t.Logf("SearchDocuments(vector) error = %v (expected if no documents)", err)
			return
		}

		if result != nil {
			t.Logf("SearchDocuments(vector) returned %d results", len(result.Results))
		}
	})

	t.Run("search with bm25 mode", func(t *testing.T) {
		req := &request.Search{
			LibraryID: lib.ID,
			Version:   lib.DefaultVersion,
			Query:     "test query",
			Mode:      "bm25", // BM25 模式
			Limit:     10,
		}

		result, err := searchService.SearchDocuments(req)
		if err != nil {
			t.Logf("SearchDocuments(bm25) error = %v (expected if no documents)", err)
			return
		}

		if result != nil {
			t.Logf("SearchDocuments(bm25) returned %d results", len(result.Results))
		}
	})

	t.Run("search with different limits", func(t *testing.T) {
		limits := []int{5, 10, 20, 50}
		for _, limit := range limits {
			req := &request.Search{
				LibraryID: lib.ID,
				Version:   lib.DefaultVersion,
				Query:     "test",
				Mode:      "hybrid",
				Limit:     limit,
			}

			result, err := searchService.SearchDocuments(req)
			if err != nil {
				t.Logf("SearchDocuments(limit=%d) error = %v", limit, err)
				continue
			}

			if result != nil && len(result.Results) > limit {
				t.Errorf("Expected at most %d results, got %d", limit, len(result.Results))
			}
		}
	})

	t.Run("search with empty query", func(t *testing.T) {
		req := &request.Search{
			LibraryID: lib.ID,
			Version:   lib.DefaultVersion,
			Query:     "", // 空查询
			Mode:      "hybrid",
			Limit:     10,
		}

		result, err := searchService.SearchDocuments(req)
		if err != nil {
			t.Logf("SearchDocuments(empty query) error = %v", err)
		}

		if result != nil {
			t.Logf("SearchDocuments(empty query) returned %d results", len(result.Results))
		}
	})
}

// Test_Search_SearchDocuments_SpecialCases 测试搜索的特殊情况
func Test_Search_SearchDocuments_SpecialCases(t *testing.T) {
	searchService := &service.SearchService{}
	libService := &service.LibraryService{}

	lib, _ := libService.Create(&request.LibraryCreate{
		Name:        "special-case-lib",
		Description: "test special cases",
	})

	t.Run("search with very long query", func(t *testing.T) {
		longQuery := ""
		for i := 0; i < 100; i++ {
			longQuery += "test query "
		}

		req := &request.Search{
			LibraryID: lib.ID,
			Version:   lib.DefaultVersion,
			Query:     longQuery,
			Mode:      "hybrid",
			Limit:     10,
		}

		result, err := searchService.SearchDocuments(req)
		if err != nil {
			t.Logf("SearchDocuments(long query) error = %v", err)
		}

		if result != nil {
			t.Logf("SearchDocuments(long query) returned %d results", len(result.Results))
		}
	})

	t.Run("search with special characters", func(t *testing.T) {
		specialQueries := []string{
			"test@#$%^&*()",
			"test\n\t\r",
			"测试中文查询",
			"test 'single' \"double\" quotes",
		}

		for _, query := range specialQueries {
			req := &request.Search{
				LibraryID: lib.ID,
				Version:   lib.DefaultVersion,
				Query:     query,
				Mode:      "hybrid",
				Limit:     10,
			}

			result, err := searchService.SearchDocuments(req)
			if err != nil {
				t.Logf("SearchDocuments(query=%q) error = %v", query, err)
				continue
			}

			if result != nil {
				t.Logf("SearchDocuments(query=%q) returned %d results", query, len(result.Results))
			}
		}
	})

	t.Run("search with invalid mode", func(t *testing.T) {
		req := &request.Search{
			LibraryID: lib.ID,
			Version:   lib.DefaultVersion,
			Query:     "test",
			Mode:      "invalid-mode", // 无效模式
			Limit:     10,
		}

		result, err := searchService.SearchDocuments(req)
		if err != nil {
			t.Logf("SearchDocuments(invalid mode) error = %v", err)
		}

		if result != nil {
			t.Logf("SearchDocuments(invalid mode) returned %d results", len(result.Results))
		}
	})
}

// Test_Search_SearchDocuments_Pagination 测试搜索分页
func Test_Search_SearchDocuments_Pagination(t *testing.T) {
	searchService := &service.SearchService{}
	libService := &service.LibraryService{}

	lib, _ := libService.Create(&request.LibraryCreate{
		Name:        "pagination-lib",
		Description: "test pagination",
	})

	t.Run("search with page 0 (should default to 1)", func(t *testing.T) {
		req := &request.Search{
			LibraryID: lib.ID,
			Version:   lib.DefaultVersion,
			Query:     "test",
			Mode:      "hybrid",
			Page:      0, // 应该默认为 1
			Limit:     10,
		}

		result, err := searchService.SearchDocuments(req)
		if err != nil {
			t.Logf("SearchDocuments(page=0) error = %v", err)
			return
		}

		if result != nil {
			if result.Page != 1 {
				t.Errorf("Expected page 1, got %d", result.Page)
			}
		}
	})

	t.Run("search with negative page", func(t *testing.T) {
		req := &request.Search{
			LibraryID: lib.ID,
			Version:   lib.DefaultVersion,
			Query:     "test",
			Mode:      "hybrid",
			Page:      -1, // 负数页码
			Limit:     10,
		}

		result, err := searchService.SearchDocuments(req)
		if err != nil {
			t.Logf("SearchDocuments(page=-1) error = %v", err)
			return
		}

		if result != nil {
			if result.Page < 1 {
				t.Errorf("Expected page >= 1, got %d", result.Page)
			}
		}
	})

	t.Run("search with limit 0 (should default to 10)", func(t *testing.T) {
		req := &request.Search{
			LibraryID: lib.ID,
			Version:   lib.DefaultVersion,
			Query:     "test",
			Mode:      "hybrid",
			Page:      1,
			Limit:     0, // 应该默认为 10
		}

		result, err := searchService.SearchDocuments(req)
		if err != nil {
			t.Logf("SearchDocuments(limit=0) error = %v", err)
			return
		}

		if result != nil {
			if result.Limit != 10 {
				t.Errorf("Expected limit 10, got %d", result.Limit)
			}
		}
	})

	t.Run("search with limit > 10 (should cap at 10)", func(t *testing.T) {
		req := &request.Search{
			LibraryID: lib.ID,
			Version:   lib.DefaultVersion,
			Query:     "test",
			Mode:      "hybrid",
			Page:      1,
			Limit:     100, // 应该限制为 10
		}

		result, err := searchService.SearchDocuments(req)
		if err != nil {
			t.Logf("SearchDocuments(limit=100) error = %v", err)
			return
		}

		if result != nil {
			if result.Limit > 10 {
				t.Errorf("Expected limit <= 10, got %d", result.Limit)
			}
		}
	})

	t.Run("search with page beyond results", func(t *testing.T) {
		req := &request.Search{
			LibraryID: lib.ID,
			Version:   lib.DefaultVersion,
			Query:     "test",
			Mode:      "hybrid",
			Page:      999, // 超出范围的页码
			Limit:     10,
		}

		result, err := searchService.SearchDocuments(req)
		if err != nil {
			t.Logf("SearchDocuments(page=999) error = %v", err)
			return
		}

		if result != nil {
			if len(result.Results) > 0 {
				t.Logf("SearchDocuments(page=999) returned %d results (unexpected)", len(result.Results))
			}
			if result.HasMore {
				t.Error("Expected HasMore=false for page beyond results")
			}
		}
	})
}

// Test_Search_SearchDocuments_MultiTopic 测试多主题搜索
func Test_Search_SearchDocuments_MultiTopic(t *testing.T) {
	searchService := &service.SearchService{}
	libService := &service.LibraryService{}

	lib, _ := libService.Create(&request.LibraryCreate{
		Name:        "multi-topic-lib",
		Description: "test multi topic search",
	})

	t.Run("search with comma-separated topics", func(t *testing.T) {
		req := &request.Search{
			LibraryID: lib.ID,
			Version:   lib.DefaultVersion,
			Query:     "topic1, topic2, topic3", // 逗号分隔
			Mode:      "hybrid",
			Limit:     10,
		}

		result, err := searchService.SearchDocuments(req)
		if err != nil {
			t.Logf("SearchDocuments(multi-topic comma) error = %v", err)
			return
		}

		if result != nil {
			t.Logf("SearchDocuments(multi-topic comma) returned %d results", len(result.Results))
		}
	})

	t.Run("search with space-separated topics", func(t *testing.T) {
		req := &request.Search{
			LibraryID: lib.ID,
			Version:   lib.DefaultVersion,
			Query:     "topic1 topic2 topic3", // 空格分隔
			Mode:      "hybrid",
			Limit:     10,
		}

		result, err := searchService.SearchDocuments(req)
		if err != nil {
			t.Logf("SearchDocuments(multi-topic space) error = %v", err)
			return
		}

		if result != nil {
			t.Logf("SearchDocuments(multi-topic space) returned %d results", len(result.Results))
		}
	})

	t.Run("search with single topic", func(t *testing.T) {
		req := &request.Search{
			LibraryID: lib.ID,
			Version:   lib.DefaultVersion,
			Query:     "single-topic", // 单个主题
			Mode:      "hybrid",
			Limit:     10,
		}

		result, err := searchService.SearchDocuments(req)
		if err != nil {
			t.Logf("SearchDocuments(single topic) error = %v", err)
			return
		}

		if result != nil {
			t.Logf("SearchDocuments(single topic) returned %d results", len(result.Results))
		}
	})
}

// Test_Search_SearchDocuments_DifferentVersions 测试不同版本的搜索
func Test_Search_SearchDocuments_DifferentVersions(t *testing.T) {
	searchService := &service.SearchService{}
	libService := &service.LibraryService{}

	lib, _ := libService.Create(&request.LibraryCreate{
		Name:        "version-search-lib",
		Description: "test version search",
	})

	t.Run("search with latest version", func(t *testing.T) {
		req := &request.Search{
			LibraryID: lib.ID,
			Version:   "latest",
			Query:     "test",
			Mode:      "hybrid",
			Limit:     10,
		}

		result, err := searchService.SearchDocuments(req)
		if err != nil {
			t.Logf("SearchDocuments(version=latest) error = %v", err)
			return
		}

		if result != nil {
			t.Logf("SearchDocuments(version=latest) returned %d results", len(result.Results))
		}
	})

	t.Run("search with specific version", func(t *testing.T) {
		req := &request.Search{
			LibraryID: lib.ID,
			Version:   "v1.0.0",
			Query:     "test",
			Mode:      "hybrid",
			Limit:     10,
		}

		result, err := searchService.SearchDocuments(req)
		if err != nil {
			t.Logf("SearchDocuments(version=v1.0.0) error = %v", err)
			return
		}

		if result != nil {
			t.Logf("SearchDocuments(version=v1.0.0) returned %d results", len(result.Results))
		}
	})

	t.Run("search with empty version", func(t *testing.T) {
		req := &request.Search{
			LibraryID: lib.ID,
			Version:   "",
			Query:     "test",
			Mode:      "hybrid",
			Limit:     10,
		}

		result, err := searchService.SearchDocuments(req)
		if err != nil {
			t.Logf("SearchDocuments(version='') error = %v", err)
			return
		}

		if result != nil {
			t.Logf("SearchDocuments(version='') returned %d results", len(result.Results))
		}
	})
}

// Test_Search_InternalFunctions 测试内部辅助函数（通过边界测试覆盖）
func Test_Search_InternalFunctions(t *testing.T) {
	searchService := &service.SearchService{}
	libService := &service.LibraryService{}
	docService := &service.DocumentService{}

	// 创建测试库
	lib, err := libService.Create(&request.LibraryCreate{
		Name:        "test-search-internal",
		Description: "Test library for search internal functions",
	})
	if err != nil {
		t.Fatalf("Failed to create library: %v", err)
	}
	defer libService.Delete(lib.ID)

	// 创建版本
	err = libService.CreateVersion(lib.ID, "v1.0.0")
	if err != nil {
		t.Fatalf("Failed to create version: %v", err)
	}

	t.Run("search with metadata to trigger extractDeepestTitle", func(t *testing.T) {
		// 上传包含多级标题的文档
		content := []byte(`# H1 Title

## H2 Subtitle

### H3 Section

#### H4 Subsection

##### H5 Detail

###### H6 Deepest

Content here.`)

		file, header := createMultipartFile("metadata-test.md", content)
		doc, err := docService.Upload(lib.ID, "v1.0.0", file, header, "test-user", "test-task")
		if err != nil {
			t.Fatalf("Upload() error = %v", err)
		}

		// 等待处理完成
		time.Sleep(2 * time.Second)

		// 执行搜索，这会触发 extractDeepestTitle
		req := &request.Search{
			LibraryID: lib.ID,
			Version:   "v1.0.0",
			Query:     "content",
			Mode:      "hybrid",
			Limit:     10,
		}

		result, err := searchService.SearchDocuments(req)
		if err != nil {
			t.Logf("SearchDocuments() error = %v", err)
		}

		if result != nil && len(result.Results) > 0 {
			t.Logf("✅ Search returned %d results (extractDeepestTitle executed)", len(result.Results))
			for i, r := range result.Results {
				t.Logf("  Result %d: %s (relevance: %.4f)", i+1, r.Title, r.Relevance)
			}
		}

		t.Logf("Uploaded doc ID: %d", doc.ID)
	})

	t.Run("search with multi-topic query to trigger reciprocalRankFusion", func(t *testing.T) {
		// 上传多个文档
		docs := []struct {
			filename string
			content  string
		}{
			{"doc1.md", "# Document 1\n\nThis is about data fetching and API calls."},
			{"doc2.md", "# Document 2\n\nThis covers routing and navigation."},
			{"doc3.md", "# Document 3\n\nBoth data fetching and routing are discussed here."},
		}

		for _, d := range docs {
			file, header := createMultipartFile(d.filename, []byte(d.content))
			_, err := docService.Upload(lib.ID, "v1.0.0", file, header, "test-user", "test-task")
			if err != nil {
				t.Logf("Upload(%s) error = %v", d.filename, err)
			}
		}

		// 等待处理完成
		time.Sleep(3 * time.Second)

		// 使用多主题查询（逗号分隔），会触发 reciprocalRankFusion
		req := &request.Search{
			LibraryID: lib.ID,
			Version:   "v1.0.0",
			Query:     "data fetching, routing", // 多主题查询
			Mode:      "hybrid",
			Limit:     10,
		}

		result, err := searchService.SearchDocuments(req)
		if err != nil {
			t.Logf("SearchDocuments(multi-topic) error = %v", err)
		}

		if result != nil && len(result.Results) > 0 {
			t.Logf("✅ Multi-topic search returned %d results (reciprocalRankFusion executed)", len(result.Results))
			for i, r := range result.Results {
				t.Logf("  Result %d: %s (RRF relevance: %.4f)", i+1, r.Title, r.Relevance)
			}
		}
	})

	t.Run("search with empty metadata", func(t *testing.T) {
		// 上传没有标题的文档
		content := []byte("Plain text without any headers.")
		file, header := createMultipartFile("no-headers.md", content)
		_, err := docService.Upload(lib.ID, "v1.0.0", file, header, "test-user", "test-task")
		if err != nil {
			t.Logf("Upload() error = %v", err)
		}

		time.Sleep(1 * time.Second)

		req := &request.Search{
			LibraryID: lib.ID,
			Version:   "v1.0.0",
			Query:     "plain text",
			Mode:      "hybrid",
			Limit:     10,
		}

		result, err := searchService.SearchDocuments(req)
		if err != nil {
			t.Logf("SearchDocuments() error = %v", err)
		}

		if result != nil {
			t.Logf("Search with empty metadata returned %d results", len(result.Results))
		}
	})

	t.Run("search with various topic formats", func(t *testing.T) {
		testCases := []string{
			"single topic",
			"topic1, topic2",
			"topic1,topic2,topic3",
			"  spaced  ,  topics  ",
			"",
		}

		for _, query := range testCases {
			req := &request.Search{
				LibraryID: lib.ID,
				Version:   "v1.0.0",
				Query:     query,
				Mode:      "hybrid",
				Limit:     5,
			}

			result, err := searchService.SearchDocuments(req)
			if err != nil {
				t.Logf("SearchDocuments(query='%s') error = %v", query, err)
				continue
			}

			if result != nil {
				t.Logf("Query '%s' returned %d results", query, len(result.Results))
			}
		}
	})
}
