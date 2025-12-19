package response

import (
	"github.com/gin-gonic/gin"
)

// ============================================================
// 版本刷新专用 SSE 响应
// ============================================================

// RefreshStatus 版本刷新状态（用于 SSE 推送）
type RefreshStatus struct {
	DocID    uint   `json:"doc_id,omitempty"`    // 文档 ID
	DocTitle string `json:"doc_title,omitempty"` // 文档标题
	Stage    string `json:"stage"`               // 阶段: started, doc_processing, doc_completed, doc_failed, all_completed, error
	Current  int    `json:"current"`             // 当前第几个
	Total    int    `json:"total"`               // 总共几个
	Message  string `json:"message"`             // 消息
}

// RefreshSSEWriter 版本刷新专用 SSE 写入器
type RefreshSSEWriter struct {
	*SSEWriter // 继承通用 SSE 写入器
}

// NewRefreshSSEWriter 创建版本刷新专用 SSE 写入器
func NewRefreshSSEWriter(c *gin.Context) (*RefreshSSEWriter, bool) {
	base, ok := NewSSEWriter(c)
	if !ok {
		return nil, false
	}
	return &RefreshSSEWriter{SSEWriter: base}, true
}

// SendStatus 发送刷新状态事件（通用）
func (s *RefreshSSEWriter) SendStatus(status RefreshStatus) {
	if status.Stage == "error" {
		s.Send(SSE_ERROR, status.Message, status)
	} else {
		s.SendSuccess(status.Message, status)
	}
}

// SendStarted 发送开始事件
func (s *RefreshSSEWriter) SendStarted(total int, message string) {
	s.SendSuccess(message, RefreshStatus{
		Stage:   "started",
		Current: 0,
		Total:   total,
		Message: message,
	})
}

// SendDocProcessing 发送文档处理中事件
func (s *RefreshSSEWriter) SendDocProcessing(docID uint, docTitle string, current, total int) {
	s.SendSuccess("正在处理: "+docTitle, RefreshStatus{
		DocID:    docID,
		DocTitle: docTitle,
		Stage:    "doc_processing",
		Current:  current,
		Total:    total,
		Message:  "正在处理: " + docTitle,
	})
}

// SendDocCompleted 发送文档完成事件
func (s *RefreshSSEWriter) SendDocCompleted(docID uint, docTitle string, current, total int) {
	s.SendSuccess("处理完成: "+docTitle, RefreshStatus{
		DocID:    docID,
		DocTitle: docTitle,
		Stage:    "doc_completed",
		Current:  current,
		Total:    total,
		Message:  "处理完成: " + docTitle,
	})
}

// SendDocFailed 发送文档失败事件
func (s *RefreshSSEWriter) SendDocFailed(docID uint, docTitle string, current, total int, errMsg string) {
	s.Send(SSE_ERROR, errMsg, RefreshStatus{
		DocID:    docID,
		DocTitle: docTitle,
		Stage:    "doc_failed",
		Current:  current,
		Total:    total,
		Message:  errMsg,
	})
}

// SendAllCompleted 发送全部完成事件
func (s *RefreshSSEWriter) SendAllCompleted(total int, message string) {
	s.SendSuccess(message, RefreshStatus{
		Stage:   "all_completed",
		Current: total,
		Total:   total,
		Message: message,
	})
}
