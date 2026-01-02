package test_test

import (
	"fmt"
	"mime/multipart"
	"net/http/httptest"
	"os"
	"testing"
	"time"

	"go-mcp-context/internal/model/request"
	"go-mcp-context/internal/model/response"
	"go-mcp-context/internal/service"
	"go-mcp-context/internal/transport"
	"go-mcp-context/pkg/utils"

	"github.com/gin-gonic/gin"
)

// mockResponseWriter 模拟 ResponseWriter
type mockResponseWriter struct {
	responses []*response.MCPResponse
	errors    []*response.MCPError
	closed    bool
}

func newMockResponseWriter() *mockResponseWriter {
	return &mockResponseWriter{
		responses: make([]*response.MCPResponse, 0),
		errors:    make([]*response.MCPError, 0),
	}
}

func (m *mockResponseWriter) WriteResponse(resp *response.MCPResponse) error {
	m.responses = append(m.responses, resp)
	return nil
}

func (m *mockResponseWriter) WriteError(err *response.MCPError, id interface{}) error {
	m.errors = append(m.errors, err)
	return nil
}

func (m *mockResponseWriter) Close() error {
	m.closed = true
	return nil
}

// Test_MCPHandler_Initialize 测试 initialize 方法
func Test_MCPHandler_Initialize(t *testing.T) {
	handler := service.NewMCPHandler()
	writer := newMockResponseWriter()

	t.Run("initialize request", func(t *testing.T) {
		gin.SetMode(gin.TestMode)
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)

		req := &transport.RequestContext{
			Transport: transport.TransportHTTP,
			Method:    "initialize",
			Params:    map[string]interface{}{},
			ID:        1,
			GinCtx:    c,
		}

		err := handler.ProcessRequest(req, writer)
		if err != nil {
			t.Fatalf("ProcessRequest() error = %v", err)
		}

		if len(writer.responses) != 1 {
			t.Fatalf("Expected 1 response, got %d", len(writer.responses))
		}

		resp := writer.responses[0]
		if resp.JSONRPC != "2.0" {
			t.Errorf("Expected JSONRPC 2.0, got %s", resp.JSONRPC)
		}

		result, ok := resp.Result.(map[string]interface{})
		if !ok {
			t.Fatal("Expected result to be map[string]interface{}")
		}

		if result["protocolVersion"] != "2025-11-25" {
			t.Errorf("Expected protocolVersion 2025-11-25, got %v", result["protocolVersion"])
		}

		if _, ok := result["capabilities"]; !ok {
			t.Error("Expected capabilities in result")
		}

		if _, ok := result["serverInfo"]; !ok {
			t.Error("Expected serverInfo in result")
		}
	})
}

// Test_MCPHandler_Initialized 测试 initialized 通知
func Test_MCPHandler_Initialized(t *testing.T) {
	handler := service.NewMCPHandler()
	writer := newMockResponseWriter()

	t.Run("initialized notification", func(t *testing.T) {
		gin.SetMode(gin.TestMode)
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)

		req := &transport.RequestContext{
			Transport: transport.TransportHTTP,
			Method:    "initialized",
			Params:    map[string]interface{}{},
			ID:        nil, // 通知没有 ID
			GinCtx:    c,
		}

		err := handler.ProcessRequest(req, writer)
		if err != nil {
			t.Fatalf("ProcessRequest() error = %v", err)
		}
	})
}

