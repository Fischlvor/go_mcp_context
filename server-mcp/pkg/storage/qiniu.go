package storage

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"

	"github.com/qiniu/go-sdk/v7/auth/qbox"
	"github.com/qiniu/go-sdk/v7/storage"
)

// QiniuStorage 七牛云存储实现
type QiniuStorage struct {
	accessKey     string
	secretKey     string
	bucket        string
	domain        string
	useHTTPS      bool
	useCdnDomains bool
	mac           *qbox.Mac
	cfg           *storage.Config
}

// QiniuConfig 七牛云配置
type QiniuConfig struct {
	AccessKey     string
	SecretKey     string
	Bucket        string
	Domain        string
	Zone          string // z0, z1, z2, na0, as0
	UseHTTPS      bool
	UseCdnDomains bool
}

// NewQiniuStorage 创建七牛云存储实例
func NewQiniuStorage(cfg QiniuConfig) *QiniuStorage {
	mac := qbox.NewMac(cfg.AccessKey, cfg.SecretKey)
	storageCfg := &storage.Config{
		UseHTTPS:      cfg.UseHTTPS,
		UseCdnDomains: cfg.UseCdnDomains,
	}

	// 设置区域
	switch cfg.Zone {
	case "z0", "ZoneHuadong":
		storageCfg.Zone = &storage.ZoneHuadong
	case "z1", "ZoneHuabei":
		storageCfg.Zone = &storage.ZoneHuabei
	case "z2", "ZoneHuanan":
		storageCfg.Zone = &storage.ZoneHuanan
	case "na0", "ZoneBeimei":
		storageCfg.Zone = &storage.ZoneBeimei
	case "as0", "ZoneXinjiapo":
		storageCfg.Zone = &storage.ZoneXinjiapo
	default:
		storageCfg.Zone = &storage.ZoneHuadong
	}

	return &QiniuStorage{
		accessKey:     cfg.AccessKey,
		secretKey:     cfg.SecretKey,
		bucket:        cfg.Bucket,
		domain:        cfg.Domain,
		useHTTPS:      cfg.UseHTTPS,
		useCdnDomains: cfg.UseCdnDomains,
		mac:           mac,
		cfg:           storageCfg,
	}
}

// Upload 上传文件
func (q *QiniuStorage) Upload(ctx context.Context, key string, reader io.Reader, size int64, contentType string) (*UploadResult, error) {
	// 覆盖上传策略
	putPolicy := storage.PutPolicy{
		Scope: fmt.Sprintf("%s:%s", q.bucket, key),
	}
	if contentType != "" {
		putPolicy.MimeLimit = contentType
	}
	upToken := putPolicy.UploadToken(q.mac)

	formUploader := storage.NewFormUploader(q.cfg)
	ret := storage.PutRet{}
	putExtra := storage.PutExtra{}

	err := formUploader.Put(ctx, &ret, upToken, key, reader, size, &putExtra)
	if err != nil {
		return nil, fmt.Errorf("qiniu upload failed: %w", err)
	}

	return &UploadResult{
		Key:         ret.Key,
		URL:         q.GetPublicURL(ret.Key),
		Size:        size,
		ContentType: contentType,
		ETag:        ret.Hash,
		UploadedAt:  time.Now(),
	}, nil
}

// Download 下载文件
func (q *QiniuStorage) Download(ctx context.Context, key string) (io.ReadCloser, error) {
	// 通过公开 URL 下载文件
	url := q.GetPublicURL(key)

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return nil, fmt.Errorf("create request failed: %w", err)
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("download failed: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		resp.Body.Close()
		return nil, fmt.Errorf("download failed: status %d", resp.StatusCode)
	}

	return resp.Body, nil
}

// Delete 删除文件
func (q *QiniuStorage) Delete(ctx context.Context, key string) error {
	bucketManager := storage.NewBucketManager(q.mac, q.cfg)
	err := bucketManager.Delete(q.bucket, key)
	if err != nil {
		return fmt.Errorf("qiniu delete failed: %w", err)
	}
	return nil
}

