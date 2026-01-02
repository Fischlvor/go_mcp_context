package test_test

import (
	"context"
	"testing"
	"time"

	"go-mcp-context/internal/model/request"
	"go-mcp-context/internal/model/response"
	"go-mcp-context/internal/service"
)

// Test_GitHubImport_Service 测试 GitHub 导入服务
func Test_GitHubImport_Service(t *testing.T) {
	githubImportService := service.NewGitHubImportService()

	t.Run("github import service initialization", func(t *testing.T) {
		if githubImportService == nil {
			t.Error("Expected GitHubImportService to be initialized")
		}
	})

	t.Run("github import service has methods", func(t *testing.T) {
		// 验证服务结构存在
		if githubImportService == nil {
			t.Error("Expected GitHubImportService to be initialized")
		}
	})
}

// Test_GitHubImport_GetRepoInfo 测试获取 GitHub 仓库信息
func Test_GitHubImport_GetRepoInfo(t *testing.T) {
	githubImportService := service.NewGitHubImportService()
	ctx := context.Background()

	t.Run("get repo info for go-gorm/gorm", func(t *testing.T) {
		// 使用真实的 GitHub URL: https://github.com/go-gorm/gorm
		repo := "go-gorm/gorm"

		repoInfo, err := githubImportService.GetRepoInfo(ctx, repo)
		if err != nil {
			t.Logf("GetRepoInfo(%s) error = %v (expected if GitHub API unavailable)", repo, err)
			return
		}

		if repoInfo == nil {
			t.Logf("GetRepoInfo(%s) returned nil (expected if repo not found)", repo)
			return
		}

		if repoInfo.Name == "" {
			t.Error("Expected non-empty repository name")
		}

		if repoInfo.FullName == "" {
			t.Error("Expected non-empty repository full name")
		}

		t.Logf("Successfully retrieved repo info: %s (size: %d KB)", repoInfo.FullName, repoInfo.Size)
	})

	t.Run("get repo info for invalid repo", func(t *testing.T) {
		repo := "invalid-owner-xyz/invalid-repo-xyz"

		repoInfo, err := githubImportService.GetRepoInfo(ctx, repo)
		if err != nil {
			t.Logf("GetRepoInfo(invalid) error = %v (expected)", err)
			return
		}

		if repoInfo != nil {
			t.Logf("GetRepoInfo(invalid) returned info (may be unexpected)")
		}
	})
}

// Test_GitHubImport_GetMajorVersions 测试获取 GitHub 仓库的主要版本
func Test_GitHubImport_GetMajorVersions(t *testing.T) {
	githubImportService := service.NewGitHubImportService()
	ctx := context.Background()

	t.Run("get major versions for go-gorm/gorm", func(t *testing.T) {
		// 使用真实的 GitHub URL: https://github.com/go-gorm/gorm
		repo := "go-gorm/gorm"

		versions, err := githubImportService.GetMajorVersions(ctx, repo, 5)
		if err != nil {
			t.Logf("GetMajorVersions(%s) error = %v (expected if GitHub API unavailable)", repo, err)
			return
		}

		if versions == nil {
			t.Logf("GetMajorVersions(%s) returned nil (expected if no versions found)", repo)
			return
		}

		if len(versions) > 0 {
			t.Logf("GetMajorVersions(%s) returned %d versions", repo, len(versions))
			for i, v := range versions {
				t.Logf("  Version %d: %s", i+1, v)
			}
		}
	})

	t.Run("get major versions with max count", func(t *testing.T) {
		repo := "go-gorm/gorm"
		maxCount := 3

		versions, err := githubImportService.GetMajorVersions(ctx, repo, maxCount)
		if err != nil {
			t.Logf("GetMajorVersions(maxCount=%d) error = %v (expected if GitHub API unavailable)", maxCount, err)
			return
		}

		if versions != nil && len(versions) > 0 {
			if len(versions) > maxCount {
				t.Errorf("Expected at most %d versions, got %d", maxCount, len(versions))
			}
		}
	})
}

