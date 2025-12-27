// @title go-mcp-context API
// @version 1.0.0
// @description 私有化文档检索服务 - 为企业内网的 AI IDE 提供实时、准确的技术文档和代码示例
// @termsOfService http://swagger.io/terms/
// @contact.name API Support
// @contact.url http://www.swagger.io/support
// @license.name Apache 2.0
// @license.url http://www.apache.org/licenses/LICENSE-2.0.html
// @host 10.21.71.19:8090
// @basePath /
// @schemes http https
// @securityDefinitions.apikey MCP_API_KEY
// @in header
// @name MCP_API_KEY
// @description API Key for MCP protocol calls
// @securityDefinitions.apikey JWTAuth
// @in header
// @name Authorization
// @description JWT token in the format "Bearer {token}"
package main

import (
	"net/http"
	"os"

	_ "go-mcp-context/docs"
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

	initialize.InitBufferedWriters()        // 初始化缓冲写入器（活动日志、统计、MCP日志）
	defer initialize.CloseBufferedWriters() // 关闭时刷新缓冲区

	global.Redis = initialize.ConnectRedis()
	defer global.Redis.Close()

	global.Cache = initialize.InitCache() // 初始化通用缓存服务
	global.Embedding = initialize.InitEmbedding()
	initialize.InitStorage() // 初始化存储服务
	initialize.InitLLM()     // 初始化 LLM 服务

	// 加载 SSO 公钥
	if err := middleware.LoadSSOPublicKey(global.Config.SSO.PublicKeyPath); err != nil {
		global.Log.Error("加载 SSO 公钥失败", zap.Error(err))
	}

	flag.InitFlag()

	core.RunServer()
}
