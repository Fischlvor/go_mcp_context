package response

import (
	"encoding/json"
	"fmt"

	"github.com/gin-gonic/gin"
)

// SSEResponse SSE 统一响应格式
type SSEResponse struct {
	Code int         `json:"code"` // 0=成功，非0=错误
	Msg  string      `json:"msg"`  // 消息
	Data interface{} `json:"data"` // 数据
}

// SSE 状态码
const (
	SSE_SUCCESS = 0
	SSE_ERROR   = 1
)

// SSEWriter SSE 写入器（通用）
type SSEWriter struct {
	writer  gin.ResponseWriter
	flusher interface{ Flush() }
}

// NewSSEWriter 创建 SSE 写入器
func NewSSEWriter(c *gin.Context) (*SSEWriter, bool) {
	// 设置 SSE 响应头
	c.Header("Content-Type", "text/event-stream")
	c.Header("Cache-Control", "no-cache")
	c.Header("Connection", "keep-alive")
	c.Header("Access-Control-Allow-Origin", "*")

	flusher, ok := c.Writer.(interface{ Flush() })
	if !ok {
		return nil, false
	}

	return &SSEWriter{
		writer:  c.Writer,
		flusher: flusher,
	}, true
}

// Send 发送 SSE 事件（统一格式）
func (s *SSEWriter) Send(code int, msg string, data interface{}) {
	event := SSEResponse{
		Code: code,
		Msg:  msg,
		Data: data,
	}
	jsonData, _ := json.Marshal(event)
	fmt.Fprintf(s.writer, "data: %s\n\n", jsonData)
	s.flusher.Flush()
}

// SendSuccess 发送成功事件
func (s *SSEWriter) SendSuccess(msg string, data interface{}) {
	s.Send(SSE_SUCCESS, msg, data)
}

// SendError 发送错误事件
func (s *SSEWriter) SendError(msg string) {
	s.Send(SSE_ERROR, msg, nil)
}

// SendErrorWithData 发送带数据的错误事件
func (s *SSEWriter) SendErrorWithData(msg string, data interface{}) {
	s.Send(SSE_ERROR, msg, data)
}
