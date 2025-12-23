package service

type ServiceGroup struct {
	LibraryService
	DocumentService
	SearchService
	MCPService
	ApiKeyService
	ActivityLogService
	StatsService
}

var ServiceGroupApp = new(ServiceGroup)
