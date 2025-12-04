package service

type ServiceGroup struct {
	LibraryService
	DocumentService
	SearchService
	MCPService
	ApiKeyService
}

var ServiceGroupApp = new(ServiceGroup)
