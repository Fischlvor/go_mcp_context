package middleware

import (
	"bytes"
	"encoding/json"
	"io"
	"time"

	"go-mcp-context/internal/model/database"
	"go-mcp-context/pkg/bufferedwriter/mcplog"
	"go-mcp-context/pkg/utils"

	"github.com/gin-gonic/gin"
)

// methodToFuncName MCP方法名到函数名的映射
var methodToFuncName = map[string]string{
	"initialize":                database.MCPFuncInitialize,
	"notifications/initialized": database.MCPFuncInitialized,
	"initialized":               database.MCPFuncInitialized,
	"tools/list":                database.MCPFuncToolsList,
	"tools/call":                database.MCPFuncToolsCall,
	"resources/list":            database.MCPFuncResourcesList,
	"resources/templates/list":  database.MCPFuncResourceTemplatesList,
	"search-libraries":          database.MCPFuncSearchLibraries,
	"get-library-docs":          database.MCPFuncGetLibraryDocs,
}

// MCPLogMiddleware MCP调用日志中间件
func MCPLogMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 只对MCP相关路径记录日志
		if !isMCPPath(c.Request.URL.Path) {
			c.Next()
			return
		}

		startTime := time.Now()
		actorID := utils.GetUUID(c).String()

		// 读取请求体
		var bodyBytes []byte
		if c.Request.Body != nil {
			bodyBytes, _ = io.ReadAll(c.Request.Body)
			// 重新设置请求体，供后续处理使用
			c.Request.Body = io.NopCloser(bytes.NewBuffer(bodyBytes))
		}

		// 解析请求体获取method
		var reqBody map[string]interface{}
		method := "unknown"
		if len(bodyBytes) > 0 {
			if err := json.Unmarshal(bodyBytes, &reqBody); err == nil {
				if m, ok := reqBody["method"].(string); ok {
					method = m
				}
			}
		}

		// 设置上下文信息
		c.Set("mcp_start_time", startTime)
		c.Set("mcp_actor_id", actorID)
		c.Set("mcp_method", method)
		c.Set("mcp_body", reqBody)
		c.Set("mcp_client_info", c.Request.UserAgent())
		c.Set("mcp_client_ip", c.ClientIP())

		// 执行请求
		c.Next()

		// 统一记录日志
		logMCPCall(c, startTime, actorID, method, reqBody)
	}
}

// isMCPPath 判断是否为MCP相关路径
func isMCPPath(path string) bool {
	mcpPaths := []string{
		"/mcp",
		"/mcp/health",
	}

	for _, mcpPath := range mcpPaths {
		if path == mcpPath {
			return true
		}
	}
	return false
}

// logMCPCall 记录MCP调用日志
func logMCPCall(c *gin.Context, startTime time.Time, actorID, method string, reqBody map[string]interface{}) {
	latencyMs := int(time.Since(startTime).Milliseconds())

	// 获取响应状态
	status := "success"
	errorMsg := ""
	if c.Writer.Status() >= 400 {
		status = "error"
		if len(c.Errors) > 0 {
			errorMsg = c.Errors.String()
		}
	}

	// 获取结果数量（如果有的话）
	resultCount := 0
	if result, exists := c.Get("mcp_result_count"); exists {
		if count, ok := result.(int); ok {
			resultCount = count
		}
	}

	// 获取库ID（如果有的话）
	var libraryID *uint
	if libID, exists := c.Get("mcp_library_id"); exists {
		if id, ok := libID.(uint); ok {
			libraryID = &id
		}
	}

	// 获取正确的funcName
	funcName := method // 默认使用method
	if mappedName, exists := methodToFuncName[method]; exists {
		funcName = mappedName
	}

	// 记录日志
	logEntry := &mcplog.LogEntry{
		ActorID:     actorID,
		FuncName:    funcName,
		LibraryID:   libraryID,
		Params:      reqBody, // 整个请求体存到params
		ResultCount: resultCount,
		LatencyMs:   latencyMs,
		Status:      status,
		ErrorMsg:    errorMsg,
		CreatedAt:   time.Now(),
	}

	mcplog.Log(logEntry)
}
