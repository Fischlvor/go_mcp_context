// Package bufferedwriter 提供通用的异步缓冲写入器
//
// 使用示例:
//
//	// 1. 实现 Flusher 接口
//	type MyFlusher struct{}
//	func (f *MyFlusher) Flush(batch []any) error { ... }
//	func (f *MyFlusher) Name() string { return "my-writer" }
//
//	// 2. 创建缓冲区
//	buf := bufferedwriter.New(&MyFlusher{}, bufferedwriter.Config{
//	    Size:     100,
//	    Batch:    10,
//	    Interval: 5 * time.Second,
//	})
//
//	// 3. 写入数据
//	buf.Write(myEntry)
//
//	// 4. 关闭时刷新
//	buf.Close()
package bufferedwriter

import (
	"log"
	"sync"
	"time"
)

// Config 缓冲区配置
type Config struct {
	Size     int           // 缓冲区大小（channel 容量）
	Batch    int           // 批量写入数量阈值
	Interval time.Duration // 定时刷新间隔
}

// DefaultConfig 默认配置
var DefaultConfig = Config{
	Size:     1000,
	Batch:    50,
	Interval: 2 * time.Second,
}

// Writer 写入器接口，由具体实现提供
type Writer[T any] interface {
	// WriteBatch 批量写入数据
	WriteBatch(batch []T) error
	// Close 关闭写入器
	Close() error
}

// Buffer 通用异步缓冲区
type Buffer[T any] struct {
	writer  Writer[T]
	name    string
	config  Config
	ch      chan T
	done    chan struct{}
	wg      sync.WaitGroup
	closed  bool
	closeMu sync.Mutex
}

// New 创建异步缓冲区
func New[T any](name string, writer Writer[T], config Config) *Buffer[T] {
	b := &Buffer[T]{
		writer: writer,
		name:   name,
		config: config,
		ch:     make(chan T, config.Size),
		done:   make(chan struct{}),
	}
	b.start()
	return b
}

// start 启动后台处理协程
func (b *Buffer[T]) start() {
	b.wg.Add(1)
	go func() {
		defer b.wg.Done()

		batch := make([]T, 0, b.config.Batch)
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
func (b *Buffer[T]) flush(batch []T) {
	if len(batch) == 0 {
		return
	}

	if err := b.writer.WriteBatch(batch); err != nil {
		log.Printf("[%s] Failed to flush %d entries: %v", b.name, len(batch), err)
	}
}

// Write 写入数据（缓冲区满时阻塞等待）
func (b *Buffer[T]) Write(entry T) bool {
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
		log.Printf("[%s] Buffer full, dropping entry", b.name)
		return false
	}
}

// Close 关闭缓冲区
func (b *Buffer[T]) Close() error {
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

	return nil
}
