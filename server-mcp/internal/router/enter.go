package router

type RouterGroup struct {
	BaseRouter
	LibraryRouter
	DocumentRouter
	SearchRouter
	MCPRouter
	AuthRouter
	UserRouter
	ApiKeyRouter
	ActivityLogRouter
	StatsRouter
}

var RouterGroupApp = new(RouterGroup)
