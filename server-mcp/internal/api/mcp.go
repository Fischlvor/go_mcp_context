package api

import (
	"fmt"
	"net/http"

	"go-mcp-context/internal/model/request"
	"go-mcp-context/internal/model/response"
	"go-mcp-context/pkg/global"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type MCPApi struct{}

// Health 返回 MCP 服务健康状态 (已废弃，使用HandleRequest统一处理)
// func (m *MCPApi) Health(c *gin.Context) {
// 	c.JSON(http.StatusOK, gin.H{
// 		"status":  "ok",
// 		"version": "1.0.0",
// 	})
// }

// ListTools 返回可用的 MCP 工具列表 (已废弃，使用handleToolsList统一处理)
// func (m *MCPApi) ListTools(c *gin.Context) {
// 	tools := []response.MCPToolDefinition{
// 		{
// 			Name:        "search-libraries",
// 			Description: "Search for documentation libraries by name. Returns matching libraries with metadata.",
// 			InputSchema: map[string]interface{}{
// 				"type": "object",
// 				"properties": map[string]interface{}{
// 					"libraryName": map[string]interface{}{
// 						"type":        "string",
// 						"description": "The name of the library to search for",
// 					},
// 				},
// 				"required": []string{"libraryName"},
// 			},
// 		},
// 		{
// 			Name:        "get-library-docs",
// 			Description: "Get documentation from libraries. If libraryId is provided, search within that library; otherwise search across all libraries.",
// 			InputSchema: map[string]interface{}{
// 				"type": "object",
// 				"properties": map[string]interface{}{
// 					"libraryId": map[string]interface{}{
// 						"type":        "integer",
// 						"description": "The database ID of the library (optional, from search-libraries result). If not provided, search across all libraries.",
// 					},
// 					"version": map[string]interface{}{
// 						"type":        "string",
// 						"description": "The version of the library (optional, empty means all versions)",
// 					},
// 					"topic": map[string]interface{}{
// 						"type":        "string",
// 						"description": "The topic or query to search for. Supports multiple topics separated by comma (e.g. 'routing, middleware')",
// 					},
// 					"mode": map[string]interface{}{
// 						"type":        "string",
// 						"enum":        []string{"code", "info"},
// 						"description": "Filter by content type: 'code' for code examples, 'info' for documentation (optional, empty means all types)",
// 					},
// 					"page": map[string]interface{}{
// 						"type":        "integer",
// 						"description": "Page number (1-10)",
// 						"minimum":     1,
// 						"maximum":     10,
// 						"default":     1,
// 					},
// 				},
// 				"required": []string{"topic"},
// 			},
// 		},
// 	}
//
// 	c.JSON(http.StatusOK, gin.H{
// 		"tools": tools,
// 	})
// }

// HandleRequest 处理 MCP JSON-RPC 请求
func (m *MCPApi) HandleRequest(c *gin.Context) {
	// 解析JSON-RPC请求
	var req request.MCPRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, response.MCPResponse{
			JSONRPC: "2.0",
			ID:      nil,
			Error: &response.MCPError{
				Code:    -32700,
				Message: "Parse error",
			},
		})
		return
	}

	method := req.Method
	id := req.ID

	// 根据方法分发请求
	switch method {
	case "initialize":
		m.handleInitialize(c, req.Params, req.ID)
	case "notifications/initialized", "initialized":
		m.handleInitialized(c)
	case "tools/list":
		m.handleToolsList(c, id)
	case "tools/call":
		m.handleToolsCall(c, req.Params, req.ID)
	case "resources/list":
		m.handleResourcesList(c, id)
	case "resources/templates/list":
		m.handleResourceTemplatesList(c, id)
	case "search-libraries":
		m.handleSearchLibraries(c, req.Params, req.ID)
	case "get-library-docs":
		m.handleGetLibraryDocs(c, req.Params, req.ID)
	default:
		global.Log.Warn("未知的MCP方法", zap.String("method", method))
		c.JSON(404, gin.H{
			"jsonrpc": "2.0",
			"id":      id,
			"error": gin.H{
				"code":    -32601,
				"message": "Method not found",
				"data":    fmt.Sprintf("Unknown method: %s", method),
			},
		})
	}
}

// MCPToolResult 工具调用结果
type MCPToolResult struct {
	Result      interface{} // 成功时的结果
	ResultCount int         // 结果数量（用于日志）
	LibraryID   *uint       // 关联库 ID（用于日志）
	Error       *response.MCPError
}

