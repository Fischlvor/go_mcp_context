package service

import (
	"time"

	dbmodel "go-mcp-context/internal/model/database"
	"go-mcp-context/internal/model/request"
	"go-mcp-context/internal/model/response"
	"go-mcp-context/pkg/global"
)

type LibraryService struct{}

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
	library := &dbmodel.Library{
		Name:        req.Name,
		Version:     req.Version,
		Description: req.Description,
		Metadata:    dbmodel.JSON(req.Metadata),
		Status:      "active",
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
	library.Version = req.Version
	library.Description = req.Description
	library.Metadata = dbmodel.JSON(req.Metadata)

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

	// 同时删除关联的文档和 chunks
	global.DB.Model(&dbmodel.Document{}).
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

// GetByNameVersion 根据名称和版本获取库
func (s *LibraryService) GetByNameVersion(name, version string) (*dbmodel.Library, error) {
	var library dbmodel.Library
	if err := global.DB.Where("name = ? AND version = ? AND status = ?", name, version, "active").
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

	// 统计文档数
	var docCount int64
	global.DB.Model(&dbmodel.Document{}).
		Where("library_id = ? AND status = ?", id, "active").
		Count(&docCount)

	// 统计 chunk 数和 token 数
	var stats struct {
		ChunkCount int64 `gorm:"column:chunk_count"`
		TokenCount int64 `gorm:"column:token_count"`
	}
	global.DB.Model(&dbmodel.Document{}).
		Select("COALESCE(SUM(chunk_count), 0) as chunk_count, COALESCE(SUM(token_count), 0) as token_count").
		Where("library_id = ? AND status = ?", id, "active").
		Scan(&stats)

	return &response.LibraryInfo{
		ID:            library.ID,
		Name:          library.Name,
		Version:       library.Version,
		Description:   library.Description,
		DocumentCount: int(docCount),
		ChunkCount:    int(stats.ChunkCount),
		TokenCount:    int(stats.TokenCount),
		Status:        library.Status,
		CreatedAt:     library.CreatedAt,
		UpdatedAt:     library.UpdatedAt,
	}, nil
}

// ListWithStats 获取库列表（带统计信息）
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
	if err := db.Offset(offset).Limit(pageSize).Find(&libraries).Error; err != nil {
		return nil, err
	}

	// 转换为带统计信息的响应
	result := make([]response.LibraryInfo, len(libraries))
	for i, lib := range libraries {
		// 统计文档数
		var docCount int64
		global.DB.Model(&dbmodel.Document{}).
			Where("library_id = ? AND status = ?", lib.ID, "active").
			Count(&docCount)

		// 统计 chunk 数和 token 数
		var stats struct {
			ChunkCount int64 `gorm:"column:chunk_count"`
			TokenCount int64 `gorm:"column:token_count"`
		}
		global.DB.Model(&dbmodel.Document{}).
			Select("COALESCE(SUM(chunk_count), 0) as chunk_count, COALESCE(SUM(token_count), 0) as token_count").
			Where("library_id = ? AND status = ?", lib.ID, "active").
			Scan(&stats)

		result[i] = response.LibraryInfo{
			ID:            lib.ID,
			Name:          lib.Name,
			Version:       lib.Version,
			Description:   lib.Description,
			DocumentCount: int(docCount),
			ChunkCount:    int(stats.ChunkCount),
			TokenCount:    int(stats.TokenCount),
			Status:        lib.Status,
			CreatedAt:     lib.CreatedAt,
			UpdatedAt:     lib.UpdatedAt,
		}
	}

	return &response.PageResult{
		List:     result,
		Total:    total,
		Page:     page,
		PageSize: pageSize,
	}, nil
}
