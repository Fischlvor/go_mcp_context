package utils

import (
	"math/rand"
	"sync"
	"time"

	"github.com/oklog/ulid/v2"
)

var (
	entropyMu sync.Mutex
	entropy   = rand.New(rand.NewSource(time.Now().UnixNano()))
)

// GenerateTaskID 生成唯一任务ID（ULID格式）
// 返回26字符的 ULID 字符串，按时间可排序
func GenerateTaskID() string {
	entropyMu.Lock()
	id := ulid.MustNew(ulid.Timestamp(time.Now()), entropy)
	entropyMu.Unlock()
	return id.String()
}
