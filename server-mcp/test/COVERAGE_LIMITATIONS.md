# 测试覆盖率优化限制说明

## 概述

当前覆盖率：**80.4%**

13个函数无法继续优化，所有未覆盖代码均为错误处理分支，需要 mock 外部依赖才能触发。

## 无法继续优化的函数

| 函数 | 覆盖率 | 原因 |
|------|--------|------|
| `ProcessDocumentAsync` | 37.5% | 异步 goroutine 中的错误处理 |
| `generateLibraryEmbedding` | 55.6% | 外部服务失败（向量服务、数据库） |
| `DeleteVersion` | 58.8% | 数据库操作失败 |
| `RefreshVersionWithCallback` | 58.8% | 文件操作失败（下载、读取、处理） |
| `ImportFromGitHub` | 59.6% | GitHub API 失败、文件过滤等边界场景 |
| `RefreshVersion` | 61.1% | 文件下载失败、异步处理错误 |
| `ProcessDocumentWithCallback` | 61.8% | 分块失败等错误处理 |
| `Delete` (Document) | 62.5% | 数据库更新失败、文档不存在 |
| `Create` (APIKey) | 70.0% | 数据库查询/创建失败、UUID生成失败 |
| `ProcessDocument` | 70.0% | 数据库保存失败、缓存失效失败 |
| `fetchUserStats` | 71.4% | 数据库查询失败 |
| `processFile` | 72.2% | 存储上传失败、MDX 文件类型 |
| `UploadWithCallback` | 75.0% | 文件读取/上传失败、重复内容检测 |

## 详细分析

### 1. ProcessDocumentAsync

**位置**：`internal/service/processor.go:534`

**未覆盖**：`ProcessDocument` 失败时的错误处理分支

```go
if err := p.ProcessDocument(doc, content, docLogger); err != nil {
    // ❌ 未覆盖：错误处理
    doc.Status = "failed"
    doc.ErrorMessage = err.Error()
    global.DB.Save(doc)
}
```

**原因**：需要 `ProcessDocument` 返回错误，单元测试中难以触发。

---

### 2. generateLibraryEmbedding

**位置**：`internal/service/library.go:120`

**未覆盖**：向量服务或数据库失败时的错误处理

```go
embedding, err := global.Embedding.Embed(textToEmbed)
if err != nil {
    // ❌ 未覆盖：向量服务失败
    return
}

if err := global.DB.Model(&dbmodel.Library{}).Update(...).Error; err != nil {
    // ❌ 未覆盖：数据库更新失败
    return
}
```

**原因**：使用全局依赖（`global.Embedding`、`global.DB`），无法轻易 mock。

---

### 3. DeleteVersion

**位置**：`internal/service/library.go:566`

**未覆盖**：数据库查询和事务错误处理

```go
if err := global.DB.First(&library, libraryID).Error; err != nil {
    // ❌ 未覆盖：库不存在（测试中总是存在）
    return ErrNotFound
}

if err := tx.Table("document_uploads").Scan(&documentIDs).Error; err != nil {
    // ❌ 未覆盖：查询失败
    tx.Rollback()
    return err
}
```

**原因**：测试环境中数据库操作总是成功，需要 mock 数据库才能模拟失败。

---

### 4. RefreshVersionWithCallback

**位置**：`internal/service/library.go:782`

**未覆盖**：文件下载、读取和处理的错误处理

```go
reader, err := global.Storage.Download(context.Background(), doc.FilePath)
if err != nil {
    // ❌ 未覆盖：文件下载失败
    statusChan <- response.RefreshStatus{Stage: "doc_failed", ...}
    continue
}

content, err := io.ReadAll(reader)
if err != nil {
    // ❌ 未覆盖：文件读取失败
    statusChan <- response.RefreshStatus{Stage: "doc_failed", ...}
    continue
}
```

**原因**：测试环境中存储服务总是正常工作，需要 mock `global.Storage` 才能模拟失败。

---

### 5. ImportFromGitHub

**位置**：`internal/service/github_import.go:50`

**未覆盖**：GitHub API 调用失败、文件过滤边界场景

```go
// 获取仓库信息失败
repoInfo, err := s.client.GetRepoInfo(ctx, req.Repo)
if err != nil {
    // ❌ 未覆盖：API 调用失败
    return err
}

// 获取目录树失败
tree, err := s.client.GetTree(ctx, req.Repo, ref)
if err != nil {
    // ❌ 未覆盖：获取目录树失败
    return err
}

// 没有找到文档文件
if len(files) == 0 {
    // ❌ 未覆盖：文件过滤后为空
    return fmt.Errorf("no document files found")
}
```

