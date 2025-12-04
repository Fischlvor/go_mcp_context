package initialize

import (
	"go-mcp-context/internal/middleware"
	"go-mcp-context/internal/router"
	"go-mcp-context/pkg/global"

	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/cookie"
	"github.com/gin-gonic/gin"
)

// InitRouter 初始化路由
func InitRouter() *gin.Engine {
	r := gin.New()
	r.Use(middleware.GinLogger(), middleware.GinRecovery(true))

	// 初始化 Session 中间件（用于存储 refresh_token）
	store := cookie.NewStore([]byte(global.Config.SSO.SessionsSecret))
	store.Options(sessions.Options{
		Path:     "/",
		MaxAge:   7 * 24 * 3600, // 7 天
		HttpOnly: true,
		Secure:   global.Config.System.Env == "release",
	})
	r.Use(sessions.Sessions("mcp_session", store))

	routerGroup := router.RouterGroupApp

	// API 公开路由（无需认证）- 健康检查
	publicRouter := r.Group("/api")
	{
		routerGroup.InitBaseRouter(publicRouter) // 健康检查等
	}

	// API v1 公开路由（无需认证）- 查询类接口 + 认证接口
	v1Public := r.Group("/api/v1")
	{
		routerGroup.InitAuthRouter(v1Public)           // 认证相关（SSO 登录、回调、登出）
		routerGroup.InitLibraryPublicRouter(v1Public)  // GET 库列表、详情
		routerGroup.InitDocumentPublicRouter(v1Public) // GET 文档详情
		routerGroup.InitSearchPublicRouter(v1Public)   // POST 搜索
	}

	// API v1 私有路由（需要 SSO JWT 认证）- 增删改接口
	v1Private := r.Group("/api/v1")
	v1Private.Use(middleware.SSOJWTAuth())
	{
		routerGroup.InitUserRouter(v1Private)     // 用户相关（获取用户信息）
		routerGroup.InitLibraryRouter(v1Private)  // POST/PUT/DELETE 库
		routerGroup.InitDocumentRouter(v1Private) // POST/DELETE 文档
		routerGroup.InitApiKeyRouter(v1Private)   // API Key 管理（CRUD）
	}

	// MCP routes（需要 API Key 认证）- IDE 调用
	mcp := r.Group("")
	mcp.Use(middleware.APIKeyAuth())
	{
		routerGroup.InitMCPRouter(mcp)
	}

	return r
}
