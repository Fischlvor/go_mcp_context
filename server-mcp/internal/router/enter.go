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
}

var RouterGroupApp = new(RouterGroup)
