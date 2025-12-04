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
}

var RouterGroupApp = new(RouterGroup)
