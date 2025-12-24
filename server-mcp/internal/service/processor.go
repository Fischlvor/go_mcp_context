package service

import (
	"context"
	"fmt"
	"log"
	"regexp"
	"strings"
	"sync"

	dbmodel "go-mcp-context/internal/model/database"
	"go-mcp-context/internal/model/response"
	"go-mcp-context/pkg/bufferedwriter/actlog"
	"go-mcp-context/pkg/global"
	"go-mcp-context/pkg/llm"

	"github.com/pgvector/pgvector-go"
	"github.com/pkoukk/tiktoken-go"
)

// 分块配置常量
const (
	DefaultChunkSize   = 512 // 默认块大小（tokens）
	MinChunkSize       = 50  // 最小块大小（tokens），低于此值合并
	MaxEnrichBatchSize = 10  // LLM Enrich 批处理大小
)

// DocumentProcessor 文档处理器
type DocumentProcessor struct{}

// processDocumentCore 文档处理核心逻辑（解析、分块、Enrich、生成 Embedding）
// 返回处理好的 chunks 和总 token 数，不写入数据库
func (p *DocumentProcessor) processDocumentCore(doc *dbmodel.DocumentUpload, content []byte, actLogger *actlog.TaskLogger) ([]*dbmodel.DocumentChunk, int, error) {
	// 1. 解析文档内容
	text, err := p.parseDocument(doc.FileType, content)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to parse document: %w", err)
	}

	// 2. Pre-Chunking: 清理文档（移除徽章、空白等）
	text = p.preProcessMarkdown(text)

	// 3. 分块
	chunks := p.chunkText(text, doc.ID, doc.LibraryID, doc.Version, doc.FilePath)
	if len(chunks) == 0 {
		return nil, 0, nil
	}
	actLogger.Info(actlog.EventDocChunk, fmt.Sprintf("分块完成: %s (%d 块)", doc.Title, len(chunks)))

	// 4. LLM Enrich: 生成 Title 和 Description
	if global.LLM != nil {
		actLogger.Info(actlog.EventDocEnrich, fmt.Sprintf("AI 增强中: %s (%d 块)", doc.Title, len(chunks)))
		if err := p.enrichChunks(chunks); err != nil {
			log.Printf("[Processor] WARNING: LLM enrich failed: %v, using fallback", err)
		}
	}

	// 5. 批量生成 Embedding
	actLogger.Info(actlog.EventDocEmbed, fmt.Sprintf("生成 Embedding: %s (%d 块)", doc.Title, len(chunks)))
	texts := make([]string, len(chunks))
	for i, chunk := range chunks {
		texts[i] = chunk.ChunkText
	}

	embeddings, err := global.Embedding.EmbedBatch(texts)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to generate embeddings: %w", err)
	}

	// 6. 设置 Embedding
	totalTokens := 0
	for i := range chunks {
		if i < len(embeddings) {
			chunks[i].Embedding = pgvector.NewVector(embeddings[i])
		}
		totalTokens += chunks[i].Tokens
	}

	return chunks, totalTokens, nil
}

// ProcessDocument 处理文档（带任务日志器）
// actLogger 应已设置好 document target
func (p *DocumentProcessor) ProcessDocument(doc *dbmodel.DocumentUpload, content []byte, actLogger *actlog.TaskLogger) error {
	actLogger.Info(actlog.EventDocParse, fmt.Sprintf("处理文档: %s", doc.Title))

	// 调用核心处理逻辑（传入 logger 记录详细步骤）
	chunks, totalTokens, err := p.processDocumentCore(doc, content, actLogger)
	if err != nil {
		actLogger.Error(actlog.EventDocFailed, fmt.Sprintf("处理失败: %s - %s", doc.Title, err.Error()))
		return err
	}
	if len(chunks) == 0 {
		return nil
	}

	// 批量插入
	if err := global.DB.CreateInBatches(chunks, 100).Error; err != nil {
		actLogger.Error(actlog.EventDocFailed, fmt.Sprintf("保存失败: %s - %s", doc.Title, err.Error()))
		return fmt.Errorf("failed to save chunks: %w", err)
	}

	// 更新文档状态和统计信息
	doc.Status = "completed"
	doc.ChunkCount = len(chunks)
	doc.TokenCount = totalTokens
	if err := global.DB.Save(doc).Error; err != nil {
		return fmt.Errorf("failed to update document status: %w", err)
	}

	// 失效该库版本的搜索缓存（Version Tag 模式）
	searchService := &SearchService{}
	if err := searchService.InvalidateLibraryCache(doc.LibraryID, doc.Version); err != nil {
		log.Printf("[Processor] WARNING: Failed to invalidate cache for library %d version %s: %v", doc.LibraryID, doc.Version, err)
	}

	actLogger.Info(actlog.EventDocComplete, fmt.Sprintf("处理完成: %s (%d 块, %d tokens)", doc.Title, len(chunks), totalTokens))

	return nil
}

