package service

import (
	"context"
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"sort"
	"strings"
	"sync"
	"time"

	dbmodel "go-mcp-context/internal/model/database"
	"go-mcp-context/internal/model/request"
	"go-mcp-context/internal/model/response"
	"go-mcp-context/pkg/cache"
	"go-mcp-context/pkg/global"

	"github.com/pgvector/pgvector-go"
)

const (
	// 搜索结果缓存 TTL
	SearchCacheTTL = 24 * time.Hour
	// 搜索结果缓存 key 前缀
	SearchCachePrefix = "search:topic:"
)

type SearchService struct{}

// InvalidateLibraryCache 失效指定库版本的搜索缓存
// 在文档处理完成后调用，使旧的搜索结果缓存失效
func (s *SearchService) InvalidateLibraryCache(libraryID uint, version string) error {
	if global.Cache == nil {
		return nil
	}

	tag := s.buildSearchCacheTag(libraryID, version)
	return global.Cache.InvalidateTags([]string{tag})
}

// 混合搜索权重 - 使用RRF算法
const (
	VectorRRFWeight = 0.7 // 向量搜索RRF权重
	BM25RRFWeight   = 0.3 // BM25搜索RRF权重
	HotWeight       = 0.2 // 热度权重（保持不变）

	// RRF 常量
	RRFConstant = 60 // Elasticsearch 默认值，较高值让低排名文档也有影响力
)

// searchCandidate 搜索候选项（内部使用）
type searchCandidate struct {
	Chunk       dbmodel.DocumentChunk
	DocTitle    string // 文档标题
	VectorScore float64
	BM25Score   float64
	HotScore    float64
	FinalScore  float64
}

// SearchDocuments 搜索文档
func (s *SearchService) SearchDocuments(req *request.Search) (*response.SearchResult, error) {
	ctx := context.Background()

	// 参数默认值
	page := req.Page
	limit := req.Limit
	if page <= 0 {
		page = 1
	}
	if limit <= 0 {
		limit = 10
	}
	if limit > 10 {
		limit = 10
	}

	// 拆分 topic（支持逗号、空格分隔）
	topics := splitTopics(req.Query)

	var candidates []searchCandidate
	var err error

	if len(topics) <= 1 {
		// 单个 topic，使用混合RRF搜索
		candidates, err = s.searchSingleTopic(ctx, req, req.Query)
	} else {
		// 多个 topic，并行搜索 + RRF 合并
		candidates, err = s.searchMultiTopicsWithRRF(ctx, req, topics)
	}
	if err != nil {
		return nil, err
	}

	// 5. 分页返回
	total := len(candidates)
	start := (page - 1) * limit
	end := start + limit
	if start >= total {
		return &response.SearchResult{
			Results: []response.SearchResultItem{},
			Total:   int64(total),
			Page:    page,
			Limit:   limit,
			HasMore: false,
		}, nil
	}
	if end > total {
		end = total
	}

	results := make([]response.SearchResultItem, 0, end-start)
	chunkIDs := make([]uint, 0, end-start)
	for _, c := range candidates[start:end] {
		item := response.SearchResultItem{
			ChunkID:     c.Chunk.ID,
			UploadID:    c.Chunk.UploadID,
			LibraryID:   c.Chunk.LibraryID,
			Version:     c.Chunk.Version,
			Mode:        c.Chunk.ChunkType,   // code 或 info
			Title:       c.Chunk.Title,       // code mode: LLM 生成, info mode: headers 层级
			Description: c.Chunk.Description, // code mode: LLM 生成, info mode: 空
			Source:      c.Chunk.Source,
			Language:    c.Chunk.Language, // code mode: 代码语言, info mode: 空
			Code:        c.Chunk.Code,     // code mode: 代码内容, info mode: 空
			Tokens:      c.Chunk.Tokens,
			Relevance:   c.FinalScore,
		}
		// info 类型的块返回 content（chunk_text）
		if c.Chunk.ChunkType == "info" {
			item.Content = c.Chunk.ChunkText
		}
		results = append(results, item)
		chunkIDs = append(chunkIDs, c.Chunk.ID)
	}

	// 异步更新 access_count
	if len(chunkIDs) > 0 {
		go func(ids []uint) {
			global.DB.Table("document_chunks").
				Where("id IN ?", ids).
				UpdateColumn("access_count", global.DB.Raw("access_count + 1"))
		}(chunkIDs)
	}

	return &response.SearchResult{
		Results: results,
		Total:   int64(total),
		Page:    page,
		Limit:   limit,
		HasMore: end < total,
	}, nil
}

