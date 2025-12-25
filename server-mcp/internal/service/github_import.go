package service

import (
	"context"
	"fmt"
	"path/filepath"
	"strconv"
	"strings"
	"sync"

	dbmodel "go-mcp-context/internal/model/database"
	"go-mcp-context/internal/model/request"
	"go-mcp-context/internal/model/response"
	"go-mcp-context/pkg/bufferedwriter/actlog"
	"go-mcp-context/pkg/github"
	"go-mcp-context/pkg/global"
	"go-mcp-context/pkg/utils"
)

// GitHubImportService GitHub 导入服务
type GitHubImportService struct {
	client *github.Client
}

// NewGitHubImportService 创建 GitHub 导入服务
func NewGitHubImportService() *GitHubImportService {
	return &GitHubImportService{
		client: github.NewClient(),
	}
}

// GetRepoInfo 获取仓库信息（代理到 client）
func (s *GitHubImportService) GetRepoInfo(ctx context.Context, repo string) (*github.Repo, error) {
	return s.client.GetRepoInfo(ctx, repo)
}

// GetMajorVersions 获取每个大版本的最新版本（代理到 client）
func (s *GitHubImportService) GetMajorVersions(ctx context.Context, repo string, maxCount int) ([]string, error) {
	return s.client.GetMajorVersions(ctx, repo, maxCount)
}

// sendProgress 安全发送进度（channel 可为 nil）
func sendProgress(ch chan response.GitHubImportProgress, progress response.GitHubImportProgress) {
	if ch != nil {
		ch <- progress
	}
}

