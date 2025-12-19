package embedding

import (
	"context"
	"crypto/md5"
	"encoding/hex"
	"encoding/json"
	"time"

	"github.com/redis/go-redis/v9"
	"go.uber.org/zap"
)

const (
	embeddingCachePrefix = "embedding:query:"
	embeddingCacheTTL    = 24 * time.Hour // 缓存 24 小时
)

// CachedEmbeddingService 带缓存的 Embedding 服务包装器
type CachedEmbeddingService struct {
	inner  EmbeddingService // 被包装的实际 embedding 服务
	redis  *redis.Client    // Redis 客户端
	logger *zap.Logger      // 日志
	ctx    context.Context  // 上下文
}

// NewCachedEmbeddingService 创建带缓存的 Embedding 服务
func NewCachedEmbeddingService(inner EmbeddingService, redis *redis.Client, logger *zap.Logger) *CachedEmbeddingService {
	return &CachedEmbeddingService{
		inner:  inner,
		redis:  redis,
		logger: logger,
		ctx:    context.Background(),
	}
}

// Embed 生成 embedding（带缓存）
func (c *CachedEmbeddingService) Embed(text string) ([]float32, error) {
	// 1. 生成缓存 key（使用 text 的 MD5）
	hash := md5.Sum([]byte(text))
	cacheKey := embeddingCachePrefix + hex.EncodeToString(hash[:])

	// 2. 尝试从 Redis 获取缓存
	if c.redis != nil {
		cached, err := c.redis.Get(c.ctx, cacheKey).Result()
		if err == nil && cached != "" {
			// 缓存命中，解析 JSON
			var vector []float32
			if err := json.Unmarshal([]byte(cached), &vector); err == nil {
				if c.logger != nil {
					c.logger.Debug("Embedding cache hit", zap.String("text_prefix", truncate(text, 50)))
				}
				return vector, nil
			}
		}
	}

	// 3. 缓存未命中，调用实际的 embedding 服务
	if c.logger != nil {
		c.logger.Debug("Embedding cache miss, calling API", zap.String("text_prefix", truncate(text, 50)))
	}
	vector, err := c.inner.Embed(text)
	if err != nil {
		return nil, err
	}

	// 4. 存入 Redis 缓存
	if c.redis != nil {
		if data, err := json.Marshal(vector); err == nil {
			c.redis.Set(c.ctx, cacheKey, string(data), embeddingCacheTTL)
		}
	}

	return vector, nil
}

// EmbedBatch 批量生成 embedding（带缓存）
func (c *CachedEmbeddingService) EmbedBatch(texts []string) ([][]float32, error) {
	results := make([][]float32, len(texts))
	var uncachedTexts []string
	var uncachedIndices []int

	// 1. 检查缓存
	for i, text := range texts {
		hash := md5.Sum([]byte(text))
		cacheKey := embeddingCachePrefix + hex.EncodeToString(hash[:])

		if c.redis != nil {
			cached, err := c.redis.Get(c.ctx, cacheKey).Result()
			if err == nil && cached != "" {
				var vector []float32
				if err := json.Unmarshal([]byte(cached), &vector); err == nil {
					results[i] = vector
					continue
				}
			}
		}
		// 缓存未命中
		uncachedTexts = append(uncachedTexts, text)
		uncachedIndices = append(uncachedIndices, i)
	}

	// 2. 批量调用未缓存的
	if len(uncachedTexts) > 0 {
		if c.logger != nil {
			c.logger.Debug("Embedding batch cache miss",
				zap.Int("cached", len(texts)-len(uncachedTexts)),
				zap.Int("uncached", len(uncachedTexts)))
		}

		vectors, err := c.inner.EmbedBatch(uncachedTexts)
		if err != nil {
			return nil, err
		}

		// 3. 填充结果并缓存
		for j, vector := range vectors {
			idx := uncachedIndices[j]
			results[idx] = vector

			// 存入缓存
			if c.redis != nil {
				text := uncachedTexts[j]
				hash := md5.Sum([]byte(text))
				cacheKey := embeddingCachePrefix + hex.EncodeToString(hash[:])
				if data, err := json.Marshal(vector); err == nil {
					c.redis.Set(c.ctx, cacheKey, string(data), embeddingCacheTTL)
				}
			}
		}
	}

	return results, nil
}

// GetDimension 返回 embedding 维度
func (c *CachedEmbeddingService) GetDimension() int {
	return c.inner.GetDimension()
}

// GetModelName 返回模型名称
func (c *CachedEmbeddingService) GetModelName() string {
	return c.inner.GetModelName()
}

// GetMaxBatchSize 返回最大批量大小
func (c *CachedEmbeddingService) GetMaxBatchSize() int {
	return c.inner.GetMaxBatchSize()
}

// truncate 截断字符串
func truncate(s string, maxLen int) string {
	if len(s) <= maxLen {
		return s
	}
	return s[:maxLen] + "..."
}