// Test_MCPHandler_ToolsList 测试 tools/list 方法
func Test_MCPHandler_ToolsList(t *testing.T) {
	handler := service.NewMCPHandler()
	writer := newMockResponseWriter()

	t.Run("list tools", func(t *testing.T) {
		gin.SetMode(gin.TestMode)
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)

		req := &transport.RequestContext{
			Transport: transport.TransportHTTP,
			Method:    "tools/list",
			Params:    map[string]interface{}{},
			ID:        2,
			GinCtx:    c,
		}

		err := handler.ProcessRequest(req, writer)
		if err != nil {
			t.Fatalf("ProcessRequest() error = %v", err)
		}

		if len(writer.responses) != 1 {
			t.Fatalf("Expected 1 response, got %d", len(writer.responses))
		}

		resp := writer.responses[0]
		result, ok := resp.Result.(map[string]interface{})
		if !ok {
			t.Fatalf("Expected result to be map[string]interface{}, got %T", resp.Result)
		}

		toolsRaw, exists := result["tools"]
		if !exists {
			t.Fatal("Expected 'tools' key in result")
		}

		// tools 可能是 []map[string]interface{} 类型
		var toolsCount int
		switch v := toolsRaw.(type) {
		case []interface{}:
			toolsCount = len(v)
		case []map[string]interface{}:
			toolsCount = len(v)
		default:
			t.Fatalf("Expected tools to be []interface{} or []map[string]interface{}, got %T", toolsRaw)
		}

		if toolsCount != 2 {
			t.Errorf("Expected 2 tools, got %d", toolsCount)
		}
	})
}

// Test_MCPHandler_ToolsCall 测试 tools/call 方法
func Test_MCPHandler_ToolsCall(t *testing.T) {
	handler := service.NewMCPHandler()

	t.Run("call search-libraries tool", func(t *testing.T) {
		writer := newMockResponseWriter()
		gin.SetMode(gin.TestMode)
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)

		// 创建测试库
		libService := &service.LibraryService{}
		lib, err := libService.Create(&request.LibraryCreate{
			Name:        "test-mcp-lib",
			Description: "Test library for MCP",
		})
		if err != nil {
			t.Fatalf("Failed to create library: %v", err)
		}
		defer libService.Delete(lib.ID)

		req := &transport.RequestContext{
			Transport: transport.TransportHTTP,
			Method:    "tools/call",
			Params: map[string]interface{}{
				"name": "search-libraries",
				"arguments": map[string]interface{}{
					"libraryName": "test",
				},
			},
			ID:     3,
			GinCtx: c,
		}

		err = handler.ProcessRequest(req, writer)
		if err != nil {
			t.Fatalf("ProcessRequest() error = %v", err)
		}

		if len(writer.responses) != 1 {
			t.Fatalf("Expected 1 response, got %d", len(writer.responses))
		}
	})

	t.Run("call get-library-docs tool", func(t *testing.T) {
		writer := newMockResponseWriter()
		gin.SetMode(gin.TestMode)
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)

		// 创建测试库
		libService := &service.LibraryService{}
		lib, err := libService.Create(&request.LibraryCreate{
			Name:        "test-get-docs-lib",
			Description: "Test library for get-library-docs",
		})
		if err != nil {
			t.Fatalf("Failed to create library: %v", err)
		}
		defer libService.Delete(lib.ID)

		req := &transport.RequestContext{
			Transport: transport.TransportHTTP,
			Method:    "tools/call",
			Params: map[string]interface{}{
				"name": "get-library-docs",
				"arguments": map[string]interface{}{
					"libraryId": float64(lib.ID),
					"topic":     "test",
					"version":   lib.DefaultVersion,
					"mode":      "code",
					"page":      float64(1),
				},
			},
			ID:     4,
			GinCtx: c,
		}

		err = handler.ProcessRequest(req, writer)
		if err != nil {
			t.Fatalf("ProcessRequest() error = %v", err)
		}

		if len(writer.responses) != 1 {
			t.Fatalf("Expected 1 response, got %d", len(writer.responses))
		}
	})

	t.Run("call unknown tool", func(t *testing.T) {
		writer := newMockResponseWriter()
		gin.SetMode(gin.TestMode)
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)

		req := &transport.RequestContext{
			Transport: transport.TransportHTTP,
			Method:    "tools/call",
			Params: map[string]interface{}{
				"name":      "unknown-tool",
				"arguments": map[string]interface{}{},
			},
			ID:     5,
			GinCtx: c,
		}

		err := handler.ProcessRequest(req, writer)
		if err != nil {
			t.Fatalf("ProcessRequest() error = %v", err)
		}

		if len(writer.errors) != 1 {
			t.Fatalf("Expected 1 error, got %d", len(writer.errors))
		}
	})
}

