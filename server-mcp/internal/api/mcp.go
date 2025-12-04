package api

import (
	"net/http"

	"go-mcp-context/internal/model/request"
	"go-mcp-context/internal/model/response"

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
					"libraryID": map[string]interface{}{
						"type":        "string",
						"description": "The ID of the library (format: name/version)",
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
				"required": []string{"libraryID", "topic"},
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

// handleToolCall 处理工具调用
func handleToolCall(c *gin.Context, req request.MCPRequest) {
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

	switch toolName {
	case "search-libraries":
		handleSearchLibraries(c, req.ID, arguments)
	case "get-library-docs":
		handleGetLibraryDocs(c, req.ID, arguments)
	default:
		c.JSON(http.StatusOK, response.MCPResponse{
			JSONRPC: "2.0",
			ID:      req.ID,
			Error: &response.MCPError{
				Code:    -32602,
				Message: "Unknown tool: " + toolName,
			},
		})
	}
}

// handleSearchLibraries 处理 search-libraries 工具
func handleSearchLibraries(c *gin.Context, id interface{}, args map[string]interface{}) {
	libraryName, _ := args["libraryName"].(string)
	if libraryName == "" {
		c.JSON(http.StatusOK, response.MCPResponse{
			JSONRPC: "2.0",
			ID:      id,
			Error: &response.MCPError{
				Code:    -32602,
				Message: "Invalid params: libraryName is required",
			},
		})
		return
	}

	// 调用 service 层
	req := &request.MCPSearchLibraries{LibraryName: libraryName}
	result, err := mcpService.SearchLibraries(req)
	if err != nil {
		c.JSON(http.StatusOK, response.MCPResponse{
			JSONRPC: "2.0",
			ID:      id,
			Error: &response.MCPError{
				Code:    -32603,
				Message: "Internal error: " + err.Error(),
			},
		})
		return
	}

	c.JSON(http.StatusOK, response.MCPResponse{
		JSONRPC: "2.0",
		ID:      id,
		Result:  result,
	})
}

// handleGetLibraryDocs 处理 get-library-docs 工具
func handleGetLibraryDocs(c *gin.Context, id interface{}, args map[string]interface{}) {
	libraryID, _ := args["libraryID"].(string)
	topic, _ := args["topic"].(string)
	mode, _ := args["mode"].(string)
	page := 1
	if p, ok := args["page"].(float64); ok {
		page = int(p)
	}

	if libraryID == "" || topic == "" {
		c.JSON(http.StatusOK, response.MCPResponse{
			JSONRPC: "2.0",
			ID:      id,
			Error: &response.MCPError{
				Code:    -32602,
				Message: "Invalid params: libraryID and topic are required",
			},
		})
		return
	}

	// 调用 service 层
	req := &request.MCPGetLibraryDocs{
		LibraryID: libraryID,
		Topic:     topic,
		Mode:      mode,
		Page:      page,
	}
	result, err := mcpService.GetLibraryDocs(req)
	if err != nil {
		c.JSON(http.StatusOK, response.MCPResponse{
			JSONRPC: "2.0",
			ID:      id,
			Error: &response.MCPError{
				Code:    -32603,
				Message: "Internal error: " + err.Error(),
			},
		})
		return
	}

	c.JSON(http.StatusOK, response.MCPResponse{
		JSONRPC: "2.0",
		ID:      id,
		Result:  result,
	})
}