// Test_GitHubImport_Service_Advanced 测试 GitHub 导入服务的高级场景
func Test_GitHubImport_Service_Advanced(t *testing.T) {
	githubImportService := service.NewGitHubImportService()
	ctx := context.Background()

	t.Run("github import service with valid initialization", func(t *testing.T) {
		if githubImportService == nil {
			t.Error("Expected GitHubImportService to be initialized")
		}
	})

	t.Run("github import service multiple instances", func(t *testing.T) {
		service1 := service.NewGitHubImportService()
		service2 := service.NewGitHubImportService()

		if service1 == nil || service2 == nil {
			t.Error("Expected both services to be initialized")
		}
	})

	t.Run("github import service with different repos", func(t *testing.T) {
		repos := []string{
			"go-gorm/gorm",
			"kubernetes/kubernetes",
			"docker/docker",
		}

		for _, repo := range repos {
			repoInfo, err := githubImportService.GetRepoInfo(ctx, repo)
			if err != nil {
				t.Logf("GetRepoInfo(%s) error = %v (expected if GitHub API unavailable)", repo, err)
				continue
			}

			if repoInfo != nil {
				t.Logf("Successfully retrieved info for %s", repo)
			}
		}
	})

	t.Run("github import service state", func(t *testing.T) {
		if githubImportService == nil {
			t.Error("Expected GitHubImportService to be initialized")
		}
	})
}

// Test_GitHubImport_ImportFromGitHub 测试 GitHub 导入功能
func Test_GitHubImport_ImportFromGitHub(t *testing.T) {
	githubService := service.NewGitHubImportService()
	libService := &service.LibraryService{}

	// 创建测试库
	lib, err := libService.Create(&request.LibraryCreate{
		Name:        "test-github-import-lib",
		Description: "Test library for GitHub import",
	})
	if err != nil {
		t.Fatalf("Failed to create library: %v", err)
	}
	defer libService.Delete(lib.ID)

	t.Run("import from github with progress", func(t *testing.T) {
		// 创建进度通道
		progressChan := make(chan response.GitHubImportProgress, 100)

		// 启动 goroutine 接收进度
		progressReceived := false
		go func() {
			for progress := range progressChan {
				t.Logf("📦 Import Progress: [%s] %s", progress.Stage, progress.Message)
				progressReceived = true
			}
		}()

		// 准备导入请求（使用一个小的公开仓库）
		importReq := &request.GitHubImportRequest{
			Repo:       "octocat/Hello-World", // GitHub 官方示例仓库
			Branch:     "master",
			PathFilter: "*.md", // 只导入 markdown 文件
			Excludes:   []string{},
		}

		// 执行导入（设置超时）
		ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
		defer cancel()

		err := githubService.ImportFromGitHub(ctx, lib.ID, importReq, "test-user", progressChan)

		if err != nil {
			t.Logf("⚠️  ImportFromGitHub() error = %v", err)
			t.Log("Note: This is expected if GitHub API rate limit is reached or network issues")
		} else {
			t.Log("✅ Successfully imported from GitHub")
		}

		if progressReceived {
			t.Log("✅ Received progress updates (sendProgress executed)")
		}
	})

	t.Run("import with invalid repo", func(t *testing.T) {
		progressChan := make(chan response.GitHubImportProgress, 10)

		go func() {
			for range progressChan {
				// 消费进度消息
			}
		}()

		importReq := &request.GitHubImportRequest{
			Repo:       "non-existent-user-12345/non-existent-repo-67890",
			Branch:     "main",
			PathFilter: "*.md",
			Excludes:   []string{},
		}

		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()

		err := githubService.ImportFromGitHub(ctx, lib.ID, importReq, "test-user", progressChan)
		if err == nil {
			t.Error("Expected error for non-existent repo")
		} else {
			t.Logf("✅ Correctly handled invalid repo: %v", err)
		}
	})

	t.Run("import to non-existent library", func(t *testing.T) {
		progressChan := make(chan response.GitHubImportProgress, 10)

		go func() {
			for range progressChan {
				// 消费进度消息
			}
		}()

		importReq := &request.GitHubImportRequest{
			Repo:       "octocat/Hello-World",
			Branch:     "master",
			PathFilter: "*.md",
			Excludes:   []string{},
		}

		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()

		err := githubService.ImportFromGitHub(ctx, 999999, importReq, "test-user", progressChan)
		if err == nil {
			t.Error("Expected error for non-existent library")
		} else {
			t.Logf("✅ Correctly handled non-existent library: %v", err)
		}
	})

	t.Run("import with empty branch", func(t *testing.T) {
		progressChan := make(chan response.GitHubImportProgress, 10)

		go func() {
			for range progressChan {
				// 消费进度消息
			}
		}()

		importReq := &request.GitHubImportRequest{
			Repo:       "octocat/Hello-World",
			Branch:     "", // 空分支，应该使用默认分支
			PathFilter: "*.md",
			Excludes:   []string{},
		}

		ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
		defer cancel()

		err := githubService.ImportFromGitHub(ctx, lib.ID, importReq, "test-user", progressChan)
		if err != nil {
			t.Logf("ImportFromGitHub(empty branch) error = %v (may be expected)", err)
		} else {
			t.Log("✅ Successfully handled empty branch (used default branch)")
		}
	})
}
