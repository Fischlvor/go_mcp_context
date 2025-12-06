package embedding

import (
	"context"
	"net/http"

	"github.com/sashabaranov/go-openai"
)

// OpenAIProxyEmbedding 第三方 OpenAI 兼容 API 的 Embedding 服务
type OpenAIProxyEmbedding struct {
	client    *openai.Client
	model     openai.EmbeddingModel
	dimension int
}

// NewOpenAIProxyEmbedding 创建第三方代理 Embedding 服务
func NewOpenAIProxyEmbedding(apiKey string, baseURL string, model string, dimension int) *OpenAIProxyEmbedding {
	config := openai.DefaultConfig(apiKey)
	config.BaseURL = baseURL
	// 不走系统代理，直连
	config.HTTPClient = &http.Client{
		Transport: &http.Transport{
			Proxy: nil,
		},
	}
	client := openai.NewClientWithConfig(config)

	var embeddingModel openai.EmbeddingModel
	switch model {
	case "text-embedding-3-large":
		embeddingModel = openai.LargeEmbedding3
	default:
		embeddingModel = openai.SmallEmbedding3
	}

	return &OpenAIProxyEmbedding{
		client:    client,
		model:     embeddingModel,
		dimension: dimension,
	}
}

// Embed generates embedding for a single text
func (e *OpenAIProxyEmbedding) Embed(text string) ([]float32, error) {
	if text == "" {
		return nil, ErrEmptyInput
	}

	resp, err := e.client.CreateEmbeddings(context.Background(), openai.EmbeddingRequest{
		Input: []string{text},
		Model: e.model,
	})
	if err != nil {
		return nil, err
	}

	if len(resp.Data) == 0 {
		return nil, ErrEmptyInput
	}

	return resp.Data[0].Embedding, nil
}

// EmbedBatch generates embeddings for multiple texts
func (e *OpenAIProxyEmbedding) EmbedBatch(texts []string) ([][]float32, error) {
	if len(texts) == 0 {
		return nil, ErrEmptyInput
	}

	if len(texts) > e.GetMaxBatchSize() {
		return nil, ErrBatchTooLarge
	}

	resp, err := e.client.CreateEmbeddings(context.Background(), openai.EmbeddingRequest{
		Input: texts,
		Model: e.model,
	})
	if err != nil {
		return nil, err
	}

	embeddings := make([][]float32, len(resp.Data))
	for i, data := range resp.Data {
		embeddings[i] = data.Embedding
	}

	return embeddings, nil
}

// GetDimension returns the embedding dimension
func (e *OpenAIProxyEmbedding) GetDimension() int {
	return e.dimension
}

// GetModelName returns the model name
func (e *OpenAIProxyEmbedding) GetModelName() string {
	return string(e.model)
}

// GetMaxBatchSize returns the maximum batch size
func (e *OpenAIProxyEmbedding) GetMaxBatchSize() int {
	return 2048
}
