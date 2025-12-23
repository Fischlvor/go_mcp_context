package actlog

import (
	"log"
	"sync"
	"time"
)

// BufferConfig 缓冲区配置
type BufferConfig struct {
	Size     int           // 缓冲区大小
	Batch    int           // 批量写入数量
	Interval time.Duration // 刷新间隔
}

// DefaultBufferConfig 默认配置
var DefaultBufferConfig = BufferConfig{
	Size:     1000,
	Batch:    50,
	Interval: 2 * time.Second,
}

// Buffer 异步缓冲区
type Buffer struct {
	writer  Writer
	config  BufferConfig
	ch      chan *LogEntry
	done    chan struct{}
	wg      sync.WaitGroup
	closed  bool
	closeMu sync.Mutex
}

// NewBuffer 创建异步缓冲区
func NewBuffer(writer Writer, config BufferConfig) *Buffer {
	b := &Buffer{
		writer: writer,
		config: config,
		ch:     make(chan *LogEntry, config.Size),
		done:   make(chan struct{}),
	}
	b.start()
	return b
}

// start 启动后台处理协程
func (b *Buffer) start() {
	b.wg.Add(1)
	go func() {
		defer b.wg.Done()

		batch := make([]*LogEntry, 0, b.config.Batch)
		ticker := time.NewTicker(b.config.Interval)
		defer ticker.Stop()

		for {
			select {
			case entry, ok := <-b.ch:
				if !ok {
					// channel 关闭，刷新剩余数据
					b.flush(batch)
					return
				}
				batch = append(batch, entry)
				if len(batch) >= b.config.Batch {
					b.flush(batch)
					batch = batch[:0]
				}

			case <-ticker.C:
				// 定时刷新
				if len(batch) > 0 {
					b.flush(batch)
					batch = batch[:0]
				}

			case <-b.done:
				// 收到关闭信号，刷新剩余数据
				b.flush(batch)
				return
			}
		}
	}()
}

// flush 刷新缓冲区
func (b *Buffer) flush(batch []*LogEntry) {
	if len(batch) == 0 {
		return
	}

	if err := b.writer.WriteBatch(batch); err != nil {
		log.Printf("[actlog] Failed to flush %d entries: %v", len(batch), err)
	}
}

// Write 写入日志（非阻塞）
func (b *Buffer) Write(entry *LogEntry) bool {
	b.closeMu.Lock()
	if b.closed {
		b.closeMu.Unlock()
		return false
	}
	b.closeMu.Unlock()

	select {
	case b.ch <- entry:
		return true
	default:
		// 缓冲区满，丢弃
		log.Printf("[actlog] Buffer full, dropping log: %s", entry.Message)
		return false
	}
}

// Close 关闭缓冲区
func (b *Buffer) Close() error {
	b.closeMu.Lock()
	if b.closed {
		b.closeMu.Unlock()
		return nil
	}
	b.closed = true
	b.closeMu.Unlock()

	close(b.done)
	b.wg.Wait()
	close(b.ch)

	return b.writer.Close()
}
