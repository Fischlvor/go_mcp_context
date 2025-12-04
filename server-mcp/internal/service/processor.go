package service

import (
	"fmt"
	"strings"
	"unicode/utf8"

	dbmodel "go-mcp-context/internal/model/database"
	"go-mcp-context/pkg/global"

	"github.com/pgvector/pgvector-go"
)

// DocumentProcessor 文档处理器
type DocumentProcessor struct{}

// ProcessDocument 处理文档（解析、分块、生成 Embedding）
func (p *DocumentProcessor) ProcessDocument(doc *dbmodel.Document, content []byte) error {
	// 1. 解析文档内容
	text, err := p.parseDocument(doc.FileType, content)
	if err != nil {
		return fmt.Errorf("failed to parse document: %w", err)
	}

	// 2. 分块
	chunks := p.chunkText(text, doc.ID, doc.LibraryID)
	if len(chunks) == 0 {
		return nil
	}

	// 3. 批量生成 Embedding
	texts := make([]string, len(chunks))
	for i, chunk := range chunks {
		texts[i] = chunk.ChunkText
	}

	embeddings, err := global.Embedding.EmbedBatch(texts)
	if err != nil {
		return fmt.Errorf("failed to generate embeddings: %w", err)
	}

	// 4. 设置 Embedding 并存储
	for i := range chunks {
		if i < len(embeddings) {
			chunks[i].Embedding = pgvector.NewVector(embeddings[i])
		}
	}

	// 批量插入
	if err := global.DB.CreateInBatches(chunks, 100).Error; err != nil {
		return fmt.Errorf("failed to save chunks: %w", err)
	}

	// 5. 更新文档状态
	doc.Status = "processed"
	if err := global.DB.Save(doc).Error; err != nil {
		return fmt.Errorf("failed to update document status: %w", err)
	}

	return nil
}

// parseDocument 解析文档内容
func (p *DocumentProcessor) parseDocument(fileType string, content []byte) (string, error) {
	switch fileType {
	case "markdown":
		return string(content), nil
	case "pdf":
		// TODO: 使用 PDF 解析库
		return string(content), nil
	case "docx":
		// TODO: 使用 DOCX 解析库
		return string(content), nil
	case "swagger":
		// TODO: 解析 Swagger/OpenAPI
		return string(content), nil
	default:
		return string(content), nil
	}
}

// chunkText 文本分块
func (p *DocumentProcessor) chunkText(text string, documentID, libraryID uint) []*dbmodel.DocumentChunk {
	chunkSize := global.Config.Chunker.ChunkSize
	overlap := global.Config.Chunker.Overlap

	if chunkSize <= 0 {
		chunkSize = 512
	}
	if overlap <= 0 {
		overlap = 50
	}

	// 按段落分割
	paragraphs := strings.Split(text, "\n\n")
	var chunks []*dbmodel.DocumentChunk
	var currentChunk strings.Builder
	var currentTokens int
	chunkIndex := 0

	for _, para := range paragraphs {
		para = strings.TrimSpace(para)
		if para == "" {
			continue
		}

		paraTokens := p.estimateTokens(para)

		// 如果当前段落加上已有内容超过限制，先保存当前块
		if currentTokens+paraTokens > chunkSize && currentTokens > 0 {
			chunk := p.createChunk(currentChunk.String(), chunkIndex, documentID, libraryID, currentTokens)
			chunks = append(chunks, chunk)
			chunkIndex++

			// 保留 overlap 部分
			overlapText := p.getOverlapText(currentChunk.String(), overlap)
			currentChunk.Reset()
			currentChunk.WriteString(overlapText)
			currentTokens = p.estimateTokens(overlapText)
		}

		if currentChunk.Len() > 0 {
			currentChunk.WriteString("\n\n")
		}
		currentChunk.WriteString(para)
		currentTokens += paraTokens
	}

	// 保存最后一个块
	if currentChunk.Len() > 0 {
		chunk := p.createChunk(currentChunk.String(), chunkIndex, documentID, libraryID, currentTokens)
		chunks = append(chunks, chunk)
	}

	return chunks
}

// createChunk 创建文档块
func (p *DocumentProcessor) createChunk(text string, index int, documentID, libraryID uint, tokens int) *dbmodel.DocumentChunk {
	chunkType := p.detectChunkType(text)

	return &dbmodel.DocumentChunk{
		DocumentID: documentID,
		LibraryID:  libraryID,
		ChunkIndex: index,
		ChunkText:  text,
		Tokens:     tokens,
		ChunkType:  chunkType,
		Status:     "active",
	}
}

// detectChunkType 检测块类型
func (p *DocumentProcessor) detectChunkType(text string) string {
	// 检测是否包含代码块
	hasCode := strings.Contains(text, "```") ||
		strings.Contains(text, "    ") || // 缩进代码
		strings.Contains(text, "func ") ||
		strings.Contains(text, "function ") ||
		strings.Contains(text, "class ") ||
		strings.Contains(text, "def ")

	// 检测是否主要是说明文字
	hasInfo := strings.Contains(text, "# ") ||
		strings.Contains(text, "## ") ||
		len(text) > 200 && !hasCode

	if hasCode && hasInfo {
		return "mixed"
	} else if hasCode {
		return "code"
	}
	return "info"
}

// estimateTokens 估算 token 数（简单实现：字符数/4）
func (p *DocumentProcessor) estimateTokens(text string) int {
	return utf8.RuneCountInString(text) / 4
}

// getOverlapText 获取重叠文本
func (p *DocumentProcessor) getOverlapText(text string, overlapTokens int) string {
	runes := []rune(text)
	overlapChars := overlapTokens * 4
	if overlapChars >= len(runes) {
		return text
	}
	return string(runes[len(runes)-overlapChars:])
}

// ProcessDocumentAsync 异步处理文档
func (p *DocumentProcessor) ProcessDocumentAsync(doc *dbmodel.Document, content []byte) {
	go func() {
		if err := p.ProcessDocument(doc, content); err != nil {
			// 更新文档状态为失败
			doc.Status = "failed"
			global.DB.Save(doc)
			// TODO: 记录错误日志
		}
	}()
}