// ProcessDocumentForRefresh 处理文档用于刷新（无感知更新）
// 返回生成的 chunks，不直接写入数据库，由调用方控制事务
// batchVersion: 批次版本号，用于原子切换
func (p *DocumentProcessor) ProcessDocumentForRefresh(doc *dbmodel.DocumentUpload, content []byte, batchVersion int64, actLogger *actlog.TaskLogger) ([]*dbmodel.DocumentChunk, int, error) {
	actLogger.Info(actlog.EventDocParse, fmt.Sprintf("处理文档: %s", doc.Title))
	chunks, totalTokens, err := p.processDocumentCore(doc, content, actLogger)
	if err != nil {
		return nil, 0, err
	}
	if len(chunks) == 0 {
		return nil, 0, nil
	}

	// 设置 BatchVersion 和 pending 状态
	for i := range chunks {
		chunks[i].BatchVersion = batchVersion
		chunks[i].Status = "pending" // 待激活状态
	}

	return chunks, totalTokens, nil
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

// MarkdownSection 带元数据的 Markdown 段落
type MarkdownSection struct {
	Content string            // 段落内容
	Headers map[string]string // 标题层级 {"h1": "Title", "h2": "Section", "h3": "Subsection"}
}

// chunkText 文本分块（Markdown 语义分块，带标题元数据）
// 参考 LangChain MarkdownHeaderTextSplitter 的设计
func (p *DocumentProcessor) chunkText(text string, uploadID, libraryID uint, version, source string) []*dbmodel.DocumentChunk {
	chunkSize := global.Config.Chunker.ChunkSize
	if chunkSize <= 0 {
		chunkSize = 512
	}

	// 按 Markdown 标题分割成 sections（带元数据）
	sections := p.splitMarkdownWithMetadata(text)

	var chunks []*dbmodel.DocumentChunk
	chunkIndex := 0

	for _, section := range sections {
		content := strings.TrimSpace(section.Content)
		if content == "" {
			continue
		}

		sectionTokens := p.countTokens(content)

		// 如果 section 小于 chunkSize，直接作为一个 chunk
		if sectionTokens <= chunkSize {
			chunk := p.createChunkWithMetadata(content, chunkIndex, uploadID, libraryID, version, source, sectionTokens, section.Headers)
			chunks = append(chunks, chunk)
			chunkIndex++
			continue
		}

		// section 超过 chunkSize，按段落分割（继承标题元数据）
		subChunks := p.splitLargeSectionWithMetadata(content, chunkSize, chunkIndex, uploadID, libraryID, version, source, section.Headers)
		chunks = append(chunks, subChunks...)
		chunkIndex += len(subChunks)
	}

	log.Printf("[Chunker] Created %d chunks from document", len(chunks))
	return chunks
}

// splitMarkdownWithMetadata 按 Markdown 标题分割，并提取标题层级元数据
// 注意：代码块内的 # 不是标题，需要排除
func (p *DocumentProcessor) splitMarkdownWithMetadata(text string) []MarkdownSection {
	var sections []MarkdownSection

	// 先标记代码块位置，避免把代码块内的 # 当作标题
	codeBlockRanges := p.findCodeBlockRanges(text)

	// 正则匹配标题行
	headingRegex := regexp.MustCompile(`(?m)^(#{1,6})\s+(.*)$`)

	// 当前标题层级状态
	currentHeaders := make(map[string]string)

	// 找到所有标题（排除代码块内的）
	allMatches := headingRegex.FindAllStringSubmatchIndex(text, -1)
	var matches [][]int
	for _, match := range allMatches {
		if !p.isInCodeBlock(match[0], codeBlockRanges) {
			matches = append(matches, match)
		}
	}

	if len(matches) == 0 {
		// 没有标题，整个文档作为一个 section
		return []MarkdownSection{{Content: text, Headers: make(map[string]string)}}
	}

	// 第一个标题前的内容
	if matches[0][0] > 0 {
		before := strings.TrimSpace(text[:matches[0][0]])
		if before != "" {
			sections = append(sections, MarkdownSection{
				Content: before,
				Headers: copyHeaders(currentHeaders),
			})
		}
	}

	// 处理每个标题及其内容
	for i, match := range matches {
		// 提取标题级别和文本
		level := len(text[match[2]:match[3]]) // # 的数量
		title := strings.TrimSpace(text[match[4]:match[5]])

		// 更新标题层级（清除更低级别的标题）
		headerKey := fmt.Sprintf("h%d", level)
		currentHeaders[headerKey] = title
		for l := level + 1; l <= 6; l++ {
			delete(currentHeaders, fmt.Sprintf("h%d", l))
		}

		// 获取内容范围（从标题行之后开始，不包含标题行本身）
		// match[1] 是标题行的结尾位置
		contentStart := match[1]
		// 跳过标题行后的换行符
		for contentStart < len(text) && (text[contentStart] == '\n' || text[contentStart] == '\r') {
			contentStart++
		}

		var contentEnd int
		if i+1 < len(matches) {
			contentEnd = matches[i+1][0]
		} else {
			contentEnd = len(text)
		}

		// 提取内容（不含标题行）
		content := ""
		if contentStart < contentEnd {
			content = strings.TrimSpace(text[contentStart:contentEnd])
		}

		// 只有当内容不为空时才创建 section
		// 空标题（只有标题没有内容）的 headers 会传递给下一个有内容的 section
		if content != "" {
			sections = append(sections, MarkdownSection{
				Content: content,
				Headers: copyHeaders(currentHeaders),
			})
		}
	}

	return sections
}

// findCodeBlockRanges 找到所有代码块的位置范围
func (p *DocumentProcessor) findCodeBlockRanges(text string) [][2]int {
	var ranges [][2]int
	codeBlockRegex := regexp.MustCompile("(?s)```.*?```")
	matches := codeBlockRegex.FindAllStringIndex(text, -1)
	for _, match := range matches {
		ranges = append(ranges, [2]int{match[0], match[1]})
	}
	return ranges
}

// isInCodeBlock 检查位置是否在代码块内
func (p *DocumentProcessor) isInCodeBlock(pos int, ranges [][2]int) bool {
	for _, r := range ranges {
		if pos >= r[0] && pos < r[1] {
			return true
		}
	}
	return false
}

// copyHeaders 复制标题 map
func copyHeaders(headers map[string]string) map[string]string {
	result := make(map[string]string)
	for k, v := range headers {
		result[k] = v
	}
	return result
}

// createChunkWithMetadata 创建带元数据的文档块
// Code Mode: 设置 Language、Code，Metadata 在 enrichChunks 中清空
// Info Mode: Language、Code、Description 为空，Metadata 存 headers 层级
func (p *DocumentProcessor) createChunkWithMetadata(text string, index int, uploadID, libraryID uint, version, source string, tokens int, headers map[string]string) *dbmodel.DocumentChunk {
	chunkType := p.detectChunkType(text)

	// 生成标题（使用标题层级构建）
	title := p.buildTitleFromHeaders(headers)

	chunk := &dbmodel.DocumentChunk{
		LibraryID:  libraryID,
		UploadID:   uploadID,
		Version:    version,
		ChunkIndex: index,
		Title:      title,
		Source:     source,
		ChunkText:  text,
		Tokens:     tokens,
		ChunkType:  chunkType,
		Status:     "active",
	}

	if chunkType == "code" {
		// Code Mode: 设置 Language 和 Code，Metadata 为空
		chunk.Language = p.detectLanguage(text)
		chunk.Code = p.extractCodeBlock(text)
		// Metadata 不需要存储，LLM 用 Title（headers 层级格式）作为上下文
	} else {
		// Info Mode: Language、Code、Description 为空
		// Metadata 存 headers 层级（用于调试和追溯）
		metadata := make(dbmodel.JSON)
		for k, v := range headers {
			metadata[k] = v
		}
		chunk.Metadata = metadata
	}

	return chunk
}

// splitLargeSectionWithMetadata 分割超大 section（保持代码块完整）
func (p *DocumentProcessor) splitLargeSectionWithMetadata(text string, chunkSize int, startIndex int, uploadID, libraryID uint, version, source string, headers map[string]string) []*dbmodel.DocumentChunk {
	var chunks []*dbmodel.DocumentChunk

	// 先拆成原子单元（代码块作为整体，其他按段落）
	atoms := p.splitIntoAtoms(text)

	var currentChunk strings.Builder
	var currentTokens int
	chunkIndex := startIndex

	for _, atom := range atoms {
		atom = strings.TrimSpace(atom)
		if atom == "" {
			continue
		}

		atomTokens := p.countTokens(atom)

		// 如果单个原子块就超过 chunkSize，单独作为一个 chunk（不再切分）
		if atomTokens > chunkSize {
			// 先保存当前累积的内容
			if currentChunk.Len() > 0 {
				chunk := p.createChunkWithMetadata(currentChunk.String(), chunkIndex, uploadID, libraryID, version, source, currentTokens, headers)
				chunks = append(chunks, chunk)
				chunkIndex++
				currentChunk.Reset()
				currentTokens = 0
			}
			// 大原子块单独成 chunk
			chunk := p.createChunkWithMetadata(atom, chunkIndex, uploadID, libraryID, version, source, atomTokens, headers)
			chunks = append(chunks, chunk)
			chunkIndex++
			continue
		}

		// 正常累积
		if currentTokens+atomTokens > chunkSize && currentTokens > 0 {
			chunk := p.createChunkWithMetadata(currentChunk.String(), chunkIndex, uploadID, libraryID, version, source, currentTokens, headers)
			chunks = append(chunks, chunk)
			chunkIndex++
			currentChunk.Reset()
			currentTokens = 0
		}

		if currentChunk.Len() > 0 {
			currentChunk.WriteString("\n\n")
		}
		currentChunk.WriteString(atom)
		currentTokens += atomTokens
	}

	if currentChunk.Len() > 0 {
		chunk := p.createChunkWithMetadata(currentChunk.String(), chunkIndex, uploadID, libraryID, version, source, currentTokens, headers)
		chunks = append(chunks, chunk)
	}

	return chunks
}

// extractCodeBlock 提取代码块内容
func (p *DocumentProcessor) extractCodeBlock(text string) string {
	codeBlockRegex := regexp.MustCompile("(?s)```(?:\\w+)?\\n?(.*?)```")
	matches := codeBlockRegex.FindAllStringSubmatch(text, -1)
	if len(matches) == 0 {
		return ""
	}
	// 返回第一个代码块的内容
	return strings.TrimSpace(matches[0][1])
}

// detectLanguage 检测代码语言
func (p *DocumentProcessor) detectLanguage(text string) string {
	// 从代码块标记中提取语言
	langRegex := regexp.MustCompile("```(\\w+)")
	matches := langRegex.FindStringSubmatch(text)
	if len(matches) > 1 {
		return matches[1]
	}
	return "markdown"
}

// buildTitleFromHeaders 从标题层级构建标题
func (p *DocumentProcessor) buildTitleFromHeaders(headers map[string]string) string {
	var parts []string
	for i := 1; i <= 6; i++ {
		key := fmt.Sprintf("h%d", i)
		if v, ok := headers[key]; ok && v != "" {
			parts = append(parts, v)
		}
	}
	if len(parts) == 0 {
		return ""
	}
	return strings.Join(parts, " > ")
}

// splitIntoAtoms 将文本拆成原子单元（代码块完整保留，其他按段落分）
func (p *DocumentProcessor) splitIntoAtoms(text string) []string {
	var atoms []string

	// 找到所有代码块
	codeBlockRegex := regexp.MustCompile("(?s)```.*?```")
	codeBlocks := codeBlockRegex.FindAllStringIndex(text, -1)

	if len(codeBlocks) == 0 {
		// 没有代码块，直接按段落分
		return strings.Split(text, "\n\n")
	}

	lastEnd := 0
	for _, block := range codeBlocks {
		// 代码块前的文本，按段落分
		if block[0] > lastEnd {
			before := text[lastEnd:block[0]]
			for _, para := range strings.Split(before, "\n\n") {
				para = strings.TrimSpace(para)
				if para != "" {
					atoms = append(atoms, para)
				}
			}
		}
		// 代码块作为整体
		codeBlock := strings.TrimSpace(text[block[0]:block[1]])
		if codeBlock != "" {
			atoms = append(atoms, codeBlock)
		}
		lastEnd = block[1]
	}

	// 最后一个代码块后的文本
	if lastEnd < len(text) {
		after := text[lastEnd:]
		for _, para := range strings.Split(after, "\n\n") {
			para = strings.TrimSpace(para)
			if para != "" {
				atoms = append(atoms, para)
			}
		}
	}

	return atoms
}

// detectChunkType 检测块类型（只有 code 和 info 两种）
func (p *DocumentProcessor) detectChunkType(text string) string {
	// 检测是否包含代码块（``` 包裹的代码）
	hasCodeBlock := strings.Contains(text, "```")

	if hasCodeBlock {
		return "code"
	}
	return "info"
}

// countTokens 使用 tiktoken 准确计算 token 数
func (p *DocumentProcessor) countTokens(text string) int {
	// 使用 cl100k_base 编码（GPT-4, text-embedding-3-small 使用的编码）
	enc, err := tiktoken.GetEncoding("cl100k_base")
	if err != nil {
		// 降级到简单估算
		return len(text) / 4
	}
	tokens := enc.Encode(text, nil, nil)
	return len(tokens)
}

// ProcessDocumentAsync 异步处理文档（单文档上传）
// actLogger: 已配置好的任务日志器
func (p *DocumentProcessor) ProcessDocumentAsync(doc *dbmodel.DocumentUpload, content []byte, docLogger *actlog.TaskLogger) {
	go func() {
		if err := p.ProcessDocument(doc, content, docLogger); err != nil {
			log.Printf("[Processor] ERROR processing document %s: %v", doc.Title, err)
			doc.Status = "failed"
			doc.ErrorMessage = err.Error()
			global.DB.Save(doc)
			docLogger.Error(actlog.EventDocFailed, fmt.Sprintf("处理失败: %s - %s", doc.Title, err.Error()))
		} else {
			// 单文档上传任务完成
			docLogger.Success(actlog.EventDocComplete, fmt.Sprintf("处理完成: %s", doc.Title))
		}
	}()
}

// ProcessDocumentWithCallback 处理文档（带状态回调和任务日志器）
// actLogger 应已设置好 document target
func (p *DocumentProcessor) ProcessDocumentWithCallback(doc *dbmodel.DocumentUpload, content []byte, statusChan chan response.ProcessStatus, actLogger *actlog.TaskLogger) {
	defer close(statusChan)

	log.Printf("[Processor] Starting to process document: %s (ID: %d)", doc.Title, doc.ID)
	actLogger.Info(actlog.EventDocParse, fmt.Sprintf("开始处理文档: %s", doc.Title))

	// 1. 解析文档
	statusChan <- response.ProcessStatus{Stage: "parsing", Progress: 5, Message: "正在解析文档...", Status: "processing"}
	text, err := p.parseDocument(doc.FileType, content)
	if err != nil {
		statusChan <- response.ProcessStatus{Stage: "failed", Progress: 0, Message: "解析失败: " + err.Error(), Status: "failed"}
		actLogger.Error(actlog.EventDocFailed, fmt.Sprintf("解析失败: %s - %s", doc.Title, err.Error()))
		doc.Status = "failed"
		doc.ErrorMessage = err.Error()
		global.DB.Save(doc)
		return
	}

	// 2. Pre-Chunking: 清理文档
	statusChan <- response.ProcessStatus{Stage: "preprocessing", Progress: 10, Message: "正在预处理文档...", Status: "processing"}
	text = p.preProcessMarkdown(text)

	// 3. 分块
	statusChan <- response.ProcessStatus{Stage: "chunking", Progress: 20, Message: "正在分块...", Status: "processing"}
	chunks := p.chunkText(text, doc.ID, doc.LibraryID, doc.Version, doc.FilePath)
	if len(chunks) == 0 {
		statusChan <- response.ProcessStatus{Stage: "failed", Progress: 0, Message: "分块失败：无有效内容", Status: "failed"}
		actLogger.Error(actlog.EventDocFailed, fmt.Sprintf("分块失败: %s - 无有效内容", doc.Title))
		doc.Status = "failed"
		doc.ErrorMessage = "no valid chunks"
		global.DB.Save(doc)
		return
	}
	actLogger.Info(actlog.EventDocChunk, fmt.Sprintf("分块完成: %s (%d 块)", doc.Title, len(chunks)))

	// 4. LLM Enrich: 生成 Title 和 Description
	if global.LLM != nil {
		statusChan <- response.ProcessStatus{Stage: "enriching", Progress: 35, Message: fmt.Sprintf("正在 AI 增强（%d 块）...", len(chunks)), Status: "processing"}
		actLogger.Info(actlog.EventDocEnrich, fmt.Sprintf("AI 增强中: %s (%d 块)", doc.Title, len(chunks)))
		if err := p.enrichChunks(chunks); err != nil {
			log.Printf("[Processor] WARNING: LLM enrich failed: %v, using fallback", err)
		}
	}

	// 5. 生成 Embedding
	statusChan <- response.ProcessStatus{Stage: "embedding", Progress: 60, Message: fmt.Sprintf("正在生成 Embedding（%d 块）...", len(chunks)), Status: "processing"}
	actLogger.Info(actlog.EventDocEmbed, fmt.Sprintf("生成 Embedding: %s (%d 块)", doc.Title, len(chunks)))
	var texts []string
	for _, chunk := range chunks {
		texts = append(texts, chunk.ChunkText)
	}

	embeddings, err := global.Embedding.EmbedBatch(texts)
	if err != nil {
		statusChan <- response.ProcessStatus{Stage: "failed", Progress: 0, Message: "Embedding 生成失败: " + err.Error(), Status: "failed"}
		actLogger.Error(actlog.EventDocFailed, fmt.Sprintf("Embedding 失败: %s - %s", doc.Title, err.Error()))
		doc.Status = "failed"
		doc.ErrorMessage = err.Error()
		global.DB.Save(doc)
		return
	}

	// 6. 保存
	statusChan <- response.ProcessStatus{Stage: "saving", Progress: 85, Message: "正在保存...", Status: "processing"}
	for i, chunk := range chunks {
		if i < len(embeddings) {
			chunk.Embedding = pgvector.NewVector(embeddings[i])
		}
	}

	if err := global.DB.CreateInBatches(chunks, 100).Error; err != nil {
		statusChan <- response.ProcessStatus{Stage: "failed", Progress: 0, Message: "保存失败: " + err.Error(), Status: "failed"}
		doc.Status = "failed"
		doc.ErrorMessage = err.Error()
		global.DB.Save(doc)
		return
	}

	// 7. 计算统计信息
	totalTokens := 0
	for _, chunk := range chunks {
		totalTokens += chunk.Tokens
	}

	// 8. 完成 - 更新文档状态和统计
	doc.Status = "completed"
	doc.ChunkCount = len(chunks)
	doc.TokenCount = totalTokens
	if err := global.DB.Save(doc).Error; err != nil {
		log.Printf("[Processor] Failed to save document stats: %v", err)
	}

	// 9. 失效该库版本的搜索缓存（Version Tag 模式）
	searchService := &SearchService{}
	if err := searchService.InvalidateLibraryCache(doc.LibraryID, doc.Version); err != nil {
		log.Printf("[Processor] WARNING: Failed to invalidate cache for library %d version %s: %v", doc.LibraryID, doc.Version, err)
	}

	log.Printf("[Processor] Successfully processed document: %s (chunks: %d, tokens: %d)", doc.Title, len(chunks), totalTokens)
	// SSE 版本单文档上传任务完成
	actLogger.Success(actlog.EventDocComplete, fmt.Sprintf("上传完成: %s (%d 块, %d tokens)", doc.Title, len(chunks), totalTokens))
	statusChan <- response.ProcessStatus{Stage: "completed", Progress: 100, Message: "处理完成", Status: "completed"}
}

// ============================================================================
// Pre-Chunking: 文档预处理
// ============================================================================

// preProcessMarkdown 预处理 Markdown 文档
// 移除徽章、空白行、HTML 标签等无效内容
func (p *DocumentProcessor) preProcessMarkdown(text string) string {
	// 1. 移除 Markdown 徽章 [![...](...)
	badgeRegex := regexp.MustCompile(`\[!\[.*?\]\(.*?\)\]\(.*?\)`)
	text = badgeRegex.ReplaceAllString(text, "")

	// 2. 移除单独的图片标签 ![...](...)
	// 但保留带说明的图片
	standaloneImgRegex := regexp.MustCompile(`(?m)^!\[.*?\]\(.*?\)\s*$`)
	text = standaloneImgRegex.ReplaceAllString(text, "")

	// 3. 移除 HTML 注释
	htmlCommentRegex := regexp.MustCompile(`<!--[\s\S]*?-->`)
	text = htmlCommentRegex.ReplaceAllString(text, "")

	// 4. 移除 HTML img 标签
	htmlImgRegex := regexp.MustCompile(`<img[^>]*>`)
	text = htmlImgRegex.ReplaceAllString(text, "")

	// 5. 移除连续的空行（保留最多一个）
	multipleNewlines := regexp.MustCompile(`\n{3,}`)
	text = multipleNewlines.ReplaceAllString(text, "\n\n")

	// 6. 移除 --- 分隔线
	hrRegex := regexp.MustCompile(`(?m)^---+\s*$`)
	text = hrRegex.ReplaceAllString(text, "")

	return strings.TrimSpace(text)
}

// ============================================================================
// LLM Enrich: 为文档块生成 Title 和 Description
// ============================================================================

// enrichChunks 使用 LLM 为 code 类型的文档块生成 Title 和 Description
// Info 类型的块保持 Title 为 headers 层级，不调用 LLM
// 使用 5 个 worker 并发加速处理
func (p *DocumentProcessor) enrichChunks(chunks []*dbmodel.DocumentChunk) error {
	ctx := context.Background()

	// 筛选需要富化的 code chunks
	type enrichTask struct {
		idx   int
		chunk *dbmodel.DocumentChunk
	}

	var tasks []enrichTask
	for i, chunk := range chunks {
		if chunk.ChunkType == "info" {
			log.Printf("[Processor] Info chunk %d: keeping headers as title: %s", i, chunk.Title)
			continue
		}
		tasks = append(tasks, enrichTask{idx: i, chunk: chunk})
	}

	if len(tasks) == 0 {
		return nil
	}

	// Worker Pool: 5 个 worker
	const workerCount = 5
	taskChan := make(chan enrichTask, len(tasks))
	var wg sync.WaitGroup

	// 启动 workers
	for w := 0; w < workerCount; w++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for task := range taskChan {
				input := llm.EnrichInput{
					Content: task.chunk.ChunkText,
					Headers: task.chunk.Title,
				}

				output, err := global.LLM.Enrich(ctx, input)
				if err != nil {
					log.Printf("[Processor] Enrich failed for chunk %d: %v", task.idx, err)
					continue
				}

				if output.Title != "" {
					task.chunk.Title = output.Title
				}
				if output.Description != "" {
					task.chunk.Description = output.Description
				}

				log.Printf("[Processor] Enriched code chunk %d: %s", task.idx, output.Title)
			}
		}()
	}

	// 发送任务
	for _, task := range tasks {
		taskChan <- task
	}
	close(taskChan)

	wg.Wait()
	return nil
}
