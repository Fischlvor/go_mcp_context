package main

import (
	"go-mcp-context/internal/initialize"
	"go-mcp-context/internal/middleware"
	"go-mcp-context/pkg/core"
	"go-mcp-context/pkg/global"
	"go-mcp-context/scripts/flag"

	"go.uber.org/zap"
)

func main() {
	global.Config = core.InitConf()
	global.Log = core.InitLogger()

	global.DB = initialize.InitGorm()
	global.Redis = initialize.ConnectRedis()
	global.Embedding = initialize.InitEmbedding()

	// 加载 SSO 公钥
	if err := middleware.LoadSSOPublicKey(global.Config.SSO.PublicKeyPath); err != nil {
		global.Log.Error("加载 SSO 公钥失败", zap.Error(err))
	}

	defer global.Redis.Close()

	flag.InitFlag()

	core.RunServer()
}
