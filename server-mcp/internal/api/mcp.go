package api

import (
	"net/http"
	"time"

	"go-mcp-context/internal/model/request"
	"go-mcp-context/internal/model/response"
	"go-mcp-context/pkg/bufferedwriter/mcplog"
	"go-mcp-context/pkg/utils"

	"github.com/gin-gonic/gin"
)

type MCPApi struct{}

// Health 返回 MCP 服务健康状态
func (m *MCPApi) Health(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"status":  "ok",
		"version": "1.0.0",
	})
}

// ListTools 返回可用的 MCP 工具列表
func (m *MCPApi) ListTools(c *gin.Context) {
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
			Description: "Get documentation from a specific library. Returns relevant code snippets and documentation.",
			InputSchema: map[string]interface{}{
				"type": "object",
				"properties": map[string]interface{}{
					"libraryId": map[string]interface{}{
						"type":        "integer",
						"description": "The database ID of the library (from search-libraries result)",
					},
					"version": map[string]interface{}{
						"type":        "string",
						"description": "The version of the library (optional, defaults to defaultVersion)",
					},
					"topic": map[string]interface{}{
						"type":        "string",
						"description": "The topic or query to search for",
					},
					"mode": map[string]interface{}{
						"type":        "string",
						"enum":        []string{"code", "info"},
						"description": "Filter by content type: 'code' for code examples, 'info' for documentation",
					},
					"page": map[string]interface{}{
						"type":        "integer",
						"description": "Page number (1-10)",
						"minimum":     1,
						"maximum":     10,
						"default":     1,
					},
				},
				"required": []string{"libraryId", "topic"},
			},
		},
	}

	c.JSON(http.StatusOK, gin.H{
		"tools": tools,
	})
}

// HandleRequest 处理 MCP JSON-RPC 请求
func (m *MCPApi) HandleRequest(c *gin.Context) {
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

	// 验证 JSON-RPC 版本
	if req.JSONRPC != "2.0" {
		c.JSON(http.StatusBadRequest, response.MCPResponse{
			JSONRPC: "2.0",
			ID:      req.ID,
			Error: &response.MCPError{
				Code:    -32600,
				Message: "Invalid Request: jsonrpc must be 2.0",
			},
		})
		return
	}

	// 路由到对应处理器
	switch req.Method {
	case "tools/call":
		handleToolCall(c, req)
	case "tools/list":
		m.ListTools(c)
	default:
		c.JSON(http.StatusOK, response.MCPResponse{
			JSONRPC: "2.0",
			ID:      req.ID,
			Error: &response.MCPError{
				Code:    -32601,
				Message: "Method not found: " + req.Method,
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

// handleToolCall 处理工具调用（统一记录日志）
func handleToolCall(c *gin.Context, req request.MCPRequest) {
	startTime := time.Now()
	actorID := utils.GetUUID(c).String()

	toolName, ok := req.Params["name"].(string)
	if !ok {
		c.JSON(http.StatusOK, response.MCPResponse{
			JSONRPC: "2.0",
			ID:      req.ID,
			Error: &response.MCPError{
				Code:    -32602,
				Message: "Invalid params: missing tool name",
			},
		})
		return
	}

	arguments, _ := req.Params["arguments"].(map[string]interface{})

	// 调用具体工具，返回统一结果
	var toolResult MCPToolResult
	switch toolName {
	case "search-libraries":
		toolResult = doSearchLibraries(arguments)
	case "get-library-docs":
		toolResult = doGetLibraryDocs(arguments)
	default:
		c.JSON(http.StatusOK, response.MCPResponse{
			JSONRPC: "2.0",
			ID:      req.ID,
			Error: &response.MCPError{
				Code:    -32602,
				Message: "Unknown tool: " + toolName,
			},
		})
		return
	}

	latencyMs := int(time.Since(startTime).Milliseconds())

	// 统一记录日志
	logEntry := &mcplog.LogEntry{
		ActorID:     actorID,
		FuncName:    toolName,
		LibraryID:   toolResult.LibraryID,
		Params:      arguments, // 直接存储请求参数
		ResultCount: toolResult.ResultCount,
		LatencyMs:   latencyMs,
		Status:      "success",
	}
	if toolResult.Error != nil {
		logEntry.Status = "error"
		logEntry.ErrorMsg = toolResult.Error.Message
	}
	mcplog.Log(logEntry)

	// 统一返回响应
	if toolResult.Error != nil {
		c.JSON(http.StatusOK, response.MCPResponse{
			JSONRPC: "2.0",
			ID:      req.ID,
			Error:   toolResult.Error,
		})
		return
	}

	c.JSON(http.StatusOK, response.MCPResponse{
		JSONRPC: "2.0",
		ID:      req.ID,
		Result:  toolResult.Result,
	})
}

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

	if libraryID == 0 || topic == "" {
		return MCPToolResult{
			Error: &response.MCPError{
				Code:    -32602,
				Message: "Invalid params: libraryId and topic are required",
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
