package service

import (
	"fmt"
	"strings"

	dbmodel "go-mcp-context/internal/model/database"
	"go-mcp-context/internal/model/request"
	"go-mcp-context/internal/model/response"
	"go-mcp-context/pkg/global"

	"github.com/agnivade/levenshtein"
)

type MCPService struct {
	searchService *SearchService
}

// NewMCPService 创建 MCP 服务
func NewMCPService() *MCPService {
	return &MCPService{
		searchService: &SearchService{},
	}
}

// SearchLibraries 搜索库（MCP 工具）
func (s *MCPService) SearchLibraries(req *request.MCPSearchLibraries) (*response.MCPSearchLibrariesResult, error) {
	var libraries []dbmodel.Library

	// 前缀匹配 + 模糊匹配
	err := global.DB.Where("status = ? AND name ILIKE ?", "active", req.LibraryName+"%").
		Order("name ASC").
		Limit(10).
		Find(&libraries).Error

	if err != nil {
		return nil, err
	}

	// 如果前缀匹配结果不足，尝试包含匹配
	if len(libraries) < 5 {
		var moreLibraries []dbmodel.Library
		global.DB.Where("status = ? AND name ILIKE ? AND name NOT ILIKE ?",
			"active", "%"+req.LibraryName+"%", req.LibraryName+"%").
			Order("name ASC").
			Limit(10 - len(libraries)).
			Find(&moreLibraries)
		libraries = append(libraries, moreLibraries...)
	}

	// 转换为响应格式并计算匹配分数
	result := &response.MCPSearchLibrariesResult{
		Libraries: make([]response.MCPLibraryInfo, 0, len(libraries)),
	}

	for _, lib := range libraries {
		// 统计文档片段数
		var snippetCount int64
		global.DB.Model(&dbmodel.DocumentChunk{}).
			Where("library_id = ? AND status = ?", lib.ID, "active").
			Count(&snippetCount)

		// 计算匹配分数
		score := calculateMatchScore(req.LibraryName, lib.Name)

		result.Libraries = append(result.Libraries, response.MCPLibraryInfo{
			ID:          fmt.Sprintf("%s/%s", lib.Name, lib.Version),
			Name:        lib.Name,
			Version:     lib.Version,
			Description: lib.Description,
			Snippets:    int(snippetCount),
			Score:       score,
		})
	}

	return result, nil
}

// GetLibraryDocs 获取库文档（MCP 工具）
func (s *MCPService) GetLibraryDocs(req *request.MCPGetLibraryDocs) (*response.MCPGetLibraryDocsResult, error) {
	// 解析 libraryID (格式: name/version)
	parts := strings.SplitN(req.LibraryID, "/", 2)
	if len(parts) != 2 {
		return nil, ErrInvalidParams
	}
	name, version := parts[0], parts[1]

	// 查找库
	libraryService := &LibraryService{}
	library, err := libraryService.GetByNameVersion(name, version)
	if err != nil {
		return nil, ErrNotFound
	}

	// 分页参数
	page := req.Page
	if page < 1 || page > 10 {
		page = 1
	}
	limit := 10 // MCP 每页固定 10 条

	// 如果有 topic，执行混合搜索
	if req.Topic != "" {
		searchResult, err := s.searchService.SearchDocuments(&request.Search{
			LibraryID: library.ID,
			Query:     req.Topic,
			Mode:      req.Mode,
			Page:      page,
			Limit:     limit,
		})
		if err != nil {
			return nil, err
		}

		documents := make([]response.MCPDocumentChunk, 0, len(searchResult.Results))
		for _, r := range searchResult.Results {
			documents = append(documents, response.MCPDocumentChunk{
				Content:   r.Content,
				ChunkType: r.ChunkType,
				Score:     r.Score,
			})
		}

		return &response.MCPGetLibraryDocsResult{
			Documents: documents,
			Page:      page,
			HasMore:   searchResult.HasMore,
		}, nil
	}

	// 没有 topic，返回库的所有文档块（按热度排序）
	documentService := &DocumentService{}
	chunks, total, err := documentService.GetChunksByLibrary(library.ID, req.Mode, page, limit)
	if err != nil {
		return nil, err
	}

	documents := make([]response.MCPDocumentChunk, 0, len(chunks))
	for _, chunk := range chunks {
		documents = append(documents, response.MCPDocumentChunk{
			Content:   chunk.ChunkText,
			ChunkType: chunk.ChunkType,
			Score:     1.0, // 无搜索时默认分数
		})
	}

	return &response.MCPGetLibraryDocsResult{
		Documents: documents,
		Page:      page,
		HasMore:   int64(page*limit) < total,
	}, nil
}

// calculateMatchScore 计算名称匹配分数
func calculateMatchScore(query, name string) float64 {
	query = strings.ToLower(query)
	name = strings.ToLower(name)

	// 完全匹配
	if query == name {
		return 1.0
	}

	// 前缀匹配
	if strings.HasPrefix(name, query) {
		return 0.9
	}

	// 包含匹配
	if strings.Contains(name, query) {
		return 0.8
	}

	// Levenshtein 相似度
	maxLen := len(query)
	if len(name) > maxLen {
		maxLen = len(name)
	}
	if maxLen == 0 {
		return 0
	}

	distance := levenshtein.ComputeDistance(query, name)
	similarity := 1.0 - float64(distance)/float64(maxLen)

	if similarity < 0 {
		return 0
	}
	return similarity * 0.7 // 最高 0.7 分
}
