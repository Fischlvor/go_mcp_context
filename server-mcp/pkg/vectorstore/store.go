package vectorstore

import "errors"

// ErrNotFound is returned when the requested item is not found
var ErrNotFound = errors.New("not found")

// VectorStore defines the interface for vector storage and search
type VectorStore interface {
	// Store stores chunks with their embeddings
	Store(libraryID uint, chunks []ChunkData) error

	// Search performs vector similarity search
	Search(query []float32, libraryID uint, limit int) ([]SearchResult, error)

	// HybridSearch performs combined vector and keyword search
	HybridSearch(query []float32, keywords string, libraryID uint, limit int) ([]SearchResult, error)

	// Delete removes all chunks for a library
	Delete(libraryID uint) error

	// DeleteByDocument removes all chunks for a document
	DeleteByDocument(documentID uint) error

	// UpdateAccessCount increments the access count for chunks
	UpdateAccessCount(chunkIDs []uint) error

	// GetStats returns storage statistics
	GetStats(libraryID uint) (*StoreStats, error)
}

// ChunkData represents a chunk to be stored
type ChunkData struct {
	DocumentID uint
	LibraryID  uint
	ChunkIndex int
	ChunkText  string
	Tokens     int
	Embedding  []float32
	ChunkType  string
	Metadata   map[string]interface{}
}

// SearchResult represents a search result
type SearchResult struct {
	ChunkID     uint
	DocumentID  uint
	LibraryID   uint
	ChunkText   string
	ChunkType   string
	Score       float64 // Similarity score (0-1)
	BM25Score   float64 // BM25 score for hybrid search
	AccessCount int
	Metadata    map[string]interface{}
}

// StoreStats represents storage statistics
type StoreStats struct {
	TotalChunks    int64
	TotalTokens    int64
	TotalDocuments int64
}
