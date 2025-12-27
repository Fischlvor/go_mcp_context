package service

import (
	"encoding/json"
	"fmt"
	"go-mcp-context/internal/model/request"
	"go-mcp-context/internal/model/response"
	"go-mcp-context/internal/transport"
	"go-mcp-context/pkg/global"
	"net/url"
	"strings"

	"go.uber.org/zap"
)

// MCPHandler MCP统一处理器
// 负责协议无关的业务逻辑处理，调用MCPService执行具体业务
type MCPHandler struct {
	mcpService *MCPService
}

// NewMCPHandler 创建MCP处理器
func NewMCPHandler() *MCPHandler {
	return &MCPHandler{
		mcpService: NewMCPService(),
	}
}

// ProcessRequest 统一处理MCP请求
// 根据method分发到对应的处理方法，通过ResponseWriter返回响应
func (h *MCPHandler) ProcessRequest(req *transport.RequestContext, writer transport.ResponseWriter) error {
	global.Log.Info("处理MCP请求",
		zap.String("transport", string(req.Transport)),
		zap.String("method", req.Method),
	)

	switch req.Method {
	case "initialize":
		return h.handleInitialize(req, writer)

	case "notifications/initialized", "initialized":
		return h.handleInitialized(req, writer)

	case "tools/list":
		return h.handleToolsList(req, writer)

	case "tools/call":
		return h.handleToolsCall(req, writer)

	case "resources/list":
		return h.handleResourcesList(req, writer)

	case "resources/templates/list":
		return h.handleResourceTemplatesList(req, writer)

	case "resources/read":
		return h.handleResourcesRead(req, writer)

	default:
		// 未知方法
		return writer.WriteError(&response.MCPError{
			Code:    -32601,
			Message: "Method not found: " + req.Method,
		}, req.ID)
	}
}

// handleInitialize 处理initialize请求
func (h *MCPHandler) handleInitialize(req *transport.RequestContext, writer transport.ResponseWriter) error {
	result := map[string]interface{}{
		"protocolVersion": "2025-11-25",
		"capabilities": map[string]interface{}{
			"tools": map[string]interface{}{
				"listChanged": true,
			},
			"resources": map[string]interface{}{
				"subscribe":   true,
				"listChanged": true,
			},
			"logging": map[string]interface{}{},
		},
		"serverInfo": map[string]interface{}{
			"name":    "go-mcp-context",
			"version": "1.0.0",
		},
	}

	// 统计结果数量（initialize返回的是capabilities和serverInfo，计为1）
	req.GinCtx.Set("mcp_result_count", 1)

	resp := &response.MCPResponse{
		JSONRPC: "2.0",
		ID:      req.ID,
		Result:  result,
	}

	return writer.WriteResponse(resp)
}

// handleInitialized 处理initialized通知
func (h *MCPHandler) handleInitialized(req *transport.RequestContext, writer transport.ResponseWriter) error {
	// 通知不需要响应
	global.Log.Info("收到MCP初始化完成通知")
	return nil
}

// handleToolsList 处理tools/list请求
func (h *MCPHandler) handleToolsList(req *transport.RequestContext, writer transport.ResponseWriter) error {
	tools := []map[string]interface{}{
		{
			"name":        "search-libraries",
			"description": "Search for documentation libraries by name. Returns matching libraries with metadata including available versions. Use this method to discover libraries and get their version information (versions array and defaultVersion) before calling get-library-docs.",
			"inputSchema": map[string]interface{}{
				"type": "object",
				"properties": map[string]interface{}{
					"libraryName": map[string]interface{}{
						"type":        "string",
						"description": "The name of the library to search for",
					},
				},
				"required": []string{"libraryName"},
			},
		},
		{
			"name":        "get-library-docs",
			"description": "Get documentation for a specific library. Requires libraryId, topic, and version. Supports comma-separated topics for multi-topic search.",
			"inputSchema": map[string]interface{}{
				"type": "object",
				"properties": map[string]interface{}{
					"libraryId": map[string]interface{}{
						"type":        "integer",
						"description": "Library ID (required)",
					},
					"topic": map[string]interface{}{
						"type":        "string",
						"description": "Documentation topic (required, supports comma-separated topics like 'overview,api,examples')",
					},
					"version": map[string]interface{}{
						"type":        "string",
						"description": "Library version (required)",
					},
					"mode": map[string]interface{}{
						"type":        "string",
						"description": "Search mode: info or code",
						"enum":        []string{"info", "code"},
					},
					"page": map[string]interface{}{
						"type":        "integer",
						"description": "Page number (1-10)",
					},
				},
				"required": []string{"libraryId", "topic", "version"},
			},
		},
	}

	// 统计结果数量（返回2个工具）
	req.GinCtx.Set("mcp_result_count", len(tools))

	resp := &response.MCPResponse{
		JSONRPC: "2.0",
		ID:      req.ID,
		Result:  map[string]interface{}{"tools": tools},
	}

	return writer.WriteResponse(resp)
}

