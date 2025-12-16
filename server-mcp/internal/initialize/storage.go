package initialize

import (
	"go-mcp-context/pkg/global"
	"go-mcp-context/pkg/storage"
)

// InitStorage 初始化存储服务
func InitStorage() {
	storageType := global.Config.System.StorageType
	if storageType == "" {
		storageType = "local" // 默认使用本地存储
	}

	switch storageType {
	case "qiniu":
		cfg := global.Config.Qiniu
		if cfg.AccessKey == "" || cfg.SecretKey == "" {
			global.Log.Warn("Qiniu storage not configured, falling back to local storage")
			global.Storage = storage.NewLocalStorage()
			return
		}

		global.Storage = storage.NewQiniuStorage(storage.QiniuConfig{
			AccessKey:     cfg.AccessKey,
			SecretKey:     cfg.SecretKey,
			Bucket:        cfg.Bucket,
			Domain:        cfg.Domain,
			Zone:          cfg.Zone,
			UseHTTPS:      cfg.UseHTTPS,
			UseCdnDomains: cfg.UseCdnDomains,
		})

		global.Log.Info("Qiniu storage initialized successfully")

	default: // "local"
		global.Storage = storage.NewLocalStorage()
		global.Log.Info("Local storage initialized successfully")
	}
}
