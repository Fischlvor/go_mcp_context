package storage

import (
	"context"
	"io"
	"time"
)

// Storage 文件存储接口
type Storage interface {
	// 核心方法
	Upload(ctx context.Context, key string, reader io.Reader, size int64, contentType string) (*UploadResult, error)
	Download(ctx context.Context, key string) (io.ReadCloser, error)
	Delete(ctx context.Context, key string) error

	// 批量操作
	DeleteByPrefix(ctx context.Context, prefix string) error // 删除某个版本的所有文件
	ListByPrefix(ctx context.Context, prefix string) ([]FileInfo, error)

	// 元数据
	GetFileInfo(ctx context.Context, key string) (*FileInfo, error)
	Exists(ctx context.Context, key string) (bool, error)

	// URL 生成
	GetPublicURL(key string) string                                                     // 获取公开访问 URL
	GetSignedURL(ctx context.Context, key string, expiry time.Duration) (string, error) // 获取签名 URL

	// 健康检查
	Health(ctx context.Context) error
}

// UploadResult 上传结果
type UploadResult struct {
	Key         string    // 存储 Key
	URL         string    // 访问 URL
	Size        int64     // 文件大小
	ContentType string    // 内容类型
	ETag        string    // 文件 Hash
	UploadedAt  time.Time // 上传时间
}

// FileInfo 文件信息
type FileInfo struct {
	Key          string    // 存储 Key
	Size         int64     // 文件大小
	ContentType  string    // 内容类型
	LastModified time.Time // 最后修改时间
	ETag         string    // 文件 Hash
}
