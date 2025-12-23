package service

import (
	"context"
	"fmt"
	"io"
	"log"
	"regexp"
	"strconv"
	"sync"
	"time"

	dbmodel "go-mcp-context/internal/model/database"
	"go-mcp-context/internal/model/request"
	"go-mcp-context/internal/model/response"
	"go-mcp-context/pkg/bufferedwriter/actlog"
	"go-mcp-context/pkg/global"
	"go-mcp-context/pkg/utils"

	"github.com/lib/pq"
)

type LibraryService struct{}

// ValidateVersion 验证版本格式（Semantic Versioning）
// 支持的格式：
// - v1.0.0, v1.2.3, v2.0.0
// - v1.0.0-alpha, v1.0.0-beta, v1.0.0-rc.1
// - 1.0.0（不带 v 前缀）
func (s *LibraryService) ValidateVersion(version string) error {
	if version == "" {
		return ErrInvalidParams
	}

	if len(version) > 50 {
		return ErrInvalidParams
	}

	// Semantic Versioning 正则表达式
	// 支持: v1.0.0, 1.0.0, v1.0.0-alpha, v1.0.0-beta.1, v1.0.0-rc.1 等
	pattern := `^v?(\d+)\.(\d+)\.(\d+)(-[a-zA-Z0-9]+(\.[a-zA-Z0-9]+)*)?(\+[a-zA-Z0-9]+(\.[a-zA-Z0-9]+)*)?$`
	matched, err := regexp.MatchString(pattern, version)
	if err != nil || !matched {
		return ErrInvalidParams
	}

	return nil
}

// List 获取库列表
func (s *LibraryService) List(req *request.LibraryList) (*response.PageResult, error) {
	var libraries []dbmodel.Library
	var total int64

	db := global.DB.Model(&dbmodel.Library{})

	// 条件过滤
	if req.Name != nil && *req.Name != "" {
		db = db.Where("name LIKE ?", "%"+*req.Name+"%")
	}
	if req.Status != nil && *req.Status != "" {
		db = db.Where("status = ?", *req.Status)
	} else {
		db = db.Where("status = ?", "active")
	}

	// 计算总数
	if err := db.Count(&total).Error; err != nil {
		return nil, err
	}

	// 分页
	page := req.Page
	pageSize := req.PageSize
	if page <= 0 {
		page = 1
	}
	if pageSize <= 0 {
		pageSize = 10
	}

	offset := (page - 1) * pageSize
	if err := db.Offset(offset).Limit(pageSize).Find(&libraries).Error; err != nil {
		return nil, err
	}

	return &response.PageResult{
		List:     libraries,
		Total:    total,
		Page:     page,
		PageSize: pageSize,
	}, nil
}

// Create 创建库
func (s *LibraryService) Create(req *request.LibraryCreate) (*dbmodel.Library, error) {
	// 默认 source_type 为 local
	sourceType := req.SourceType
	if sourceType == "" {
		sourceType = "local"
	}

	// 默认版本
	defaultVersion := req.DefaultVersion
	if defaultVersion == "" {
		defaultVersion = "default"
	}

	library := &dbmodel.Library{
		Name:           req.Name,
		Description:    req.Description,
		SourceType:     sourceType,
		SourceURL:      req.SourceURL,
		Status:         "active",
		DefaultVersion: defaultVersion,
		Versions:       []string{}, // versions 只存正常版本，不包含 default
		CreatedBy:      req.CreatedBy,
	}

	if err := global.DB.Create(library).Error; err != nil {
		return nil, err
	}

	return library, nil
}

// GetByID 根据 ID 获取库
func (s *LibraryService) GetByID(id uint) (*dbmodel.Library, error) {
	var library dbmodel.Library
	if err := global.DB.First(&library, id).Error; err != nil {
		return nil, err
	}
	return &library, nil
}

// getBySourceURL 根据 SourceURL 获取库（内部使用）
func (s *LibraryService) getBySourceURL(sourceURL string) (*dbmodel.Library, error) {
	var library dbmodel.Library
	if err := global.DB.Where("source_url = ?", sourceURL).First(&library).Error; err != nil {
		return nil, err
	}
	return &library, nil
}

