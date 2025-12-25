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
	// GenerateLibraryTitle 为库生成简短友好的名称
	GenerateLibraryTitle(ctx context.Context, repoName, description string) (string, error)
}

// EnrichInput Enrich 输入
type EnrichInput struct {
	Content string // 原始内容（不含标题行）
	Headers string // 标题层级（如：Gin Web Framework > Getting Started > Installation）
}

// EnrichOutput Enrich 输出（Context7 风格）
type EnrichOutput struct {
	Title       string `json:"title"`       // 简洁标题（5-15字）
	Description string `json:"description"` // 描述（50-150字）
}
