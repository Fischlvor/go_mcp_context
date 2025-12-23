package api

import (
	"strconv"

	"go-mcp-context/internal/model/response"

	"github.com/gin-gonic/gin"
)

type ActivityLogApi struct{}

// List 获取库的活动日志
// @Summary 获取活动日志列表
// @Tags ActivityLog
// @Param libraryId query int true "库ID"
// @Param limit query int false "返回数量，默认50，最大100"
// @Success 200 {object} response.Response
// @Router /api/v1/logs [get]
func (a *ActivityLogApi) List(c *gin.Context) {
	// 解析库ID
	libraryIDStr := c.Query("libraryId")
	if libraryIDStr == "" {
		response.FailWithMessage("缺少 libraryId 参数", c)
		return
	}
	libraryID, err := strconv.ParseUint(libraryIDStr, 10, 64)
	if err != nil {
		response.FailWithMessage("无效的库ID", c)
		return
	}

	// 获取最新任务的日志
	result, err := activityLogService.ListByLatestTask(uint(libraryID))
	if err != nil {
		response.FailWithMessage("获取日志失败: "+err.Error(), c)
		return
	}

	response.OkWithData(gin.H{
		"logs":    result.Logs,
		"task_id": result.TaskID,
		"status":  result.Status,
	}, c)
}