// ImportFromGitHub 从 GitHub 导入文档（progressChan 可为 nil）
func (s *GitHubImportService) ImportFromGitHub(ctx context.Context, libraryID uint, req *request.GitHubImportRequest, actorID string, progressChan chan response.GitHubImportProgress) error {
	if progressChan != nil {
		defer close(progressChan)
	}

	// 使用传入的任务 ID，或生成新的
	taskID := req.TaskID
	if taskID == "" {
		taskID = utils.GenerateTaskID()
	}

	// 1. 检查库是否存在
	var library dbmodel.Library
	if err := global.DB.First(&library, libraryID).Error; err != nil {
		sendProgress(progressChan, response.GitHubImportProgress{Stage: "failed", Message: "库不存在"})
		return ErrNotFound
	}

	// 2. 确定 ref（branch 或 tag）
	ref := req.Branch
	if req.Tag != "" {
		ref = req.Tag
	}
	if ref == "" {
		// 获取默认分支
		repoInfo, err := s.client.GetRepoInfo(ctx, req.Repo)
		if err != nil {
			sendProgress(progressChan, response.GitHubImportProgress{Stage: "failed", Message: "获取仓库信息失败: " + err.Error()})
			return err
		}
		ref = repoInfo.DefaultBranch
	}

	sendProgress(progressChan, response.GitHubImportProgress{Stage: "fetching_tree", Message: fmt.Sprintf("获取目录树: %s@%s", req.Repo, ref)})

	// 3. 获取目录树
	tree, err := s.client.GetTree(ctx, req.Repo, ref)
	if err != nil {
		sendProgress(progressChan, response.GitHubImportProgress{Stage: "failed", Message: "获取目录树失败: " + err.Error()})
		return err
	}

	// 4. 过滤文件
	files := s.client.FilterTree(tree, req.PathFilter, req.Excludes)
	if len(files) == 0 {
		sendProgress(progressChan, response.GitHubImportProgress{Stage: "failed", Message: "没有找到文档文件"})
		return fmt.Errorf("no document files found")
	}

	sendProgress(progressChan, response.GitHubImportProgress{
		Stage:   "downloading",
		Total:   len(files),
		Current: 0,
		Message: fmt.Sprintf("找到 %d 个文档文件", len(files)),
	})

	// 5. 确定版本名
	version := req.Version
	if version == "" {
		if req.Tag != "" {
			version = req.Tag
		} else {
			version = "latest"
		}
	}

	// 创建任务日志器（开始日志已在 API 层同步写入）
	actLogger := actlog.NewTaskLogger(libraryID, taskID, version).
		WithTarget("version", version).
		WithActor(actorID)

	actLogger.Info(actlog.EventGHImportDownload, fmt.Sprintf("找到 %d 个文档文件", len(files)))

	// 5.1 检查版本是否已存在
	versionExists := false
	for _, v := range library.Versions {
		if v == version {
			versionExists = true
			break
		}
	}
	// 如果版本已存在且不是默认版本，拒绝重复导入
	if versionExists {
		sendProgress(progressChan, response.GitHubImportProgress{Stage: "failed", Message: fmt.Sprintf("版本 %s 已存在", version)})
		return ErrVersionExists
	}

	// 6. 获取仓库大小，选择下载方式
	repoInfo, _ := s.client.GetRepoInfo(ctx, req.Repo)
	repoSizeKB := 0
	if repoInfo != nil {
		repoSizeKB = repoInfo.Size
	}

	// 阈值：100MB = 100 * 1024 KB
	const tarballThresholdKB = 100 * 1024
	useTarball := repoSizeKB >= tarballThresholdKB

	processor := &DocumentProcessor{}
	successCount := 0
	failCount := 0

	if useTarball {
		// === 大仓库：使用 tarball 流式下载 ===
		actLogger.Info(actlog.EventGHImportDownload, fmt.Sprintf("大仓库（%dMB），使用 tarball 流式下载", repoSizeKB/1024))
		sendProgress(progressChan, response.GitHubImportProgress{
			Stage:   "downloading",
			Message: fmt.Sprintf("大仓库（%dMB），使用 tarball 流式下载", repoSizeKB/1024),
		})

		// 构建文件过滤器（基于已过滤的文件列表）
		allowedPaths := make(map[string]bool)
		for _, f := range files {
			allowedPaths[f.Path] = true
		}
		filter := func(path string) bool {
			return allowedPaths[path]
		}

		fileChan, errChan := s.client.DownloadTarballFiles(ctx, req.Repo, ref, filter)

		// 处理文件
		for file := range fileChan {
			successCount, failCount = s.processFile(ctx, file.Path, file.Content, library, version, taskID, actLogger, req, processor, progressChan, successCount, failCount, len(files))
		}

		// 检查错误
		if err := <-errChan; err != nil {
			sendProgress(progressChan, response.GitHubImportProgress{
				Stage:   "warning",
				Message: fmt.Sprintf("tarball 下载出错: %s", err.Error()),
			})
		}
	} else {
		// === 小仓库：使用多 API 并行下载 ===
		actLogger.Info(actlog.EventGHImportDownload, fmt.Sprintf("开始下载: %d 个文件", len(files)))
		sendProgress(progressChan, response.GitHubImportProgress{
			Stage:   "downloading",
			Message: fmt.Sprintf("开始下载 %s@%s（%d 个文件）", req.Repo, version, len(files)),
		})

		type downloadResult struct {
			path    string
			content []byte
			err     error
		}

		results := make(chan downloadResult, len(files))
		var wg sync.WaitGroup
		semaphore := make(chan struct{}, 10)

		for _, item := range files {
			wg.Add(1)
			go func(item github.TreeItem) {
				defer wg.Done()
				semaphore <- struct{}{}
				defer func() { <-semaphore }()

				content, err := s.client.DownloadFile(ctx, req.Repo, ref, item.Path)
				results <- downloadResult{path: item.Path, content: content, err: err}
			}(item)
		}

		go func() {
			wg.Wait()
			close(results)
		}()

		for result := range results {
			if result.err != nil {
				failCount++
				actLogger.Warning(actlog.EventGHImportDownload, fmt.Sprintf("下载失败: %s", result.path))
				sendProgress(progressChan, response.GitHubImportProgress{
					Stage:    "downloading",
					Current:  successCount + failCount,
					Total:    len(files),
					FileName: result.path,
					Message:  fmt.Sprintf("下载失败: %s - %s", result.path, result.err.Error()),
				})
				continue
			}
			successCount, failCount = s.processFile(ctx, result.path, result.content, library, version, taskID, actLogger, req, processor, progressChan, successCount, failCount, len(files))
		}
	}

	// 8. 有成功文件时，创建版本并更新 source 信息
	if successCount > 0 {
		// 创建版本（如果不存在）
		if !versionExists && version != library.DefaultVersion {
			libService := &LibraryService{}
			if err := libService.CreateVersion(libraryID, version); err != nil {
				if err == ErrVersionExists {
					sendProgress(progressChan, response.GitHubImportProgress{
						Stage:   "info",
						Message: fmt.Sprintf("版本 %s 已存在", version),
					})
				} else {
					// 版本创建失败，但文件已上传，记录警告
					sendProgress(progressChan, response.GitHubImportProgress{
						Stage:   "warning",
						Message: fmt.Sprintf("版本创建失败: %s，但文件已上传", err.Error()),
					})
				}
			} else {
				sendProgress(progressChan, response.GitHubImportProgress{
					Stage:   "info",
					Message: fmt.Sprintf("版本 %s 创建成功", version),
				})
			}
		}

		// 更新 source 信息
		global.DB.Model(&library).Updates(map[string]interface{}{
			"source_type": "github",
			"source_url":  req.Repo,
		})
	}

	// 记录导入完成
	if failCount == 0 {
		actLogger.Success(actlog.EventGHImportComplete, fmt.Sprintf("导入完成: %s@%s (成功 %d)", req.Repo, version, successCount))
	} else {
		actLogger.Warning(actlog.EventGHImportComplete, fmt.Sprintf("导入完成: %s@%s (成功 %d, 失败 %d)", req.Repo, version, successCount, failCount))
	}

	sendProgress(progressChan, response.GitHubImportProgress{
		Stage:   "completed",
		Current: successCount + failCount,
		Total:   len(files),
		Message: fmt.Sprintf("导入完成：成功 %d，失败 %d", successCount, failCount),
	})

	return nil
}

