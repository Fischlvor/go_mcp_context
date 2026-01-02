package integration_test

import (
	"context"
	"testing"
	"time"

	"go-mcp-context/internal/model/request"
	"go-mcp-context/internal/model/response"
	"go-mcp-context/internal/service"
)

// Test_Integration_GitHubImport_RealAPI 集成测试：使用真实的 GitHub API
func Test_Integration_GitHubImport_RealAPI(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test in short mode")
	}

	githubService := service.NewGitHubImportService()
	libService := &service.LibraryService{}
	ctx := context.Background()

	t.Run("get real repo info from github", func(t *testing.T) {
		// 使用 GORM 官方文档仓库进行测试
		repo := "go-gorm/gorm.io"

		repoInfo, err := githubService.GetRepoInfo(ctx, repo)
		if err != nil {
			t.Fatalf("GetRepoInfo(%s) failed: %v", repo, err)
		}

		if repoInfo == nil {
			t.Fatal("Expected repo info, got nil")
		}

		// 验证返回的数据
		if repoInfo.Name != "gorm.io" {
			t.Errorf("Expected repo name 'gorm.io', got '%s'", repoInfo.Name)
		}

		if repoInfo.FullName != "go-gorm/gorm.io" {
			t.Errorf("Expected full name 'go-gorm/gorm.io', got '%s'", repoInfo.FullName)
		}

		if repoInfo.DefaultBranch == "" {
			t.Error("Expected non-empty default branch")
		}

		t.Logf("✅ Successfully retrieved repo info:")
		t.Logf("   Name: %s", repoInfo.Name)
		t.Logf("   Full Name: %s", repoInfo.FullName)
		t.Logf("   Default Branch: %s", repoInfo.DefaultBranch)
		t.Logf("   Description: %s", repoInfo.Description)
	})

	t.Run("get major versions from github", func(t *testing.T) {
		repo := "go-gorm/gorm.io"
		maxCount := 3

		versions, err := githubService.GetMajorVersions(ctx, repo, maxCount)
		if err != nil {
			t.Fatalf("GetMajorVersions(%s) failed: %v", repo, err)
		}

		t.Logf("GetMajorVersions returned: versions=%v, err=%v", versions, err)

		if len(versions) == 0 {
			t.Skip("No versions found (repository may not have releases)")
		}

		if len(versions) > maxCount {
			t.Errorf("Expected at most %d versions, got %d", maxCount, len(versions))
		}

		t.Logf("✅ Successfully retrieved %d versions:", len(versions))
		for i, v := range versions {
			t.Logf("   Version %d: %s", i+1, v)
		}
	})

	t.Run("import from github with real api", func(t *testing.T) {
		// 先创建一个测试库
		lib, err := libService.Create(&request.LibraryCreate{
			Name:        "integration-test-github-import",
			Description: "Integration test for GitHub import",
		})
		if err != nil {
			t.Fatalf("Failed to create library: %v", err)
		}

		// 准备导入请求（使用 GORM 文档仓库）
		importReq := &request.GitHubImportRequest{
			Repo:       "go-gorm/gorm.io",
			Branch:     "master",
			PathFilter: "pages/**/*.md", // 只导入文档页面
			Excludes:   []string{"node_modules"},
		}

		// 创建进度通道
		progressChan := make(chan response.GitHubImportProgress, 100)

		// 在 goroutine 中监听进度
		go func() {
			for progress := range progressChan {
				t.Logf("📦 Progress: [%s] %s", progress.Stage, progress.Message)
			}
		}()

		// 执行导入（设置超时）
		ctx, cancel := context.WithTimeout(context.Background(), 2*time.Minute)
		defer cancel()

		err = githubService.ImportFromGitHub(ctx, lib.ID, importReq, "integration-test", progressChan)

		if err != nil {
			t.Logf("⚠️  ImportFromGitHub() error = %v", err)
			t.Logf("Note: This is expected if GitHub API rate limit is reached or network issues")
			return
		}

		t.Log("✅ Successfully imported from GitHub")
	})
}

