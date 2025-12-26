# Redis 缓存架构指南

## 概述

go-mcp-context 项目使用 Redis 作为缓存层来显著提高频繁访问的文档查询的搜索性能。本文档详细介绍了项目中的缓存架构、实现方式和最佳实践。

## 缓存架构

### 1. 缓存层次结构

```
┌─────────────────────────────────────────────────────────────┐
│                    应用层 (API)                              │
├─────────────────────────────────────────────────────────────┤
│                   服务层 (Service)                           │
├─────────────────────────────────────────────────────────────┤
│                   缓存层 (Cache)                             │
│  ┌─────────────────┐  ┌─────────────────┐  ┌──────────────┐ │
│  │   搜索结果缓存   │  │  Embedding缓存  │  │  Tag版本缓存  │ │
│  │  (24小时TTL)    │  │  (24小时TTL)    │  │   (永久)     │ │
│  └─────────────────┘  └─────────────────┘  └──────────────┘ │
├─────────────────────────────────────────────────────────────┤
│                   Redis 存储                                │
└─────────────────────────────────────────────────────────────┘
```

### 2. 缓存类型

#### 2.1 搜索结果缓存
- **用途**: 缓存文档搜索的结果，避免重复的向量搜索和BM25搜索
- **TTL**: 24小时
- **失效策略**: Tag版本失效机制

#### 2.2 Embedding 缓存
- **用途**: 缓存文本的向量表示，避免重复调用OpenAI API
- **TTL**: 24小时
- **Key格式**: `embedding:query:{md5_hash}`

#### 2.3 Tag版本缓存
- **用途**: 实现O(1)批量缓存失效
- **TTL**: 永久（通过版本号递增实现失效）
- **Key格式**: `tag:version:{tag_name}`

## 核心接口设计

### 1. 基础缓存接口

```go
// Cache 基础缓存接口
type Cache interface {
    // Get 从缓存获取值，反序列化到 dest
    Get(key string, dest interface{}) error
    // Set 设置缓存值，ttl 为过期时间
    Set(key string, value interface{}, ttl time.Duration) error
    // Delete 删除指定 key
    Delete(key string) error
    // Exists 检查 key 是否存在
    Exists(key string) (bool, error)
    // Clear 清除指定前缀的所有缓存
    Clear(prefix string) error
    // Close 关闭缓存连接
    Close() error
}
```

### 2. Tag感知缓存接口

```go
// TagAwareCache 带 Tag 版本的缓存接口
// 通过 Version Tag 模式实现 O(1) 缓存失效
type TagAwareCache interface {
    Cache
    
    // GetTagVersion 获取 tag 的当前版本号（不存在返回 0）
    GetTagVersion(tag string) (int64, error)
    
    // InvalidateTags 递增 tag 版本号，使关联的缓存失效
    InvalidateTags(tags []string) error
}
```

## Redis 实现

### 1. Redis 配置

```yaml
# config.yaml
redis:
  address: "localhost:6379"
  password: ""
  db: 3

cache:
  ttl: 24h
  prefix: "mcp:"
```

### 2. Redis 客户端初始化

```go
// internal/initialize/redis.go
func InitRedis() {
    config := global.Config.Redis
    
    client := redis.NewClient(&redis.Options{
        Addr:     config.Address,
        Password: config.Password,
        DB:       config.DB,
    })
    
    // 测试连接
    ctx := context.Background()
    if err := client.Ping(ctx).Err(); err != nil {
        global.Log.Fatal("Redis connection failed", zap.Error(err))
    }
    
    // 创建缓存实例
    global.Cache = cache.NewRedisCacheWithClient(client, global.Config.Cache.Prefix)
    global.Log.Info("Redis initialized successfully")
}
```

### 3. RedisCache 实现

```go
// pkg/cache/redis.go
type RedisCache struct {
    client *redis.Client
    prefix string
    ctx    context.Context
}

// Get 从缓存获取值，反序列化到 dest
func (c *RedisCache) Get(key string, dest interface{}) error {
    data, err := c.client.Get(c.ctx, c.prefix+key).Bytes()
    if err == redis.Nil {
        return ErrCacheMiss
    }
    if err != nil {
        return err
    }
    return json.Unmarshal(data, dest)
}

// Set 设置缓存值（带 TTL）
func (c *RedisCache) Set(key string, value interface{}, ttl time.Duration) error {
    data, err := json.Marshal(value)
    if err != nil {
        return err
    }
    return c.client.Set(c.ctx, c.prefix+key, data, ttl).Err()
}
```

## 搜索结果缓存

### 1. 缓存键设计

```go
// 搜索缓存 key 格式
// search:topic:{library_id}:{version}:{mode}:{topic_hash}
func (s *SearchService) buildSearchCacheKey(libraryID uint, version, mode, topic string) string {
    hash := md5.Sum([]byte(topic))
    topicHash := hex.EncodeToString(hash[:])
    return fmt.Sprintf("%s%d:%s:%s:%s", SearchCachePrefix, libraryID, version, mode, topicHash)
}

// 搜索缓存 tag 格式
// library:{library_id}:{version}
func (s *SearchService) buildSearchCacheTag(libraryID uint, version string) string {
    return fmt.Sprintf("library:%d:%s", libraryID, version)
}
```