// InitFromGitHub 从 GitHub URL 初始化创建库
// 流程：解析 URL -> 验证连通性 -> 检查重复 -> 创建库
// 返回：库信息、仓库默认分支、错误
func (s *LibraryService) InitFromGitHub(ctx context.Context, githubURL string, createdBy string) (*dbmodel.Library, string, error) {
	// 1. 解析 GitHub URL
	repo, err := utils.ParseGitHubURL(githubURL)
	if err != nil {
		return nil, "", fmt.Errorf("无效的 GitHub URL: %w", err)
	}

	// 2. 验证仓库连通性
	githubService := NewGitHubImportService()
	repoInfo, err := githubService.GetRepoInfo(ctx, repo)
	if err != nil {
		return nil, "", fmt.Errorf("无法访问仓库: %w", err)
	}

	// 3. 检查是否已存在
	existingLib, _ := s.getBySourceURL(repo)
	if existingLib != nil {
		return nil, "", fmt.Errorf("该库已存在: %s (ID: %d)", existingLib.Name, existingLib.ID)
	}

	// 4. 创建库（默认版本为 latest）
	repoName := utils.ExtractRepoName(repo)
	library, err := s.Create(&request.LibraryCreate{
		Name:           repoName,
		Description:    repoInfo.Description,
		SourceType:     "github",
		SourceURL:      repo,
		DefaultVersion: "latest",
		CreatedBy:      createdBy,
	})
	if err != nil {
		return nil, "", fmt.Errorf("创建库失败: %w", err)
	}

	return library, repoInfo.DefaultBranch, nil
}

// Update 更新库
func (s *LibraryService) Update(id uint, req *request.LibraryCreate) (*dbmodel.Library, error) {
	var library dbmodel.Library
	if err := global.DB.First(&library, id).Error; err != nil {
		return nil, err
	}

	library.Name = req.Name
	library.Description = req.Description
	library.SourceType = req.SourceType
	library.SourceURL = req.SourceURL

	if err := global.DB.Save(&library).Error; err != nil {
		return nil, err
	}

	return &library, nil
}

// Delete 删除库（软删除）
func (s *LibraryService) Delete(id uint) error {
	now := time.Now()
	result := global.DB.Model(&dbmodel.Library{}).
		Where("id = ? AND deleted_at IS NULL", id).
		Updates(map[string]interface{}{"status": "deleted", "deleted_at": now})

	if result.Error != nil {
		return result.Error
	}

	if result.RowsAffected == 0 {
		return ErrNotFound
	}

	// 同时删除关联的文档上传记录和 chunks
	global.DB.Model(&dbmodel.DocumentUpload{}).
		Where("library_id = ? AND deleted_at IS NULL", id).
		Updates(map[string]interface{}{"status": "deleted", "deleted_at": now})
	global.DB.Model(&dbmodel.DocumentChunk{}).
		Where("library_id = ? AND deleted_at IS NULL", id).
		Updates(map[string]interface{}{"status": "deleted", "deleted_at": now})

	return nil
}

// SearchByName 根据名称模糊搜索库
func (s *LibraryService) SearchByName(name string) ([]dbmodel.Library, error) {
	var libraries []dbmodel.Library

	// 前缀匹配优先
	err := global.DB.Where("status = ? AND name LIKE ?", "active", name+"%").
		Order("name ASC").
		Limit(10).
		Find(&libraries).Error

	if err != nil {
		return nil, err
	}

	return libraries, nil
}

// GetByName 根据名称获取库
func (s *LibraryService) GetByName(name string) (*dbmodel.Library, error) {
	var library dbmodel.Library
	if err := global.DB.Where("name = ? AND status = ?", name, "active").
		First(&library).Error; err != nil {
		return nil, ErrNotFound
	}
	return &library, nil
}

