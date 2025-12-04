package service

import (
	"context"
	"fmt"
	"sort"

	dbmodel "go-mcp-context/internal/model/database"
	"go-mcp-context/internal/model/request"
	"go-mcp-context/internal/model/response"
	"go-mcp-context/pkg/global"

	"github.com/pgvector/pgvector-go"
)

type SearchService struct{}

// 重排序权重
const (
	VectorWeight = 0.5 // 向量相似度权重
	BM25Weight   = 0.3 // BM25 权重
	HotWeight    = 0.2 // 热度权重
)

// searchCandidate 搜索候选项（内部使用）
type searchCandidate struct {
	Chunk       dbmodel.DocumentChunk
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
	if limit > 50 {
		limit = 50
	}

	// 1. 生成查询向量
	queryVector, err := global.Embedding.Embed(req.Query)
	if err != nil {
		return nil, fmt.Errorf("failed to generate embedding: %w", err)
	}

	// 2. 执行向量搜索 (Top-50)
	vectorResults, err := s.vectorSearch(ctx, req.LibraryID, queryVector, req.Mode, 50)
	if err != nil {
		return nil, fmt.Errorf("vector search failed: %w", err)
	}

	// 3. 执行 BM25 关键词搜索 (Top-50)
	bm25Results, err := s.bm25Search(ctx, req.LibraryID, req.Query, req.Mode, 50)
	if err != nil {
		return nil, fmt.Errorf("bm25 search failed: %w", err)
	}

	// 4. 合并去重并重排序
	candidates := s.mergeAndRerank(vectorResults, bm25Results)

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
	for _, c := range candidates[start:end] {
		results = append(results, response.SearchResultItem{
			ChunkID:     c.Chunk.ID,
			DocumentID:  c.Chunk.DocumentID,
			LibraryID:   c.Chunk.LibraryID,
			Content:     c.Chunk.ChunkText,
			ChunkType:   c.Chunk.ChunkType,
			Score:       c.FinalScore,
			VectorScore: c.VectorScore,
			BM25Score:   c.BM25Score,
		})
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
func (s *SearchService) vectorSearch(ctx context.Context, libraryID uint, queryVector []float32, mode string, limit int) ([]searchCandidate, error) {
	var chunks []struct {
		dbmodel.DocumentChunk
		Distance float64 `gorm:"column:distance"`
	}

	query := global.DB.Table("document_chunks").
		Select("document_chunks.*, embedding <=> ? as distance", pgvector.NewVector(queryVector)).
		Where("library_id = ? AND status = ?", libraryID, "active").
		Order("distance ASC").
		Limit(limit)

	if mode != "" {
		query = query.Where("chunk_type = ?", mode)
	}

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
			VectorScore: similarity,
		}
	}

	return results, nil
}

// bm25Search BM25 关键词搜索
func (s *SearchService) bm25Search(ctx context.Context, libraryID uint, query string, mode string, limit int) ([]searchCandidate, error) {
	var chunks []struct {
		dbmodel.DocumentChunk
		Rank float64 `gorm:"column:rank"`
	}

	// 使用 PostgreSQL 全文搜索
	sqlQuery := global.DB.Table("document_chunks").
		Select("document_chunks.*, ts_rank(to_tsvector('simple', chunk_text), plainto_tsquery('simple', ?)) as rank", query).
		Where("library_id = ? AND status = ?", libraryID, "active").
		Where("to_tsvector('simple', chunk_text) @@ plainto_tsquery('simple', ?)", query).
		Order("rank DESC").
		Limit(limit)

	if mode != "" {
		sqlQuery = sqlQuery.Where("chunk_type = ?", mode)
	}

	if err := sqlQuery.Find(&chunks).Error; err != nil {
		return nil, err
	}

	results := make([]searchCandidate, len(chunks))
	for i, c := range chunks {
		results[i] = searchCandidate{
			Chunk:     c.DocumentChunk,
			BM25Score: c.Rank,
		}
	}

	return results, nil
}

// mergeAndRerank 合并去重并重排序
func (s *SearchService) mergeAndRerank(vectorResults, bm25Results []searchCandidate) []searchCandidate {
	// 用 map 合并去重
	candidateMap := make(map[uint]*searchCandidate)

	// 获取最大热度用于归一化
	var maxAccessCount int
	for _, c := range vectorResults {
		if c.Chunk.AccessCount > maxAccessCount {
			maxAccessCount = c.Chunk.AccessCount
		}
	}
	for _, c := range bm25Results {
		if c.Chunk.AccessCount > maxAccessCount {
			maxAccessCount = c.Chunk.AccessCount
		}
	}
	if maxAccessCount == 0 {
		maxAccessCount = 1 // 避免除零
	}

	// 合并向量搜索结果
	for _, c := range vectorResults {
		candidate := c
		candidate.HotScore = float64(c.Chunk.AccessCount) / float64(maxAccessCount)
		candidateMap[c.Chunk.ID] = &candidate
	}

	// 合并 BM25 结果
	for _, c := range bm25Results {
		if existing, ok := candidateMap[c.Chunk.ID]; ok {
			existing.BM25Score = c.BM25Score
		} else {
			candidate := c
			candidate.HotScore = float64(c.Chunk.AccessCount) / float64(maxAccessCount)
			candidateMap[c.Chunk.ID] = &candidate
		}
	}

	// 归一化 BM25 分数
	var maxBM25 float64
	for _, c := range candidateMap {
		if c.BM25Score > maxBM25 {
			maxBM25 = c.BM25Score
		}
	}
	if maxBM25 > 0 {
		for _, c := range candidateMap {
			c.BM25Score = c.BM25Score / maxBM25
		}
	}

	// 计算最终分数
	candidates := make([]searchCandidate, 0, len(candidateMap))
	for _, c := range candidateMap {
		c.FinalScore = VectorWeight*c.VectorScore + BM25Weight*c.BM25Score + HotWeight*c.HotScore
		candidates = append(candidates, *c)
	}

	// 按最终分数降序排序
	sort.Slice(candidates, func(i, j int) bool {
		return candidates[i].FinalScore > candidates[j].FinalScore
	})

	return candidates
}