// handleToolCall 处理工具调用（统一记录日志）(已废弃，使用handleSearchLibrariesFromMap和handleGetLibraryDocsFromMap)
// func handleToolCall(c *gin.Context, req request.MCPRequest) {
// 	startTime := time.Now()
// 	actorID := utils.GetUUID(c).String()
//
// 	toolName, ok := req.Params["name"].(string)
// 	if !ok {
// 		c.JSON(http.StatusOK, response.MCPResponse{
// 			JSONRPC: "2.0",
// 			ID:      req.ID,
// 			Error: &response.MCPError{
// 				Code:    -32602,
// 				Message: "Invalid params: missing tool name",
// 			},
// 		})
// 		return
// 	}
//
// 	arguments, _ := req.Params["arguments"].(map[string]interface{})
//
// 	// 调用具体工具，返回统一结果
// 	var toolResult MCPToolResult
// 	switch toolName {
// 	case "search-libraries":
// 		toolResult = doSearchLibraries(arguments)
// 	case "get-library-docs":
// 		toolResult = doGetLibraryDocs(arguments)
// 	default:
// 		c.JSON(http.StatusOK, response.MCPResponse{
// 			JSONRPC: "2.0",
// 			ID:      req.ID,
// 			Error: &response.MCPError{
// 				Code:    -32602,
// 				Message: "Unknown tool: " + toolName,
// 			},
// 		})
// 		return
// 	}
//
// 	latencyMs := int(time.Since(startTime).Milliseconds())
//
// 	// 统一记录日志
// 	logEntry := &mcplog.LogEntry{
// 		ActorID:     actorID,
// 		FuncName:    toolName,
// 		LibraryID:   toolResult.LibraryID,
// 		Params:      arguments, // 直接存储请求参数
// 		ResultCount: toolResult.ResultCount,
// 		LatencyMs:   latencyMs,
// 		Status:      "success",
// 	}
// 	if toolResult.Error != nil {
// 		logEntry.Status = "error"
// 		logEntry.ErrorMsg = toolResult.Error.Message
// 	}
// 	mcplog.Log(logEntry)
//
// 	// 统一返回响应
// 	if toolResult.Error != nil {
// 		c.JSON(http.StatusOK, response.MCPResponse{
// 			JSONRPC: "2.0",
// 			ID:      req.ID,
// 			Error:   toolResult.Error,
// 		})
// 		return
// 	}
//
// 	c.JSON(http.StatusOK, response.MCPResponse{
// 		JSONRPC: "2.0",
// 		ID:      req.ID,
// 		Result:  toolResult.Result,
// 	})
// }

// doSearchLibraries 执行 search-libraries 工具
func doSearchLibraries(args map[string]interface{}) MCPToolResult {
	libraryName, _ := args["libraryName"].(string)
	if libraryName == "" {
		return MCPToolResult{
			Error: &response.MCPError{
				Code:    -32602,
				Message: "Invalid params: libraryName is required",
			},
		}
	}

	// 调用 service 层
	req := &request.MCPSearchLibraries{LibraryName: libraryName}
	result, err := mcpService.SearchLibraries(req)
	if err != nil {
		return MCPToolResult{
			Error: &response.MCPError{
				Code:    -32603,
				Message: "Internal error: " + err.Error(),
			},
		}
	}

	return MCPToolResult{
		Result:      result,
		ResultCount: len(result.Libraries),
	}
}

// doGetLibraryDocs 执行 get-library-docs 工具
func doGetLibraryDocs(args map[string]interface{}) MCPToolResult {
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

	if topic == "" {
		return MCPToolResult{
			Error: &response.MCPError{
				Code:    -32602,
				Message: "Invalid params: topic is required",
			},
		}
	}

	// 验证 mode 值（允许为空，表示搜索所有类型）
	if mode != "" && mode != "code" && mode != "info" {
		return MCPToolResult{
			Error: &response.MCPError{
				Code:    -32602,
				Message: "Invalid params: mode must be 'code' or 'info'",
			},
		}
	}

	// 调用 service 层
	req := &request.MCPGetLibraryDocs{
		LibraryID: libraryID,
		Version:   version,
		Topic:     topic,
		Mode:      mode,
		Page:      page,
	}
	result, err := mcpService.GetLibraryDocs(req)
	if err != nil {
		return MCPToolResult{
			LibraryID: &libraryID,
			Error: &response.MCPError{
				Code:    -32603,
				Message: "Internal error: " + err.Error(),
			},
		}
	}

	return MCPToolResult{
		Result:      result,
		ResultCount: len(result.Documents),
		LibraryID:   &result.LibraryID,
	}
}

// handleInitialize 处理MCP初始化请求
func (m *MCPApi) handleInitialize(c *gin.Context, params map[string]interface{}, id interface{}) {
	global.Log.Info("处理MCP初始化请求", zap.Any("params", params))
	c.JSON(200, gin.H{
		"jsonrpc": "2.0",
		"id":      id,
		"result": gin.H{
			"protocolVersion": "2025-11-25",
			"capabilities": gin.H{
				"tools": gin.H{
					"listChanged": true,
				},
				"resources": gin.H{
					"subscribe":   true,
					"listChanged": true,
				},
				"logging": gin.H{},
			},
			"serverInfo": gin.H{
				"name":    "go-mcp-context",
				"version": "1.0.0",
			},
		},
	})
}

// handleInitialized 处理MCP初始化完成通知
func (m *MCPApi) handleInitialized(c *gin.Context) {
	global.Log.Info("收到MCP初始化完成通知")
	c.Status(200) // 通知不需要响应
}