// vectorSearch 向量搜索
func (s *SearchService) vectorSearch(ctx context.Context, libraryID uint, queryVector []float32, mode string, version string, limit int) ([]searchCandidate, error) {
	var chunks []struct {
		dbmodel.DocumentChunk
		Distance float64 `gorm:"column:distance"`
		DocTitle string  `gorm:"column:doc_title"`
	}

	query := global.DB.Model(&dbmodel.DocumentChunk{}).
		Select("document_chunks.*, document_uploads.title as doc_title, embedding <=> ? as distance", pgvector.NewVector(queryVector)).
		Joins("LEFT JOIN document_uploads ON document_uploads.id = document_chunks.upload_id").
		Where("document_chunks.status = ? AND document_chunks.deleted_at IS NULL", "active").
		Where("document_chunks.library_id = ?", libraryID).
		Where("document_chunks.version = ?", version)

	// mode 过滤：code 搜索 code 类型，info 搜索 info 类型
	if mode == "code" {
		query = query.Where("document_chunks.chunk_type = ?", "code")
	} else if mode == "info" {
		query = query.Where("document_chunks.chunk_type = ?", "info")
	}

	query = query.
		Order("distance ASC").
		Limit(limit)

	if err := query.Find(&chunks).Error; err != nil {
		return nil, err
	}

	results := make([]searchCandidate, len(chunks))
	for i, c := range chunks {
		// 将距离转换为相似度 (1 - distance)，距离越小相似度越高
		similarity := 1.0 - c.Distance
		if similarity < 0 {
			similarity = 0
		}
		results[i] = searchCandidate{
			Chunk:       c.DocumentChunk,
			DocTitle:    c.DocTitle,
			VectorScore: similarity,
		}
	}

	return results, nil
}

// bm25Search BM25 关键词搜索
func (s *SearchService) bm25Search(ctx context.Context, libraryID uint, query string, mode string, version string, limit int) ([]searchCandidate, error) {
	var chunks []struct {
		dbmodel.DocumentChunk
		Rank     float64 `gorm:"column:rank"`
		DocTitle string  `gorm:"column:doc_title"`
	}

	// 使用 PostgreSQL 全文搜索
	sqlQuery := global.DB.Model(&dbmodel.DocumentChunk{}).
		Select("document_chunks.*, document_uploads.title as doc_title, ts_rank(document_chunks.chunk_tsvector_simple, plainto_tsquery('simple', ?)) as rank", query).
		Joins("LEFT JOIN document_uploads ON document_uploads.id = document_chunks.upload_id").
		Where("document_chunks.status = ? AND document_chunks.deleted_at IS NULL", "active").
		Where("document_chunks.library_id = ?", libraryID).
		Where("document_chunks.version = ?", version).
		Where("document_chunks.chunk_tsvector_simple @@ plainto_tsquery('simple', ?)", query).
		Order("rank DESC").
		Limit(limit)

	// mode 过滤：code 搜索 code 类型，info 搜索 info 类型
	if mode == "code" {
		sqlQuery = sqlQuery.Where("document_chunks.chunk_type = ?", "code")
	} else if mode == "info" {
		sqlQuery = sqlQuery.Where("document_chunks.chunk_type = ?", "info")
	}

	if err := sqlQuery.Find(&chunks).Error; err != nil {
		return nil, err
	}

	results := make([]searchCandidate, len(chunks))
	for i, c := range chunks {
		results[i] = searchCandidate{
			Chunk:     c.DocumentChunk,
			DocTitle:  c.DocTitle,
			BM25Score: c.Rank,
		}
	}

	return results, nil
}

// mergeAndRerank 合并去重并重排序 - 使用RRF算法
func (s *SearchService) mergeAndRerank(vectorResults, bm25Results []searchCandidate) []searchCandidate {
	// 对得分进行归一化（可选，RRF主要基于排名）
	s.normalizeScoresMinMax(vectorResults, "vector")
	s.normalizeScoresMinMax(bm25Results, "bm25")

	// 使用RRF算法合并结果
	return s.hybridRRF(vectorResults, bm25Results)
}

// extractDeepestTitle 从 Metadata 提取最深层级的标题
// ⚠️ 已注释：此函数从未被使用（死代码）
/*
func extractDeepestTitle(metadata dbmodel.JSON) string {
	if metadata == nil {
		return ""
	}

	// 从 h6 到 h1 找最深层级的标题
	for level := 6; level >= 1; level-- {
		key := fmt.Sprintf("h%d", level)
		if title, ok := metadata[key]; ok {
			if titleStr, ok := title.(string); ok && titleStr != "" {
				return titleStr
			}
		}
	}

	return ""
}
*/

// splitTopics 拆分 topic 字符串（按逗号分隔）
// 例如 "data fetching, routing" -> ["data fetching", "routing"]
func splitTopics(query string) []string {
	parts := strings.Split(query, ",")
	var topics []string

	for _, part := range parts {
		part = strings.TrimSpace(part)
		if part != "" {
			topics = append(topics, part)
		}
	}

	return topics
}

