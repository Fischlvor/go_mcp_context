package service

import (
	"strings"

	dbmodel "go-mcp-context/internal/model/database"
	"go-mcp-context/internal/model/request"
	"go-mcp-context/internal/model/response"
	"go-mcp-context/pkg/bufferedwriter/stats"
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

		// 获取版本列表
		versions := []string(lib.Versions)
		if len(versions) == 0 {
			defaultVer := lib.DefaultVersion
			if defaultVer == "" {
				defaultVer = "latest"
			}
			versions = []string{defaultVer}
		}

		// 默认版本
		defaultVersion := lib.DefaultVersion
		if defaultVersion == "" {
			defaultVersion = "latest"
		}

		result.Libraries = append(result.Libraries, response.MCPLibraryInfo{
			LibraryID:      lib.ID,
			Name:           lib.Name,
			Versions:       versions,
			DefaultVersion: defaultVersion,
			Description:    lib.Description,
			Snippets:       int(snippetCount),
			Score:          score,
		})
	}

	// 统计 MCP 调用（全局统计，不关联具体库）
	stats.Increment(dbmodel.MetricMCPSearchLibraries, 1)

	return result, nil
}

// GetLibraryDocs 获取库文档（MCP 工具）
// 支持两种模式：
// 1. 指定 libraryID：在特定库中搜索
// 2. 不指定 libraryID（为 0）：全局搜索所有库
func (s *MCPService) GetLibraryDocs(req *request.MCPGetLibraryDocs) (*response.MCPGetLibraryDocsResult, error) {
	// 分页参数
	page := req.Page
	if page < 1 || page > 10 {
		page = 1
	}
	limit := 10 // MCP 每页固定 10 条

	version := req.Version

	// 如果指定了 libraryID，验证库是否存在
	var libraryID uint
	if req.LibraryID > 0 {
		libraryService := &LibraryService{}
		library, err := libraryService.GetByID(req.LibraryID)
		if err != nil {
			return nil, ErrNotFound
		}
		libraryID = library.ID
	}

	// 执行搜索（libraryID 为 0 时全局搜索）
	searchResult, err := s.searchService.SearchDocuments(&request.Search{
		LibraryID: libraryID, // 0 表示全局搜索
		Query:     req.Topic,
		Mode:      req.Mode,
		Version:   version,
		Page:      page,
		Limit:     limit,
	})
	if err != nil {
		return nil, err
	}

	documents := make([]response.MCPDocumentChunk, 0, len(searchResult.Results))
	for _, r := range searchResult.Results {
		doc := response.MCPDocumentChunk{
			Title:       r.Title,
			Description: r.Description, // code mode 有值，info mode 为空
			Source:      r.Source,
			Version:     r.Version,
			Mode:        r.Mode,
			Language:    r.Language, // code mode 有值，info mode 为空
			Code:        r.Code,     // code mode 有值，info mode 为空
			Tokens:      r.Tokens,
			Relevance:   r.Relevance,
		}
		// info 模式才返回 content（chunk_text）
		if r.Mode == "info" {
			doc.Content = r.Content
		}
		documents = append(documents, doc)
	}

	// 统计 MCP 调用（如果有指定库）
	if libraryID > 0 {
		stats.IncrementWithLibrary(libraryID, dbmodel.MetricMCPGetLibraryDocs, 1)
	}

	return &response.MCPGetLibraryDocsResult{
		LibraryID: libraryID,
		Documents: documents,
		Page:      page,
		HasMore:   searchResult.HasMore,
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
