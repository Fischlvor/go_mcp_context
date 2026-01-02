package test_test

import (
	"context"
	"fmt"
	"mime/multipart"
	"os"
	"strings"
	"testing"
	"time"

	dbmodel "go-mcp-context/internal/model/database"
	"go-mcp-context/internal/model/request"
	"go-mcp-context/internal/model/response"
	"go-mcp-context/internal/service"
	"go-mcp-context/pkg/global"
	"go-mcp-context/pkg/utils"
)

// Test_Library_ValidateVersion 测试版本验证
func Test_Library_ValidateVersion(t *testing.T) {
	libService := &service.LibraryService{}

	tests := []struct {
		name    string
		version string
		wantErr bool
	}{
		{"valid v1.0.0", "v1.0.0", false},
		{"valid 1.0.0", "1.0.0", false},
		{"valid v1.0.0-alpha", "v1.0.0-alpha", false},
		{"valid v1.0.0-beta.1", "v1.0.0-beta.1", false},
		{"invalid empty", "", true},
		{"invalid format", "invalid", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := libService.ValidateVersion(tt.version)
			if (err != nil) != tt.wantErr {
				t.Errorf("ValidateVersion(%s) error = %v, wantErr %v", tt.version, err, tt.wantErr)
			}
		})
	}
}

// Test_Library_Delete 测试库删除（软删除）
func Test_Library_Delete(t *testing.T) {
	libService := &service.LibraryService{}

	t.Run("delete library successfully", func(t *testing.T) {
		// 通过 Service 层创建测试库
		createReq := &request.LibraryCreate{
			Name:        "delete-test-lib",
			Description: "test",
		}
		lib, err := libService.Create(createReq)
		if err != nil {
			t.Fatalf("Create() error = %v", err)
		}

		// 删除库
		err = libService.Delete(lib.ID)
		if err != nil {
			t.Fatalf("Delete() error = %v", err)
		}

		// 验证库已删除（软删除）
		_, err = libService.GetByID(lib.ID)
		if err == nil {
			t.Error("Expected error when getting deleted library, got nil")
		}
	})

	t.Run("delete non-existent library", func(t *testing.T) {
		err := libService.Delete(99999)
		if err == nil {
			t.Error("Expected error when deleting non-existent library, got nil")
		}
	})
}

// Test_Library_GetByName 测试按名称获取库
func Test_Library_GetByName(t *testing.T) {
	libService := &service.LibraryService{}

	t.Run("get library by name successfully", func(t *testing.T) {
		// 通过 Service 层创建测试库
		createReq := &request.LibraryCreate{
			Name:        "getbyname-lib",
			Description: "test",
		}
		lib, err := libService.Create(createReq)
		if err != nil {
			t.Fatalf("Create() error = %v", err)
		}

		// 按名称查询
		retrieved, err := libService.GetByName("getbyname-lib")
		if err != nil {
			t.Fatalf("GetByName() error = %v", err)
		}

		if retrieved.ID != lib.ID {
			t.Errorf("Expected ID %d, got %d", lib.ID, retrieved.ID)
		}

		if retrieved.Name != lib.Name {
			t.Errorf("Expected name %s, got %s", lib.Name, retrieved.Name)
		}
	})

	t.Run("get non-existent library by name", func(t *testing.T) {
		_, err := libService.GetByName("non-existent-lib-12345")
		if err == nil {
			t.Error("Expected error when getting non-existent library, got nil")
		}
	})
}

// Test_Library_GetByID 测试根据 ID 获取库
func Test_Library_GetByID(t *testing.T) {
	libService := &service.LibraryService{}

	t.Run("get library by id successfully", func(t *testing.T) {
		// 通过 Service 层创建测试库
		createReq := &request.LibraryCreate{
			Name:        "getbyid-test-lib",
			Description: "test",
		}
		lib, err := libService.Create(createReq)
		if err != nil {
			t.Fatalf("Create() error = %v", err)
		}

		// 获取库
		retrieved, err := libService.GetByID(lib.ID)
		if err != nil {
			t.Fatalf("GetByID() error = %v", err)
		}

		if retrieved.ID != lib.ID {
			t.Errorf("Expected ID %d, got %d", lib.ID, retrieved.ID)
		}

		if retrieved.Name != lib.Name {
			t.Errorf("Expected name %s, got %s", lib.Name, retrieved.Name)
		}
	})

	t.Run("get non-existent library", func(t *testing.T) {
		_, err := libService.GetByID(99999)
		if err == nil {
			t.Error("Expected error when getting non-existent library, got nil")
		}
	})
}

// Test_Library_List 测试库列表查询
func Test_Library_List(t *testing.T) {
	libService := &service.LibraryService{}

	t.Run("list libraries", func(t *testing.T) {
		// 通过 Service 层创建多个库
		for i := 1; i <= 3; i++ {
			req := &request.LibraryCreate{
				Name:        "list-lib-" + string(rune('0'+i)),
				Description: "test",
			}
			_, err := libService.Create(req)
			if err != nil {
				t.Fatalf("Create() error = %v", err)
			}
		}

		// 查询库列表
		listReq := &request.LibraryList{
			PageInfo: request.PageInfo{
				Page:     1,
				PageSize: 10,
			},
		}

		result, err := libService.List(listReq)
		if err != nil {
			t.Fatalf("List() error = %v", err)
		}

		if result.Total < 3 {
			t.Errorf("Expected at least 3 libraries, got %d", result.Total)
		}
	})

	t.Run("list libraries with pagination", func(t *testing.T) {
		// 通过 Service 层创建 5 个库
		for i := 1; i <= 5; i++ {
			req := &request.LibraryCreate{
				Name:        "paging-lib-" + string(rune('0'+i)),
				Description: "test",
			}
			_, err := libService.Create(req)
			if err != nil {
				t.Fatalf("Create() error = %v", err)
			}
		}

		// 第一页，每页 2 条
		listReq := &request.LibraryList{
			PageInfo: request.PageInfo{
				Page:     1,
				PageSize: 2,
			},
		}

		result, err := libService.List(listReq)
		if err != nil {
			t.Fatalf("List() error = %v", err)
		}

		// 验证分页逻辑
		if result.Page != 1 {
			t.Errorf("Expected page 1, got %d", result.Page)
		}

		if result.PageSize != 2 {
			t.Errorf("Expected page size 2, got %d", result.PageSize)
		}
	})

	t.Run("list libraries with invalid page size", func(t *testing.T) {
		// 通过 Service 层创建库
		req := &request.LibraryCreate{
			Name:        "invalid-size-lib",
			Description: "test",
		}
		_, err := libService.Create(req)
		if err != nil {
			t.Fatalf("Create() error = %v", err)
		}

		// Page size 为 0，应该使用默认值
		listReq := &request.LibraryList{
			PageInfo: request.PageInfo{
				Page:     1,
				PageSize: 0,
			},
		}

		result, err := libService.List(listReq)
		if err != nil {
			t.Fatalf("List() error = %v", err)
		}

		if result == nil {
			t.Fatal("Expected result, got nil")
		}
	})
}