// Test_MCPHandler_ResourcesList 测试 resources/list 方法
func Test_MCPHandler_ResourcesList(t *testing.T) {
	handler := service.NewMCPHandler()
	writer := newMockResponseWriter()

	t.Run("list resources", func(t *testing.T) {
		gin.SetMode(gin.TestMode)
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)

		req := &transport.RequestContext{
			Transport: transport.TransportHTTP,
			Method:    "resources/list",
			Params:    map[string]interface{}{},
			ID:        6,
			GinCtx:    c,
		}

		err := handler.ProcessRequest(req, writer)
		if err != nil {
			t.Fatalf("ProcessRequest() error = %v", err)
		}

		if len(writer.responses) != 1 {
			t.Fatalf("Expected 1 response, got %d", len(writer.responses))
		}
	})
}

// Test_MCPHandler_ResourcesRead 测试 resources/read 方法
func Test_MCPHandler_ResourcesRead(t *testing.T) {
	handler := service.NewMCPHandler()

	t.Run("read library resource", func(t *testing.T) {
		writer := newMockResponseWriter()
		gin.SetMode(gin.TestMode)
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)

		// 创建测试库
		libService := &service.LibraryService{}
		lib, err := libService.Create(&request.LibraryCreate{
			Name:        "test-resource-lib",
			Description: "Test library for resources",
		})
		if err != nil {
			t.Fatalf("Failed to create library: %v", err)
		}
		defer libService.Delete(lib.ID)

		// 使用正确的 URI 格式: go-mcp-context:///library/{libraryId}
		uri := fmt.Sprintf("go-mcp-context:///library/%d", lib.ID)
		req := &transport.RequestContext{
			Transport: transport.TransportHTTP,
			Method:    "resources/read",
			Params: map[string]interface{}{
				"uri": uri,
			},
			ID:     7,
			GinCtx: c,
		}

		err = handler.ProcessRequest(req, writer)
		if err != nil {
			t.Fatalf("ProcessRequest() error = %v", err)
		}

		if len(writer.responses) != 1 {
			t.Fatalf("Expected 1 response, got %d", len(writer.responses))
		}
	})

	t.Run("read docs resource successfully", func(t *testing.T) {
		writer := newMockResponseWriter()
		gin.SetMode(gin.TestMode)
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)

		// 创建测试库
		libService := &service.LibraryService{}
		lib, err := libService.Create(&request.LibraryCreate{
			Name:        "test-docs-resource-lib",
			Description: "Test library for docs resources",
		})
		if err != nil {
			t.Fatalf("Failed to create library: %v", err)
		}
		defer libService.Delete(lib.ID)

		// 上传文档
		docService := &service.DocumentService{}
		content := []byte("# Test Document\n\nThis is a test document for docs resource reading.")
		tmpFile, err := os.CreateTemp("", "docs-resource-*.md")
		if err != nil {
			t.Fatalf("Failed to create temp file: %v", err)
		}
		defer os.Remove(tmpFile.Name())

		if _, err := tmpFile.Write(content); err != nil {
			t.Fatalf("Failed to write temp file: %v", err)
		}
		tmpFile.Seek(0, 0)

		fileHeader := &multipart.FileHeader{
			Filename: "docs-resource.md",
			Size:     int64(len(content)),
		}

		taskID := utils.GenerateTaskID()
		_, err = docService.Upload(lib.ID, lib.DefaultVersion, tmpFile, fileHeader, "test-user", taskID)
		if err != nil {
			t.Logf("Upload() error = %v (may fail if storage not available)", err)
			return
		}

		// 等待文档处理
		time.Sleep(5 * time.Second)

		// 读取文档资源
		uri := fmt.Sprintf("go-mcp-context:///docs/chunk/%d/%s/test", lib.ID, lib.DefaultVersion)
		req := &transport.RequestContext{
			Transport: transport.TransportHTTP,
			Method:    "resources/read",
			Params: map[string]interface{}{
				"uri": uri,
			},
			ID:     8,
			GinCtx: c,
		}

		err = handler.ProcessRequest(req, writer)
		if err != nil {
			t.Fatalf("ProcessRequest() error = %v", err)
		}

		if len(writer.responses) != 1 {
			t.Fatalf("Expected 1 response, got %d", len(writer.responses))
		}

		// 验证响应包含 contents
		resp := writer.responses[0]
		result, ok := resp.Result.(map[string]interface{})
		if !ok {
			t.Fatal("Expected result to be map[string]interface{}")
		}

		if _, exists := result["contents"]; !exists {
			t.Error("Expected 'contents' key in result")
		}

		t.Log("✅ Successfully read docs resource")
	})
}