// handleToolsCall 处理tools/call请求
func (h *MCPHandler) handleToolsCall(req *transport.RequestContext, writer transport.ResponseWriter) error {
	// 提取工具名称
	toolName, ok := req.Params["name"].(string)
	if !ok {
		return writer.WriteError(&response.MCPError{
			Code:    -32602,
			Message: "Invalid params: missing tool name",
		}, req.ID)
	}

	// 提取工具参数
	arguments, _ := req.Params["arguments"].(map[string]interface{})

	// 根据工具名称分发
	switch toolName {
	case "search-libraries":
		return h.handleSearchLibraries(arguments, req, writer)

	case "get-library-docs":
		return h.handleGetLibraryDocs(arguments, req, writer)

	default:
		return writer.WriteError(&response.MCPError{
			Code:    -32602,
			Message: "Unknown tool: " + toolName,
		}, req.ID)
	}
}

// handleSearchLibraries 处理search-libraries工具调用
func (h *MCPHandler) handleSearchLibraries(args map[string]interface{}, req *transport.RequestContext, writer transport.ResponseWriter) error {
	// 提取参数
	libraryName, _ := args["libraryName"].(string)
	if libraryName == "" {
		return writer.WriteError(&response.MCPError{
			Code:    -32602,
			Message: "Invalid params: libraryName is required",
		}, req.ID)
	}

	// 调用service层
	searchReq := &request.MCPSearchLibraries{LibraryName: libraryName}
	result, err := h.mcpService.SearchLibraries(searchReq)
	if err != nil {
		return writer.WriteError(&response.MCPError{
			Code:    -32603,
			Message: "Internal error: " + err.Error(),
		}, req.ID)
	}

	// 设置结果信息到context，供中间件记录日志
	req.GinCtx.Set("mcp_result_count", len(result.Libraries))

	// 转换为MCP规范的格式
	// MCP期望的格式: { content: [{ type: "text", text: "..." }] }
	resultJSON, _ := json.Marshal(result)
	mcpResult := map[string]interface{}{
		"content": []map[string]interface{}{
			{
				"type": "text",
				"text": string(resultJSON),
			},
		},
	}

	resp := &response.MCPResponse{
		JSONRPC: "2.0",
		ID:      req.ID,
		Result:  mcpResult,
	}

	return writer.WriteResponse(resp)
}

// handleGetLibraryDocs 处理get-library-docs工具调用
func (h *MCPHandler) handleGetLibraryDocs(args map[string]interface{}, req *transport.RequestContext, writer transport.ResponseWriter) error {
	// 提取参数
	var libraryID uint
	if id, ok := args["libraryId"].(float64); ok {
		libraryID = uint(id)
	}
	version, _ := args["version"].(string)
	topic, _ := args["topic"].(string)
	mode, _ := args["mode"].(string)
	page := 1
	if p, ok := args["page"].(float64); ok {
		page = int(p)
	}

	// 参数验证
	if topic == "" {
		return writer.WriteError(&response.MCPError{
			Code:    -32602,
			Message: "Invalid params: topic is required",
		}, req.ID)
	}

	// 调用service层
	docsReq := &request.MCPGetLibraryDocs{
		LibraryID: libraryID,
		Topic:     topic,
		Version:   version,
		Mode:      mode,
		Page:      page,
	}
	result, err := h.mcpService.GetLibraryDocs(docsReq)
	if err != nil {
		return writer.WriteError(&response.MCPError{
			Code:    -32603,
			Message: "Internal error: " + err.Error(),
		}, req.ID)
	}

	// 设置结果信息到context，供中间件记录日志
	req.GinCtx.Set("mcp_result_count", len(result.Documents))
	if result.LibraryID > 0 {
		req.GinCtx.Set("mcp_library_id", result.LibraryID)
	}

	// 转换为MCP规范的格式
	// MCP期望的格式: { content: [{ type: "text", text: "..." }] }
	resultJSON, _ := json.Marshal(result)
	mcpResult := map[string]interface{}{
		"content": []map[string]interface{}{
			{
				"type": "text",
				"text": string(resultJSON),
			},
		},
	}

	resp := &response.MCPResponse{
		JSONRPC: "2.0",
		ID:      req.ID,
		Result:  mcpResult,
	}

	return writer.WriteResponse(resp)
}