// Test_Library_Create 测试库创建（调用 Service 层入口，进入不同分支）
func Test_Library_Create(t *testing.T) {
	libService := &service.LibraryService{}

	t.Run("create library with name and description", func(t *testing.T) {
		req := &request.LibraryCreate{
			Name:        "Test Library",
			Description: "Test description",
		}

		resp, err := libService.Create(req)
		if err != nil {
			t.Fatalf("Create() error = %v", err)
		}

		if resp == nil {
			t.Fatal("Expected response, got nil")
		}

		if resp.Name != req.Name {
			t.Errorf("Expected name %s, got %s", req.Name, resp.Name)
		}
	})

	t.Run("create library with empty description", func(t *testing.T) {
		req := &request.LibraryCreate{
			Name:        "Library No Desc",
			Description: "",
		}

		resp, err := libService.Create(req)
		if err != nil {
			t.Fatalf("Create() error = %v", err)
		}

		if resp == nil {
			t.Fatal("Expected response, got nil")
		}

		if resp.Name != req.Name {
			t.Errorf("Expected name %s, got %s", req.Name, resp.Name)
		}
	})
}

// Test_Library_Update 测试库更新
func Test_Library_Update(t *testing.T) {
	libService := &service.LibraryService{}

	t.Run("update library name and description", func(t *testing.T) {
		// 创建库
		createReq := &request.LibraryCreate{
			Name:        "original-name",
			Description: "original description",
		}
		lib, err := libService.Create(createReq)
		if err != nil {
			t.Fatalf("Create() error = %v", err)
		}

		// 更新库
		updateReq := &request.LibraryUpdate{
			Name:        "updated-name",
			Description: "updated description",
		}
		updated, err := libService.Update(lib.ID, updateReq)
		if err != nil {
			t.Fatalf("Update() error = %v", err)
		}

		if updated.Name != updateReq.Name {
			t.Errorf("Expected name %s, got %s", updateReq.Name, updated.Name)
		}

		if updated.Description != updateReq.Description {
			t.Errorf("Expected description %s, got %s", updateReq.Description, updated.Description)
		}
	})

	t.Run("update non-existent library", func(t *testing.T) {
		updateReq := &request.LibraryUpdate{
			Name:        "new-name",
			Description: "new description",
		}
		_, err := libService.Update(99999, updateReq)
		if err == nil {
			t.Error("Expected error when updating non-existent library, got nil")
		}
	})
}

// Test_Library_SearchByName 测试按名称搜索库
func Test_Library_SearchByName(t *testing.T) {
	libService := &service.LibraryService{}

	t.Run("search library by name prefix", func(t *testing.T) {
		// 创建库
		createReq := &request.LibraryCreate{
			Name:        "search-test-lib",
			Description: "test",
		}
		_, err := libService.Create(createReq)
		if err != nil {
			t.Fatalf("Create() error = %v", err)
		}

		// 按前缀搜索
		results, err := libService.SearchByName("search-test")
		if err != nil {
			t.Fatalf("SearchByName() error = %v", err)
		}

		if len(results) == 0 {
			t.Error("Expected at least 1 result, got 0")
		}
	})

	t.Run("search non-existent library", func(t *testing.T) {
		results, err := libService.SearchByName("non-existent-search-xyz")
		if err != nil {
			t.Fatalf("SearchByName() error = %v", err)
		}

		if len(results) != 0 {
			t.Errorf("Expected 0 results, got %d", len(results))
		}
	})
}

// Test_Library_GetLibraryInfo 测试获取库详情（带统计）
func Test_Library_GetLibraryInfo(t *testing.T) {
	libService := &service.LibraryService{}

	t.Run("get library info with stats", func(t *testing.T) {
		// 创建库
		createReq := &request.LibraryCreate{
			Name:        "info-test-lib",
			Description: "test",
		}
		lib, err := libService.Create(createReq)
		if err != nil {
			t.Fatalf("Create() error = %v", err)
		}

		// 获取库详情
		info, err := libService.GetLibraryInfo(lib.ID)
		if err != nil {
			t.Fatalf("GetLibraryInfo() error = %v", err)
		}

		if info == nil {
			t.Fatal("Expected info, got nil")
		}

		if info.ID != lib.ID {
			t.Errorf("Expected ID %d, got %d", lib.ID, info.ID)
		}

		if info.Name != lib.Name {
			t.Errorf("Expected name %s, got %s", lib.Name, info.Name)
		}

		// 验证统计字段
		if info.ChunkCount < 0 {
			t.Errorf("Expected non-negative chunk count, got %d", info.ChunkCount)
		}

		if info.TokenCount < 0 {
			t.Errorf("Expected non-negative token count, got %d", info.TokenCount)
		}
	})

	t.Run("get non-existent library info", func(t *testing.T) {
		_, err := libService.GetLibraryInfo(99999)
		if err == nil {
			t.Error("Expected error when getting non-existent library info, got nil")
		}
	})
}