// Test_MCPHandler_UnknownMethod 测试未知方法
func Test_MCPHandler_UnknownMethod(t *testing.T) {
	handler := service.NewMCPHandler()
	writer := newMockResponseWriter()

	t.Run("unknown method", func(t *testing.T) {
		gin.SetMode(gin.TestMode)
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)

		req := &transport.RequestContext{
			Transport: transport.TransportHTTP,
			Method:    "unknown/method",
			Params:    map[string]interface{}{},
			ID:        999,
			GinCtx:    c,
		}

		err := handler.ProcessRequest(req, writer)
		if err != nil {
			t.Fatalf("ProcessRequest() error = %v", err)
		}

		if len(writer.errors) != 1 {
			t.Fatalf("Expected 1 error, got %d", len(writer.errors))
		}

		if writer.errors[0].Code != -32601 {
			t.Errorf("Expected error code -32601, got %d", writer.errors[0].Code)
		}
	})
}

// Test_MCPHandler_ToolsCall_EdgeCases 测试 tools/call 边界情况
func Test_MCPHandler_ToolsCall_EdgeCases(t *testing.T) {
	handler := service.NewMCPHandler()

	t.Run("missing tool name", func(t *testing.T) {
		writer := newMockResponseWriter()
		gin.SetMode(gin.TestMode)
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)

		req := &transport.RequestContext{
			Transport: transport.TransportHTTP,
			Method:    "tools/call",
			Params: map[string]interface{}{
				// 缺少 name 参数
				"arguments": map[string]interface{}{},
			},
			ID:     10,
			GinCtx: c,
		}

		err := handler.ProcessRequest(req, writer)
		if err != nil {
			t.Fatalf("ProcessRequest() error = %v", err)
		}

		if len(writer.errors) != 1 {
			t.Fatalf("Expected 1 error, got %d", len(writer.errors))
		}

		if writer.errors[0].Code != -32602 {
			t.Errorf("Expected error code -32602, got %d", writer.errors[0].Code)
		}
	})

	t.Run("search-libraries with empty libraryName", func(t *testing.T) {
		writer := newMockResponseWriter()
		gin.SetMode(gin.TestMode)
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)

		req := &transport.RequestContext{
			Transport: transport.TransportHTTP,
			Method:    "tools/call",
			Params: map[string]interface{}{
				"name": "search-libraries",
				"arguments": map[string]interface{}{
					"libraryName": "", // 空字符串
				},
			},
			ID:     11,
			GinCtx: c,
		}

		err := handler.ProcessRequest(req, writer)
		if err != nil {
			t.Fatalf("ProcessRequest() error = %v", err)
		}

		if len(writer.errors) != 1 {
			t.Fatalf("Expected 1 error, got %d", len(writer.errors))
		}

		if writer.errors[0].Code != -32602 {
			t.Errorf("Expected error code -32602, got %d", writer.errors[0].Code)
		}
	})

	t.Run("get-library-docs with empty topic", func(t *testing.T) {
		writer := newMockResponseWriter()
		gin.SetMode(gin.TestMode)
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)

		req := &transport.RequestContext{
			Transport: transport.TransportHTTP,
			Method:    "tools/call",
			Params: map[string]interface{}{
				"name": "get-library-docs",
				"arguments": map[string]interface{}{
					"libraryId": float64(1),
					"topic":     "", // 空 topic
				},
			},
			ID:     12,
			GinCtx: c,
		}

		err := handler.ProcessRequest(req, writer)
		if err != nil {
			t.Fatalf("ProcessRequest() error = %v", err)
		}

		if len(writer.errors) != 1 {
			t.Fatalf("Expected 1 error, got %d", len(writer.errors))
		}

		if writer.errors[0].Code != -32602 {
			t.Errorf("Expected error code -32602, got %d", writer.errors[0].Code)
		}
	})
}

