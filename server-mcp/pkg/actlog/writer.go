package actlog

// Writer 日志写入器接口
// 支持扩展到不同的写入目标（数据库、文件、消息队列等）
type Writer interface {
	// Write 写入单条日志
	Write(entry *LogEntry) error

	// WriteBatch 批量写入日志
	WriteBatch(entries []*LogEntry) error

	// Close 关闭写入器
	Close() error
}