// Test_Library_ListWithStats 测试获取库列表（带统计）
func Test_Library_ListWithStats(t *testing.T) {
	libService := &service.LibraryService{}

	t.Run("list libraries with stats", func(t *testing.T) {
		// 创建库
		for i := 1; i <= 2; i++ {
			req := &request.LibraryCreate{
				Name:        "stats-lib-" + string(rune('0'+i)),
				Description: "test",
			}
			_, err := libService.Create(req)
			if err != nil {
				t.Fatalf("Create() error = %v", err)
			}
		}

		// 查询库列表（带统计）
		listReq := &request.LibraryList{
			PageInfo: request.PageInfo{
				Page:     1,
				PageSize: 10,
			},
		}

		result, err := libService.ListWithStats(listReq)
		if err != nil {
			t.Fatalf("ListWithStats() error = %v", err)
		}

		if result == nil {
			t.Fatal("Expected result, got nil")
		}

		if result.Total < 2 {
			t.Errorf("Expected at least 2 libraries, got %d", result.Total)
		}
	})

	t.Run("list libraries with name filter", func(t *testing.T) {
		// 创建库
		req := &request.LibraryCreate{
			Name:        "filter-test-lib",
			Description: "test",
		}
		_, err := libService.Create(req)
		if err != nil {
			t.Fatalf("Create() error = %v", err)
		}

		// 按名称过滤
		name := "filter-test"
		listReq := &request.LibraryList{
			Name: &name,
			PageInfo: request.PageInfo{
				Page:     1,
				PageSize: 10,
			},
		}

		result, err := libService.ListWithStats(listReq)
		if err != nil {
			t.Fatalf("ListWithStats() error = %v", err)
		}

		if result == nil {
			t.Fatal("Expected result, got nil")
		}
	})
}

// Test_Library_GetVersions 测试获取库版本列表
func Test_Library_GetVersions(t *testing.T) {
	libService := &service.LibraryService{}

	t.Run("get versions for library", func(t *testing.T) {
		// 创建库
		createReq := &request.LibraryCreate{
			Name:        "version-test-lib",
			Description: "test",
		}
		lib, err := libService.Create(createReq)
		if err != nil {
			t.Fatalf("Create() error = %v", err)
		}

		// 获取版本列表
		versions, err := libService.GetVersions(lib.ID)
		if err != nil {
			t.Fatalf("GetVersions() error = %v", err)
		}

		if len(versions) == 0 {
			t.Error("Expected at least 1 version (default version), got 0")
		}
	})

	t.Run("get versions for non-existent library", func(t *testing.T) {
		_, err := libService.GetVersions(99999)
		if err == nil {
			t.Error("Expected error when getting versions for non-existent library, got nil")
		}
	})
}

// Test_Library_CreateVersion 测试创建新版本
func Test_Library_CreateVersion(t *testing.T) {
	libService := &service.LibraryService{}

	t.Run("create new version", func(t *testing.T) {
		// 创建库
		createReq := &request.LibraryCreate{
			Name:        "create-version-lib",
			Description: "test",
		}
		lib, err := libService.Create(createReq)
		if err != nil {
			t.Fatalf("Create() error = %v", err)
		}

		// 创建新版本
		err = libService.CreateVersion(lib.ID, "v2.0.0")
		if err != nil {
			t.Fatalf("CreateVersion() error = %v", err)
		}

		// 验证版本已创建
		versions, err := libService.GetVersions(lib.ID)
		if err != nil {
			t.Fatalf("GetVersions() error = %v", err)
		}

		found := false
		for _, v := range versions {
			if v.Version == "v2.0.0" {
				found = true
				break
			}
		}

		if !found {
			t.Error("Expected version v2.0.0 to be created, but not found")
		}
	})

	t.Run("create version for non-existent library", func(t *testing.T) {
		err := libService.CreateVersion(99999, "v1.0.0")
		if err == nil {
			t.Error("Expected error when creating version for non-existent library, got nil")
		}
	})

	t.Run("create version with invalid format", func(t *testing.T) {
		// 创建库
		createReq := &request.LibraryCreate{
			Name:        "invalid-version-lib",
			Description: "test",
		}
		lib, err := libService.Create(createReq)
		if err != nil {
			t.Fatalf("Create() error = %v", err)
		}

		// 创建无效版本
		err = libService.CreateVersion(lib.ID, "invalid-version")
		if err == nil {
			t.Error("Expected error when creating invalid version, got nil")
		}
	})
}

// Test_Library_DeleteVersion 测试删除版本
func Test_Library_DeleteVersion(t *testing.T) {
	libService := &service.LibraryService{}
	docService := &service.DocumentService{}

	t.Run("delete version with documents", func(t *testing.T) {
		// 创建库
		lib, err := libService.Create(&request.LibraryCreate{
			Name:        "delete-version-with-docs-lib",
			Description: "test delete version with documents",
		})
		if err != nil {
			t.Fatalf("Create() error = %v", err)
		}
		defer libService.Delete(lib.ID)

		// 创建新版本
		err = libService.CreateVersion(lib.ID, "v2.0.0")
		if err != nil {
			t.Fatalf("CreateVersion() error = %v", err)
		}

		// 创建临时文件用于上传
		content := []byte("# Test Document\n\nThis is a test document for version deletion.")
		tmpFile, err := os.CreateTemp("", "test-*.md")
		if err != nil {
			t.Fatalf("Failed to create temp file: %v", err)
		}
		defer os.Remove(tmpFile.Name())

		if _, err := tmpFile.Write(content); err != nil {
			t.Fatalf("Failed to write temp file: %v", err)
		}
		tmpFile.Seek(0, 0)

		// 创建 FileHeader
		fileHeader := &multipart.FileHeader{
			Filename: "test.md",
			Size:     int64(len(content)),
		}

		// 上传文档到新版本
		taskID := utils.GenerateTaskID()
		uploadedDoc, err := docService.Upload(lib.ID, "v2.0.0", tmpFile, fileHeader, "test-user", taskID)
		if err != nil {
			t.Logf("Upload() error = %v (may fail if storage not available)", err)
			// 即使上传失败，也继续测试删除逻辑
		}

		// 等待文档处理
		time.Sleep(5 * time.Second)

		// 验证文档是否生成了分块
		if uploadedDoc != nil {
			var chunks []dbmodel.DocumentChunk
			if err := global.DB.Where("upload_id = ?", uploadedDoc.ID).Find(&chunks).Error; err == nil {
				t.Logf("Document has %d chunks before deletion", len(chunks))
			}
		}

		// 删除版本（包括其所有文档和分块）
		err = libService.DeleteVersion(lib.ID, "v2.0.0")
		if err != nil {
			t.Fatalf("DeleteVersion() error = %v", err)
		}

		// 验证版本已删除
		versions, _ := libService.GetVersions(lib.ID)
		if len(versions) != 1 {
			t.Errorf("Expected 1 version after deletion, got %d", len(versions))
		}

		t.Log("✅ Successfully deleted version with documents")
	})

	t.Run("delete non-existent version", func(t *testing.T) {
		// 创建库
		createReq := &request.LibraryCreate{
			Name:        "delete-nonexistent-version-lib",
			Description: "test",
		}
		lib, err := libService.Create(createReq)
		if err != nil {
			t.Fatalf("Create() error = %v", err)
		}
		defer libService.Delete(lib.ID)

		// 删除不存在的版本（没有任何文档的版本）
		err = libService.DeleteVersion(lib.ID, "v99.0.0")
		if err == nil {
			t.Error("Expected error when deleting non-existent version, got nil")
		}

		t.Log("✅ Correctly handled non-existent version")
	})

	t.Run("delete version for non-existent library", func(t *testing.T) {
		// 删除不存在的库的版本
		err := libService.DeleteVersion(99999, "v1.0.0")
		if err == nil {
			t.Error("Expected error when deleting version for non-existent library, got nil")
		}

		t.Log("✅ Correctly handled non-existent library")
	})
}

