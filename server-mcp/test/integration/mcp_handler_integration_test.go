package integration_test

import (
	"encoding/json"
	"testing"

	"go-mcp-context/internal/model/request"
	"go-mcp-context/internal/model/response"
	"go-mcp-context/internal/service"
	"go-mcp-context/internal/transport"

	"github.com/gin-gonic/gin"
)

// mockResponseWriter 模拟 ResponseWriter 用于测试
type mockResponseWriter struct {
	responses []interface{}
	errors    []*response.MCPError
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
	return nil
}

// Test_Integration_MCPHandler_Initialize 集成测试：MCP 初始化
func Test_Integration_MCPHandler_Initialize(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test in short mode")
	}

	handler := service.NewMCPHandler()
	writer := &mockResponseWriter{}

	t.Run("handle initialize request", func(t *testing.T) {
		// 创建模拟的 Gin 上下文
		ginCtx, _ := gin.CreateTestContext(nil)

		req := &transport.RequestContext{
			Transport: transport.TransportHTTP,
			Method:    "initialize",
			ID:        1,
			GinCtx:    ginCtx,
		}

		err := handler.ProcessRequest(req, writer)
		if err != nil {
			t.Fatalf("ProcessRequest(initialize) failed: %v", err)
		}

		if len(writer.responses) == 0 {
			t.Fatal("Expected response, got none")
		}

		// 验证响应内容
		resp := writer.responses[0].(*response.MCPResponse)
		if resp.JSONRPC != "2.0" {
			t.Errorf("Expected JSONRPC 2.0, got %s", resp.JSONRPC)
		}

		if resp.ID != 1 {
			t.Errorf("Expected ID 1, got %v", resp.ID)
		}

		// 验证 result 包含必要的字段
		result, ok := resp.Result.(map[string]interface{})
		if !ok {
			t.Fatal("Expected result to be map[string]interface{}")
		}

		if _, ok := result["protocolVersion"]; !ok {
			t.Error("Expected protocolVersion in result")
		}

		if _, ok := result["capabilities"]; !ok {
			t.Error("Expected capabilities in result")
		}

		if _, ok := result["serverInfo"]; !ok {
			t.Error("Expected serverInfo in result")
		}

		t.Log("✅ Initialize request handled successfully")
	})
}

// Test_Integration_MCPHandler_ToolsList 集成测试：工具列表
func Test_Integration_MCPHandler_ToolsList(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test in short mode")
	}

	handler := service.NewMCPHandler()
	writer := &mockResponseWriter{}

	t.Run("handle tools/list request", func(t *testing.T) {
		ginCtx, _ := gin.CreateTestContext(nil)

		req := &transport.RequestContext{
			Transport: transport.TransportHTTP,
			Method:    "tools/list",
			ID:        2,
			GinCtx:    ginCtx,
		}

		err := handler.ProcessRequest(req, writer)
		if err != nil {
			t.Fatalf("ProcessRequest(tools/list) failed: %v", err)
		}

		if len(writer.responses) == 0 {
			t.Fatal("Expected response, got none")
		}

		resp := writer.responses[0].(*response.MCPResponse)
		if resp.JSONRPC != "2.0" {
			t.Errorf("Expected JSONRPC 2.0, got %s", resp.JSONRPC)
		}

		// 验证工具列表
		t.Logf("Response Result type: %T, value: %+v", resp.Result, resp.Result)

		result, ok := resp.Result.(map[string]interface{})
		if !ok {
			t.Fatalf("Expected result to be map[string]interface{}, got %T", resp.Result)
		}

		t.Logf("Result map: %+v", result)

		// tools 可能是 []interface{} 或 []map[string]interface{}
		toolsRaw, exists := result["tools"]
		if !exists {
			t.Fatal("Expected 'tools' field in result")
		}

		// 尝试转换为 []interface{}
		var toolsCount int
		switch v := toolsRaw.(type) {
		case []interface{}:
			toolsCount = len(v)
		case []map[string]interface{}:
			toolsCount = len(v)
		default:
			t.Fatalf("Expected tools to be array, got %T", toolsRaw)
		}

		if toolsCount == 0 {
			t.Error("Expected at least one tool")
		}

		t.Logf("✅ Tools list returned %d tools", toolsCount)
	})
}