// GetLibraryInfo 获取库详情（带统计信息）
func (s *LibraryService) GetLibraryInfo(id uint) (*response.LibraryInfo, error) {
	var library dbmodel.Library
	if err := global.DB.First(&library, id).Error; err != nil {
		return nil, err
	}

	// 统计文档上传数
	var docCount int64
	global.DB.Model(&dbmodel.DocumentUpload{}).
		Where("library_id = ? AND status = ?", id, "completed").
		Count(&docCount)

	// 统计 chunk 数和 token 数
	var stats struct {
		ChunkCount int64 `gorm:"column:chunk_count"`
		TokenCount int64 `gorm:"column:token_count"`
	}
	global.DB.Model(&dbmodel.DocumentUpload{}).
		Select("COALESCE(SUM(chunk_count), 0) as chunk_count, COALESCE(SUM(token_count), 0) as token_count").
		Where("library_id = ? AND status = ?", id, "completed").
		Scan(&stats)

	return &response.LibraryInfo{
		ID:             library.ID,
		Name:           library.Name,
		DefaultVersion: library.DefaultVersion,
		Versions:       library.Versions,
		SourceType:     library.SourceType,
		SourceURL:      library.SourceURL,
		Description:    library.Description,
		DocumentCount:  int(docCount),
		ChunkCount:     int(stats.ChunkCount),
		TokenCount:     int(stats.TokenCount),
		Status:         library.Status,
		CreatedAt:      library.CreatedAt,
		UpdatedAt:      library.UpdatedAt,
	}, nil
}

// ListWithStats 获取库列表（带统计信息，返回精简字段）
func (s *LibraryService) ListWithStats(req *request.LibraryList) (*response.PageResult, error) {
	var libraries []dbmodel.Library
	var total int64

	db := global.DB.Model(&dbmodel.Library{})

	// 条件过滤
	if req.Name != nil && *req.Name != "" {
		db = db.Where("name LIKE ?", "%"+*req.Name+"%")
	}
	if req.Status != nil && *req.Status != "" {
		db = db.Where("status = ?", *req.Status)
	} else {
		db = db.Where("status = ?", "active")
	}

	// 计算总数
	if err := db.Count(&total).Error; err != nil {
		return nil, err
	}

	// 分页
	page := req.Page
	pageSize := req.PageSize
	if page <= 0 {
		page = 1
	}
	if pageSize <= 0 {
		pageSize = 10
	}

	offset := (page - 1) * pageSize

	// 排序
	sort := "updated_at DESC" // 默认按更新时间
	if req.Sort != nil && *req.Sort == "popular" {
		// 按 MCP 调用次数排序（LEFT JOIN statistics）
		db = db.Select("libraries.*, COALESCE(s.metric_value, 0) as popularity").
			Joins("LEFT JOIN statistics s ON s.library_id = libraries.id AND s.metric_name = ?", dbmodel.MetricMCPGetLibraryDocs)
		sort = "popularity DESC, updated_at DESC"
	}

	if err := db.Order(sort).Offset(offset).Limit(pageSize).Find(&libraries).Error; err != nil {
		return nil, err
	}

	// 转换为精简的列表响应
	result := make([]response.LibraryListItem, len(libraries))
	for i, lib := range libraries {
		// 统计 chunk 数和 token 数
		var stats struct {
			ChunkCount int64 `gorm:"column:chunk_count"`
			TokenCount int64 `gorm:"column:token_count"`
		}
		global.DB.Model(&dbmodel.DocumentUpload{}).
			Select("COALESCE(SUM(chunk_count), 0) as chunk_count, COALESCE(SUM(token_count), 0) as token_count").
			Where("library_id = ? AND status = ?", lib.ID, "completed").
			Scan(&stats)

		result[i] = response.LibraryListItem{
			ID:             lib.ID,
			Name:           lib.Name,
			SourceType:     lib.SourceType,
			SourceURL:      lib.SourceURL,
			DefaultVersion: lib.DefaultVersion,
			TokenCount:     int(stats.TokenCount),
			ChunkCount:     int(stats.ChunkCount),
			UpdatedAt:      lib.UpdatedAt,
		}
	}

	return &response.PageResult{
		List:     result,
		Total:    total,
		Page:     page,
		PageSize: pageSize,
	}, nil
}

