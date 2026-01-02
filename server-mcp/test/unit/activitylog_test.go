package test_test

import (
	"testing"
	"time"

	"go-mcp-context/internal/service"
	"go-mcp-context/pkg/global"

	dbmodel "go-mcp-context/internal/model/database"
	"go-mcp-context/internal/model/request"
)

// TestActivityLogListByLatestTask 测试获取最新任务日志
func Test_ActivityLog_ListByLatestTask(t *testing.T) {
	activityLogService := &service.ActivityLogService{}

	t.Run("get latest task logs", func(t *testing.T) {
		// 通过 Service 层创建库
		libService := &service.LibraryService{}
		libReq := &request.LibraryCreate{
			Name:        "actlog-test-lib",
			Description: "test",
		}
		lib, err := libService.Create(libReq)
		if err != nil {
			t.Fatalf("Create library error = %v", err)
		}
		taskID := "test-task-001"

		// 创建多条日志
		logs := []dbmodel.ActivityLog{
			{
				LibraryID: lib.ID,
				TaskID:    taskID,
				Status:    "processing",
				Message:   "Task started",
			},
			{
				LibraryID: lib.ID,
				TaskID:    taskID,
				Status:    "processing",
				Message:   "Processing documents",
			},
			{
				LibraryID: lib.ID,
				TaskID:    taskID,
				Status:    "success",
				Message:   "Task completed",
			},
		}

		for _, log := range logs {
			global.DB.Create(&log)
		}

		// 获取最新任务日志
		result, err := activityLogService.ListByLatestTask(lib.ID)
		if err != nil {
			t.Fatalf("ListByLatestTask() error = %v", err)
		}

		if result == nil {
			t.Fatal("Expected result, got nil")
		}

		if len(result.Logs) < 3 {
			t.Errorf("Expected at least 3 logs, got %d", len(result.Logs))
		}

		if result.TaskID != taskID {
			t.Errorf("Expected task ID %s, got %s", taskID, result.TaskID)
		}

		if result.Status != "complete" {
			t.Errorf("Expected status 'complete', got '%s'", result.Status)
		}
	})

	t.Run("get processing task logs", func(t *testing.T) {
		libService := &service.LibraryService{}
		lib, _ := libService.Create(&request.LibraryCreate{Name: "actlog-test-lib-2", Description: "test"})
		taskID := "test-task-002"

		// 创建日志，最后一条是 processing
		logs := []dbmodel.ActivityLog{
			{
				LibraryID: lib.ID,
				TaskID:    taskID,
				Status:    "processing",
				Message:   "Task started",
			},
			{
				LibraryID: lib.ID,
				TaskID:    taskID,
				Status:    "processing",
				Message:   "Still processing",
			},
		}

		for _, log := range logs {
			global.DB.Create(&log)
		}

		result, err := activityLogService.ListByLatestTask(lib.ID)
		if err != nil {
			t.Fatalf("ListByLatestTask() error = %v", err)
		}

		if result.Status != "processing" {
			t.Errorf("Expected status 'processing', got '%s'", result.Status)
		}
	})

	t.Run("get empty logs", func(t *testing.T) {
		libService := &service.LibraryService{}
		lib, _ := libService.Create(&request.LibraryCreate{Name: "actlog-test-lib-3", Description: "test"})

		result, err := activityLogService.ListByLatestTask(lib.ID)
		if err != nil {
			t.Fatalf("ListByLatestTask() error = %v", err)
		}

		if len(result.Logs) != 0 {
			t.Errorf("Expected 0 logs, got %d", len(result.Logs))
		}

		if result.Status != "complete" {
			t.Errorf("Expected status 'complete' for empty logs, got '%s'", result.Status)
		}
	})
}