// processFile 处理单个文件（上传存储 + 创建文档记录 + 异步处理）
func (s *GitHubImportService) processFile(
	ctx context.Context,
	filePath string,
	content []byte,
	library dbmodel.Library,
	version string,
	taskID string,
	actLogger *actlog.TaskLogger,
	req *request.GitHubImportRequest,
	processor *DocumentProcessor,
	progressChan chan response.GitHubImportProgress,
	successCount, failCount, total int,
) (int, int) {
	// 计算存储路径
	storagePath := filePath
	if req.PathFilter != "" {
		storagePath = strings.TrimPrefix(storagePath, req.PathFilter)
	}
	storagePath = strings.TrimPrefix(storagePath, "/")

	// 生成存储 Key
	libDir := sanitizeFileName(library.Name)
	versionDir := sanitizeFileName(version)
	key := filepath.Join(global.Config.Qiniu.PathPrefix, libDir, versionDir, storagePath)

	// 根据扩展名设置 MIME 类型
	mimeType := "text/markdown"
	ext := strings.ToLower(filepath.Ext(filePath))
	if ext == ".mdx" {
		mimeType = "text/mdx"
	}

	// 上传到存储
	uploadResult, err := global.Storage.Upload(ctx, key, strings.NewReader(string(content)), int64(len(content)), mimeType)
	if err != nil {
		failCount++
		actLogger.Warning(actlog.EventGHImportDownload, fmt.Sprintf("上传失败: %s", filePath))
		sendProgress(progressChan, response.GitHubImportProgress{
			Stage:    "downloading",
			Current:  successCount + failCount,
			Total:    total,
			FileName: filePath,
			Message:  fmt.Sprintf("上传失败: %s - %s", filePath, err.Error()),
		})
		return successCount, failCount
	}

	// 创建文档记录
	doc := &dbmodel.DocumentUpload{
		LibraryID:   library.ID,
		Version:     version,
		Title:       filepath.Base(filePath),
		FilePath:    uploadResult.Key,
		FileType:    getFileType(filepath.Ext(filePath)),
		FileSize:    int64(len(content)),
		ContentHash: uploadResult.ETag,
		Status:      "processing",
	}

	if err := global.DB.Create(doc).Error; err != nil {
		failCount++
		return successCount, failCount
	}

	// 同步处理文档，并推送处理状态
	statusChan := make(chan response.ProcessStatus, 10)
	docLogger := actLogger.WithTarget("document", strconv.FormatUint(uint64(doc.ID), 10))
	go processor.ProcessDocumentWithCallback(doc, content, statusChan, docLogger, false) // GitHub 导入是中间步骤

	// 转发处理状态到 progressChan
	processingFailed := false
	for status := range statusChan {
		sendProgress(progressChan, response.GitHubImportProgress{
			Stage:    status.Stage,
			Current:  successCount + failCount,
			Total:    total,
			FileName: filePath,
			Message:  fmt.Sprintf("[%s] %s", filePath, status.Message),
		})
		if status.Stage == "completed" {
			break
		}
		if status.Stage == "failed" {
			processingFailed = true
			break
		}
	}

	if processingFailed {
		failCount++
	} else {
		successCount++
	}
	return successCount, failCount
}