// GetVersions 获取库的所有版本（用于上传时选择）
func (s *LibraryService) GetVersions(libraryID uint) ([]response.VersionInfo, error) {
	// 检查库是否存在
	var library dbmodel.Library
	if err := global.DB.First(&library, libraryID).Error; err != nil {
		return nil, ErrNotFound
	}

	// 查询每个版本的统计信息
	type versionStats struct {
		Version     string    `gorm:"column:version"`
		TokenCount  int       `gorm:"column:token_count"`
		ChunkCount  int       `gorm:"column:chunk_count"`
		LastUpdated time.Time `gorm:"column:last_updated"`
	}
	var stats []versionStats
	if err := global.DB.Table("document_uploads").
		Select("version, COALESCE(SUM(token_count), 0) as token_count, COALESCE(SUM(chunk_count), 0) as chunk_count, MAX(updated_at) as last_updated").
		Where("library_id = ? AND status != ?", libraryID, "deleted").
		Group("version").
		Find(&stats).Error; err != nil {
		return nil, err
	}

	// 构建 version -> stats 映射
	statsMap := make(map[string]versionStats)
	for _, stat := range stats {
		statsMap[stat.Version] = stat
	}

	var versions []response.VersionInfo

	// 先添加 default_version
	if stat, ok := statsMap[library.DefaultVersion]; ok {
		versions = append(versions, response.VersionInfo{
			Version:     library.DefaultVersion,
			TokenCount:  stat.TokenCount,
			ChunkCount:  stat.ChunkCount,
			LastUpdated: stat.LastUpdated,
		})
	} else {
		versions = append(versions, response.VersionInfo{
			Version:     library.DefaultVersion,
			TokenCount:  0,
			ChunkCount:  0,
			LastUpdated: library.UpdatedAt,
		})
	}

	// 再添加 versions 数组中的所有版本（倒序）
	for i := len(library.Versions) - 1; i >= 0; i-- {
		v := library.Versions[i]
		if stat, ok := statsMap[v]; ok {
			versions = append(versions, response.VersionInfo{
				Version:     v,
				TokenCount:  stat.TokenCount,
				ChunkCount:  stat.ChunkCount,
				LastUpdated: stat.LastUpdated,
			})
		} else {
			versions = append(versions, response.VersionInfo{
				Version:     v,
				TokenCount:  0,
				ChunkCount:  0,
				LastUpdated: library.UpdatedAt,
			})
		}
	}

	return versions, nil
}

// CreateVersion 创建新版本
func (s *LibraryService) CreateVersion(libraryID uint, version string) error {
	// 自动添加 v 前缀（如果没有的话）
	if !regexp.MustCompile(`^v`).MatchString(version) {
		version = "v" + version
	}

	// 验证版本格式
	if err := s.ValidateVersion(version); err != nil {
		return err
	}

	// 检查库是否存在
	var library dbmodel.Library
	if err := global.DB.First(&library, libraryID).Error; err != nil {
		return ErrNotFound
	}

	// 检查版本是否已存在于 versions 数组
	versionInArray := false
	for _, v := range library.Versions {
		if v == version {
			versionInArray = true
			break
		}
	}

	if versionInArray {
		return ErrVersionExists
	}

	// 检查 document_uploads 表是否有该版本的文档
	var count int64
	if err := global.DB.Table("document_uploads").
		Where("library_id = ? AND version = ?", libraryID, version).
		Count(&count).Error; err != nil {
		return err
	}

	// 添加版本到 versions 数组（即使文档已存在，也要确保版本在数组中）
	library.Versions = append(library.Versions, version)

	// 保存到数据库
	if err := global.DB.Model(&library).Update("versions", library.Versions).Error; err != nil {
		return err
	}

	return nil
}

