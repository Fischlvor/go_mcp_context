package database

import (
	"time"
)

// MCP 函数名常量
const (
	MCPFuncSearchLibraries = "search_libraries"
	MCPFuncGetLibraryDocs  = "get_library_docs"
)

// MCPCallLog MCP 调用日志
type MCPCallLog struct {
	ID        uint   `json:"id" gorm:"primaryKey"`
	ActorID   string `json:"actor_id" gorm:"size:36;not null;index"`  // 调用者 UUID
	FuncName  string `json:"func_name" gorm:"size:64;not null;index"` // search_libraries / get_library_docs
	LibraryID *uint  `json:"library_id,omitempty" gorm:"index"`       // 关联库（get_library_docs 有）

	// 请求参数（JSON 格式）
	Params string `json:"params,omitempty" gorm:"type:jsonb"` // 请求参数 JSON

	// 响应摘要
	ResultCount int    `json:"result_count" gorm:"default:0"`           // 返回结果数
	LatencyMs   int    `json:"latency_ms" gorm:"default:0"`             // 响应时间(ms)
	Status      string `json:"status" gorm:"size:16;default:'success'"` // success/error
	ErrorMsg    string `json:"error_msg,omitempty" gorm:"size:500"`     // 错误信息

	CreatedAt time.Time `json:"created_at" gorm:"autoCreateTime;index"`
}

func (MCPCallLog) TableName() string {
	return "mcp_call_logs"
}
