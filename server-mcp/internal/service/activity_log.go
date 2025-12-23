package service

import (
	"go-mcp-context/internal/model/database"
	"go-mcp-context/pkg/global"
)

// ActivityLogService 活动日志服务（仅用于 API 查询）
// 写入日志请使用 pkg/actlog 包
type ActivityLogService struct{}

// ActivityLogResult 活动日志查询结果
type ActivityLogResult struct {
	Logs   []database.ActivityLog `json:"logs"`
	TaskID string                 `json:"task_id,omitempty"`
	Status string                 `json:"status"` // "complete" or "processing"
}

// ListByLatestTask 获取库的最新任务日志（单次查询优化）
func (s *ActivityLogService) ListByLatestTask(libraryID uint) (*ActivityLogResult, error) {
	// 使用子查询一次获取最新任务的所有日志
	var logs []database.ActivityLog
	err := global.DB.
		Where("library_id = ? AND task_id = (SELECT task_id FROM activity_logs WHERE library_id = ? AND task_id IS NOT NULL AND task_id != '' ORDER BY created_at DESC LIMIT 1)", libraryID, libraryID).
		Order("created_at ASC").
		Find(&logs).Error

	if err != nil || len(logs) == 0 {
		// 没有任务日志，返回空结果
		return &ActivityLogResult{
			Logs:   []database.ActivityLog{},
			Status: "complete",
		}, nil
	}

	// 判断任务状态：最新日志（最后一条）是 info 说明还在处理中
	status := "complete"
	lastLog := logs[len(logs)-1]
	if lastLog.Status == "info" {
		status = "processing"
	}

	return &ActivityLogResult{
		Logs:   logs,
		TaskID: logs[0].TaskID,
		Status: status,
	}, nil
}

// List 获取库的活动日志列表（所有日志）
func (s *ActivityLogService) List(libraryID uint, limit int) ([]database.ActivityLog, error) {
	if limit <= 0 || limit > 100 {
		limit = 50
	}

	var logs []database.ActivityLog
	err := global.DB.
		Where("library_id = ?", libraryID).
		Order("created_at DESC").
		Limit(limit).
		Find(&logs).Error

	return logs, err
}
