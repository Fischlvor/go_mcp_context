package embedding

import "errors"

// ErrEmptyInput is returned when the input text is empty
var ErrEmptyInput = errors.New("input text is empty")

// ErrBatchTooLarge is returned when the batch size exceeds the limit
var ErrBatchTooLarge = errors.New("batch size exceeds limit")

// EmbeddingService defines the interface for embedding generation
type EmbeddingService interface {
	// Embed generates embeddings for a single text
	Embed(text string) ([]float32, error)

	// EmbedBatch generates embeddings for multiple texts
	EmbedBatch(texts []string) ([][]float32, error)

	// GetDimension returns the embedding dimension
	GetDimension() int

	// GetModelName returns the model name
	GetModelName() string

	// GetMaxBatchSize returns the maximum batch size
	GetMaxBatchSize() int
}

// EmbeddingResult contains the result of embedding generation
type EmbeddingResult struct {
	Embedding []float32
	Tokens    int
}