// Test_MCPHandler_ResourceTemplatesList 测试 resources/templates/list
func Test_MCPHandler_ResourceTemplatesList(t *testing.T) {
	handler := service.NewMCPHandler()
	writer := newMockResponseWriter()

	t.Run("list resource templates", func(t *testing.T) {
		gin.SetMode(gin.TestMode)
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)

		req := &transport.RequestContext{
			Transport: transport.TransportHTTP,
			Method:    "resources/templates/list",
			Params:    map[string]interface{}{},
			ID:        13,
			GinCtx:    c,
		}

		err := handler.ProcessRequest(req, writer)
		if err != nil {
			t.Fatalf("ProcessRequest() error = %v", err)
		}

		if len(writer.responses) != 1 {
			t.Fatalf("Expected 1 response, got %d", len(writer.responses))
		}

		resp := writer.responses[0]
		result, ok := resp.Result.(map[string]interface{})
		if !ok {
			t.Fatal("Expected result to be map[string]interface{}")
		}

		templates, ok := result["resourceTemplates"].([]map[string]interface{})
		if !ok {
			t.Fatal("Expected resourceTemplates to be []map[string]interface{}")
		}

		if len(templates) != 2 {
			t.Errorf("Expected 2 templates, got %d", len(templates))
		}
	})
}

// Test_MCPHandler_ResourcesRead_EdgeCases 测试 resources/read 边界情况
func Test_MCPHandler_ResourcesRead_EdgeCases(t *testing.T) {
	handler := service.NewMCPHandler()

	t.Run("missing uri parameter", func(t *testing.T) {
		writer := newMockResponseWriter()
		gin.SetMode(gin.TestMode)
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)

		req := &transport.RequestContext{
			Transport: transport.TransportHTTP,
			Method:    "resources/read",
			Params:    map[string]interface{}{}, // 缺少 uri
			ID:        14,
			GinCtx:    c,
		}

		err := handler.ProcessRequest(req, writer)
		if err != nil {
			t.Fatalf("ProcessRequest() error = %v", err)
		}

		if len(writer.errors) != 1 {
			t.Fatalf("Expected 1 error, got %d", len(writer.errors))
		}

		if writer.errors[0].Code != -32602 {
			t.Errorf("Expected error code -32602, got %d", writer.errors[0].Code)
		}
	})

	t.Run("empty uri", func(t *testing.T) {
		writer := newMockResponseWriter()
		gin.SetMode(gin.TestMode)
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)

		req := &transport.RequestContext{
			Transport: transport.TransportHTTP,
			Method:    "resources/read",
			Params: map[string]interface{}{
				"uri": "", // 空 URI
			},
			ID:     15,
			GinCtx: c,
		}

		err := handler.ProcessRequest(req, writer)
		if err != nil {
			t.Fatalf("ProcessRequest() error = %v", err)
		}

		if len(writer.errors) != 1 {
			t.Fatalf("Expected 1 error, got %d", len(writer.errors))
		}
	})

	t.Run("invalid uri format", func(t *testing.T) {
		writer := newMockResponseWriter()
		gin.SetMode(gin.TestMode)
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)

		req := &transport.RequestContext{
			Transport: transport.TransportHTTP,
			Method:    "resources/read",
			Params: map[string]interface{}{
				"uri": "invalid://uri:with:colons", // 无效 URI
			},
			ID:     16,
			GinCtx: c,
		}

		err := handler.ProcessRequest(req, writer)
		if err != nil {
			t.Fatalf("ProcessRequest() error = %v", err)
		}

		if len(writer.errors) != 1 {
			t.Fatalf("Expected 1 error, got %d", len(writer.errors))
		}
	})

	t.Run("empty path segments", func(t *testing.T) {
		writer := newMockResponseWriter()
		gin.SetMode(gin.TestMode)
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)

		req := &transport.RequestContext{
			Transport: transport.TransportHTTP,
			Method:    "resources/read",
			Params: map[string]interface{}{
				"uri": "go-mcp-context:///", // 空路径
			},
			ID:     17,
			GinCtx: c,
		}

		err := handler.ProcessRequest(req, writer)
		if err != nil {
			t.Fatalf("ProcessRequest() error = %v", err)
		}

		if len(writer.errors) != 1 {
			t.Fatalf("Expected 1 error, got %d", len(writer.errors))
		}
	})

	t.Run("unknown resource type", func(t *testing.T) {
		writer := newMockResponseWriter()
		gin.SetMode(gin.TestMode)
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)

		req := &transport.RequestContext{
			Transport: transport.TransportHTTP,
			Method:    "resources/read",
			Params: map[string]interface{}{
				"uri": "go-mcp-context:///unknown/resource",
			},
			ID:     18,
			GinCtx: c,
		}

		err := handler.ProcessRequest(req, writer)
		if err != nil {
			t.Fatalf("ProcessRequest() error = %v", err)
		}

		if len(writer.errors) != 1 {
			t.Fatalf("Expected 1 error, got %d", len(writer.errors))
		}

		if writer.errors[0].Code != -32602 {
			t.Errorf("Expected error code -32602, got %d", writer.errors[0].Code)
		}
	})
}

