package api

import "go-mcp-context/internal/service"

type ApiGroup struct {
	LibraryApi
	DocumentApi
	SearchApi
	MCPApi
	AuthApi
	UserApi
	ApiKeyApi
	ActivityLogApi
	StatsApi
}

var ApiGroupApp = new(ApiGroup)

var libraryService = service.ServiceGroupApp.LibraryService
var documentService = service.ServiceGroupApp.DocumentService
var searchService = service.ServiceGroupApp.SearchService
var mcpService = service.ServiceGroupApp.MCPService
var apiKeyService = service.ServiceGroupApp.ApiKeyService
var activityLogService = service.ServiceGroupApp.ActivityLogService
var statsService = service.ServiceGroupApp.StatsService
