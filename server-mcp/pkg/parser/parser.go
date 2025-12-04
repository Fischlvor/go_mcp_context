package parser

import "errors"

// ErrUnsupportedFormat is returned when the document format is not supported
var ErrUnsupportedFormat = errors.New("unsupported document format")

// DocumentParser defines the interface for document parsing
type DocumentParser interface {
	// Parse extracts text content from a file
	Parse(filePath string) (string, error)

	// ParseBytes extracts text content from bytes
	ParseBytes(data []byte) (string, error)

	// GetFormat returns the format type this parser handles
	GetFormat() string

	// SupportedExtensions returns the file extensions this parser supports
	SupportedExtensions() []string

	// CanParse checks if this parser can handle the given file
	CanParse(filePath string) bool
}

// ParseResult contains the result of parsing a document
type ParseResult struct {
	Text     string            // Extracted text content
	Metadata map[string]string // Document metadata
	Format   string            // Document format
}