// Test_Integration_MCPHandler_ToolsCall 集成测试：工具调用
func Test_Integration_MCPHandler_ToolsCall(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test in short mode")
	}

	handler := service.NewMCPHandler()
	libService := &service.LibraryService{}

	// 先创建一个测试库
	lib, err := libService.Create(&request.LibraryCreate{
		Name:        "integration-test-mcp-tools",
		Description: "Integration test for MCP tools",
	})
	if err != nil {
		t.Fatalf("Failed to create library: %v", err)
	}

	t.Run("handle tools/call search-libraries", func(t *testing.T) {
		writer := &mockResponseWriter{}
		ginCtx, _ := gin.CreateTestContext(nil)

		// 构造工具调用参数
		params := map[string]interface{}{
			"name": "search-libraries",
			"arguments": map[string]interface{}{
				"libraryName": "integration-test",
			},
		}

		paramsJSON, _ := json.Marshal(params)

		req := &transport.RequestContext{
			Transport: transport.TransportHTTP,
			Method:    "tools/call",
			ID:        3,
			Params:    params,
			GinCtx:    ginCtx,
		}
		_ = paramsJSON // 避免未使用警告

		err := handler.ProcessRequest(req, writer)
		if err != nil {
			t.Fatalf("ProcessRequest(tools/call) failed: %v", err)
		}

		if len(writer.responses) == 0 {
			t.Fatal("Expected response, got none")
		}

		resp := writer.responses[0].(*response.MCPResponse)
		if resp.JSONRPC != "2.0" {
			t.Errorf("Expected JSONRPC 2.0, got %s", resp.JSONRPC)
		}

		t.Logf("✅ Tool call executed successfully, library ID: %d", lib.ID)
	})
}

// Test_Integration_MCPHandler_ErrorHandling 集成测试：错误处理
func Test_Integration_MCPHandler_ErrorHandling(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test in short mode")
	}

	handler := service.NewMCPHandler()
	writer := &mockResponseWriter{}

	t.Run("handle unknown method", func(t *testing.T) {
		ginCtx, _ := gin.CreateTestContext(nil)

		req := &transport.RequestContext{
			Transport: transport.TransportHTTP,
			Method:    "unknown/method",
			ID:        999,
			GinCtx:    ginCtx,
		}

		err := handler.ProcessRequest(req, writer)
		if err != nil {
			t.Fatalf("ProcessRequest(unknown) failed: %v", err)
		}

		if len(writer.errors) == 0 {
			t.Fatal("Expected error response, got none")
		}

		mcpErr := writer.errors[0]
		if mcpErr.Code != -32601 {
			t.Errorf("Expected error code -32601, got %d", mcpErr.Code)
		}

		if mcpErr.Message == "" {
			t.Error("Expected error message")
		}

		t.Logf("✅ Unknown method correctly handled with error: %s", mcpErr.Message)
	})
}

// Test_Integration_MCPHandler_ResourcesList 集成测试：资源列表
func Test_Integration_MCPHandler_ResourcesList(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test in short mode")
	}

	handler := service.NewMCPHandler()
	writer := &mockResponseWriter{}

	t.Run("handle resources/list request", func(t *testing.T) {
		ginCtx, _ := gin.CreateTestContext(nil)

		req := &transport.RequestContext{
			Transport: transport.TransportHTTP,
			Method:    "resources/list",
			ID:        4,
			GinCtx:    ginCtx,
		}

		err := handler.ProcessRequest(req, writer)
		if err != nil {
			t.Fatalf("ProcessRequest(resources/list) failed: %v", err)
		}

		if len(writer.responses) == 0 {
			t.Fatal("Expected response, got none")
		}

		resp := writer.responses[0].(*response.MCPResponse)
		if resp.JSONRPC != "2.0" {
			t.Errorf("Expected JSONRPC 2.0, got %s", resp.JSONRPC)
		}

		t.Log("✅ Resources list request handled successfully")
	})
}