// DeleteVersion 删除版本及其所有文档和分块
func (s *LibraryService) DeleteVersion(libraryID uint, version string) error {
	// 检查库是否存在
	var library dbmodel.Library
	if err := global.DB.First(&library, libraryID).Error; err != nil {
		return ErrNotFound
	}

	// 检查版本是否存在
	var count int64
	if err := global.DB.Table("document_uploads").
		Where("library_id = ? AND version = ?", libraryID, version).
		Count(&count).Error; err != nil {
		return err
	}

	if count == 0 {
		return ErrNotFound
	}

	// 开始事务
	tx := global.DB.Begin()

	// 获取该版本的所有文档 ID
	var documentIDs []uint
	if err := tx.Table("document_uploads").
		Select("id").
		Where("library_id = ? AND version = ?", libraryID, version).
		Scan(&documentIDs).Error; err != nil {
		tx.Rollback()
		return err
	}

	// 删除分块
	if err := tx.Where("upload_id IN ?", documentIDs).
		Delete(&dbmodel.DocumentChunk{}).Error; err != nil {
		tx.Rollback()
		return err
	}

	// 删除文档
	if err := tx.Where("library_id = ? AND version = ?", libraryID, version).
		Delete(&dbmodel.DocumentUpload{}).Error; err != nil {
		tx.Rollback()
		return err
	}

	// 从 library.Versions 数组中移除该版本
	newVersions := make([]string, 0, len(library.Versions))
	for _, v := range library.Versions {
		if v != version {
			newVersions = append(newVersions, v)
		}
	}

	// 更新 library
	updates := map[string]interface{}{
		"versions": pq.StringArray(newVersions),
	}

	// 如果删除的是默认版本，更新默认版本
	if library.DefaultVersion == version {
		if len(newVersions) > 0 {
			updates["default_version"] = newVersions[0]
		} else {
			updates["default_version"] = ""
		}
	}

	if err := tx.Model(&library).Updates(updates).Error; err != nil {
		tx.Rollback()
		return err
	}

	if err := tx.Commit().Error; err != nil {
		return err
	}

	return nil
}

