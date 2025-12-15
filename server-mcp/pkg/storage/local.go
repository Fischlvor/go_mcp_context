package storage

import (
	"context"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"time"
)

// LocalStorage 本地文件存储实现
type LocalStorage struct {
	basePath string // 基础路径，默认为 "uploads"
}

// NewLocalStorage 创建本地存储实例
func NewLocalStorage() *LocalStorage {
	return &LocalStorage{
		basePath: "uploads",
	}
}

// Upload 上传文件到本地存储
func (l *LocalStorage) Upload(ctx context.Context, key string, reader io.Reader, size int64, contentType string) (*UploadResult, error) {
	// 生成实际存储路径
	actualPath := filepath.Join(l.basePath, key)

	// 创建目录
	dir := filepath.Dir(actualPath)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return nil, fmt.Errorf("failed to create directory: %w", err)
	}

	// 创建文件
	file, err := os.Create(actualPath)
	if err != nil {
		return nil, fmt.Errorf("failed to create file: %w", err)
	}
	defer file.Close()

	// 复制内容
	if _, err := io.Copy(file, reader); err != nil {
		os.Remove(actualPath) // 删除失败的文件
		return nil, fmt.Errorf("failed to write file: %w", err)
	}

	return &UploadResult{
		Key:         key,
		URL:         "",
		Size:        size,
		ContentType: contentType,
		ETag:        "",
		UploadedAt:  time.Now(),
	}, nil
}

// Download 下载文件（本地存储直接返回文件句柄）
func (l *LocalStorage) Download(ctx context.Context, key string) (io.ReadCloser, error) {
	actualPath := filepath.Join(l.basePath, key)
	file, err := os.Open(actualPath)
	if err != nil {
		return nil, fmt.Errorf("failed to open file: %w", err)
	}
	return file, nil
}

// Delete 删除文件
func (l *LocalStorage) Delete(ctx context.Context, key string) error {
	actualPath := filepath.Join(l.basePath, key)
	if err := os.Remove(actualPath); err != nil {
		return fmt.Errorf("failed to delete file: %w", err)
	}
	return nil
}

// DeleteByPrefix 按前缀删除文件（删除整个目录）
func (l *LocalStorage) DeleteByPrefix(ctx context.Context, prefix string) error {
	actualPath := filepath.Join(l.basePath, prefix)
	if err := os.RemoveAll(actualPath); err != nil {
		return fmt.Errorf("failed to delete directory: %w", err)
	}
	return nil
}

// ListByPrefix 按前缀列出文件
func (l *LocalStorage) ListByPrefix(ctx context.Context, prefix string) ([]FileInfo, error) {
	actualPath := filepath.Join(l.basePath, prefix)

	var files []FileInfo
	err := filepath.Walk(actualPath, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if !info.IsDir() {
			// 计算相对路径（相对于 basePath）
			relPath, _ := filepath.Rel(l.basePath, path)
			files = append(files, FileInfo{
				Key:          relPath,
				Size:         info.Size(),
				ContentType:  "",
				LastModified: info.ModTime(),
				ETag:         "",
			})
		}
		return nil
	})

	if err != nil && !os.IsNotExist(err) {
		return nil, fmt.Errorf("failed to list files: %w", err)
	}

	return files, nil
}

// GetFileInfo 获取文件信息
func (l *LocalStorage) GetFileInfo(ctx context.Context, key string) (*FileInfo, error) {
	actualPath := filepath.Join(l.basePath, key)
	info, err := os.Stat(actualPath)
	if err != nil {
		return nil, fmt.Errorf("failed to get file info: %w", err)
	}

	return &FileInfo{
		Key:          key,
		Size:         info.Size(),
		ContentType:  "",
		LastModified: info.ModTime(),
		ETag:         "",
	}, nil
}

// Exists 检查文件是否存在
func (l *LocalStorage) Exists(ctx context.Context, key string) (bool, error) {
	actualPath := filepath.Join(l.basePath, key)
	_, err := os.Stat(actualPath)
	if err == nil {
		return true, nil
	}
	if os.IsNotExist(err) {
		return false, nil
	}
	return false, err
}

// GetPublicURL 获取公开访问 URL（本地存储不提供 URL）
func (l *LocalStorage) GetPublicURL(key string) string {
	return ""
}

// GetSignedURL 获取签名 URL（本地存储不支持）
func (l *LocalStorage) GetSignedURL(ctx context.Context, key string, expiry time.Duration) (string, error) {
	return "", fmt.Errorf("local storage does not support signed URLs")
}

// Health 健康检查
func (l *LocalStorage) Health(ctx context.Context) error {
	// 检查基础路径是否可访问
	if err := os.MkdirAll(l.basePath, 0755); err != nil {
		return fmt.Errorf("local storage health check failed: %w", err)
	}
	return nil
}

// 确保 LocalStorage 实现了 Storage 接口
var _ Storage = (*LocalStorage)(nil)
