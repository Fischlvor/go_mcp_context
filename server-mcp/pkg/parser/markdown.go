package parser

import (
	"os"
	"path/filepath"
	"strings"

	"github.com/gomarkdown/markdown"
	"github.com/gomarkdown/markdown/html"
	mdparser "github.com/gomarkdown/markdown/parser"
)

// MarkdownParser implements DocumentParser for Markdown files
type MarkdownParser struct{}

// NewMarkdownParser creates a new MarkdownParser
func NewMarkdownParser() *MarkdownParser {
	return &MarkdownParser{}
}

// Parse extracts text content from a Markdown file
func (p *MarkdownParser) Parse(filePath string) (string, error) {
	data, err := os.ReadFile(filePath)
	if err != nil {
		return "", err
	}
	return p.ParseBytes(data)
}

// ParseBytes extracts text content from Markdown bytes
func (p *MarkdownParser) ParseBytes(data []byte) (string, error) {
	// Create markdown parser with extensions
	extensions := mdparser.CommonExtensions | mdparser.AutoHeadingIDs | mdparser.NoEmptyLineBeforeBlock
	parser := mdparser.NewWithExtensions(extensions)

	// Parse markdown to AST
	doc := parser.Parse(data)

	// Render to HTML (we'll strip tags to get plain text)
	htmlFlags := html.CommonFlags | html.HrefTargetBlank
	opts := html.RendererOptions{Flags: htmlFlags}
	renderer := html.NewRenderer(opts)

	htmlContent := markdown.Render(doc, renderer)

	// Strip HTML tags to get plain text
	text := stripHTMLTags(string(htmlContent))

	return strings.TrimSpace(text), nil
}

// GetFormat returns the format type
func (p *MarkdownParser) GetFormat() string {
	return "markdown"
}

// SupportedExtensions returns supported file extensions
func (p *MarkdownParser) SupportedExtensions() []string {
	return []string{".md", ".markdown", ".mdown", ".mkd"}
}

// CanParse checks if this parser can handle the given file
func (p *MarkdownParser) CanParse(filePath string) bool {
	ext := strings.ToLower(filepath.Ext(filePath))
	for _, supported := range p.SupportedExtensions() {
		if ext == supported {
			return true
		}
	}
	return false
}

// stripHTMLTags removes HTML tags from a string
func stripHTMLTags(s string) string {
	var result strings.Builder
	inTag := false

	for _, r := range s {
		switch {
		case r == '<':
			inTag = true
		case r == '>':
			inTag = false
			result.WriteRune(' ')
		case !inTag:
			result.WriteRune(r)
		}
	}

	// Clean up multiple spaces
	text := result.String()
	for strings.Contains(text, "  ") {
		text = strings.ReplaceAll(text, "  ", " ")
	}

	return text
}
