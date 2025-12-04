package service

import (
	"crypto/sha256"
	"encoding/hex"
	"io"
	"mime/multipart"
	"path/filepath"

	dbmodel "go-mcp-context/internal/model/database"
	"go-mcp-context/internal/model/request"
	"go-mcp-context/internal/model/response"
	"go-mcp-context/pkg/global"
)

type DocumentService struct{}

// List 获取文档列表
func (s *DocumentService) List(req *request.DocumentList) (*response.PageResult, error) {
	var documents []dbmodel.Document
	var total int64

	db := global.DB.Model(&dbmodel.Document{})

	// 条件过滤
	if req.LibraryID != nil && *req.LibraryID > 0 {
		db = db.Where("library_id = ?", *req.LibraryID)
	}
	if req.Title != nil && *req.Title != "" {
		db = db.Where("title LIKE ?", "%"+*req.Title+"%")
	}
	if req.FileType != nil && *req.FileType != "" {
		db = db.Where("file_type = ?", *req.FileType)
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
	if err := db.Offset(offset).Limit(pageSize).Find(&documents).Error; err != nil {
		return nil, err
	}

	return &response.PageResult{
		List:     documents,
		Total:    total,
		Page:     page,
		PageSize: pageSize,
	}, nil
}

// Upload 上传文档
func (s *DocumentService) Upload(libraryID uint, file multipart.File, header *multipart.FileHeader) (*dbmodel.Document, error) {
	// 检查库是否存在
	var library dbmodel.Library
	if err := global.DB.First(&library, libraryID).Error; err != nil {
		return nil, ErrNotFound
	}

	// 读取文件内容
	content, err := io.ReadAll(file)
	if err != nil {
		return nil, err
	}

	// 计算内容哈希
	hash := sha256.Sum256(content)
	contentHash := hex.EncodeToString(hash[:])

	// 检查是否已存在相同内容的文档
	var existingDoc dbmodel.Document
	if err := global.DB.Where("library_id = ? AND content_hash = ? AND status = ?",
		libraryID, contentHash, "active").First(&existingDoc).Error; err == nil {
		return nil, ErrAlreadyExists
	}

	// 确定文件类型
	ext := filepath.Ext(header.Filename)
	fileType := getFileType(ext)
	if fileType == "" {
		return nil, ErrInvalidParams
	}

	// 创建文档记录
	doc := &dbmodel.Document{
		LibraryID:   libraryID,
		Title:       header.Filename,
		FilePath:    header.Filename,
		FileType:    fileType,
		FileSize:    int64(len(content)),
		ContentHash: contentHash,
		Status:      "active",
	}

	if err := global.DB.Create(doc).Error; err != nil {
		return nil, err
	}

	// 异步处理文档（解析、分块、生成 Embedding）
	processor := &DocumentProcessor{}
	processor.ProcessDocumentAsync(doc, content)

	return doc, nil
}

// GetByID 根据 ID 获取文档
func (s *DocumentService) GetByID(id uint) (*dbmodel.Document, error) {
	var doc dbmodel.Document
	if err := global.DB.First(&doc, id).Error; err != nil {
		return nil, ErrNotFound
	}
	return &doc, nil
}

// Delete 删除文档（软删除）
func (s *DocumentService) Delete(id uint) error {
	result := global.DB.Model(&dbmodel.Document{}).
		Where("id = ?", id).
		Update("status", "deleted")

	if result.Error != nil {
		return result.Error
	}

	if result.RowsAffected == 0 {
		return ErrNotFound
	}

	// 同时删除关联的 chunks
	global.DB.Model(&dbmodel.DocumentChunk{}).
		Where("document_id = ?", id).
		Update("status", "deleted")

	return nil
}

// getFileType 根据扩展名返回文件类型
func getFileType(ext string) string {
	switch ext {
	case ".md", ".markdown":
		return "markdown"
	case ".pdf":
		return "pdf"
	case ".docx":
		return "docx"
	case ".json", ".yaml", ".yml":
		return "swagger"
	default:
		return ""
	}
}

// GetChunksByLibrary 获取库的文档块（按热度排序）
func (s *DocumentService) GetChunksByLibrary(libraryID uint, mode string, page, limit int) ([]dbmodel.DocumentChunk, int64, error) {
	var chunks []dbmodel.DocumentChunk
	var total int64

	query := global.DB.Model(&dbmodel.DocumentChunk{}).
		Where("library_id = ? AND status = ?", libraryID, "active")

	if mode != "" {
		query = query.Where("chunk_type = ?", mode)
	}

	// 计算总数
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// 分页查询（按热度排序）
	offset := (page - 1) * limit
	if err := query.Order("access_count DESC").
		Offset(offset).
		Limit(limit).
		Find(&chunks).Error; err != nil {
		return nil, 0, err
	}

	return chunks, total, nil
}
