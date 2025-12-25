package api

import (
	"go-mcp-context/internal/model/request"
	"go-mcp-context/internal/model/response"

	"github.com/gin-gonic/gin"
)

type SearchApi struct{}

// Search 搜索文档
func (s *SearchApi) Search(c *gin.Context) {
	var req request.Search
	if err := c.ShouldBindJSON(&req); err != nil {
		response.FailWithMessage("参数错误: "+err.Error(), c)
		return
	}

	if req.LibraryID == 0 {
		response.FailWithMessage("library_id 必须大于 0", c)
		return
	}

	if req.Version == "" {
		response.FailWithMessage("version 为必填参数", c)
		return
	}

	if req.Page <= 0 {
		req.Page = 1
	}
	if req.Limit <= 0 || req.Limit > 10 {
		req.Limit = 10
	}

	result, err := searchService.SearchDocuments(&req)
	if err != nil {
		response.FailWithMessage("搜索失败: "+err.Error(), c)
		return
	}

	response.OkWithData(result, c)
}
