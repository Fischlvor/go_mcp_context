# 混合搜索架构文档

## 概述

本文档描述了 go-mcp-context 项目中的混合搜索系统架构，该系统结合了向量搜索和BM25关键词搜索，使用RRF（Reciprocal Rank Fusion）算法进行结果融合，为MCP文档检索提供高质量的搜索结果。

## 核心架构

### 搜索流程

```
用户查询 + 搜索参数
├─ libraryID (必需)
├─ version (必需)
├─ mode (可选: code/info)
└─ query (必需，但某些API支持空字符串浏览模式)
    ↓
查询预处理（topic分割）
    ↓
搜索结果缓存检查
├─ 缓存Key: search:topic:{libraryID}:{version}:{mode}:{topic_hash}
├─ 命中缓存 → 直接返回结果
└─ 缓存未命中 ↓
    ↓
并行召回阶段（限定库范围）
├─ 向量搜索 (Top-50)
│  ├─ Embedding缓存检查
│  ├─ 缓存未命中 → 生成查询向量
│  └─ pgvector相似度搜索 (WHERE library_id={libraryID} AND version={version})
└─ BM25搜索 (Top-50)
   └─ PostgreSQL全文搜索 (WHERE library_id={libraryID} AND version={version})
    ↓
得分归一化阶段
├─ Min-Max归一化向量得分
└─ Min-Max归一化BM25得分
    ↓
RRF融合排序
├─ 基于排名的RRF算法
├─ 向量权重: 0.7
├─ BM25权重: 0.3
└─ 热度权重: 0.2
    ↓
搜索结果缓存存储 & 分页返回结果
```

## 技术实现

### 1. 搜索入口

**主要方法**: `SearchDocuments(req *request.Search)`

**处理逻辑**:
- 单查询: 直接使用混合RRF搜索
- 多查询: 并行搜索各topic后RRF合并

### 2. 并行召回

**向量搜索**: `vectorSearch()`
- 使用pgvector进行cosine相似度计算
- 召回Top-50相关文档
- 得分范围: [0, 1]

**BM25搜索**: `bm25Search()`
- 使用PostgreSQL全文搜索
- ts_rank BM25评分算法
- 召回Top-50相关文档
- 得分范围: 动态变化

### 3. 得分归一化

**Min-Max归一化**: `normalizeScoresMinMax()`
- 将不同搜索方式的得分归一化到[0,1]范围
- 公式: `(score - min) / (max - min)`
- 避免某种搜索方式占主导地位

### 4. RRF融合算法

**混合RRF**: `hybridRRF()`
- 基于排名而非绝对得分进行融合
- RRF公式: `score = Σ weight / (rank + k)`
- RRF常量k = 60 (Elasticsearch默认值)

**权重配置**:
```go
VectorRRFWeight = 0.7  // 向量搜索权重
BM25RRFWeight   = 0.3  // BM25搜索权重
HotWeight       = 0.2  // 热度权重
```

### 5. 缓存机制

**多层缓存**:
- Embedding缓存: 查询向量生成结果
- 搜索结果缓存: 完整搜索结果 (24小时TTL)
- Tag感知失效: 库版本更新时自动失效相关缓存

## 关键优势

### 1. 解决得分冲突
- **问题**: 向量搜索和BM25搜索得分范围不一致
- **解决**: RRF基于排名融合，不依赖绝对得分值
- **效果**: 避免某种搜索方式占主导地位

### 2. 提高搜索质量
- **语义理解**: 向量搜索捕获语义相似性
- **精确匹配**: BM25搜索处理关键词精确匹配
- **智能融合**: RRF算法平衡两种搜索方式的优势

### 3. 性能优化
- **并行召回**: 向量搜索和BM25搜索并行执行
- **缓存加速**: 多层缓存减少重复计算
- **轻量级融合**: RRF算法计算简单高效

### 4. 业界验证
- **Elasticsearch**: 8.8+版本默认使用RRF
- **OpenSearch**: 2.11+版本提供RRF支持
- **Azure AI Search**: RRF是混合查询的标准算法

## 配置参数

### 搜索参数
```go
// 召回数量
VectorTopK = 50    // 向量搜索召回数
BM25TopK   = 50    // BM25搜索召回数

// RRF参数
RRFConstant = 60   // RRF算法常量k

// 权重配置
VectorRRFWeight = 0.7  // 向量搜索权重
BM25RRFWeight   = 0.3  // BM25搜索权重
HotWeight       = 0.2  // 热度权重
```

### 缓存配置
```go
SearchCacheTTL = 24 * time.Hour  // 搜索结果缓存时长
```

## 性能指标

### 搜索延迟
- 向量搜索: ~20ms
- BM25搜索: ~30ms  
- RRF融合: ~5ms
- **总延迟**: ~50ms (并行执行)

## 相关文档

- [缓存架构文档](./CACHE.md)
- [MCP接口文档](./MCP.md)
- [API接口文档](./API.md)

---

*最后更新: 2025-12-26*
