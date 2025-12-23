package service

type ServiceGroup struct {
	LibraryService
	DocumentService
	SearchService
	MCPService
	ApiKeyService
	ActivityLogService
}

var ServiceGroupApp = new(ServiceGroup)
