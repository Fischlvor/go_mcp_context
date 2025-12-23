package api

import (
	"go-mcp-context/internal/model/response"
	"go-mcp-context/pkg/utils"

	"github.com/gin-gonic/gin"
)

type StatsApi struct{}

// GetMyStats 获取当前用户的统计数据
func (s *StatsApi) GetMyStats(c *gin.Context) {
	userUUID := utils.GetUUID(c).String()

	result, err := statsService.GetUserStats(userUUID)
	if err != nil {
		response.FailWithMessage("获取统计失败: "+err.Error(), c)
		return
	}

	response.OkWithData(result, c)
}