// handleToolsList 处理工具列表请求
func (m *MCPApi) handleToolsList(c *gin.Context, id interface{}) {
	global.Log.Info("处理工具列表请求")

	tools := []response.MCPToolDefinition{
		{
			Name:        "search-libraries",
			Description: "Search for documentation libraries by name. Returns matching libraries with metadata.",
			InputSchema: map[string]interface{}{
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
			Name:        "get-library-docs",
			Description: "Get documentation from a specific library. Requires libraryId, version, and topic.",
			InputSchema: map[string]interface{}{
				"type": "object",
				"properties": map[string]interface{}{
					"libraryId": map[string]interface{}{
						"type":        "integer",
						"description": "The database ID of the library (required, from search-libraries result)",
					},
					"version": map[string]interface{}{
						"type":        "string",
						"description": "The version of the library (required)",
					},
					"topic": map[string]interface{}{
						"type":        "string",
						"description": "The topic to search for within the library (required)",
					},
					"mode": map[string]interface{}{
						"type":        "string",
						"enum":        []string{"code", "info"},
						"description": "The type of documentation to retrieve: 'code' for API reference with code examples, 'info' for conceptual guides",
					},
					"page": map[string]interface{}{
						"type":        "integer",
						"minimum":     1,
						"maximum":     10,
						"description": "Page number for pagination (1-10)",
					},
				},
				"required": []string{"libraryId", "version", "topic"},
			},
		},
	}

	c.JSON(200, gin.H{
		"jsonrpc": "2.0",
		"id":      id,
		"result": gin.H{
			"tools": tools,
		},
	})
}

// handleResourcesList 处理资源列表请求
func (m *MCPApi) handleResourcesList(c *gin.Context, id interface{}) {
	global.Log.Info("处理资源列表请求")
	c.JSON(200, gin.H{
		"jsonrpc": "2.0",
		"id":      id,
		"result": gin.H{
			"resources": []gin.H{}, // 暂时返回空列表
		},
	})
}

// handleResourceTemplatesList 处理模板列表请求
func (m *MCPApi) handleResourceTemplatesList(c *gin.Context, id interface{}) {
	global.Log.Info("处理模板列表请求")
	c.JSON(200, gin.H{
		"jsonrpc": "2.0",
		"id":      id,
		"result": gin.H{
			"resourceTemplates": []gin.H{}, // 暂时返回空列表
		},
	})
}

// handleToolsCall 处理工具调用请求
func (m *MCPApi) handleToolsCall(c *gin.Context, params map[string]interface{}, id interface{}) {
	// 获取工具名称
	toolName, ok := params["name"].(string)
	if !ok {
		c.JSON(http.StatusOK, response.MCPResponse{
			JSONRPC: "2.0",
			ID:      id,
			Error: &response.MCPError{
				Code:    -32602,
				Message: "Invalid params: missing tool name",
			},
		})
		return
	}

	// 获取工具参数
	arguments, _ := params["arguments"].(map[string]interface{})

	// 根据工具名称分发到具体的处理方法
	switch toolName {
	case "search-libraries":
		m.handleSearchLibraries(c, arguments, id)
	case "get-library-docs":
		m.handleGetLibraryDocs(c, arguments, id)
	default:
		c.JSON(http.StatusOK, response.MCPResponse{
			JSONRPC: "2.0",
			ID:      id,
			Error: &response.MCPError{
				Code:    -32602,
				Message: "Unknown tool: " + toolName,
			},
		})
	}
}

// handleSearchLibraries 处理搜索库请求
func (m *MCPApi) handleSearchLibraries(c *gin.Context, params map[string]interface{}, id interface{}) {
	// 调用工具函数
	result := doSearchLibraries(params)

	// 设置结果信息到context，供中间件记录日志
	c.Set("mcp_result_count", result.ResultCount)
	if result.LibraryID != nil {
		c.Set("mcp_library_id", *result.LibraryID)
	}

	// 返回响应
	if result.Error != nil {
		c.JSON(http.StatusOK, response.MCPResponse{
			JSONRPC: "2.0",
			ID:      id,
			Error:   result.Error,
		})
		return
	}

	c.JSON(http.StatusOK, response.MCPResponse{
		JSONRPC: "2.0",
		ID:      id,
		Result:  result.Result,
	})
}

// handleGetLibraryDocs 处理获取库文档请求
func (m *MCPApi) handleGetLibraryDocs(c *gin.Context, params map[string]interface{}, id interface{}) {
	// 调用工具函数
	result := doGetLibraryDocs(params)

	// 设置结果信息到context，供中间件记录日志
	c.Set("mcp_result_count", result.ResultCount)
	if result.LibraryID != nil {
		c.Set("mcp_library_id", *result.LibraryID)
	}

	// 返回响应
	if result.Error != nil {
		c.JSON(http.StatusOK, response.MCPResponse{
			JSONRPC: "2.0",
			ID:      id,
			Error:   result.Error,
		})
		return
	}

	c.JSON(http.StatusOK, response.MCPResponse{
		JSONRPC: "2.0",
		ID:      id,
		Result:  result.Result,
	})
}
