# 测试文档

## 📋 概述

**测试覆盖率**：81.0% ✅（目标：80%+）

**测试框架**：Go 标准测试框架 + 真实数据库环境

---

## 🔒 数据库隔离

| 环境 | 数据库 | Redis DB |
|------|--------|----------|
| **生产** | `mcp_context` | DB 3 |
| **测试** | `mcp_context_test` | DB 15 |

**保证**：测试 100% 不会影响生产数据

---

## 📁 目录结构

```
test/
├── README.md                    # 本文档
├── COVERAGE_LIMITATIONS.md      # 覆盖率限制说明（14个无法优化的函数）
├── Makefile                     # 测试命令
├── coverage.out                 # 覆盖率数据
├── all_functions_coverage.txt   # 函数覆盖率报告
├── test_log.txt                 # 测试日志
│
├── unit/                        # 单元测试（11个文件）
│   ├── setup_test.go
│   ├── library_test.go
│   ├── document_test.go
│   ├── processor_test.go
│   ├── search_test.go
│   ├── mcp_test.go
│   ├── mcp_handler_test.go
│   ├── apikey_test.go
│   ├── stats_test.go
│   ├── activitylog_test.go
│   └── github_import_test.go
│
└── integration/                 # 集成测试
    ├── setup_integration_test.go
    ├── github_import_integration_test.go
    └── mcp_handler_integration_test.go
```

---

## 🚀 快速开始

### 使用 Makefile（推荐）

```bash
cd test

# 运行所有测试并生成覆盖率
make all

# 只运行单元测试
make test-unit

# 运行指定测试
make test-unit TEST=Test_Library_Create

# 查看覆盖率
make show-coverage

# 查看测试日志
make show-log

# 清理生成的文件
make clean
```

### 直接使用 go test

```bash
# 运行所有单元测试
go test ./test/unit/... -v

# 运行指定测试
go test ./test/unit/... -v -run Test_Library_Create

# 生成覆盖率
go test ./test/unit/... -v -coverprofile=test/coverage.out -coverpkg=./internal/service/...

# 查看覆盖率
go tool cover -func=test/coverage.out
go tool cover -html=test/coverage.out -o test/coverage.html
```

---

## 📊 测试覆盖情况

### 测试文件

| Service | 测试文件 | 状态 |
|---------|---------|------|
| LibraryService | `library_test.go` | ✅ |
| DocumentService | `document_test.go` | ✅ |
| ProcessorService | `processor_test.go` | ✅ |
| SearchService | `search_test.go` | ✅ |
| MCPService | `mcp_test.go` | ✅ |
| MCPHandler | `mcp_handler_test.go` | ✅ |
| ApiKeyService | `apikey_test.go` | ✅ |
| StatsService | `stats_test.go` | ✅ |
| ActivityLogService | `activitylog_test.go` | ✅ |
| GitHubImportService | `github_import_test.go` | ✅ |

**总计**：10个测试文件，所有测试通过

### 覆盖率说明

- **当前覆盖率**：81.0% ✅
- **目标覆盖率**：80%+
- **统计范围**：`internal/service/...`（业务逻辑层）
- **不统计**：`pkg/`、`cmd/`、`internal/handler/`、`internal/middleware/` 等
- **原因**：业务逻辑是核心代码，其他层（如 HTTP handler、中间件）主要是框架代码和路由配置
- **无法优化函数**：14个（详见 [COVERAGE_LIMITATIONS.md](./COVERAGE_LIMITATIONS.md)）

---

## 📝 测试命名规范

### 基础格式

```
Test_{Service}_{Method}
```

### 高级测试后缀

```
Test_{Service}_{Method}_{Suffix}
```

### 示例

| 测试类型 | 命名示例 |
|---------|---------|
| 基础测试 | `Test_Library_Create` |
| 高级测试 | `Test_Library_Create_Advanced` |
| 边界测试 | `Test_Document_Delete_EdgeCases` |
| 集成测试 | `Test_Integration_GitHubImport_RealAPI` |

### 子测试命名

```go
func Test_Library_Create(t *testing.T) {
    libService := &service.LibraryService{}
    
    t.Run("create library with valid data", func(t *testing.T) {
        // 测试代码
    })
    
    t.Run("create library with empty name", func(t *testing.T) {
        // 测试代码
    })
}
```

---

## 🔧 配置说明

### 测试配置文件

**位置**：`configs/config.test.yaml`

**关键配置**：
```yaml
postgres:
  db_name: mcp_context_test  # 测试数据库

redis:
  db: 15                     # 测试 Redis DB
```

---

## ⚠️ 注意事项

### 数据库隔离

- ✅ 测试使用独立的 `mcp_context_test` 数据库
- ✅ 生产数据库 `mcp_context` 完全不受影响
- ✅ 测试完成后数据会保留（方便检查）

### Redis 隔离

- ✅ 测试使用 Redis DB 15
- ✅ 生产使用 Redis DB 3
- ✅ 完全隔离，互不影响

### Makefile 实时输出

- ✅ 使用 `stdbuf -oL -eL` 实现实时日志输出
- ✅ 不再出现日志延迟

---

## 📚 相关文档

- [COVERAGE_LIMITATIONS.md](./COVERAGE_LIMITATIONS.md) - 覆盖率限制说明
- [集成测试文档](./integration/README.md) - 集成测试说明
- [Makefile](./Makefile) - 测试命令定义

---

## ✅ 检查清单

### 运行测试前

- [ ] PostgreSQL 可访问
- [ ] Redis 可访问
- [ ] 测试配置文件存在（`configs/config.test.yaml`）

### 运行测试后

- [ ] 所有测试通过（`make all`）
- [ ] 覆盖率 ≥ 80%（`make show-coverage`）
- [ ] 生产数据库数据完整（未被修改）

---
