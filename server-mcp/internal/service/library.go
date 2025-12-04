package service

import (
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
	result := global.DB.Model(&dbmodel.Library{}).
		Where("id = ?", id).
		Update("status", "deleted")

	if result.Error != nil {
		return result.Error
	}

	if result.RowsAffected == 0 {
		return ErrNotFound
	}

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
