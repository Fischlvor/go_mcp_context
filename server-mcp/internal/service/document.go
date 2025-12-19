package service

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"io"
	"mime/multipart"
	"path/filepath"
	"regexp"
	"strings"
	"time"

	dbmodel "go-mcp-context/internal/model/database"
	"go-mcp-context/internal/model/request"
	"go-mcp-context/internal/model/response"
	"go-mcp-context/pkg/global"

	"gorm.io/gorm"
)

type DocumentService struct{}

// List 获取文档上传记录列表
func (s *DocumentService) List(req *request.DocumentList) (*response.PageResult, error) {
	var documents []dbmodel.DocumentUpload
	var total int64

	// 构建基础查询条件
	baseQuery := global.DB.Model(&dbmodel.DocumentUpload{})

	// 条件过滤
	if req.LibraryID != nil && *req.LibraryID > 0 {
		baseQuery = baseQuery.Where("library_id = ?", *req.LibraryID)

		// 版本过滤：不传 version 时使用 library 的 default_version
		if req.Version != nil && *req.Version != "" {
			baseQuery = baseQuery.Where("version = ?", *req.Version)
		} else {
			// 查询 library 的 default_version
			var library dbmodel.Library
			if err := global.DB.Select("default_version").First(&library, *req.LibraryID).Error; err == nil {
				baseQuery = baseQuery.Where("version = ?", library.DefaultVersion)
			}
		}
	}
	// 默认查询非删除状态的文档
	baseQuery = baseQuery.Where("status != ?", "deleted")

	// 计算总数（使用 Session 克隆避免影响后续查询）
	if err := baseQuery.Session(&gorm.Session{}).Count(&total).Error; err != nil {
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
	if err := baseQuery.Order("created_at DESC").Offset(offset).Limit(pageSize).Find(&documents).Error; err != nil {
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
func (s *DocumentService) Upload(libraryID uint, version string, file multipart.File, header *multipart.FileHeader) (*dbmodel.DocumentUpload, error) {
	// 检查库是否存在
	var library dbmodel.Library
	if err := global.DB.First(&library, libraryID).Error; err != nil {
		return nil, ErrNotFound
	}

	// 检查版本是否存在
	versionExists := version == library.DefaultVersion
	if !versionExists {
		for _, v := range library.Versions {
			if v == version {
				versionExists = true
				break
			}
		}
	}
	if !versionExists {
		return nil, fmt.Errorf("version %s does not exist", version)
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
	var existingDoc dbmodel.DocumentUpload
	if err := global.DB.Where("library_id = ? AND version = ? AND content_hash = ? AND status != ?",
		libraryID, version, contentHash, "deleted").First(&existingDoc).Error; err == nil {
		return nil, ErrAlreadyExists
	}

	// 确定文件类型
	ext := filepath.Ext(header.Filename)
	fileType := getFileType(ext)
	if fileType == "" {
		return nil, ErrInvalidParams
	}

	// 生成存储 Key: {path_prefix}/{lib_name}/{version}/{filename}
	libDir := sanitizeFileName(library.Name)
	versionDir := sanitizeFileName(version)
	key := filepath.Join(global.Config.Qiniu.PathPrefix, libDir, versionDir, header.Filename)

	// 使用 Storage 接口上传
	result, err := global.Storage.Upload(context.Background(), key, file, header.Size, "")
	if err != nil {
		return nil, fmt.Errorf("failed to upload file: %w", err)
	}

	// 创建文档上传记录（初始状态为 processing）
	doc := &dbmodel.DocumentUpload{
		LibraryID:   libraryID,
		Version:     version,
		Title:       header.Filename,
		FilePath:    result.Key,
		FileType:    fileType,
		FileSize:    int64(len(content)),
		ContentHash: contentHash,
		Status:      "processing",
	}

	if err := global.DB.Create(doc).Error; err != nil {
		// 如果数据库创建失败，删除已上传的文件
		global.Storage.Delete(context.Background(), result.Key)
		return nil, err
	}

	// 异步处理文档（解析、分块、生成 Embedding）
	processor := &DocumentProcessor{}
	processor.ProcessDocumentAsync(doc, content)

	return doc, nil
}

// UploadWithCallback 上传文档（带状态回调）
func (s *DocumentService) UploadWithCallback(libraryID uint, version string, file multipart.File, header *multipart.FileHeader, statusChan chan response.ProcessStatus) (*dbmodel.DocumentUpload, error) {
	// 检查库是否存在
	var library dbmodel.Library
	if err := global.DB.First(&library, libraryID).Error; err != nil {
		close(statusChan)
		return nil, ErrNotFound
	}

	// 检查版本是否存在
	versionExists := version == library.DefaultVersion
	if !versionExists {
		for _, v := range library.Versions {
			if v == version {
				versionExists = true
				break
			}
		}
	}
	if !versionExists {
		close(statusChan)
		return nil, fmt.Errorf("version %s does not exist", version)
	}

	// 读取文件内容
	content, err := io.ReadAll(file)
	if err != nil {
		close(statusChan)
		return nil, err
	}

	// 计算内容哈希
	hash := sha256.Sum256(content)
	contentHash := hex.EncodeToString(hash[:])

	// 检查是否已存在相同内容的文档
	var existingDoc dbmodel.DocumentUpload
	if err := global.DB.Where("library_id = ? AND version = ? AND content_hash = ? AND status != ?",
		libraryID, version, contentHash, "deleted").First(&existingDoc).Error; err == nil {
		close(statusChan)
		return nil, ErrAlreadyExists
	}

	// 确定文件类型
	ext := filepath.Ext(header.Filename)
	fileType := getFileType(ext)
	if fileType == "" {
		close(statusChan)
		return nil, ErrInvalidParams
	}

	// 生成存储 Key: {path_prefix}/{lib_name}/{version}/{filename}
	libDir := sanitizeFileName(library.Name)
	versionDir := sanitizeFileName(version)
	key := filepath.Join(global.Config.Qiniu.PathPrefix, libDir, versionDir, header.Filename)

	// 使用 Storage 接口上传
	result, err := global.Storage.Upload(context.Background(), key, file, header.Size, "")
	if err != nil {
		close(statusChan)
		return nil, fmt.Errorf("failed to upload file: %w", err)
	}

	// 创建文档上传记录
	doc := &dbmodel.DocumentUpload{
		LibraryID:   libraryID,
		Version:     version,
		Title:       header.Filename,
		FilePath:    result.Key,
		FileType:    fileType,
		FileSize:    int64(len(content)),
		ContentHash: contentHash,
		Status:      "processing",
	}

	if err := global.DB.Create(doc).Error; err != nil {
		global.Storage.Delete(context.Background(), result.Key)
		close(statusChan)
		return nil, err
	}

	// 异步处理文档（带状态回调）
	processor := &DocumentProcessor{}
	go processor.ProcessDocumentWithCallback(doc, content, statusChan)

	return doc, nil
}

// sanitizeFileName 清理文件名，移除不安全字符
func sanitizeFileName(name string) string {
	// 替换空格为下划线
	name = strings.ReplaceAll(name, " ", "_")
	// 只保留字母、数字、下划线、点和连字符
	reg := regexp.MustCompile(`[^a-zA-Z0-9._-]`)
	name = reg.ReplaceAllString(name, "")
	// 转小写
	return strings.ToLower(name)
}

// GetByID 根据 ID 获取文档上传记录
func (s *DocumentService) GetByID(id uint) (*dbmodel.DocumentUpload, error) {
	var doc dbmodel.DocumentUpload
	if err := global.DB.First(&doc, id).Error; err != nil {
		return nil, ErrNotFound
	}
	return &doc, nil
}

// GetLatestContent 获取库的最新文档内容（按创建时间倒序）
// 如果指定了版本，则只查询该版本的文档
func (s *DocumentService) GetLatestContent(libraryID uint, version string) (string, string, error) {
	var doc dbmodel.DocumentUpload
	query := global.DB.Where("library_id = ? AND status = ?", libraryID, "completed")

	// 如果指定了版本，添加版本过滤
	if version != "" {
		query = query.Where("version = ?", version)
	}

	if err := query.Order("created_at DESC").First(&doc).Error; err != nil {
		return "", "", ErrNotFound
	}

	// 从存储读取文件内容
	reader, err := global.Storage.Download(context.Background(), doc.FilePath)
	if err != nil {
		return doc.Title, "", nil
	}
	defer reader.Close()

	content, err := io.ReadAll(reader)
	if err != nil {
		return doc.Title, "", nil
	}

	return doc.Title, string(content), nil
}

// GetChunks 获取库的文档块
// mode: "code" 只返回代码块, "info" 只返回文档块, "" 返回全部
func (s *DocumentService) GetChunks(libraryID uint, version string, mode string) ([]dbmodel.DocumentChunk, error) {
	var chunks []dbmodel.DocumentChunk
	query := global.DB.Where("library_id = ? AND status = ?", libraryID, "active")

	// 版本过滤
	if version != "" {
		query = query.Where("version = ?", version)
	}

	// 类型过滤
	if mode == "code" {
		query = query.Where("chunk_type = ?", "code")
	} else if mode == "info" {
		query = query.Where("chunk_type = ?", "info")
	}

	if err := query.Order("chunk_index ASC").Find(&chunks).Error; err != nil {
		return nil, err
	}

	return chunks, nil
}

// Delete 删除文档上传记录（软删除）
func (s *DocumentService) Delete(id uint) error {
	now := time.Now()
	result := global.DB.Model(&dbmodel.DocumentUpload{}).
		Where("id = ? AND deleted_at IS NULL", id).
		Updates(map[string]interface{}{"status": "deleted", "deleted_at": now})

	if result.Error != nil {
		return result.Error
	}

	if result.RowsAffected == 0 {
		return ErrNotFound
	}

	// 同时删除关联的 chunks
	global.DB.Model(&dbmodel.DocumentChunk{}).
		Where("upload_id = ? AND deleted_at IS NULL", id).
		Updates(map[string]interface{}{"status": "deleted", "deleted_at": now})

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
