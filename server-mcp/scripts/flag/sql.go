package flag

import (
	"go-mcp-context/internal/initialize"
)

// SQL 初始化数据库表结构
func SQL() error {
	initialize.InitTables()
	return nil
}
