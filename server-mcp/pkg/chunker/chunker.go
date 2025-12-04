package chunker

// Chunker defines the interface for text chunking
type Chunker interface {
	// Chunk splits text into chunks
	Chunk(text string) []Chunk

	// GetChunkSize returns the chunk size in tokens
	GetChunkSize() int

	// GetOverlap returns the overlap size in tokens
	GetOverlap() int
}

// Chunk represents a text chunk
type Chunk struct {
	Index     int               // Position in the original document
	Text      string            // Chunk text content
	Tokens    int               // Estimated token count
	ChunkType string            // code, info, or mixed
	Metadata  map[string]string // Additional metadata
}
