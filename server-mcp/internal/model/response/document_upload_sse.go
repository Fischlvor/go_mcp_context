package response

import (
	"github.com/gin-gonic/gin"
)

// ============================================================
// 文档上传专用 SSE 响应
// ============================================================

// DocumentSSEData 文档上传 SSE 数据结构
type DocumentSSEData struct {
	Stage      string `json:"stage"`                 // 处理阶段: uploaded, parsing, chunking, embedding, saving, completed, failed
	Progress   int    `json:"progress"`              // 进度百分比 0-100
	Message    string `json:"message"`               // 阶段描述
	Status     string `json:"status"`                // 状态: processing, completed, failed
	DocumentID uint   `json:"document_id,omitempty"` // 文档ID
	Title      string `json:"title,omitempty"`       // 文档标题
}

// DocumentSSEWriter 文档上传专用 SSE 写入器
type DocumentSSEWriter struct {
	*SSEWriter // 继承通用 SSE 写入器
}

// NewDocumentSSEWriter 创建文档上传专用 SSE 写入器
func NewDocumentSSEWriter(c *gin.Context) (*DocumentSSEWriter, bool) {
	base, ok := NewSSEWriter(c)
	if !ok {
		return nil, false
	}
	return &DocumentSSEWriter{SSEWriter: base}, true
}

// SendProgress 发送进度事件
func (s *DocumentSSEWriter) SendProgress(stage string, progress int, message string, docID uint) {
	s.SendSuccess("success", DocumentSSEData{
		Stage:      stage,
		Progress:   progress,
		Message:    message,
		Status:     "processing",
		DocumentID: docID,
	})
}

// SendComplete 发送完成事件
func (s *DocumentSSEWriter) SendComplete(docID uint, title string) {
	s.SendSuccess("success", DocumentSSEData{
		Stage:      "completed",
		Progress:   100,
		Message:    "处理完成",
		Status:     "completed",
		DocumentID: docID,
		Title:      title,
	})
}

// SendFailed 发送失败事件
func (s *DocumentSSEWriter) SendFailed(message string, docID uint) {
	s.Send(SSE_ERROR, message, DocumentSSEData{
		Stage:      "failed",
		Progress:   0,
		Message:    message,
		Status:     "failed",
		DocumentID: docID,
	})
}