// TestActivityLogList 测试获取活动日志列表
func Test_ActivityLog_List(t *testing.T) {
	activityLogService := &service.ActivityLogService{}

	t.Run("list activity logs with limit", func(t *testing.T) {
		libService := &service.LibraryService{}
		lib, _ := libService.Create(&request.LibraryCreate{Name: "actlog-list-lib", Description: "test"})

		// 创建 10 条日志
		for i := 1; i <= 10; i++ {
			log := &dbmodel.ActivityLog{
				LibraryID: lib.ID,
				Status:    "success",
				Message:   "Log message " + string(rune('0'+i)),
			}
			global.DB.Create(log)
			time.Sleep(10 * time.Millisecond) // 确保时间顺序
		}

		// 查询前 5 条
		logs, err := activityLogService.List(lib.ID, 5)
		if err != nil {
			t.Fatalf("List() error = %v", err)
		}

		if len(logs) != 5 {
			t.Errorf("Expected 5 logs, got %d", len(logs))
		}
	})

	t.Run("list activity logs with default limit", func(t *testing.T) {
		libService := &service.LibraryService{}
		lib, _ := libService.Create(&request.LibraryCreate{Name: "actlog-list-lib-2", Description: "test"})

		// 创建 3 条日志
		for i := 1; i <= 3; i++ {
			log := &dbmodel.ActivityLog{
				LibraryID: lib.ID,
				Status:    "success",
				Message:   "Log " + string(rune('0'+i)),
			}
			global.DB.Create(log)
		}

		// 使用默认 limit（0 应该变成 50）
		logs, err := activityLogService.List(lib.ID, 0)
		if err != nil {
			t.Fatalf("List() error = %v", err)
		}

		if len(logs) != 3 {
			t.Errorf("Expected 3 logs, got %d", len(logs))
		}
	})

	t.Run("list activity logs with invalid limit", func(t *testing.T) {
		libService := &service.LibraryService{}
		lib, _ := libService.Create(&request.LibraryCreate{Name: "actlog-list-lib-3", Description: "test"})

		log := &dbmodel.ActivityLog{
			LibraryID: lib.ID,
			Status:    "success",
			Message:   "Test log",
		}
		global.DB.Create(log)

		// 使用超过最大值的 limit（应该被限制为 50）
		logs, err := activityLogService.List(lib.ID, 200)
		if err != nil {
			t.Fatalf("List() error = %v", err)
		}

		if len(logs) != 1 {
			t.Errorf("Expected 1 log, got %d", len(logs))
		}
	})

	t.Run("list empty activity logs", func(t *testing.T) {
		libService := &service.LibraryService{}
		lib, _ := libService.Create(&request.LibraryCreate{Name: "actlog-list-lib-4", Description: "test"})

		logs, err := activityLogService.List(lib.ID, 10)
		if err != nil {
			t.Fatalf("List() error = %v", err)
		}

		if len(logs) != 0 {
			t.Errorf("Expected 0 logs, got %d", len(logs))
		}
	})

	t.Run("list activity logs with limit 0", func(t *testing.T) {
		libService := &service.LibraryService{}
		lib, _ := libService.Create(&request.LibraryCreate{Name: "actlog-list-lib-5", Description: "test"})

		// 创建一些日志
		for i := 1; i <= 5; i++ {
			log := &dbmodel.ActivityLog{
				LibraryID: lib.ID,
				Status:    "success",
				Message:   "Log " + string(rune('0'+i)),
			}
			global.DB.Create(log)
		}

		// 使用 limit 0（应该使用默认值）
		logs, err := activityLogService.List(lib.ID, 0)
		if err != nil {
			t.Fatalf("List() error = %v", err)
		}

		if len(logs) != 5 {
			t.Errorf("Expected 5 logs, got %d", len(logs))
		}
	})

	t.Run("list activity logs with different statuses", func(t *testing.T) {
		libService := &service.LibraryService{}
		lib, _ := libService.Create(&request.LibraryCreate{Name: "actlog-list-lib-6", Description: "test"})

		// 创建不同状态的日志
		statuses := []string{"success", "processing", "error", "warning"}
		for _, status := range statuses {
			log := &dbmodel.ActivityLog{
				LibraryID: lib.ID,
				Status:    status,
				Message:   "Log with status: " + status,
			}
			global.DB.Create(log)
		}

		// 查询所有日志
		logs, err := activityLogService.List(lib.ID, 10)
		if err != nil {
			t.Fatalf("List() error = %v", err)
		}

		if len(logs) < 4 {
			t.Errorf("Expected at least 4 logs, got %d", len(logs))
		}
	})
}
