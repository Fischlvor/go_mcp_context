package chunker

import (
	"regexp"
	"strings"
)

// TokenBasedChunker implements Chunker using token-based splitting
type TokenBasedChunker struct {
	chunkSize int
	overlap   int
}

// NewTokenBasedChunker creates a new token-based chunker
func NewTokenBasedChunker(chunkSize, overlap int) *TokenBasedChunker {
	return &TokenBasedChunker{
		chunkSize: chunkSize,
		overlap:   overlap,
	}
}

// Chunk splits text into chunks based on token count
func (c *TokenBasedChunker) Chunk(text string) []Chunk {
	if text == "" {
		return nil
	}

	// Split text into paragraphs first
	paragraphs := splitIntoParagraphs(text)

	var chunks []Chunk
	var currentChunk strings.Builder
	currentTokens := 0
	chunkIndex := 0

	for _, para := range paragraphs {
		paraTokens := estimateTokens(para)

		// If adding this paragraph exceeds chunk size
		if currentTokens+paraTokens > c.chunkSize && currentTokens > 0 {
			// Save current chunk
			chunkText := strings.TrimSpace(currentChunk.String())
			if chunkText != "" {
				chunks = append(chunks, Chunk{
					Index:     chunkIndex,
					Text:      chunkText,
					Tokens:    currentTokens,
					ChunkType: detectChunkType(chunkText),
				})
				chunkIndex++
			}

			// Start new chunk with overlap
			currentChunk.Reset()
			currentTokens = 0

			// Add overlap from previous content if needed
			if c.overlap > 0 && len(chunks) > 0 {
				overlapText := getOverlapText(chunks[len(chunks)-1].Text, c.overlap)
				if overlapText != "" {
					currentChunk.WriteString(overlapText)
					currentChunk.WriteString("\n\n")
					currentTokens = estimateTokens(overlapText)
				}
			}
		}

		// Add paragraph to current chunk
		if currentChunk.Len() > 0 {
			currentChunk.WriteString("\n\n")
		}
		currentChunk.WriteString(para)
		currentTokens += paraTokens
	}

	// Don't forget the last chunk
	if currentChunk.Len() > 0 {
		chunkText := strings.TrimSpace(currentChunk.String())
		if chunkText != "" {
			chunks = append(chunks, Chunk{
				Index:     chunkIndex,
				Text:      chunkText,
				Tokens:    currentTokens,
				ChunkType: detectChunkType(chunkText),
			})
		}
	}

	return chunks
}

// GetChunkSize returns the chunk size
func (c *TokenBasedChunker) GetChunkSize() int {
	return c.chunkSize
}

// GetOverlap returns the overlap size
func (c *TokenBasedChunker) GetOverlap() int {
	return c.overlap
}

// splitIntoParagraphs splits text into paragraphs
func splitIntoParagraphs(text string) []string {
	// Split by double newlines or code blocks
	parts := regexp.MustCompile(`\n\s*\n`).Split(text, -1)

	var paragraphs []string
	for _, part := range parts {
		part = strings.TrimSpace(part)
		if part != "" {
			paragraphs = append(paragraphs, part)
		}
	}

	return paragraphs
}

// estimateTokens estimates the number of tokens in text
// Using a simple heuristic: ~4 characters per token for English
func estimateTokens(text string) int {
	// Count words and special characters
	words := len(strings.Fields(text))
	// Rough estimate: 1.3 tokens per word on average
	return int(float64(words) * 1.3)
}

// getOverlapText gets the last N tokens worth of text
func getOverlapText(text string, tokens int) string {
	words := strings.Fields(text)
	// Convert tokens to approximate word count
	wordCount := int(float64(tokens) / 1.3)

	if len(words) <= wordCount {
		return text
	}

	return strings.Join(words[len(words)-wordCount:], " ")
}

// detectChunkType detects if a chunk is code, info, or mixed
func detectChunkType(text string) string {
	hasCodeBlock := strings.Contains(text, "```")
	hasInlineCode := regexp.MustCompile("`[^`]+`").MatchString(text)

	// Check for API patterns
	apiPatterns := []string{
		"func ", "function ", "def ", "class ",
		"import ", "package ", "module ",
		"GET ", "POST ", "PUT ", "DELETE ",
		"@param", "@return", "@throws",
	}

	hasAPIPattern := false
	for _, pattern := range apiPatterns {
		if strings.Contains(text, pattern) {
			hasAPIPattern = true
			break
		}
	}

	if hasCodeBlock || hasAPIPattern {
		return "code"
	}

	if hasInlineCode {
		return "mixed"
	}

	return "info"
}
