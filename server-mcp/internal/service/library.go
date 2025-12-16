package service

import (
	"regexp"
	"time"

	dbmodel "go-mcp-context/internal/model/database"
	"go-mcp-context/internal/model/request"
	"go-mcp-context/internal/model/response"
	"go-mcp-context/pkg/global"
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

	library := &dbmodel.Library{
		Name:           req.Name,
		Description:    req.Description,
		SourceType:     sourceType,
		SourceURL:      req.SourceURL,
		Status:         "active",
		DefaultVersion: "default",
		Versions:       []string{}, // versions 只存正常版本，不包含 default
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
	if err := db.Order("updated_at DESC").Offset(offset).Limit(pageSize).Find(&libraries).Error; err != nil {
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

	var versions []response.VersionInfo

	// 先添加 default_version
	versions = append(versions, response.VersionInfo{
		Version:     library.DefaultVersion,
		TokenCount:  0,
		ChunkCount:  0,
		LastUpdated: library.UpdatedAt,
	})

	// 再添加 versions 数组中的所有版本（倒序）
	for i := len(library.Versions) - 1; i >= 0; i-- {
		versions = append(versions, response.VersionInfo{
			Version:     library.Versions[i],
			TokenCount:  0,
			ChunkCount:  0,
			LastUpdated: library.UpdatedAt,
		})
	}

	return versions, nil
}

// CreateVersion 创建新版本（只是标记，实际版本在上传文档时创建）
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

	// 检查版本是否已存在
	var count int64
	if err := global.DB.Table("document_uploads").
		Where("library_id = ? AND version = ?", libraryID, version).
		Count(&count).Error; err != nil {
		return err
	}

	if count > 0 {
		return ErrVersionExists
	}

	// 版本在上传文档时自动创建，这里只做验证
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
	if err := tx.Where("document_id IN ?", documentIDs).
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

	if err := tx.Commit().Error; err != nil {
		return err
	}

	return nil
}

// RefreshVersion 刷新版本（重新处理所有文档）
func (s *LibraryService) RefreshVersion(libraryID uint, version string) error {
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

	// 开始事务
	tx := global.DB.Begin()

	// 对每个文档重新处理
	for _, doc := range documents {
		// 删除该文档的旧分块
		if err := tx.Where("document_id = ?", doc.ID).
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

	// 异步处理文档（在后台队列中重新处理）
	// TODO: 将文档加入处理队列，由后台 worker 处理
	// 这里可以使用消息队列（如 Redis、RabbitMQ）或简单的 goroutine

	return nil
}