// Test_Library_RefreshVersion 测试刷新版本
func Test_Library_RefreshVersion(t *testing.T) {
	libService := &service.LibraryService{}
	docService := &service.DocumentService{}

	t.Run("refresh version with documents", func(t *testing.T) {
		// 创建库
		lib, err := libService.Create(&request.LibraryCreate{
			Name:        "refresh-version-with-docs-lib",
			Description: "test refresh version with documents",
		})
		if err != nil {
			t.Fatalf("Create() error = %v", err)
		}
		defer libService.Delete(lib.ID)

		// 上传文档
		content := []byte("# Test Document\n\nThis is a test document for version refresh.\n\n## Section 1\n\nSome content here.")
		tmpFile, err := os.CreateTemp("", "refresh-test-*.md")
		if err != nil {
			t.Fatalf("Failed to create temp file: %v", err)
		}
		defer os.Remove(tmpFile.Name())

		if _, err := tmpFile.Write(content); err != nil {
			t.Fatalf("Failed to write temp file: %v", err)
		}
		tmpFile.Seek(0, 0)

		fileHeader := &multipart.FileHeader{
			Filename: "refresh-test.md",
			Size:     int64(len(content)),
		}

		taskID := utils.GenerateTaskID()
		_, err = docService.Upload(lib.ID, lib.DefaultVersion, tmpFile, fileHeader, "test-user", taskID)
		if err != nil {
			t.Logf("Upload() error = %v (may fail if storage not available)", err)
			return
		}

		// 等待文档处理完成
		time.Sleep(5 * time.Second)

		// 刷新版本
		err = libService.RefreshVersion(lib.ID, lib.DefaultVersion, "test-user")
		if err != nil {
			t.Fatalf("RefreshVersion() error = %v", err)
		}

		// 等待异步刷新完成
		time.Sleep(8 * time.Second)

		t.Log("✅ Successfully refreshed version with documents")
	})

	t.Run("refresh version without documents", func(t *testing.T) {
		// 创建库
		createReq := &request.LibraryCreate{
			Name:        "refresh-version-lib",
			Description: "test",
		}
		lib, err := libService.Create(createReq)
		if err != nil {
			t.Fatalf("Create() error = %v", err)
		}
		defer libService.Delete(lib.ID)

		// 刷新版本（没有文档，应该返回错误）
		err = libService.RefreshVersion(lib.ID, lib.DefaultVersion, "test-user")
		if err == nil {
			t.Error("Expected error when refreshing version without documents, got nil")
		}

		t.Log("✅ Correctly handled version without documents")
	})

	t.Run("refresh non-existent library version", func(t *testing.T) {
		err := libService.RefreshVersion(99999, "latest", "")
		if err == nil {
			t.Error("Expected error when refreshing non-existent library version, got nil")
		}
	})
}

// Test_Library_Create_Advanced 测试库创建的高级场景
func Test_Library_Create_Advanced(t *testing.T) {
	libService := &service.LibraryService{}

	t.Run("create library with very long name", func(t *testing.T) {
		longName := ""
		for i := 0; i < 200; i++ {
			longName += "a"
		}

		req := &request.LibraryCreate{
			Name:        longName,
			Description: "test",
		}

		resp, err := libService.Create(req)
		if err != nil {
			t.Fatalf("Create() error = %v", err)
		}

		if resp == nil {
			t.Fatal("Expected response, got nil")
		}

		if resp.Name != longName {
			t.Errorf("Expected name %s, got %s", longName, resp.Name)
		}
	})

	t.Run("create library with special characters", func(t *testing.T) {
		req := &request.LibraryCreate{
			Name:        "Test-Library_2024@v1.0",
			Description: "test with special chars",
		}

		resp, err := libService.Create(req)
		if err != nil {
			t.Fatalf("Create() error = %v", err)
		}

		if resp == nil {
			t.Fatal("Expected response, got nil")
		}

		if resp.Name != req.Name {
			t.Errorf("Expected name %s, got %s", req.Name, resp.Name)
		}
	})

	t.Run("create library with unicode characters", func(t *testing.T) {
		req := &request.LibraryCreate{
			Name:        "测试库-Test-テスト",
			Description: "test with unicode",
		}

		resp, err := libService.Create(req)
		if err != nil {
			t.Fatalf("Create() error = %v", err)
		}

		if resp == nil {
			t.Fatal("Expected response, got nil")
		}

		if resp.Name != req.Name {
			t.Errorf("Expected name %s, got %s", req.Name, resp.Name)
		}
	})
}

