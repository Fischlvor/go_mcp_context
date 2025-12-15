package llm

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
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
const enrichPromptTemplate = `你是一个技术文档专家，负责为代码片段生成结构化的描述信息。

请分析以下代码片段，生成符合 Context7 格式的结构化输出。

## 输入信息

**文件来源：** {{.Source}}
**代码语言：** {{.Language}}
**原始内容：**
{{.Content}}

## 输出要求

请返回 JSON 格式，包含以下字段：

{
  "title": "简洁的标题，描述这段代码做什么（英文，使用动词开头，如 'Create...', 'Configure...', 'Handle...'）",
  "description": "1-3 句话的描述，说明代码的用途和使用场景（英文）",
  "content_type": "code 或 info（code 表示包含可执行代码示例，info 表示概念说明或配置指南）",
  "language": "代码语言（如 js, ts, go, python, markdown 等）"
}

## 生成规则

### Title 规则：
- 使用动词开头：Create, Configure, Handle, Implement, Define, Set up, etc.
- 简洁明了，不超过 80 个字符
- 如果是 API 文档，包含函数/方法名
- 如果是配置，说明配置什么

### Description 规则：
- 第一句说明代码做什么
- 第二句说明什么场景使用
- 第三句（可选）说明注意事项或最佳实践
- 使用第三人称描述（This code..., This configuration..., This example...）
- 不超过 200 个字符

### Content Type 规则：
- code：包含可执行的代码示例（函数、类、配置代码等）
- info：概念说明、最佳实践、架构指南、无代码的文档

请严格按照 JSON 格式返回，不要包含其他内容。`

// 确保 OpenAILLM 实现了 LLMService 接口
var _ LLMService = (*OpenAILLM)(nil)