// RefreshVersion 刷新版本（重新处理所有文档）
func (s *LibraryService) RefreshVersion(libraryID uint, version string, actorID string) error {
	// 生成任务 ID
	taskID := utils.GenerateTaskID()

	// 检查库是否存在
	var library dbmodel.Library
	if err := global.DB.First(&library, libraryID).Error; err != nil {
		return ErrNotFound
	}

	// 检查版本是否存在
	var count int64
	if err := global.DB.Table("document_uploads").
		Where("library_id = ? AND version = ?", libraryID, version).
		Count(&count).Error; err != nil {
		return err
	}

	if count == 0 {
		return ErrNotFound
	}

	// 获取该版本的所有文档
	var documents []dbmodel.DocumentUpload
	if err := global.DB.Where("library_id = ? AND version = ?", libraryID, version).
		Find(&documents).Error; err != nil {
		return err
	}

	// 创建任务日志器
	actLogger := actlog.NewTaskLogger(libraryID, taskID, version).
		WithTarget("version", version).
		WithActor(actorID)

	// 同步写入"开始"日志（确保 API 返回前日志已入库）
	startLog := &dbmodel.ActivityLog{
		LibraryID:  libraryID,
		ActorID:    actorID,
		TaskID:     taskID,
		Event:      actlog.EventVerRefresh,
		Status:     actlog.StatusInfo,
		Message:    fmt.Sprintf("开始刷新版本: %s (%d 个文档)", version, len(documents)),
		Version:    version,
		TargetType: "version",
		TargetID:   version,
	}
	global.DB.Create(startLog)

	// 开始事务
	tx := global.DB.Begin()

	// 对每个文档重新处理
	for _, doc := range documents {
		// 删除该文档的旧分块
		if err := tx.Where("upload_id = ?", doc.ID).
			Delete(&dbmodel.DocumentChunk{}).Error; err != nil {
			tx.Rollback()
			return err
		}

		// 标记文档为处理中
		if err := tx.Model(&doc).Update("status", "processing").Error; err != nil {
			tx.Rollback()
			return err
		}
	}

	if err := tx.Commit().Error; err != nil {
		return err
	}

	// 异步重新处理所有文档
	go func() {
		processor := &DocumentProcessor{}
		var wg sync.WaitGroup
		var successCount, failCount int
		var mu sync.Mutex

		for _, doc := range documents {
			wg.Add(1)
			docCopy := doc // 避免闭包问题
			go func() {
				defer wg.Done()

				// 创建文档级别日志器
				docLogger := actLogger.WithTarget("document", strconv.FormatUint(uint64(docCopy.ID), 10))

				// 从存储下载文件内容
				reader, err := global.Storage.Download(context.Background(), docCopy.FilePath)
				if err != nil {
					log.Printf("[RefreshVersion] Failed to download file %s: %v", docCopy.FilePath, err)
					docLogger.Warning(actlog.EventDocFailed, fmt.Sprintf("下载失败: %s", docCopy.Title))
					global.DB.Model(&docCopy).Update("status", "failed")
					mu.Lock()
					failCount++
					mu.Unlock()
					return
				}
				content, err := io.ReadAll(reader)
				reader.Close()
				if err != nil {
					log.Printf("[RefreshVersion] Failed to read file content %s: %v", docCopy.FilePath, err)
					docLogger.Warning(actlog.EventDocFailed, fmt.Sprintf("读取失败: %s", docCopy.Title))
					global.DB.Model(&docCopy).Update("status", "failed")
					mu.Lock()
					failCount++
					mu.Unlock()
					return
				}

				log.Printf("[RefreshVersion] Starting to reprocess document: %s (ID: %d)", docCopy.Title, docCopy.ID)
				if err := processor.ProcessDocument(&docCopy, content, docLogger); err != nil {
					log.Printf("[RefreshVersion] Failed to process document %d: %v", docCopy.ID, err)
					global.DB.Model(&docCopy).Update("status", "failed")
					mu.Lock()
					failCount++
					mu.Unlock()
					return
				}

				mu.Lock()
				successCount++
				mu.Unlock()
			}()
		}

		// 等待所有文档处理完成
		wg.Wait()

		// 记录刷新完成
		if failCount == 0 {
			actLogger.Success(actlog.EventVerRefresh, fmt.Sprintf("刷新完成: %s (成功 %d)", version, successCount))
		} else {
			actLogger.Warning(actlog.EventVerRefresh, fmt.Sprintf("刷新完成: %s (成功 %d, 失败 %d)", version, successCount, failCount))
		}

		log.Printf("[RefreshVersion] Completed reprocessing for library %d version %s: success=%d, fail=%d", libraryID, version, successCount, failCount)
	}()

	log.Printf("[RefreshVersion] Started reprocessing %d documents for library %d version %s", len(documents), libraryID, version)

	return nil
}