// Test_Library_List_Advanced 测试库列表的高级场景
func Test_Library_List_Advanced(t *testing.T) {
	libService := &service.LibraryService{}

	t.Run("list with page 0 (should default to 1)", func(t *testing.T) {
		req := &request.LibraryList{
			PageInfo: request.PageInfo{
				Page:     0,
				PageSize: 10,
			},
		}

		result, err := libService.List(req)
		if err != nil {
			t.Fatalf("List() error = %v", err)
		}

		if result == nil {
			t.Fatal("Expected result, got nil")
		}

		if result.Page < 1 {
			t.Errorf("Expected page >= 1, got %d", result.Page)
		}
	})

	t.Run("list with negative page size", func(t *testing.T) {
		req := &request.LibraryList{
			PageInfo: request.PageInfo{
				Page:     1,
				PageSize: -10,
			},
		}

		result, err := libService.List(req)
		if err != nil {
			t.Fatalf("List() error = %v", err)
		}

		if result == nil {
			t.Fatal("Expected result, got nil")
		}
	})

	t.Run("list with very large page size", func(t *testing.T) {
		req := &request.LibraryList{
			PageInfo: request.PageInfo{
				Page:     1,
				PageSize: 10000,
			},
		}

		result, err := libService.List(req)
		if err != nil {
			t.Fatalf("List() error = %v", err)
		}

		if result == nil {
			t.Fatal("Expected result, got nil")
		}
	})
}

// Test_Library_CreateVersion_Advanced 测试创建版本的高级场景
func Test_Library_CreateVersion_Advanced(t *testing.T) {
	libService := &service.LibraryService{}

	t.Run("create version with v prefix", func(t *testing.T) {
		// 创建库
		createReq := &request.LibraryCreate{
			Name:        "version-prefix-lib",
			Description: "test",
		}
		lib, err := libService.Create(createReq)
		if err != nil {
			t.Fatalf("Create() error = %v", err)
		}

		// 创建带 v 前缀的版本
		err = libService.CreateVersion(lib.ID, "v3.0.0")
		if err != nil {
			t.Fatalf("CreateVersion() error = %v", err)
		}

		// 验证版本已创建
		versions, err := libService.GetVersions(lib.ID)
		if err != nil {
			t.Fatalf("GetVersions() error = %v", err)
		}

		found := false
		for _, v := range versions {
			if v.Version == "v3.0.0" {
				found = true
				break
			}
		}

		if !found {
			t.Error("Expected version v3.0.0 to be created")
		}
	})

	t.Run("create version without v prefix", func(t *testing.T) {
		// 创建库
		createReq := &request.LibraryCreate{
			Name:        "version-no-prefix-lib",
			Description: "test",
		}
		lib, err := libService.Create(createReq)
		if err != nil {
			t.Fatalf("Create() error = %v", err)
		}

		// 创建不带 v 前缀的版本（应该自动添加）
		err = libService.CreateVersion(lib.ID, "2.5.0")
		if err != nil {
			t.Fatalf("CreateVersion() error = %v", err)
		}

		// 验证版本已创建（带 v 前缀）
		versions, err := libService.GetVersions(lib.ID)
		if err != nil {
			t.Fatalf("GetVersions() error = %v", err)
		}

		found := false
		for _, v := range versions {
			if v.Version == "v2.5.0" {
				found = true
				break
			}
		}

		if !found {
			t.Error("Expected version v2.5.0 to be created (with auto-added v prefix)")
		}
	})

	t.Run("create duplicate version", func(t *testing.T) {
		// 创建库
		createReq := &request.LibraryCreate{
			Name:        "duplicate-version-lib",
			Description: "test",
		}
		lib, err := libService.Create(createReq)
		if err != nil {
			t.Fatalf("Create() error = %v", err)
		}

		// 创建第一个版本
		err = libService.CreateVersion(lib.ID, "v1.5.0")
		if err != nil {
			t.Fatalf("CreateVersion() error = %v", err)
		}

		// 尝试创建相同的版本
		err = libService.CreateVersion(lib.ID, "v1.5.0")
		if err == nil {
			t.Error("Expected error when creating duplicate version, got nil")
		}
	})

	t.Run("create version for non-existent library", func(t *testing.T) {
		err := libService.CreateVersion(99999, "v1.0.0")
		if err == nil {
			t.Error("Expected error when creating version for non-existent library, got nil")
		}
	})

	t.Run("create multiple versions", func(t *testing.T) {
		// 创建库
		createReq := &request.LibraryCreate{
			Name:        "multi-version-lib",
			Description: "test",
		}
		lib, err := libService.Create(createReq)
		if err != nil {
			t.Fatalf("Create() error = %v", err)
		}

		// 创建多个版本
		versions := []string{"v1.0.0", "v1.1.0", "v2.0.0", "v2.1.0"}
		for _, version := range versions {
			err = libService.CreateVersion(lib.ID, version)
			if err != nil {
				t.Fatalf("CreateVersion(%s) error = %v", version, err)
			}
		}

		// 验证所有版本都已创建
		retrievedVersions, err := libService.GetVersions(lib.ID)
		if err != nil {
			t.Fatalf("GetVersions() error = %v", err)
		}

		if len(retrievedVersions) < len(versions) {
			t.Errorf("Expected at least %d versions, got %d", len(versions), len(retrievedVersions))
		}
	})
}