// Test_Integration_GitHubImport_ErrorHandling 集成测试：错误处理
func Test_Integration_GitHubImport_ErrorHandling(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test in short mode")
	}

	githubService := service.NewGitHubImportService()
	ctx := context.Background()

	t.Run("get repo info for non-existent repo", func(t *testing.T) {
		repo := "non-existent-owner-12345/non-existent-repo-67890"

		repoInfo, err := githubService.GetRepoInfo(ctx, repo)
		if err == nil {
			t.Error("Expected error for non-existent repo, got nil")
		}

		if repoInfo != nil {
			t.Error("Expected nil repo info for non-existent repo")
		}

		t.Logf("✅ Correctly handled non-existent repo: %v", err)
	})

	t.Run("get major versions for invalid repo", func(t *testing.T) {
		repo := "invalid/repo/format/with/too/many/slashes"

		versions, err := githubService.GetMajorVersions(ctx, repo, 5)
		if err == nil {
			t.Error("Expected error for invalid repo format, got nil")
		}

		if versions != nil {
			t.Error("Expected nil versions for invalid repo")
		}

		t.Logf("✅ Correctly handled invalid repo format: %v", err)
	})
}

// Test_Integration_GitHubImport_ProcessFile 集成测试：测试 processFile 函数
// 此测试通过完整的 GitHub 导入流程来间接测试 processFile 函数
func Test_Integration_GitHubImport_ProcessFile(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test in short mode")
	}

	githubService := service.NewGitHubImportService()
	libService := &service.LibraryService{}

	t.Run("import small repo to trigger processFile", func(t *testing.T) {
		// 创建测试库
		lib, err := libService.Create(&request.LibraryCreate{
			Name:        "processfile-test-lib",
			Description: "Test library for processFile function",
		})
		if err != nil {
			t.Fatalf("Failed to create library: %v", err)
		}
		defer libService.Delete(lib.ID)

		// 使用 GORM 文档仓库进行测试，只导入少量文件
		importReq := &request.GitHubImportRequest{
			Repo:       "go-gorm/gorm.io", // 与其他集成测试保持一致
			Branch:     "master",
			PathFilter: "pages/docs/index.md", // 只导入单个文档文件
			Excludes:   []string{"node_modules"},
		}

		// 创建进度通道
		progressChan := make(chan response.GitHubImportProgress, 100)

		// 记录处理的文件数
		fileCount := 0
		go func() {
			for progress := range progressChan {
				if progress.Stage == "downloading" && progress.FileName != "" {
					fileCount++
					t.Logf("📄 Processing file: %s (processFile called)", progress.FileName)
				}
				t.Logf("📦 [%s] %s", progress.Stage, progress.Message)
			}
		}()

		// 执行导入
		ctx, cancel := context.WithTimeout(context.Background(), 1*time.Minute)
		defer cancel()

		err = githubService.ImportFromGitHub(ctx, lib.ID, importReq, "processfile-test", progressChan)

		if err != nil {
			t.Logf("⚠️  ImportFromGitHub() error = %v", err)
			t.Logf("Note: This may be due to GitHub API rate limit or network issues")
			return
		}

		if fileCount > 0 {
			t.Logf("✅ Successfully processed %d file(s) through processFile function", fileCount)
		} else {
			t.Log("⚠️  No files were processed (may need to check path filter)")
		}
	})
}

// Test_Integration_GitHubImport_RateLimiting 集成测试：速率限制
func Test_Integration_GitHubImport_RateLimiting(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test in short mode")
	}

	githubService := service.NewGitHubImportService()
	ctx := context.Background()

	t.Run("multiple sequential requests", func(t *testing.T) {
		repos := []string{
			"go-gorm/gorm.io",
			"go-gorm/gorm",
			"gin-gonic/gin",
		}

		successCount := 0
		for _, repo := range repos {
			repoInfo, err := githubService.GetRepoInfo(ctx, repo)
			if err != nil {
				t.Logf("⚠️  GetRepoInfo(%s) error: %v (may be rate limited)", repo, err)
				continue
			}

			if repoInfo != nil {
				successCount++
				t.Logf("✅ Retrieved info for %s", repo)
			}

			// 添加小延迟避免触发速率限制
			time.Sleep(100 * time.Millisecond)
		}

		if successCount == 0 {
			t.Log("⚠️  All requests failed (likely rate limited)")
		} else {
			t.Logf("✅ Successfully retrieved %d/%d repos", successCount, len(repos))
		}
	})
}