// handleResourcesList 处理resources/list请求
// 从数据库查询所有可用的库，动态生成资源列表
func (h *MCPHandler) handleResourcesList(req *transport.RequestContext, writer transport.ResponseWriter) error {
	// 从数据库查询所有库
	libraries, err := h.mcpService.GetAllLibraries()
	if err != nil {
		return writer.WriteError(&response.MCPError{
			Code:    -32603,
			Message: "Internal error: " + err.Error(),
		}, req.ID)
	}

	// 动态生成资源列表（为每个库添加一个资源）
	resources := make([]map[string]interface{}, 0, len(libraries))
	for _, lib := range libraries {
		resources = append(resources, map[string]interface{}{
			"uri":         fmt.Sprintf("go-mcp-context:///library/%d", lib.ID),
			"name":        lib.Name,
			"description": lib.Description,
			"mimeType":    "application/json",
		})
	}

	// 统计结果数量
	req.GinCtx.Set("mcp_result_count", len(resources))

	resp := &response.MCPResponse{
		JSONRPC: "2.0",
		ID:      req.ID,
		Result:  map[string]interface{}{"resources": resources},
	}

	return writer.WriteResponse(resp)
}

// handleResourceTemplatesList 处理resources/templates/list请求
// 返回资源模板列表，用于客户端动态生成资源URI
func (h *MCPHandler) handleResourceTemplatesList(req *transport.RequestContext, writer transport.ResponseWriter) error {
	// 资源模板（客户端可以填充参数生成具体的资源URI）
	templates := []map[string]interface{}{
		{
			"uriTemplate": "go-mcp-context:///library/{libraryId}",
			"name":        "Library by ID",
			"description": "Get library information by ID",
			"mimeType":    "application/json",
		},
		{
			"uriTemplate": "go-mcp-context:///docs/chunk/{libraryId}/{version}/{topic}",
			"name":        "Documentation Chunk",
			"description": "Get documentation chunk for a specific library version and topic (supports comma-separated topics like overview,api,examples)",
			"mimeType":    "text/markdown",
		},
	}

	// 统计结果数量
	req.GinCtx.Set("mcp_result_count", len(templates))
	resp := &response.MCPResponse{
		JSONRPC: "2.0",
		ID:      req.ID,
		Result:  map[string]interface{}{"resourceTemplates": templates},
	}

	return writer.WriteResponse(resp)
}

// handleResourcesRead 处理resources/read请求
// 根据URI读取对应的资源内容
func (h *MCPHandler) handleResourcesRead(req *transport.RequestContext, writer transport.ResponseWriter) error {
	// 解析params中的uri参数
	uri, ok := req.Params["uri"].(string)
	if !ok || uri == "" {
		return writer.WriteError(&response.MCPError{
			Code:    -32602,
			Message: "Missing or invalid uri parameter",
		}, req.ID)
	}

	// 支持两种URI格式（使用三个斜杠表示无host）：
	// 1. go-mcp-context:///library/{libraryId} - 获取库的基本信息
	// 2. go-mcp-context:///docs/chunk/{libraryId}/{version}/{topic} - 获取库的文档块
	//    topic支持逗号分隔的多个topic，如 overview,api,examples

	// 解析URI，提取path部分
	parsedURL, err := url.Parse(uri)
	if err != nil {
		return writer.WriteError(&response.MCPError{
			Code:    -32602,
			Message: "Invalid uri format: " + uri,
		}, req.ID)
	}

	// 按 / 分割路径
	pathSegments := strings.Split(strings.TrimPrefix(parsedURL.Path, "/"), "/")

	// 过滤空的路径段
	var segments []string
	for _, seg := range pathSegments {
		if seg != "" {
			segments = append(segments, seg)
		}
	}

	// 根据第一个路径段分发到对应的处理器
	if len(segments) == 0 {
		return writer.WriteError(&response.MCPError{
			Code:    -32602,
			Message: "Invalid uri format: " + uri,
		}, req.ID)
	}

	switch segments[0] {
	case "library":
		return h.handleReadLibraryResource(segments, uri, req, writer)
	case "docs":
		return h.handleReadDocsResource(segments, uri, req, writer)
	default:
		return writer.WriteError(&response.MCPError{
			Code:    -32602,
			Message: "Unknown resource type: " + segments[0],
		}, req.ID)
	}
}