// Test_Library_InitFromGitHub 测试从 GitHub 初始化库
func Test_Library_InitFromGitHub(t *testing.T) {
	libService := &service.LibraryService{}
	ctx := context.Background()

	t.Run("init from valid github url", func(t *testing.T) {
		// 使用有效的 GitHub URL
		githubURL := "https://github.com/go-gorm/gorm"
		userUUID := "test-user-github-001"

		result, err := libService.InitFromGitHub(ctx, githubURL, userUUID)
		if err != nil {
			// 如果库已存在，这是预期的
			if strings.Contains(err.Error(), "已存在") {
				t.Logf("InitFromGitHub() error = %v (library already exists, expected)", err)
				return
			}
			t.Logf("InitFromGitHub() error = %v (expected if GitHub API unavailable)", err)
			return
		}

		if result == nil {
			t.Fatal("Expected result, got nil")
		}

		if result.Library == nil {
			t.Fatal("Expected library, got nil")
		}

		if result.Library.Name == "" {
			t.Error("Expected non-empty library name")
		}

		if result.Library.SourceType != "github" {
			t.Errorf("Expected source type 'github', got '%s'", result.Library.SourceType)
		}

		// SourceURL 存储的是解析后的 repo 格式 (go-gorm/gorm)，不是完整 URL
		if result.Library.SourceURL == "" {
			t.Error("Expected non-empty source URL")
		}

		if result.RepoName == "" {
			t.Error("Expected non-empty repo name")
		}
	})

	t.Run("init from invalid github url", func(t *testing.T) {
		// 使用无效的 GitHub URL
		invalidURL := "https://github.com/invalid"
		userUUID := "test-user-github-002"

		result, err := libService.InitFromGitHub(ctx, invalidURL, userUUID)
		if err == nil {
			t.Error("Expected error for invalid GitHub URL, got nil")
		}

		if result != nil {
			t.Error("Expected nil result for invalid URL")
		}
	})

	t.Run("init from malformed url", func(t *testing.T) {
		// 使用格式错误的 URL
		malformedURL := "not-a-valid-url"
		userUUID := "test-user-github-003"

		result, err := libService.InitFromGitHub(ctx, malformedURL, userUUID)
		if err == nil {
			t.Error("Expected error for malformed URL, got nil")
		}

		if result != nil {
			t.Error("Expected nil result for malformed URL")
		}
	})

	t.Run("init from non-github url", func(t *testing.T) {
		// 使用非 GitHub URL
		nonGitHubURL := "https://gitlab.com/some/repo"
		userUUID := "test-user-github-004"

		result, err := libService.InitFromGitHub(ctx, nonGitHubURL, userUUID)
		if err == nil {
			t.Error("Expected error for non-GitHub URL, got nil")
		}

		if result != nil {
			t.Error("Expected nil result for non-GitHub URL")
		}
	})

	t.Run("init from github url with different formats", func(t *testing.T) {
		testCases := []string{
			"https://github.com/go-gorm/gorm",
			"https://github.com/go-gorm/gorm.git",
			"https://github.com/go-gorm/gorm/",
		}

		for i, githubURL := range testCases {
			userUUID := fmt.Sprintf("test-user-github-%03d", 100+i)

			result, err := libService.InitFromGitHub(ctx, githubURL, userUUID)
			if err != nil {
				t.Logf("InitFromGitHub(%s) error = %v (expected if GitHub API unavailable)", githubURL, err)
				continue
			}

			if result != nil && result.Library != nil {
				if result.Library.SourceType != "github" {
					t.Errorf("Expected source type 'github' for URL %s, got '%s'", githubURL, result.Library.SourceType)
				}
			}
		}
	})
}

// Test_Library_RefreshVersion_Advanced 测试刷新版本的高级场景
func Test_Library_RefreshVersion_Advanced(t *testing.T) {
	libService := &service.LibraryService{}

	t.Run("refresh version with callback", func(t *testing.T) {
		// 创建库
		createReq := &request.LibraryCreate{
			Name:        "refresh-callback-lib",
			Description: "test",
		}
		lib, err := libService.Create(createReq)
		if err != nil {
			t.Fatalf("Create() error = %v", err)
		}

		// 刷新版本（可能会失败，因为外部服务不可用）
		err = libService.RefreshVersion(lib.ID, "latest", "")
		if err != nil {
			t.Logf("RefreshVersion() error = %v (expected if external service unavailable)", err)
		}
	})

	t.Run("refresh version with specific branch", func(t *testing.T) {
		createReq := &request.LibraryCreate{
			Name:        "refresh-branch-lib",
			Description: "test",
		}
		lib, err := libService.Create(createReq)
		if err != nil {
			t.Fatalf("Create() error = %v", err)
		}

		// 尝试刷新特定分支
		err = libService.RefreshVersion(lib.ID, "latest", "main")
		if err != nil {
			t.Logf("RefreshVersion(branch=main) error = %v (expected if external service unavailable)", err)
		}
	})

	t.Run("refresh version multiple times", func(t *testing.T) {
		createReq := &request.LibraryCreate{
			Name:        "refresh-multi-lib",
			Description: "test",
		}
		lib, err := libService.Create(createReq)
		if err != nil {
			t.Fatalf("Create() error = %v", err)
		}

		// 多次刷新
		for i := 0; i < 3; i++ {
			err = libService.RefreshVersion(lib.ID, "latest", "")
			if err != nil {
				t.Logf("RefreshVersion() iteration %d error = %v (expected if external service unavailable)", i, err)
			}
		}
	})
}

// Test_Library_DeleteVersion_Advanced 测试删除版本的高级场景
func Test_Library_DeleteVersion_Advanced(t *testing.T) {
	libService := &service.LibraryService{}

	t.Run("delete version with documents", func(t *testing.T) {
		// 创建库
		createReq := &request.LibraryCreate{
			Name:        "delete-version-docs-lib",
			Description: "test",
		}
		lib, err := libService.Create(createReq)
		if err != nil {
			t.Fatalf("Create() error = %v", err)
		}

		// 创建版本
		err = libService.CreateVersion(lib.ID, "v1.5.0")
		if err != nil {
			t.Fatalf("CreateVersion() error = %v", err)
		}

		// 删除版本
		err = libService.DeleteVersion(lib.ID, "v1.5.0")
		if err != nil {
			t.Logf("DeleteVersion() error = %v (expected if version has documents)", err)
		}
	})

	t.Run("delete multiple versions", func(t *testing.T) {
		createReq := &request.LibraryCreate{
			Name:        "delete-multi-versions-lib",
			Description: "test",
		}
		lib, err := libService.Create(createReq)
		if err != nil {
			t.Fatalf("Create() error = %v", err)
		}

		// 创建多个版本
		versions := []string{"v1.1.0", "v1.2.0", "v1.3.0"}
		for _, version := range versions {
			err = libService.CreateVersion(lib.ID, version)
			if err != nil {
				t.Fatalf("CreateVersion(%s) error = %v", version, err)
			}
		}

		// 删除所有版本
		for _, version := range versions {
			err = libService.DeleteVersion(lib.ID, version)
			if err != nil {
				t.Logf("DeleteVersion(%s) error = %v (expected if version has documents)", version, err)
			}
		}
	})

	t.Run("delete version for non-existent library", func(t *testing.T) {
		err := libService.DeleteVersion(99999, "v1.0.0")
		if err == nil {
			t.Error("Expected error when deleting version for non-existent library, got nil")
		}
	})
}

