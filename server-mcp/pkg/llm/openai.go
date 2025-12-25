package llm

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"strings"
	"text/template"

	"github.com/sashabaranov/go-openai"
)

// OpenAILLM OpenAI LLM 实现
type OpenAILLM struct {
	client      *openai.Client
	model       string
	maxTokens   int
	temperature float32
}

// OpenAIConfig OpenAI 配置
type OpenAIConfig struct {
	APIKey      string
	BaseURL     string
	Model       string
	MaxTokens   int
	Temperature float32
}

// NewOpenAILLM 创建 OpenAI LLM 实例
func NewOpenAILLM(cfg OpenAIConfig) *OpenAILLM {
	config := openai.DefaultConfig(cfg.APIKey)
	if cfg.BaseURL != "" {
		config.BaseURL = cfg.BaseURL
	}

	// 默认值
	model := cfg.Model
	if model == "" {
		model = "gpt-4o-mini"
	}
	maxTokens := cfg.MaxTokens
	if maxTokens == 0 {
		maxTokens = 500
	}
	temperature := cfg.Temperature
	if temperature == 0 {
		temperature = 0.3
	}

	return &OpenAILLM{
		client:      openai.NewClientWithConfig(config),
		model:       model,
		maxTokens:   maxTokens,
		temperature: temperature,
	}
}

// Enrich 为文档块生成结构化描述
func (l *OpenAILLM) Enrich(ctx context.Context, input EnrichInput) (*EnrichOutput, error) {
	prompt, err := l.renderEnrichPrompt(input)
	if err != nil {
		return nil, fmt.Errorf("render prompt failed: %w", err)
	}

	resp, err := l.client.CreateChatCompletion(ctx, openai.ChatCompletionRequest{
		Model: l.model,
		Messages: []openai.ChatCompletionMessage{
			{Role: openai.ChatMessageRoleUser, Content: prompt},
		},
		MaxTokens:   l.maxTokens,
		Temperature: l.temperature,
		ResponseFormat: &openai.ChatCompletionResponseFormat{
			Type: openai.ChatCompletionResponseFormatTypeJSONObject,
		},
	})
	if err != nil {
		return nil, fmt.Errorf("openai chat completion failed: %w", err)
	}

	if len(resp.Choices) == 0 {
		return nil, fmt.Errorf("no response from openai")
	}

	var output EnrichOutput
	if err := json.Unmarshal([]byte(resp.Choices[0].Message.Content), &output); err != nil {
		return nil, fmt.Errorf("parse response failed: %w", err)
	}

	return &output, nil
}

// Chat 通用对话
func (l *OpenAILLM) Chat(ctx context.Context, prompt string) (string, error) {
	resp, err := l.client.CreateChatCompletion(ctx, openai.ChatCompletionRequest{
		Model: l.model,
		Messages: []openai.ChatCompletionMessage{
			{Role: openai.ChatMessageRoleUser, Content: prompt},
		},
		MaxTokens:   l.maxTokens,
		Temperature: l.temperature,
	})
	if err != nil {
		return "", fmt.Errorf("openai chat completion failed: %w", err)
	}

	if len(resp.Choices) == 0 {
		return "", fmt.Errorf("no response from openai")
	}

	return resp.Choices[0].Message.Content, nil
}

// renderEnrichPrompt 渲染 Enrich 提示词
func (l *OpenAILLM) renderEnrichPrompt(input EnrichInput) (string, error) {
	tmpl, err := template.New("enrich").Parse(enrichPromptTemplate)
	if err != nil {
		return "", err
	}

	var buf bytes.Buffer
	if err := tmpl.Execute(&buf, input); err != nil {
		return "", err
	}

	return buf.String(), nil
}

// enrichPromptTemplate Enrich 提示词模板
const enrichPromptTemplate = `You are a technical documentation expert. Generate a title and description for the following code snippet.

## Code Snippet
---
{{.Content}}
---

## Context
- Section hierarchy: {{.Headers}}

## Return JSON:
{
  "title": "Concise title (5-15 words) describing the core functionality",
  "description": "Description (50-150 words): 1) What it does 2) When to use it 3) Key points"
}

## Rules
- title: Clear and concise, describe the core functionality of this code
- description: Explain what it does, when to use it, and key points
- **Keep the same language as the original document** (if the doc is in English, respond in English; if in Chinese, respond in Chinese)
- Return strict JSON only, no other content`

// GenerateLibraryTitle 为库生成简短友好的名称
func (l *OpenAILLM) GenerateLibraryTitle(ctx context.Context, repoName, description string) (string, error) {
	prompt := fmt.Sprintf(`Generate a short, friendly library name (1-3 words) for this GitHub repository.

Repository: %s
Description: %s

Rules:
- Output ONLY the library name, nothing else
- Keep it short: 1-3 words maximum
- Use proper capitalization (e.g., "Gin", "Next.js", "React Router")
- If the repo name is already a good title, just capitalize it properly
- Do NOT include "Library", "Framework", "SDK" unless it's part of the official name
- Examples: "gin-gonic/gin" → "Gin", "vercel/next.js" → "Next.js", "facebook/react" → "React", "gin-gonic/examples" → "Gin Examples"

Output:`, repoName, description)

	resp, err := l.client.CreateChatCompletion(ctx, openai.ChatCompletionRequest{
		Model: l.model,
		Messages: []openai.ChatCompletionMessage{
			{Role: openai.ChatMessageRoleUser, Content: prompt},
		},
		MaxTokens:   20,  // 只需要很短的输出
		Temperature: 0.1, // 低温度，更确定性
	})
	if err != nil {
		return "", fmt.Errorf("openai chat completion failed: %w", err)
	}

	if len(resp.Choices) == 0 {
		return "", fmt.Errorf("no response from openai")
	}

	// 清理输出（去除空白和引号）
	title := resp.Choices[0].Message.Content
	title = strings.TrimSpace(title)
	title = strings.Trim(title, `"'`)

	return title, nil
}

// 确保 OpenAILLM 实现了 LLMService 接口
var _ LLMService = (*OpenAILLM)(nil)