// RefreshVersionWithCallback 刷新版本（带 SSE 状态回调，无感知更新）
// 使用批次版本号实现原子切换，确保刷新过程中检索不受影响
func (s *LibraryService) RefreshVersionWithCallback(libraryID uint, version string, actorID string, statusChan chan response.RefreshStatus) {
	defer close(statusChan)

	// 生成任务 ID 和日志器
	taskID := utils.GenerateTaskID()
	actLogger := actlog.NewTaskLogger(libraryID, taskID, version).
		WithTarget("version", version).
		WithActor(actorID)

	// 检查库是否存在
	var library dbmodel.Library
	if err := global.DB.First(&library, libraryID).Error; err != nil {
		statusChan <- response.RefreshStatus{Stage: "error", Message: "库不存在"}
		return
	}

	// 获取该版本的所有文档
	var documents []dbmodel.DocumentUpload
	if err := global.DB.Where("library_id = ? AND version = ?", libraryID, version).
		Find(&documents).Error; err != nil {
		statusChan <- response.RefreshStatus{Stage: "error", Message: "获取文档列表失败"}
		return
	}

	if len(documents) == 0 {
		statusChan <- response.RefreshStatus{Stage: "error", Message: "该版本没有文档"}
		return
	}

	total := len(documents)
	actLogger.Info(actlog.EventVerRefresh, fmt.Sprintf("开始刷新版本: %s (%d 个文档，无感知模式)", version, total))
	statusChan <- response.RefreshStatus{
		Stage:   "started",
		Current: 0,
		Total:   total,
		Message: fmt.Sprintf("开始刷新 %d 个文档（无感知模式）", total),
	}

	// 生成新的批次版本号
	batchVersion := time.Now().UnixNano()
	log.Printf("[RefreshVersion] Starting refresh with batchVersion=%d for library %d version %s", batchVersion, libraryID, version)

	// 收集所有新 chunks 和文档统计
	processor := &DocumentProcessor{}
	type docResult struct {
		doc         *dbmodel.DocumentUpload
		chunks      []*dbmodel.DocumentChunk
		totalTokens int
		err         error
	}
	results := make([]docResult, 0, total)
	successCount := 0

	// 阶段1：处理所有文档，生成新 chunks（status=pending）
	for i, doc := range documents {
		docCopy := doc // 避免闭包问题
		statusChan <- response.RefreshStatus{
			DocID:    doc.ID,
			DocTitle: doc.Title,
			Stage:    "doc_processing",
			Current:  i + 1,
			Total:    total,
			Message:  fmt.Sprintf("正在处理: %s", doc.Title),
		}

		// 从存储读取文件内容
		reader, err := global.Storage.Download(context.Background(), doc.FilePath)
		if err != nil {
			log.Printf("[RefreshVersion] Failed to download file %s: %v", doc.FilePath, err)
			results = append(results, docResult{doc: &docCopy, err: err})
			statusChan <- response.RefreshStatus{
				DocID:    doc.ID,
				DocTitle: doc.Title,
				Stage:    "doc_failed",
				Current:  i + 1,
				Total:    total,
				Message:  fmt.Sprintf("下载文件失败: %s", doc.Title),
			}
			continue
		}
		content, err := io.ReadAll(reader)
		reader.Close()
		if err != nil {
			log.Printf("[RefreshVersion] Failed to read file content %s: %v", doc.FilePath, err)
			results = append(results, docResult{doc: &docCopy, err: err})
			statusChan <- response.RefreshStatus{
				DocID:    doc.ID,
				DocTitle: doc.Title,
				Stage:    "doc_failed",
				Current:  i + 1,
				Total:    total,
				Message:  fmt.Sprintf("读取文件内容失败: %s", doc.Title),
			}
			continue
		}

		// 处理文档，生成 chunks（不写入数据库）
		docLogger := actLogger.WithTarget("document", strconv.FormatUint(uint64(docCopy.ID), 10))
		chunks, totalTokens, err := processor.ProcessDocumentForRefresh(&docCopy, content, batchVersion, docLogger)
		if err != nil {
			log.Printf("[RefreshVersion] Failed to process document %d: %v", doc.ID, err)
			results = append(results, docResult{doc: &docCopy, err: err})
			statusChan <- response.RefreshStatus{
				DocID:    doc.ID,
				DocTitle: doc.Title,
				Stage:    "doc_failed",
				Current:  i + 1,
				Total:    total,
				Message:  fmt.Sprintf("处理失败: %s", doc.Title),
			}
			continue
		}

		results = append(results, docResult{doc: &docCopy, chunks: chunks, totalTokens: totalTokens})
		successCount++
		statusChan <- response.RefreshStatus{
			DocID:    doc.ID,
			DocTitle: doc.Title,
			Stage:    "doc_completed",
			Current:  i + 1,
			Total:    total,
			Message:  fmt.Sprintf("处理完成: %s（%d chunks）", doc.Title, len(chunks)),
		}
	}

	// 如果没有成功处理任何文档，直接返回
	if successCount == 0 {
		statusChan <- response.RefreshStatus{Stage: "error", Message: "所有文档处理失败"}
		return
	}

	// 阶段2：原子切换（事务）
	statusChan <- response.RefreshStatus{
		Stage:   "switching",
		Current: total,
		Total:   total,
		Message: "正在原子切换数据...",
	}

	tx := global.DB.Begin()
	for _, result := range results {
		if result.err != nil {
			// 标记失败的文档
			tx.Model(result.doc).Update("status", "failed")
			continue
		}

		// 插入新 chunks（status=pending）
		if len(result.chunks) > 0 {
			if err := tx.CreateInBatches(result.chunks, 100).Error; err != nil {
				tx.Rollback()
				log.Printf("[RefreshVersion] Failed to insert new chunks: %v", err)
				statusChan <- response.RefreshStatus{Stage: "error", Message: "插入新数据失败"}
				return
			}
		}

		// 原子切换：旧 chunks -> deleted，新 chunks -> active
		// 1. 将旧的 active chunks 标记为 deleted
		if err := tx.Model(&dbmodel.DocumentChunk{}).
			Where("upload_id = ? AND status = ? AND batch_version < ?", result.doc.ID, "active", batchVersion).
			Update("status", "deleted").Error; err != nil {
			tx.Rollback()
			log.Printf("[RefreshVersion] Failed to mark old chunks as deleted: %v", err)
			statusChan <- response.RefreshStatus{Stage: "error", Message: "标记旧数据失败"}
			return
		}

		// 2. 将新的 pending chunks 激活
		if err := tx.Model(&dbmodel.DocumentChunk{}).
			Where("upload_id = ? AND batch_version = ? AND status = ?", result.doc.ID, batchVersion, "pending").
			Update("status", "active").Error; err != nil {
			tx.Rollback()
			log.Printf("[RefreshVersion] Failed to activate new chunks: %v", err)
			statusChan <- response.RefreshStatus{Stage: "error", Message: "激活新数据失败"}
			return
		}

		// 更新文档状态和统计
		if err := tx.Model(result.doc).Updates(map[string]interface{}{
			"status":      "completed",
			"chunk_count": len(result.chunks),
			"token_count": result.totalTokens,
		}).Error; err != nil {
			tx.Rollback()
			log.Printf("[RefreshVersion] Failed to update document status: %v", err)
			statusChan <- response.RefreshStatus{Stage: "error", Message: "更新文档状态失败"}
			return
		}
	}

	if err := tx.Commit().Error; err != nil {
		log.Printf("[RefreshVersion] Failed to commit transaction: %v", err)
		statusChan <- response.RefreshStatus{Stage: "error", Message: "事务提交失败"}
		return
	}

	// 阶段3：对标记为 deleted 的旧数据执行 GORM 软删除（设置 deleted_at）
	go func() {
		result := global.DB.Where("library_id = ? AND version = ? AND status = ?", libraryID, version, "deleted").
			Delete(&dbmodel.DocumentChunk{})
		if result.Error != nil {
			log.Printf("[RefreshVersion] WARNING: Failed to soft delete old chunks: %v", result.Error)
		} else if result.RowsAffected > 0 {
			log.Printf("[RefreshVersion] Soft deleted %d old chunks for library %d version %s", result.RowsAffected, libraryID, version)
		}
	}()

	// 失效缓存
	searchService := &SearchService{}
	if err := searchService.InvalidateLibraryCache(libraryID, version); err != nil {
		log.Printf("[RefreshVersion] WARNING: Failed to invalidate cache: %v", err)
	}

	// 记录完成日志
	if successCount == total {
		actLogger.Success(actlog.EventVerRefresh, fmt.Sprintf("刷新完成: %s (成功 %d)", version, successCount))
	} else {
		actLogger.Warning(actlog.EventVerRefresh, fmt.Sprintf("刷新完成: %s (成功 %d, 失败 %d)", version, successCount, total-successCount))
	}

	statusChan <- response.RefreshStatus{
		Stage:   "all_completed",
		Current: total,
		Total:   total,
		Message: fmt.Sprintf("全部完成，成功处理 %d/%d 个文档", successCount, total),
	}

	log.Printf("[RefreshVersion] Completed refresh for library %d version %s: %d/%d documents", libraryID, version, successCount, total)
}