// searchSingleTopic 单个 topic 搜索（带缓存）
func (s *SearchService) searchSingleTopic(ctx context.Context, req *request.Search, topic string) ([]searchCandidate, error) {
	// 生成缓存 key: search:topic:{library_id}:{version}:{mode}:{topic_hash}
	cacheKey := s.buildSearchCacheKey(req.LibraryID, req.Version, req.Mode, topic)

	// 生成缓存 tag: library:{library_id}:{version}
	cacheTag := s.buildSearchCacheTag(req.LibraryID, req.Version)

	// 定义搜索函数
	fetchFunc := s.buildSearchFunc(ctx, req, topic)

	// 使用 GetOrSetWithTags 模式：缓存 key 包含 tag version，tag 失效时旧缓存自动失效
	return cache.GetOrSetWithTags(global.Cache, cacheKey, []string{cacheTag}, SearchCacheTTL, fetchFunc)
}

// buildSearchFunc 构建搜索函数（用于 GetOrSet）
func (s *SearchService) buildSearchFunc(ctx context.Context, req *request.Search, topic string) func() ([]searchCandidate, error) {
	return func() ([]searchCandidate, error) {
		return s.executeSearch(ctx, req, topic)
	}
}

// executeSearch 执行实际的搜索逻辑
func (s *SearchService) executeSearch(ctx context.Context, req *request.Search, topic string) ([]searchCandidate, error) {
	// 1. 生成查询向量（CachedEmbeddingService 自带缓存）
	queryVector, err := global.Embedding.Embed(topic)
	if err != nil {
		return nil, fmt.Errorf("failed to generate embedding: %w", err)
	}

	// 2. 执行向量搜索 (Top-50)
	vectorResults, err := s.vectorSearch(ctx, req.LibraryID, queryVector, req.Mode, req.Version, 50)
	if err != nil {
		return nil, fmt.Errorf("vector search failed: %w", err)
	}

	// 3. 执行 BM25 关键词搜索 (Top-50)
	bm25Results, err := s.bm25Search(ctx, req.LibraryID, topic, req.Mode, req.Version, 50)
	if err != nil {
		return nil, fmt.Errorf("bm25 search failed: %w", err)
	}

	// 4. 合并去重并重排序
	return s.mergeAndRerank(vectorResults, bm25Results), nil
}

// buildSearchCacheKey 构建搜索缓存 key
// 格式: search:topic:{library_id}:{version}:{mode}:{topic_hash}
// 参数顺序与 key 格式一致
func (s *SearchService) buildSearchCacheKey(libraryID uint, version, mode, topic string) string {
	hash := md5.Sum([]byte(topic))
	topicHash := hex.EncodeToString(hash[:])
	return fmt.Sprintf("%s%d:%s:%s:%s", SearchCachePrefix, libraryID, version, mode, topicHash)
}

// buildSearchCacheTag 构建搜索缓存 tag
// 格式: library:{library_id}:{version}
// 用于在库版本更新时批量失效相关缓存
func (s *SearchService) buildSearchCacheTag(libraryID uint, version string) string {
	return fmt.Sprintf("library:%d:%s", libraryID, version)
}

// normalizeScoresMinMax Min-Max归一化得分到[0,1]范围
func (s *SearchService) normalizeScoresMinMax(candidates []searchCandidate, scoreField string) {
	if len(candidates) == 0 {
		return
	}

	// 找出最大值和最小值
	var minScore, maxScore float64
	for i, candidate := range candidates {
		var score float64
		switch scoreField {
		case "vector":
			score = candidate.VectorScore
		case "bm25":
			score = candidate.BM25Score
		default:
			continue
		}

		if i == 0 {
			minScore, maxScore = score, score
		} else {
			if score < minScore {
				minScore = score
			}
			if score > maxScore {
				maxScore = score
			}
		}
	}

	// 避免除零
	if maxScore == minScore {
		return
	}

	// 归一化
	for i := range candidates {
		var score *float64
		switch scoreField {
		case "vector":
			score = &candidates[i].VectorScore
		case "bm25":
			score = &candidates[i].BM25Score
		default:
			continue
		}
		*score = (*score - minScore) / (maxScore - minScore)
	}
}