// Test_MCPHandler_ReadLibraryResource_EdgeCases 测试 library 资源读取边界情况
func Test_MCPHandler_ReadLibraryResource_EdgeCases(t *testing.T) {
	handler := service.NewMCPHandler()

	t.Run("library uri with missing libraryId", func(t *testing.T) {
		writer := newMockResponseWriter()
		gin.SetMode(gin.TestMode)
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)

		req := &transport.RequestContext{
			Transport: transport.TransportHTTP,
			Method:    "resources/read",
			Params: map[string]interface{}{
				"uri": "go-mcp-context:///library/", // 缺少 libraryId
			},
			ID:     19,
			GinCtx: c,
		}

		err := handler.ProcessRequest(req, writer)
		if err != nil {
			t.Fatalf("ProcessRequest() error = %v", err)
		}

		if len(writer.errors) != 1 {
			t.Fatalf("Expected 1 error, got %d", len(writer.errors))
		}
	})

	t.Run("library uri with invalid libraryId format", func(t *testing.T) {
		writer := newMockResponseWriter()
		gin.SetMode(gin.TestMode)
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)

		req := &transport.RequestContext{
			Transport: transport.TransportHTTP,
			Method:    "resources/read",
			Params: map[string]interface{}{
				"uri": "go-mcp-context:///library/abc", // 非数字 ID
			},
			ID:     20,
			GinCtx: c,
		}

		err := handler.ProcessRequest(req, writer)
		if err != nil {
			t.Fatalf("ProcessRequest() error = %v", err)
		}

		if len(writer.errors) != 1 {
			t.Fatalf("Expected 1 error, got %d", len(writer.errors))
		}

		if writer.errors[0].Code != -32602 {
			t.Errorf("Expected error code -32602, got %d", writer.errors[0].Code)
		}
	})

	t.Run("library uri with zero libraryId", func(t *testing.T) {
		writer := newMockResponseWriter()
		gin.SetMode(gin.TestMode)
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)

		req := &transport.RequestContext{
			Transport: transport.TransportHTTP,
			Method:    "resources/read",
			Params: map[string]interface{}{
				"uri": "go-mcp-context:///library/0", // ID 为 0
			},
			ID:     21,
			GinCtx: c,
		}

		err := handler.ProcessRequest(req, writer)
		if err != nil {
			t.Fatalf("ProcessRequest() error = %v", err)
		}

		if len(writer.errors) != 1 {
			t.Fatalf("Expected 1 error, got %d", len(writer.errors))
		}

		if writer.errors[0].Code != -32602 {
			t.Errorf("Expected error code -32602, got %d", writer.errors[0].Code)
		}
	})

	t.Run("library uri with non-existent libraryId", func(t *testing.T) {
		writer := newMockResponseWriter()
		gin.SetMode(gin.TestMode)
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)

		req := &transport.RequestContext{
			Transport: transport.TransportHTTP,
			Method:    "resources/read",
			Params: map[string]interface{}{
				"uri": "go-mcp-context:///library/999999", // 不存在的 ID
			},
			ID:     22,
			GinCtx: c,
		}

		err := handler.ProcessRequest(req, writer)
		if err != nil {
			t.Fatalf("ProcessRequest() error = %v", err)
		}

		if len(writer.errors) != 1 {
			t.Fatalf("Expected 1 error, got %d", len(writer.errors))
		}
	})
}

