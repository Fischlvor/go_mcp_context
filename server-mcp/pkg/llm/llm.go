package llm

import (
	"context"
)

// LLMService LLM 服务接口
type LLMService interface {
	// Enrich 为文档块生成结构化描述
	Enrich(ctx context.Context, input EnrichInput) (*EnrichOutput, error)
	// Chat 通用对话（预留）
	Chat(ctx context.Context, prompt string) (string, error)
}

// EnrichInput Enrich 输入
type EnrichInput struct {
	Source   string // 文件来源路径
	Language string // 代码语言
	Content  string // 原始内容
}

// EnrichOutput Enrich 输出（Context7 风格）
type EnrichOutput struct {
	Title       string `json:"title"`        // 标题
	Description string `json:"description"`  // 描述
	ContentType string `json:"content_type"` // code 或 info
	Language    string `json:"language"`     // 代码语言
}
