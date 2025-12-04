package embedding

import (
	"context"

	"github.com/sashabaranov/go-openai"
)

// OpenAIEmbedding implements EmbeddingService using OpenAI API
type OpenAIEmbedding struct {
	client    *openai.Client
	model     openai.EmbeddingModel
	dimension int
}

// NewOpenAIEmbedding creates a new OpenAI embedding service
func NewOpenAIEmbedding(apiKey string, model string, dimension int) *OpenAIEmbedding {
	client := openai.NewClient(apiKey)

	var embeddingModel openai.EmbeddingModel
	switch model {
	case "text-embedding-3-large":
		embeddingModel = openai.LargeEmbedding3
	default:
		embeddingModel = openai.SmallEmbedding3
	}

	return &OpenAIEmbedding{
		client:    client,
		model:     embeddingModel,
		dimension: dimension,
	}
}

// Embed generates embedding for a single text
func (e *OpenAIEmbedding) Embed(text string) ([]float32, error) {
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
func (e *OpenAIEmbedding) EmbedBatch(texts []string) ([][]float32, error) {
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
func (e *OpenAIEmbedding) GetDimension() int {
	return e.dimension
}

// GetModelName returns the model name
func (e *OpenAIEmbedding) GetModelName() string {
	return string(e.model)
}

// GetMaxBatchSize returns the maximum batch size
func (e *OpenAIEmbedding) GetMaxBatchSize() int {
	return 2048 // OpenAI limit
}
