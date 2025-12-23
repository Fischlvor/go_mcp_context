package actlog

import (
	"encoding/json"
	"log"
)

// ConsoleWriter 控制台写入器（调试用）
type ConsoleWriter struct {
	prefix string
}

// NewConsoleWriter 创建控制台写入器
func NewConsoleWriter(prefix string) *ConsoleWriter {
	return &ConsoleWriter{prefix: prefix}
}

// Write 写入单条日志
func (w *ConsoleWriter) Write(entry *LogEntry) error {
	data, _ := json.Marshal(entry)
	log.Printf("%s %s", w.prefix, string(data))
	return nil
}

// WriteBatch 批量写入日志
func (w *ConsoleWriter) WriteBatch(entries []*LogEntry) error {
	for _, entry := range entries {
		if err := w.Write(entry); err != nil {
			return err
		}
	}
	return nil
}

// Close 关闭写入器
func (w *ConsoleWriter) Close() error {
	return nil
}
