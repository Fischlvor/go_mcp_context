package mcplog

import (
	"encoding/json"
	"log"
	"sync"
	"time"

	dbmodel "go-mcp-context/internal/model/database"
	"go-mcp-context/pkg/bufferedwriter"
	"go-mcp-context/pkg/global"
)

// LogEntry MCP 调用日志条目
type LogEntry struct {
	ActorID     string
	FuncName    string
	LibraryID   *uint
	Params      map[string]interface{} // 请求参数
	ResultCount int
	LatencyMs   int
	Status      string
	ErrorMsg    string
	CreatedAt   time.Time
}

// mcpLogWriter 实现 bufferedwriter.Writer 接口
type mcpLogWriter struct{}

func (w *mcpLogWriter) WriteBatch(entries []*LogEntry) error {
	if len(entries) == 0 {
		return nil
	}

	logs := make([]dbmodel.MCPCallLog, 0, len(entries))
	for _, e := range entries {
		// 序列化参数为 JSON
		var paramsJSON string
		if e.Params != nil {
			if data, err := json.Marshal(e.Params); err == nil {
				paramsJSON = string(data)
			}
		}

		logs = append(logs, dbmodel.MCPCallLog{
			ActorID:     e.ActorID,
			FuncName:    e.FuncName,
			LibraryID:   e.LibraryID,
			Params:      paramsJSON,
			ResultCount: e.ResultCount,
			LatencyMs:   e.LatencyMs,
			Status:      e.Status,
			ErrorMsg:    e.ErrorMsg,
			CreatedAt:   e.CreatedAt,
		})
	}

	return global.DB.Create(&logs).Error
}

func (w *mcpLogWriter) Close() error {
	return nil
}

var (
	buffer *bufferedwriter.Buffer[*LogEntry]
	writer *mcpLogWriter
	mu     sync.RWMutex
)

// Init 初始化 MCP 调用日志
func Init() {
	mu.Lock()
	defer mu.Unlock()

	writer = &mcpLogWriter{}
	buffer = bufferedwriter.New("mcplog", writer, bufferedwriter.Config{
		Size:     1000,
		Batch:    50,
		Interval: 2 * time.Second,
	})
	log.Println("[mcplog] Initialized")
}

// Close 关闭 MCP 调用日志
func Close() {
	mu.Lock()
	defer mu.Unlock()

	if buffer != nil {
		buffer.Close()
		buffer = nil
	}
	log.Println("[mcplog] Closed")
}

// Log 记录 MCP 调用日志
func Log(entry *LogEntry) {
	mu.RLock()
	b := buffer
	mu.RUnlock()

	if b == nil {
		log.Printf("[mcplog] Not initialized, dropping log: %s", entry.FuncName)
		return
	}

	if entry.CreatedAt.IsZero() {
		entry.CreatedAt = time.Now()
	}
	if entry.Status == "" {
		entry.Status = "success"
	}

	b.Write(entry)
}