// DeleteByPrefix 按前缀批量删除
func (q *QiniuStorage) DeleteByPrefix(ctx context.Context, prefix string) error {
	files, err := q.ListByPrefix(ctx, prefix)
	if err != nil {
		return err
	}

	if len(files) == 0 {
		return nil
	}

	// 批量删除
	bucketManager := storage.NewBucketManager(q.mac, q.cfg)
	deleteOps := make([]string, 0, len(files))
	for _, file := range files {
		deleteOps = append(deleteOps, storage.URIDelete(q.bucket, file.Key))
	}

	// 分批处理，每批最多 1000 个
	batchSize := 1000
	for i := 0; i < len(deleteOps); i += batchSize {
		end := i + batchSize
		if end > len(deleteOps) {
			end = len(deleteOps)
		}
		batch := deleteOps[i:end]
		_, err := bucketManager.Batch(batch)
		if err != nil {
			return fmt.Errorf("qiniu batch delete failed: %w", err)
		}
	}

	return nil
}

// ListByPrefix 按前缀列出文件
func (q *QiniuStorage) ListByPrefix(ctx context.Context, prefix string) ([]FileInfo, error) {
	bucketManager := storage.NewBucketManager(q.mac, q.cfg)

	var files []FileInfo
	marker := ""
	limit := 1000

	for {
		entries, _, nextMarker, hasNext, err := bucketManager.ListFiles(q.bucket, prefix, "", marker, limit)
		if err != nil {
			return nil, fmt.Errorf("qiniu list files failed: %w", err)
		}

		for _, entry := range entries {
			files = append(files, FileInfo{
				Key:          entry.Key,
				Size:         entry.Fsize,
				ContentType:  entry.MimeType,
				LastModified: time.Unix(0, entry.PutTime*100), // 七牛时间戳是 100 纳秒
				ETag:         entry.Hash,
			})
		}

		if !hasNext {
			break
		}
		marker = nextMarker
	}

	return files, nil
}

// GetFileInfo 获取文件信息
func (q *QiniuStorage) GetFileInfo(ctx context.Context, key string) (*FileInfo, error) {
	bucketManager := storage.NewBucketManager(q.mac, q.cfg)
	fileInfo, err := bucketManager.Stat(q.bucket, key)
	if err != nil {
		return nil, fmt.Errorf("qiniu stat failed: %w", err)
	}

	return &FileInfo{
		Key:          key,
		Size:         fileInfo.Fsize,
		ContentType:  fileInfo.MimeType,
		LastModified: time.Unix(0, fileInfo.PutTime*100),
		ETag:         fileInfo.Hash,
	}, nil
}

// Exists 检查文件是否存在
func (q *QiniuStorage) Exists(ctx context.Context, key string) (bool, error) {
	_, err := q.GetFileInfo(ctx, key)
	if err != nil {
		// 检查是否是 "no such file or directory" 错误
		if strings.Contains(err.Error(), "no such file") {
			return false, nil
		}
		return false, err
	}
	return true, nil
}

// GetPublicURL 获取公开访问 URL
func (q *QiniuStorage) GetPublicURL(key string) string {
	scheme := "http"
	if q.useHTTPS {
		scheme = "https"
	}
	return fmt.Sprintf("%s://%s/%s", scheme, q.domain, key)
}

// GetSignedURL 获取签名 URL（私有空间）
func (q *QiniuStorage) GetSignedURL(ctx context.Context, key string, expiry time.Duration) (string, error) {
	deadline := time.Now().Add(expiry).Unix()
	publicURL := q.GetPublicURL(key)
	privateURL := storage.MakePrivateURL(q.mac, q.domain, key, deadline)
	if privateURL == "" {
		return publicURL, nil
	}
	return privateURL, nil
}

// Health 健康检查
func (q *QiniuStorage) Health(ctx context.Context) error {
	bucketManager := storage.NewBucketManager(q.mac, q.cfg)
	// 尝试列出一个文件来验证连接
	_, _, _, _, err := bucketManager.ListFiles(q.bucket, "", "", "", 1)
	if err != nil {
		return fmt.Errorf("qiniu health check failed: %w", err)
	}
	return nil
}

// 确保 QiniuStorage 实现了 Storage 接口
var _ Storage = (*QiniuStorage)(nil)
