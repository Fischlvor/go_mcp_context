package streamable

import (
	"encoding/json"
	"net/http"

	"go-mcp-context/internal/model/response"

	"github.com/gin-gonic/gin"
)

// StreamableResponseWriter Streamable HTTP响应写入器
// 可以根据请求复杂度选择返回JSON或SSE流
type StreamableResponseWriter struct {
	ctx          *gin.Context
	shouldStream bool // 是否使用流式响应
	method       string
	params       map[string]interface{}
}

// NewStreamableResponseWriter 创建Streamable响应写入器
func NewStreamableResponseWriter(c *gin.Context) *StreamableResponseWriter {
	return &StreamableResponseWriter{
		ctx:          c,
		shouldStream: false,
	}
}

// SetRequestInfo 设置请求信息，用于判断是否需要流式响应
func (w *StreamableResponseWriter) SetRequestInfo(method string, params map[string]interface{}) {
	w.method = method
	w.params = params
	w.shouldStream = shouldUseStreaming(method, params)
}

// shouldUseStreaming 判断是否应该使用流式响应
// 根据请求的method和params判断是否需要流式推送
//
// 当前业务特点：所有查询都是一次性返回，速度快，数据量小
// 因此暂时不使用流式响应，全部返回JSON
//
// 未来如果有以下需求时可以启用流式：
// 1. 大数据量分批返回（1000+条记录）
// 2. 长时间处理任务（10秒+）
// 3. 需要实时进度反馈
func shouldUseStreaming(method string, params map[string]interface{}) bool {
	// 当前全部返回JSON，不使用流式
	return false

	// 未来需要流式时，可以使用以下逻辑：
	/*
		switch method {
		case "initialize", "tools/list", "resources/list", "resources/templates/list":
			// 简单查询，直接返回JSON
			return false

		case "tools/call":
			// 检查具体的工具调用
			if toolName, ok := params["name"].(string); ok {
				switch toolName {
				case "search-libraries":
					// 搜索可能需要较长时间，使用流式响应
					return true
				case "get-library-docs":
					// 获取文档可能需要较长时间，使用流式响应
					return true
				}
			}
			return false

		default:
			// 其他方法默认不使用流式
			return false
		}
	*/
}

// WriteResponse 写入成功响应
// 根据shouldStream决定返回JSON还是SSE流
func (w *StreamableResponseWriter) WriteResponse(resp *response.MCPResponse) error {
	if w.shouldStream {
		// 流式响应: 设置SSE响应头并推送数据
		return w.writeSSEResponse(resp)
	}

	// 直接JSON响应
	w.ctx.JSON(http.StatusOK, resp)
	return nil
}

// WriteError 写入错误响应
func (w *StreamableResponseWriter) WriteError(err *response.MCPError, id interface{}) error {
	resp := &response.MCPResponse{
		JSONRPC: "2.0",
		ID:      id,
		Error:   err,
	}

	if w.shouldStream {
		return w.writeSSEResponse(resp)
	}

	w.ctx.JSON(http.StatusOK, resp)
	return nil
}

// writeSSEResponse 写入SSE格式的响应
// 将响应数据格式化为SSE格式并推送
func (w *StreamableResponseWriter) writeSSEResponse(resp *response.MCPResponse) error {
	// 设置SSE响应头
	w.ctx.Header("Content-Type", "text/event-stream")
	w.ctx.Header("Cache-Control", "no-cache")
	w.ctx.Header("Connection", "keep-alive")
	w.ctx.Header("Access-Control-Allow-Origin", "*")

	// 可选: 发送进度消息
	// w.sendProgress("处理中...")

	// 发送最终结果
	data, err := json.Marshal(resp)
	if err != nil {
		return err
	}

	// SSE格式: data: {json}\n\n
	_, err = w.ctx.Writer.Write([]byte("data: " + string(data) + "\n\n"))
	if err != nil {
		return err
	}

	// 立即刷新，推送给客户端
	if flusher, ok := w.ctx.Writer.(http.Flusher); ok {
		flusher.Flush()
	}

	return nil
}

// sendProgress 发送进度消息 (可选功能)
// 用于在处理过程中向客户端推送进度信息
func (w *StreamableResponseWriter) sendProgress(message string) error {
	progressData := map[string]interface{}{
		"type":    "progress",
		"message": message,
	}

	data, err := json.Marshal(progressData)
	if err != nil {
		return err
	}

	_, err = w.ctx.Writer.Write([]byte("data: " + string(data) + "\n\n"))
	if err != nil {
		return err
	}

	if flusher, ok := w.ctx.Writer.(http.Flusher); ok {
		flusher.Flush()
	}

	return nil
}

// Close 关闭写入器
func (w *StreamableResponseWriter) Close() error {
	return nil
}