// handleReadLibraryResource 处理 library/{libraryId} 资源请求
func (h *MCPHandler) handleReadLibraryResource(segments []string, uri string, req *transport.RequestContext, writer transport.ResponseWriter) error {
	if len(segments) < 2 {
		return writer.WriteError(&response.MCPError{
			Code:    -32602,
			Message: "Invalid library uri format: " + uri,
		}, req.ID)
	}

	// 校验libraryId参数
	if segments[1] == "" {
		return writer.WriteError(&response.MCPError{
			Code:    -32602,
			Message: "libraryId cannot be empty",
		}, req.ID)
	}

	var libraryID uint
	_, err := fmt.Sscanf(segments[1], "%d", &libraryID)
	if err != nil {
		return writer.WriteError(&response.MCPError{
			Code:    -32602,
			Message: "Invalid libraryId format: must be a positive integer",
		}, req.ID)
	}

	if libraryID == 0 {
		return writer.WriteError(&response.MCPError{
			Code:    -32602,
			Message: "libraryId must be greater than 0",
		}, req.ID)
	}

	// 从数据库获取库信息
	library, err := h.mcpService.GetLibraryByID(libraryID)
	if err != nil {
		return writer.WriteError(&response.MCPError{
			Code:    -32603,
			Message: "Internal error: " + err.Error(),
		}, req.ID)
	}

	if library == nil {
		return writer.WriteError(&response.MCPError{
			Code:    -32602,
			Message: "Library not found: " + uri,
		}, req.ID)
	}

	// 返回库的基本信息
	result := map[string]interface{}{
		"contents": []map[string]interface{}{
			{
				"uri":      uri,
				"mimeType": "application/json",
				"text": map[string]interface{}{
					"id":          library.ID,
					"name":        library.Name,
					"description": library.Description,
					"sourceURL":   library.SourceURL,
					"status":      library.Status,
					"createdAt":   library.CreatedAt,
					"updatedAt":   library.UpdatedAt,
				},
			},
		},
	}

	// 统计结果数量（返回1个资源）
	req.GinCtx.Set("mcp_result_count", 1)

	resp := &response.MCPResponse{
		JSONRPC: "2.0",
		ID:      req.ID,
		Result:  result,
	}

	return writer.WriteResponse(resp)
}

// handleReadDocsResource 处理 docs/chunk/{libraryId}/{version}/{topic} 资源请求
func (h *MCPHandler) handleReadDocsResource(segments []string, uri string, req *transport.RequestContext, writer transport.ResponseWriter) error {
	if len(segments) < 5 || segments[1] != "chunk" {
		return writer.WriteError(&response.MCPError{
			Code:    -32602,
			Message: "Invalid docs uri format: " + uri,
		}, req.ID)
	}

	// 校验libraryId参数
	if segments[2] == "" {
		return writer.WriteError(&response.MCPError{
			Code:    -32602,
			Message: "libraryId cannot be empty",
		}, req.ID)
	}

	var libraryID uint
	_, err := fmt.Sscanf(segments[2], "%d", &libraryID)
	if err != nil {
		return writer.WriteError(&response.MCPError{
			Code:    -32602,
			Message: "Invalid libraryId format: must be a positive integer",
		}, req.ID)
	}

	if libraryID == 0 {
		return writer.WriteError(&response.MCPError{
			Code:    -32602,
			Message: "libraryId must be greater than 0",
		}, req.ID)
	}

	// 校验version参数
	version := segments[3]
	if version == "" {
		return writer.WriteError(&response.MCPError{
			Code:    -32602,
			Message: "version cannot be empty",
		}, req.ID)
	}

	// 校验topic参数
	// topic支持逗号分隔的多个topic，如 overview,api,examples
	topic := segments[4]
	if topic == "" {
		return writer.WriteError(&response.MCPError{
			Code:    -32602,
			Message: "topic cannot be empty",
		}, req.ID)
	}

	// 调用GetLibraryDocs逻辑获取文档数据
	// 因为GetLibraryDocs返回的是工具调用格式，我们需要转换为资源读取格式
	// topic支持逗号分隔的多个topic，GetLibraryDocs会自动处理
	tempReq := &request.MCPGetLibraryDocs{
		LibraryID: libraryID,
		Topic:     topic,
		Version:   version,
	}

	result, err := h.mcpService.GetLibraryDocs(tempReq)
	if err != nil {
		return writer.WriteError(&response.MCPError{
			Code:    -32603,
			Message: "Internal error: " + err.Error(),
		}, req.ID)
	}

	// 将GetLibraryDocs的结果转换为资源读取格式
	// 构建contents数组，包含所有文档
	contents := make([]map[string]interface{}, 0)

	if result != nil && len(result.Documents) > 0 {
		for _, doc := range result.Documents {
			contents = append(contents, map[string]interface{}{
				"uri":      uri,
				"mimeType": "text/markdown",
				"text":     doc.Content, // 假设Document有Content字段
			})
		}
	}

	// 如果没有文档，返回空的contents数组
	if len(contents) == 0 {
		contents = []map[string]interface{}{}
	}

	// 统计结果数量
	req.GinCtx.Set("mcp_result_count", len(contents))

	resp := &response.MCPResponse{
		JSONRPC: "2.0",
		ID:      req.ID,
		Result: map[string]interface{}{
			"contents": contents,
		},
	}

	return writer.WriteResponse(resp)
}
