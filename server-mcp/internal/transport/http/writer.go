package http

import (
	"net/http"

	"go-mcp-context/internal/model/response"

	"github.com/gin-gonic/gin"
)

// HTTPResponseWriter HTTP响应写入器
// 实现传统的HTTP JSON响应
type HTTPResponseWriter struct {
	ctx *gin.Context
}

// NewHTTPResponseWriter 创建HTTP响应写入器
func NewHTTPResponseWriter(c *gin.Context) *HTTPResponseWriter {
	return &HTTPResponseWriter{
		ctx: c,
	}
}

// WriteResponse 写入成功响应
// 直接返回JSON格式的MCP响应
func (w *HTTPResponseWriter) WriteResponse(resp *response.MCPResponse) error {
	w.ctx.JSON(http.StatusOK, resp)
	return nil
}

// WriteError 写入错误响应
// 构造包含错误信息的MCP响应并返回
func (w *HTTPResponseWriter) WriteError(err *response.MCPError, id interface{}) error {
	resp := &response.MCPResponse{
		JSONRPC: "2.0",
		ID:      id,
		Error:   err,
	}
	w.ctx.JSON(http.StatusOK, resp)
	return nil
}

// Close 关闭写入器
// HTTP响应写入器无需特殊清理
func (w *HTTPResponseWriter) Close() error {
	return nil
}