// hybridRRF 使用RRF算法合并向量搜索和BM25搜索结果
func (s *SearchService) hybridRRF(vectorResults, bm25Results []searchCandidate) []searchCandidate {
	// 构建排名映射
	vectorRanks := make(map[uint]int)
	bm25Ranks := make(map[uint]int)

	for rank, candidate := range vectorResults {
		vectorRanks[candidate.Chunk.ID] = rank + 1 // 排名从1开始
	}

	for rank, candidate := range bm25Results {
		bm25Ranks[candidate.Chunk.ID] = rank + 1 // 排名从1开始
	}

	// 收集所有候选文档
	candidateMap := make(map[uint]*searchCandidate)

	// 处理向量搜索结果
	for _, candidate := range vectorResults {
		c := candidate
		candidateMap[c.Chunk.ID] = &c
	}

	// 处理BM25搜索结果
	for _, candidate := range bm25Results {
		if existing, ok := candidateMap[candidate.Chunk.ID]; ok {
			existing.BM25Score = candidate.BM25Score
		} else {
			c := candidate
			candidateMap[c.Chunk.ID] = &c
		}
	}

	// 计算热度归一化
	var maxAccessCount int
	for _, c := range candidateMap {
		if c.Chunk.AccessCount > maxAccessCount {
			maxAccessCount = c.Chunk.AccessCount
		}
	}
	if maxAccessCount == 0 {
		maxAccessCount = 1
	}

	// 计算RRF分数
	candidates := make([]searchCandidate, 0, len(candidateMap))
	for _, candidate := range candidateMap {
		rrfScore := 0.0

		// 向量搜索贡献
		if rank, exists := vectorRanks[candidate.Chunk.ID]; exists {
			rrfScore += VectorRRFWeight / (float64(rank) + float64(RRFConstant))
		}

		// BM25搜索贡献
		if rank, exists := bm25Ranks[candidate.Chunk.ID]; exists {
			rrfScore += BM25RRFWeight / (float64(rank) + float64(RRFConstant))
		}

		// 热度贡献
		hotScore := float64(candidate.Chunk.AccessCount) / float64(maxAccessCount)
		rrfScore += HotWeight * hotScore

		candidate.FinalScore = rrfScore
		candidate.HotScore = hotScore
		candidates = append(candidates, *candidate)
	}

	// 按RRF分数降序排序
	sort.Slice(candidates, func(i, j int) bool {
		return candidates[i].FinalScore > candidates[j].FinalScore
	})

	return candidates
}

// searchMultiTopicsWithRRF 多 topic 并行搜索 + RRF 合并
func (s *SearchService) searchMultiTopicsWithRRF(ctx context.Context, req *request.Search, topics []string) ([]searchCandidate, error) {
	// 并行搜索每个 topic
	type topicResult struct {
		topic      string
		candidates []searchCandidate
		err        error
	}

	resultChan := make(chan topicResult, len(topics))
	var wg sync.WaitGroup

	for _, topic := range topics {
		wg.Add(1)
		go func(t string) {
			defer wg.Done()
			candidates, err := s.searchSingleTopic(ctx, req, t)
			resultChan <- topicResult{topic: t, candidates: candidates, err: err}
		}(topic)
	}

	// 等待所有搜索完成
	go func() {
		wg.Wait()
		close(resultChan)
	}()

	// 收集结果
	var allResults [][]searchCandidate
	for result := range resultChan {
		if result.err != nil {
			// 记录错误但继续处理其他结果
			global.Log.Warn(fmt.Sprintf("topic search failed: %s, error: %v", result.topic, result.err))
			continue
		}
		if len(result.candidates) > 0 {
			allResults = append(allResults, result.candidates)
		}
	}

	if len(allResults) == 0 {
		return nil, nil
	}

	// 使用 RRF 合并多个结果列表
	return s.reciprocalRankFusion(allResults), nil
}

// reciprocalRankFusion 使用 RRF 算法合并多个排序结果
// 公式: score(d) = Σ 1 / (k + rank(d))
// 其中 k 是常量（默认 60），rank 是文档在每个列表中的排名（从 1 开始）
func (s *SearchService) reciprocalRankFusion(resultLists [][]searchCandidate) []searchCandidate {
	// 用 map 存储每个文档的 RRF 分数
	rrfScores := make(map[uint]float64)
	candidateMap := make(map[uint]*searchCandidate)

	for _, results := range resultLists {
		for rank, candidate := range results {
			chunkID := candidate.Chunk.ID
			// RRF 公式: 1 / (k + rank)，rank 从 1 开始
			rrfScores[chunkID] += 1.0 / float64(RRFConstant+rank+1)

			// 保存候选项（如果还没有）
			if _, exists := candidateMap[chunkID]; !exists {
				c := candidate
				candidateMap[chunkID] = &c
			}
		}
	}

	// 构建最终结果
	candidates := make([]searchCandidate, 0, len(candidateMap))
	for _, candidate := range candidateMap {
		candidate.FinalScore = rrfScores[candidate.Chunk.ID]
		candidates = append(candidates, *candidate)
	}

	// 按 RRF 分数降序排序
	sort.Slice(candidates, func(i, j int) bool {
		return candidates[i].FinalScore > candidates[j].FinalScore
	})

	return candidates
}