// Test_MCPHandler_ReadDocsResource_EdgeCases 测试 docs 资源读取边界情况
func Test_MCPHandler_ReadDocsResource_EdgeCases(t *testing.T) {
	handler := service.NewMCPHandler()

	t.Run("docs uri with wrong format", func(t *testing.T) {
		writer := newMockResponseWriter()
		gin.SetMode(gin.TestMode)
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)

		req := &transport.RequestContext{
			Transport: transport.TransportHTTP,
			Method:    "resources/read",
			Params: map[string]interface{}{
				"uri": "go-mcp-context:///docs/1/v1.0.0/api", // 缺少 chunk
			},
			ID:     23,
			GinCtx: c,
		}

		err := handler.ProcessRequest(req, writer)
		if err != nil {
			t.Fatalf("ProcessRequest() error = %v", err)
		}

		if len(writer.errors) != 1 {
			t.Fatalf("Expected 1 error, got %d", len(writer.errors))
		}
	})

	t.Run("docs uri with invalid libraryId", func(t *testing.T) {
		writer := newMockResponseWriter()
		gin.SetMode(gin.TestMode)
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)

		req := &transport.RequestContext{
			Transport: transport.TransportHTTP,
			Method:    "resources/read",
			Params: map[string]interface{}{
				"uri": "go-mcp-context:///docs/chunk/abc/v1.0.0/api", // 非数字 ID
			},
			ID:     24,
			GinCtx: c,
		}

		err := handler.ProcessRequest(req, writer)
		if err != nil {
			t.Fatalf("ProcessRequest() error = %v", err)
		}

		if len(writer.errors) != 1 {
			t.Fatalf("Expected 1 error, got %d", len(writer.errors))
		}
	})

	t.Run("docs uri with zero libraryId", func(t *testing.T) {
		writer := newMockResponseWriter()
		gin.SetMode(gin.TestMode)
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)

		req := &transport.RequestContext{
			Transport: transport.TransportHTTP,
			Method:    "resources/read",
			Params: map[string]interface{}{
				"uri": "go-mcp-context:///docs/chunk/0/v1.0.0/api", // ID 为 0
			},
			ID:     25,
			GinCtx: c,
		}

		err := handler.ProcessRequest(req, writer)
		if err != nil {
			t.Fatalf("ProcessRequest() error = %v", err)
		}

		if len(writer.errors) != 1 {
			t.Fatalf("Expected 1 error, got %d", len(writer.errors))
		}
	})

	t.Run("docs uri with empty version", func(t *testing.T) {
		writer := newMockResponseWriter()
		gin.SetMode(gin.TestMode)
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)

		req := &transport.RequestContext{
			Transport: transport.TransportHTTP,
			Method:    "resources/read",
			Params: map[string]interface{}{
				"uri": "go-mcp-context:///docs/chunk/1//api", // 空 version
			},
			ID:     26,
			GinCtx: c,
		}

		err := handler.ProcessRequest(req, writer)
		if err != nil {
			t.Fatalf("ProcessRequest() error = %v", err)
		}

		if len(writer.errors) != 1 {
			t.Fatalf("Expected 1 error, got %d", len(writer.errors))
		}
	})

	t.Run("docs uri with empty topic", func(t *testing.T) {
		writer := newMockResponseWriter()
		gin.SetMode(gin.TestMode)
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)

		req := &transport.RequestContext{
			Transport: transport.TransportHTTP,
			Method:    "resources/read",
			Params: map[string]interface{}{
				"uri": "go-mcp-context:///docs/chunk/1/v1.0.0/", // 空 topic
			},
			ID:     27,
			GinCtx: c,
		}

		err := handler.ProcessRequest(req, writer)
		if err != nil {
			t.Fatalf("ProcessRequest() error = %v", err)
		}

		if len(writer.errors) != 1 {
			t.Fatalf("Expected 1 error, got %d", len(writer.errors))
		}
	})
}
