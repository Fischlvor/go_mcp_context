package main

import (
	"net/http"
	"os"

	"go-mcp-context/internal/initialize"
	"go-mcp-context/internal/middleware"
	"go-mcp-context/pkg/core"
	"go-mcp-context/pkg/global"
	"go-mcp-context/scripts/flag"

	"go.uber.org/zap"
)

func init() {
	// 禁用系统代理
	os.Unsetenv("HTTP_PROXY")
	os.Unsetenv("HTTPS_PROXY")
	os.Unsetenv("http_proxy")
	os.Unsetenv("https_proxy")
	os.Unsetenv("ALL_PROXY")
	os.Unsetenv("all_proxy")

	// 禁用 Go 默认 HTTP 客户端的代理
	http.DefaultTransport.(*http.Transport).Proxy = nil
}

func main() {
	global.Config = core.InitConf()
	global.Log = core.InitLogger()

	global.DB = initialize.InitGorm()
	global.Redis = initialize.ConnectRedis()
	global.Embedding = initialize.InitEmbedding()
	initialize.InitStorage() // 初始化存储服务
	initialize.InitLLM()     // 初始化 LLM 服务

	// 加载 SSO 公钥
	if err := middleware.LoadSSOPublicKey(global.Config.SSO.PublicKeyPath); err != nil {
		global.Log.Error("加载 SSO 公钥失败", zap.Error(err))
	}

	defer global.Redis.Close()

	flag.InitFlag()

	core.RunServer()
}