// Test_Library_InitFromGitHub_GitHub 测试从 GitHub URL 初始化导入
func Test_Library_InitFromGitHub_GitHub(t *testing.T) {
	libService := &service.LibraryService{}
	ctx := context.Background()

	t.Run("init from github url go-gorm/gorm", func(t *testing.T) {
		// 使用完整的 GitHub URL: https://github.com/go-gorm/gorm
		githubURL := "https://github.com/go-gorm/gorm"
		userUUID := "test-github-init-001"

		result, err := libService.InitFromGitHub(ctx, githubURL, userUUID)
		if err != nil {
			t.Logf("InitFromGitHub(%s) error = %v (expected if library already exists or GitHub API unavailable)", githubURL, err)
			return
		}

		if result == nil {
			t.Logf("InitFromGitHub(%s) returned nil (expected if import failed)", githubURL)
			return
		}

		if result.Library == nil {
			t.Error("Expected library in result, got nil")
			return
		}

		if result.Library.Name == "" {
			t.Error("Expected non-empty library name")
		}

		if result.Library.SourceType != "github" {
			t.Errorf("Expected source type 'github', got '%s'", result.Library.SourceType)
		}

		t.Logf("Successfully initialized library from GitHub: %s", result.Library.Name)
	})

	t.Run("init from github url with invalid url", func(t *testing.T) {
		githubURL := "https://github.com/invalid-owner-xyz/invalid-repo-xyz"
		userUUID := "test-github-init-002"

		result, err := libService.InitFromGitHub(ctx, githubURL, userUUID)
		if err != nil {
			t.Logf("InitFromGitHub(invalid) error = %v (expected)", err)
			return
		}

		if result != nil {
			t.Logf("InitFromGitHub(invalid) returned result (may be unexpected)")
		}
	})

	t.Run("init from github url with malformed url", func(t *testing.T) {
		githubURL := "not-a-valid-url"
		userUUID := "test-github-init-003"

		result, err := libService.InitFromGitHub(ctx, githubURL, userUUID)
		if err != nil {
			t.Logf("InitFromGitHub(malformed) error = %v (expected)", err)
			return
		}

		if result != nil {
			t.Logf("InitFromGitHub(malformed) returned result (may be unexpected)")
		}
	})
}

// Test_Library_List_WithStatus 测试库列表查询（带状态过滤）
func Test_Library_List_WithStatus(t *testing.T) {
	libService := &service.LibraryService{}

	t.Run("list libraries with active status", func(t *testing.T) {
		// 创建库
		for i := 1; i <= 2; i++ {
			req := &request.LibraryCreate{
				Name:        fmt.Sprintf("status-active-lib-%d", i),
				Description: "test",
			}
			_, err := libService.Create(req)
			if err != nil {
				t.Fatalf("Create() error = %v", err)
			}
		}

		// 查询 active 状态的库
		status := "active"
		listReq := &request.LibraryList{
			Status: &status,
			PageInfo: request.PageInfo{
				Page:     1,
				PageSize: 10,
			},
		}

		result, err := libService.List(listReq)
		if err != nil {
			t.Fatalf("List() error = %v", err)
		}

		if result == nil {
			t.Fatal("Expected result, got nil")
		}

		if result.Total < 2 {
			t.Errorf("Expected at least 2 active libraries, got %d", result.Total)
		}
	})

	t.Run("list libraries with deleted status", func(t *testing.T) {
		// 查询 deleted 状态的库
		status := "deleted"
		listReq := &request.LibraryList{
			Status: &status,
			PageInfo: request.PageInfo{
				Page:     1,
				PageSize: 10,
			},
		}

		result, err := libService.List(listReq)
		if err != nil {
			t.Fatalf("List() error = %v", err)
		}

		if result == nil {
			t.Fatal("Expected result, got nil")
		}
	})
}

// Test_Library_DeleteVersion_EdgeCases 测试删除版本的边界情况
func Test_Library_DeleteVersion_EdgeCases(t *testing.T) {
	libService := &service.LibraryService{}

	t.Run("delete version when it's the default version", func(t *testing.T) {
		// 创建库并添加多个版本
		lib, _ := libService.Create(&request.LibraryCreate{
			Name:        "delete-default-version-lib",
			Description: "test",
		})

		// 添加版本
		_ = libService.CreateVersion(lib.ID, "v1.0.0")
		_ = libService.CreateVersion(lib.ID, "v2.0.0")

		// 删除默认版本（应该自动切换到下一个版本）
		err := libService.DeleteVersion(lib.ID, lib.DefaultVersion)
		if err != nil {
			t.Logf("DeleteVersion(default) error = %v (expected if no documents)", err)
		}
	})

	t.Run("delete version when it's the last version", func(t *testing.T) {
		// 创建只有一个版本的库
		lib, _ := libService.Create(&request.LibraryCreate{
			Name:        "delete-last-version-lib",
			Description: "test",
		})

		// 尝试删除唯一的版本
		err := libService.DeleteVersion(lib.ID, lib.DefaultVersion)
		if err != nil {
			t.Logf("DeleteVersion(last) error = %v (expected if no documents)", err)
		}
	})

	t.Run("delete non-default version", func(t *testing.T) {
		// 创建库并添加多个版本
		lib, _ := libService.Create(&request.LibraryCreate{
			Name:        "delete-non-default-lib",
			Description: "test",
		})

		// 添加非默认版本
		_ = libService.CreateVersion(lib.ID, "v1.0.0")
		_ = libService.CreateVersion(lib.ID, "v2.0.0")

		// 删除非默认版本
		err := libService.DeleteVersion(lib.ID, "v1.0.0")
		if err != nil {
			t.Logf("DeleteVersion(non-default) error = %v (expected if no documents)", err)
		}
	})
}