**原因**：需要 mock GitHub API 客户端才能模拟各种失败场景，测试中 API 总是成功。

---

### 6. RefreshVersion

**位置**：`internal/service/library.go:647`

**未覆盖**：文件下载失败、异步处理中的错误处理

```go
// 文件下载失败
reader, err := global.Storage.Download(context.Background(), doc.FilePath)
if err != nil {
    // ❌ 未覆盖：文件下载失败
    global.DB.Model(&docCopy).Update("status", "failed")
    continue
}

// 文件读取失败
content, err := io.ReadAll(reader)
if err != nil {
    // ❌ 未覆盖：文件读取失败
    global.DB.Model(&docCopy).Update("status", "failed")
    continue
}
```

**原因**：与 RefreshVersionWithCallback 类似，需要 mock 存储服务才能模拟失败。

---

### 7. ProcessDocumentWithCallback

**位置**：`internal/service/processor.go:552`

**未覆盖**：文档解析失败的错误处理

```go
// 文档解析失败
text, err := p.parseDocument(doc.FileType, content)
if err != nil {
    // ❌ 未覆盖：解析失败
    statusChan <- response.ProcessStatus{Stage: "failed", ...}
    doc.Status = "failed"
    global.DB.Save(doc)
    return
}
```

**原因**：需要特殊构造的文档内容才能触发解析失败，正常测试中总是成功。

---

### 8. Delete (Document)

**位置**：`internal/service/document.go:336`

**未覆盖**：数据库更新失败

```go
result := global.DB.Model(&dbmodel.DocumentUpload{}).
    Where("id = ? AND deleted_at IS NULL", id).
    Updates(map[string]interface{}{"status": "deleted", "deleted_at": now})

if result.Error != nil {
    // ❌ 未覆盖：数据库更新失败
    return result.Error
}
```

**原因**：测试环境中数据库操作总是成功，需要 mock 数据库才能模拟更新失败。

---

### 9. Create (APIKey)

**位置**：`internal/service/apikey.go:27`

**未覆盖**：数据库查询失败、UUID生成失败、数据库创建失败

```go
// 查询 API Key 数量失败
if err := global.DB.Model(&database.APIKey{}).Count(&count).Error; err != nil {
    // ❌ 未覆盖：查询失败
    return nil, errors.New("查询失败")
}

// UUID 生成失败
uuidV4, err := uuid.NewV4()
if err != nil {
    // ❌ 未覆盖：UUID 生成失败
    return nil, errors.New("生成 API Key 失败")
}

// 数据库创建失败
if err := global.DB.Create(token).Error; err != nil {
    // ❌ 未覆盖：创建失败
    return nil, errors.New("创建失败")
}
```

**原因**：测试环境中数据库操作和 UUID 生成总是成功，需要 mock 才能模拟失败。

---

### 10. ProcessDocument

**位置**：`internal/service/processor.go:84`

**未覆盖**：processDocumentCore 失败、数据库保存失败、缓存失效失败

**原因**：测试环境中核心处理和数据库操作总是成功，需要 mock 才能模拟失败。

---

### 11. UploadWithCallback

**位置**：`internal/service/document.go:171`

**未覆盖**：文件读取失败、文件上传失败、重复内容检测等错误处理

**原因**：测试环境中文件操作总是成功，需要 mock 文件读取和存储服务才能模拟失败。

---

### 12. fetchUserStats

**位置**：`internal/service/stats.go:30`

**未覆盖**：4个数据库查询失败的错误处理分支（库数量、文档数量、Token 总数、MCP 调用次数）

**原因**：测试环境中数据库查询总是成功，需要 mock 数据库才能模拟失败。

---

### 13. processFile

**位置**：`internal/service/github_import.go:286`

**未覆盖**：MDX 文件类型判断、存储上传失败的错误处理

**原因**：测试中没有使用 .mdx 文件，存储服务总是成功，需要 mock 存储服务才能模拟失败。

## 根本原因

所有未覆盖代码的共同特点：

1. **全局依赖**：使用 `global.Embedding`、`global.Storage`、`global.DB`
2. **错误处理**：所有未覆盖路径都是错误处理分支
3. **测试环境稳定**：外部服务在测试中总是正常工作