### 2. 缓存使用示例

```go
// internal/service/search.go
func (s *SearchService) searchSingleTopic(ctx context.Context, req *request.Search, topic string) ([]searchCandidate, error) {
    // 生成缓存 key
    cacheKey := s.buildSearchCacheKey(req.LibraryID, req.Version, req.Mode, topic)
    
    // 生成缓存 tag
    cacheTag := s.buildSearchCacheTag(req.LibraryID, req.Version)
    
    // 定义搜索函数
    fetchFunc := func() ([]searchCandidate, error) {
        return s.executeSearch(ctx, req, topic)
    }
    
    // 使用 GetOrSetWithTags 模式：缓存 key 包含 tag version，tag 失效时旧缓存自动失效
    return cache.GetOrSetWithTags(global.Cache, cacheKey, []string{cacheTag}, SearchCacheTTL, fetchFunc)
}
```

### 3. 缓存失效机制

```go
// 文档处理完成后，失效相关搜索缓存
func (s *SearchService) InvalidateLibraryCache(libraryID uint, version string) error {
    if global.Cache == nil {
        return nil
    }
    
    tag := s.buildSearchCacheTag(libraryID, version)
    return global.Cache.InvalidateTags([]string{tag})
}
```

## Embedding 缓存

### 1. 缓存包装器

```go
// pkg/embedding/cached.go
type CachedEmbeddingService struct {
    inner  EmbeddingService // 被包装的实际 embedding 服务
    redis  *redis.Client    // Redis 客户端
    logger *zap.Logger      // 日志
    ctx    context.Context  // 上下文
}
```

### 2. 单个文本 Embedding 缓存

```go
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
                c.logger.Debug("Embedding cache hit", zap.String("text_prefix", truncate(text, 50)))
                return vector, nil
            }
        }
    }
    
    // 3. 缓存未命中，调用实际的 embedding 服务
    c.logger.Debug("Embedding cache miss, calling API", zap.String("text_prefix", truncate(text, 50)))
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
```

### 3. 批量 Embedding 缓存

```go
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
```

## Tag版本失效机制

### 1. 工作原理

Tag版本失效机制通过维护每个tag的版本号来实现O(1)批量缓存失效：

1. **缓存存储**: 缓存key包含tag的版本号
2. **版本维护**: 每个tag在Redis中维护一个版本号
3. **失效操作**: 递增tag版本号，使所有包含旧版本号的缓存key自动失效

### 2. 实现细节

```go
// 构建带 tag 版本的缓存 key
func BuildTaggedKey(c TagAwareCache, baseKey string, tags []string) (string, error) {
    if len(tags) == 0 {
        return baseKey, nil
    }
    
    key := baseKey + ":tv"
    for _, tag := range tags {
        version, err := c.GetTagVersion(tag)
        if err != nil {
            return "", err
        }
        key += fmt.Sprintf(":%d", version)
    }
    return key, nil
}

// 带 Tag 版本的缓存旁路模式
func GetOrSetWithTags[T any](c TagAwareCache, baseKey string, tags []string, ttl time.Duration, fetchFunc func() (T, error)) (T, error) {
    var result T
    
    if c == nil {
        return fetchFunc()
    }
    
    // 构建带 tag version 的 key
    key, err := BuildTaggedKey(c, baseKey, tags)
    if err != nil {
        return fetchFunc()
    }
    
    // 尝试从缓存获取
    if err := c.Get(key, &result); err == nil {
        return result, nil
    }
    
    // 缓存未命中，调用 fetchFunc
    data, err := fetchFunc()
    if err != nil {
        return result, err
    }
    
    // 写入缓存
    _ = c.Set(key, data, ttl)
    
    return data, nil
}
```

### 3. 失效示例

```go
// 示例：库版本更新后失效相关缓存
tag := "library:6:v1.0"

// 失效前：tag version = 3
// 缓存 key: "search:topic:6:v1.0:code:abc123:tv:3"

// 执行失效
cache.InvalidateTags([]string{tag})

// 失效后：tag version = 4
// 新的缓存 key: "search:topic:6:v1.0:code:abc123:tv:4"
// 旧缓存自动失效
```

## 性能优化

### 1. 缓存命中率优化

```go
// 使用合理的缓存 TTL
const (
    SearchCacheTTL    = 24 * time.Hour  // 搜索结果缓存24小时
    EmbeddingCacheTTL = 24 * time.Hour  // Embedding缓存24小时
)

// 使用 MD5 哈希减少 key 长度
func buildCacheKey(text string) string {
    hash := md5.Sum([]byte(text))
    return "embedding:query:" + hex.EncodeToString(hash[:])
}
```

### 2. 批量操作优化