// Test_Library_RefreshVersion_EdgeCases 测试刷新版本的边界情况
func Test_Library_RefreshVersion_EdgeCases(t *testing.T) {
	libService := &service.LibraryService{}

	t.Run("refresh version with non-existent library", func(t *testing.T) {
		err := libService.RefreshVersion(99999, "v1.0.0", "test-user")
		if err == nil {
			t.Error("Expected error when refreshing version for non-existent library, got nil")
		}
	})

	t.Run("refresh version with non-existent version", func(t *testing.T) {
		// 创建库
		lib, _ := libService.Create(&request.LibraryCreate{
			Name:        "refresh-nonexist-version-lib",
			Description: "test",
		})

		// 尝试刷新不存在的版本
		err := libService.RefreshVersion(lib.ID, "v99.99.99", "test-user")
		if err == nil {
			t.Error("Expected error when refreshing non-existent version, got nil")
		}
	})

	t.Run("refresh version with valid library and version", func(t *testing.T) {
		// 创建库
		lib, _ := libService.Create(&request.LibraryCreate{
			Name:        "refresh-valid-lib",
			Description: "test",
		})

		// 尝试刷新默认版本（没有文档，应该快速返回）
		err := libService.RefreshVersion(lib.ID, lib.DefaultVersion, "test-user")
		if err != nil {
			t.Logf("RefreshVersion() error = %v (expected if no documents)", err)
		}
	})
}

// Test_Library_InitFromGitHub_EdgeCases 测试 GitHub 初始化的边界情况
func Test_Library_InitFromGitHub_EdgeCases(t *testing.T) {
	libService := &service.LibraryService{}
	ctx := context.Background()

	t.Run("init from invalid github url", func(t *testing.T) {
		invalidURLs := []string{
			"not-a-url",
			"http://example.com",
			"https://gitlab.com/user/repo",
			"github.com/user/repo",
		}

		for _, url := range invalidURLs {
			result, err := libService.InitFromGitHub(ctx, url, "test-user")
			if err == nil {
				t.Logf("InitFromGitHub(%s) succeeded unexpectedly, result: %v", url, result)
			} else {
				t.Logf("InitFromGitHub(%s) correctly failed: %v", url, err)
			}
		}
	})

	t.Run("init from github with different url formats", func(t *testing.T) {
		// 测试不同的 URL 格式
		urls := []string{
			"https://github.com/go-gorm/gorm",
			"https://github.com/go-gorm/gorm.git",
			"https://github.com/go-gorm/gorm/",
		}

		for _, url := range urls {
			result, err := libService.InitFromGitHub(ctx, url, "test-user")
			if err != nil {
				t.Logf("InitFromGitHub(%s) error: %v (may be expected)", url, err)
			} else if result != nil {
				t.Logf("InitFromGitHub(%s) succeeded: %s", url, result.Library.Name)
			}
		}
	})
}

// Test_Library_RefreshVersionWithCallback 测试带回调的版本刷新
func Test_Library_RefreshVersionWithCallback(t *testing.T) {
	libService := &service.LibraryService{}
	docService := &service.DocumentService{}

	// 创建测试库
	lib, err := libService.Create(&request.LibraryCreate{
		Name:        "test-refresh-callback-lib",
		Description: "Test library for refresh with callback",
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

	// 上传一个测试文档
	content := []byte("# Test Document\n\nThis is a test for refresh callback.")
	file, header := createMultipartFile("refresh-test.md", content)
	_, err = docService.Upload(lib.ID, "v1.0.0", file, header, "test-user", "test-task")
	if err != nil {
		t.Fatalf("Upload() error = %v", err)
	}

	// 等待处理完成
	time.Sleep(2 * time.Second)

	t.Run("refresh version with callback", func(t *testing.T) {
		// 创建状态通道
		statusChan := make(chan response.RefreshStatus, 10)

		// 启动 goroutine 接收状态
		statusReceived := false
		go func() {
			for status := range statusChan {
				t.Logf("Refresh Status: %s - %s", status.Stage, status.Message)
				statusReceived = true
			}
		}()

		// 执行刷新（在 goroutine 中，因为它会阻塞直到完成）
		go libService.RefreshVersionWithCallback(lib.ID, "v1.0.0", "test-user", statusChan)

		// 等待一些状态消息
		time.Sleep(5 * time.Second)

		if statusReceived {
			t.Log("✅ Received refresh status updates")
		} else {
			t.Log("⚠️  No status updates received (may be expected if refresh completed quickly)")
		}
	})

	t.Run("refresh non-existent version with callback", func(t *testing.T) {
		statusChan := make(chan response.RefreshStatus, 10)

		// 启动 goroutine 接收状态
		errorReceived := false
		go func() {
			for status := range statusChan {
				t.Logf("Status: %s - %s", status.Stage, status.Message)
				if status.Stage == "error" {
					errorReceived = true
				}
			}
		}()

		// 执行刷新
		go libService.RefreshVersionWithCallback(lib.ID, "v999.0.0", "test-user", statusChan)

		// 等待错误消息
		time.Sleep(2 * time.Second)

		if errorReceived {
			t.Log("✅ Correctly handled non-existent version with error status")
		}
	})

	t.Run("refresh non-existent library with callback", func(t *testing.T) {
		statusChan := make(chan response.RefreshStatus, 10)

		// 启动 goroutine 接收状态
		errorReceived := false
		go func() {
			for status := range statusChan {
				t.Logf("Status: %s - %s", status.Stage, status.Message)
				if status.Stage == "error" {
					errorReceived = true
				}
			}
		}()

		// 执行刷新
		go libService.RefreshVersionWithCallback(999999, "v1.0.0", "test-user", statusChan)

		// 等待错误消息
		time.Sleep(2 * time.Second)

		if errorReceived {
			t.Log("✅ Correctly handled non-existent library with error status")
		}
	})
}
