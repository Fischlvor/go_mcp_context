package api

import (
	"go-mcp-context/internal/model/response"
	"go-mcp-context/pkg/utils"

	"github.com/gin-gonic/gin"
)

type StatsApi struct{}

// GetMyStats 获取当前用户的统计数据
// @Summary 获取统计数据
// @Description 获取当前用户的库、文档、API Key 等统计数据（需要认证）
// @Tags Stats
// @Accept json
// @Produce json
// @Security JWTAuth
// @Success 200 {object} response.Response{data=map[string]interface{}}
// @Failure 401 {object} response.Response
// @Router /api/v1/stats/my [get]
func (s *StatsApi) GetMyStats(c *gin.Context) {
	userUUID := utils.GetUUID(c).String()

	result, err := statsService.GetUserStats(userUUID)
	if err != nil {
		response.FailWithMessage("获取统计失败: "+err.Error(), c)
		return
	}

	response.OkWithData(result, c)
}