```go
// 批量检查缓存，减少网络往返
func (c *CachedEmbeddingService) EmbedBatch(texts []string) ([][]float32, error) {
    // 先批量检查所有缓存
    // 只对未命中的文本调用 API
    // 批量写入新的缓存结果
}
```

### 3. 连接池优化

```go
// Redis 连接池配置
client := redis.NewClient(&redis.Options{
    Addr:         config.Address,
    Password:     config.Password,
    DB:           config.DB,
    PoolSize:     10,              // 连接池大小
    MinIdleConns: 5,               // 最小空闲连接
    MaxRetries:   3,               // 最大重试次数
    DialTimeout:  5 * time.Second, // 连接超时
    ReadTimeout:  3 * time.Second, // 读取超时
    WriteTimeout: 3 * time.Second, // 写入超时
})
```

## 监控和调试

### 1. 缓存统计

```go
// 添加缓存统计功能
type CacheStats struct {
    Hits       int64
    Misses     int64
    Sets       int64
    Deletes    int64
    Errors     int64
}

// 在缓存操作中记录统计
func (c *RedisCache) Get(key string, dest interface{}) error {
    data, err := c.client.Get(c.ctx, c.prefix+key).Bytes()
    if err == redis.Nil {
        atomic.AddInt64(&c.stats.Misses, 1)
        return ErrCacheMiss
    }
    if err != nil {
        atomic.AddInt64(&c.stats.Errors, 1)
        return err
    }
    atomic.AddInt64(&c.stats.Hits, 1)
    return json.Unmarshal(data, dest)
}
```

### 2. 日志记录

```go
// 记录缓存操作日志
func (c *CachedEmbeddingService) Embed(text string) ([]float32, error) {
    if cached {
        c.logger.Debug("Embedding cache hit", 
            zap.String("text_prefix", truncate(text, 50)),
            zap.String("cache_key", cacheKey))
    } else {
        c.logger.Debug("Embedding cache miss, calling API", 
            zap.String("text_prefix", truncate(text, 50)),
            zap.String("cache_key", cacheKey))
    }
}
```

### 3. 健康检查

```go
// 添加缓存健康检查
func (c *RedisCache) HealthCheck() error {
    ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
    defer cancel()
    
    return c.client.Ping(ctx).Err()
}
```

## 最佳实践

### 1. 缓存键命名规范

```go
// 使用层次化的键命名
const (
    SearchCachePrefix    = "search:topic:"     // 搜索结果缓存
    EmbeddingCachePrefix = "embedding:query:"  // Embedding缓存
    TagVersionPrefix     = "tag:version:"      // Tag版本缓存
    BlacklistPrefix      = "blacklist:"       // JWT黑名单
)
```

### 2. 错误处理

```go
// 缓存操作失败不应影响主业务逻辑
func GetOrSet[T any](c Cache, key string, ttl time.Duration, fetchFunc func() (T, error)) (T, error) {
    var result T
    
    // 尝试从缓存获取，失败则继续执行业务逻辑
    if c != nil {
        if err := c.Get(key, &result); err == nil {
            return result, nil
        }
    }
    
    // 缓存未命中，调用 fetchFunc
    data, err := fetchFunc()
    if err != nil {
        return result, err
    }
    
    // 写入缓存，失败不影响返回结果
    if c != nil {
        _ = c.Set(key, data, ttl)
    }
    
    return data, nil
}
```

### 3. 缓存预热

```go
// 在应用启动时预热热门查询
func (s *SearchService) WarmupCache() {
    popularQueries := []string{
        "data fetching",
        "routing",
        "authentication",
        "state management",
    }
    
    for _, query := range popularQueries {
        go func(q string) {
            // 预热缓存
            s.searchSingleTopic(context.Background(), &request.Search{
                LibraryID: 1,
                Version:   "latest",
                Mode:      "code",
                Query:     q,
            }, q)
        }(query)
    }
}
```

### 4. 缓存清理

```go
// 定期清理过期的tag版本
func (c *RedisCache) CleanupExpiredTags() error {
    // 扫描所有 tag:version: 键
    iter := c.client.Scan(c.ctx, 0, c.prefix+tagVersionKeyPrefix+"*", 100).Iterator()
    
    for iter.Next(c.ctx) {
        key := iter.Val()
        // 检查是否有关联的缓存，如果没有则删除tag版本
        // 这里可以根据业务需求实现具体的清理逻辑
    }
    
    return iter.Err()
}
```

## 总结

go-mcp-context 项目通过精心设计的Redis缓存架构，实现了：

1. **高性能搜索**: 搜索结果缓存避免重复的向量搜索和BM25搜索
2. **成本优化**: Embedding缓存减少OpenAI API调用成本
3. **智能失效**: Tag版本机制实现O(1)批量缓存失效
4. **高可用性**: 缓存失败不影响主业务逻辑
5. **易于监控**: 完整的统计和日志记录

这套缓存架构为文档查询提供了显著的性能提升，同时保持了系统的稳定性和可维护性。
