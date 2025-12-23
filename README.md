# go-mcp-context

<div align="center">

![Go Version](https://img.shields.io/badge/Go-1.23-blue)
![Vue Version](https://img.shields.io/badge/Vue-3.5-green)
![License](https://img.shields.io/badge/license-MIT-orange)

私有化的 Context7 替代方案，为企业内网的 AI IDE 提供实时、准确的技术文档和代码示例

🌐 **在线体验**: [https://mcp.hsk423.cn](https://mcp.hsk423.cn)

</div>

---

## 📖 项目介绍

go-mcp-context 是一个私有化的文档检索服务，通过 MCP 协议为 AI IDE（如 Cursor、Windsurf、VSCode）提供企业内部技术文档的智能检索能力。

### ✨ 核心特性

- 🔌 **MCP 协议支持** - 标准 MCP 协议接口，支持 IDE 无缝集成
- 🔍 **向量检索** - 基于 PostgreSQL + pgvector 的高性能向量搜索
- 📄 **多格式文档** - 支持 Markdown、PDF、DOCX、Swagger 等格式
- 🔀 **混合搜索** - 向量相似度 + BM25 关键词搜索
- 📊 **智能重排序** - 多指标评分优化搜索结果
- 🔐 **双重认证** - SSO JWT 管理 + API Key MCP 调用
- 🎨 **现代化 UI** - Vue3 + TypeScript + TailwindCSS

## 🛠️ 技术栈

### 后端技术

| 技术 | 版本 | 说明 |
|------|------|------|
| Go | 1.23 | 主要开发语言 |
| Gin | 1.10 | Web 框架 |
| GORM | 1.25 | ORM 框架 |
| PostgreSQL | 15 | 主数据库 + pgvector |
| Redis | 6 | 缓存数据库 |
| OpenAI API | - | Embedding 生成 |
| JWT | - | 身份认证 |
| Zap | 1.27 | 日志框架 |

### 前端技术

| 技术 | 版本 | 说明 |
|------|------|------|
| Vue | 3.5 | 前端框架 |
| TypeScript | 5.x | 类型系统 |
| TailwindCSS | 3.x | CSS 框架 |
| Vite | 6.x | 构建工具 |
| Axios | 1.x | HTTP 客户端 |

### 基础设施

- **容器化**: Docker + Docker Compose
- **向量存储**: PostgreSQL + pgvector 扩展
- **认证**: SSO JWT + API Key

## 🚀 快速开始

### 环境要求

- Go 1.23+
- Docker & Docker Compose
- OpenAI API Key

### 本地开发

```bash
# 克隆项目
git clone <repo-url>
cd go-mcp-context

# 启动依赖服务
docker-compose up -d postgres redis

# 设置环境变量
export OPENAI_API_KEY=your-api-key
export JWT_SECRET=your-jwt-secret

# 运行服务
go run ./cmd/server -config configs/config.yaml
```

### Docker 部署

```bash
# 设置环境变量
export OPENAI_API_KEY=your-api-key
export JWT_SECRET=your-jwt-secret

# 启动所有服务
docker-compose up -d
```

## API 端点

### 健康检查

```
GET /health
```

### 库管理

```
GET    /api/v1/libraries                        # 获取库列表
POST   /api/v1/libraries                        # 创建库
GET    /api/v1/libraries/:id                    # 获取库详情
PUT    /api/v1/libraries/:id                    # 更新库
DELETE /api/v1/libraries/:id                    # 删除库
GET    /api/v1/libraries/github/releases        # 获取 GitHub 仓库版本列表
POST   /api/v1/libraries/:id/github/import-sse  # 从 GitHub 导入文档（SSE）
```

### 文档管理

```
GET    /api/v1/documents/list                           # 获取文档列表
GET    /api/v1/documents/detail/:id                     # 获取文档详情
GET    /api/v1/documents/chunks/:mode/:libid/*version   # 获取库的文档块 (mode: code/info, version 可选)
POST   /api/v1/documents/upload                         # 上传文档
POST   /api/v1/documents/upload-sse                     # 上传文档（SSE 实时状态）
DELETE /api/v1/documents/:id                            # 删除文档
```

### 搜索

```
POST /api/v1/search  # 搜索文档
```

### MCP 端点

```
GET  /mcp/health  # MCP 健康检查
GET  /mcp/tools   # 获取工具列表
POST /mcp         # JSON-RPC 2.0 请求
```

## MCP 工具

### search-libraries

搜索文档库。

```json
{
  "jsonrpc": "2.0",
  "id": 1,
  "method": "tools/call",
  "params": {
    "name": "search-libraries",
    "arguments": {
      "libraryName": "react"
    }
  }
}
```

### get-library-docs

获取库文档。

```json
{
  "jsonrpc": "2.0",
  "id": 2,
  "method": "tools/call",
  "params": {
    "name": "get-library-docs",
    "arguments": {
      "libraryID": "react/18.3.0",
      "topic": "useState hook",
      "mode": "code",
      "page": 1
    }
  }
}
```

## IDE 配置

在 IDE 的 MCP 配置中添加：

```json
{
  "mcpServers": {
    "go-mcp-context": {
      "type": "streamable-http",
      "url": "http://localhost:8080/mcp",
      "headers": {
        "Authorization": "Bearer YOUR_API_TOKEN"
      }
    }
  }
}
```

## 📁 项目结构

```
go-mcp-context/
├── server-mcp/               # MCP 后端服务
│   ├── cmd/                  # 主程序入口
│   ├── configs/              # 配置文件
│   ├── internal/
│   │   ├── api/                  # HTTP 处理器
│   │   │   ├── auth.go           # 认证 API
│   │   │   ├── apikey.go         # API Key 管理
│   │   │   ├── document.go       # 文档管理
│   │   │   ├── library.go        # 库管理
│   │   │   ├── mcp.go            # MCP 端点
│   │   │   ├── search.go         # 搜索 API
│   │   │   └── user.go           # 用户 API
│   │   ├── initialize/           # 初始化模块
│   │   │   ├── gorm.go           # 数据库初始化
│   │   │   ├── redis.go          # Redis 初始化
│   │   │   └── router.go         # 路由初始化
│   │   ├── middleware/           # 中间件
│   │   │   ├── api_key.go        # API Key 认证
│   │   │   └── sso_jwt.go        # SSO JWT 认证
│   │   ├── model/                # 数据模型
│   │   │   ├── database/         # 数据库模型
│   │   │   ├── request/          # 请求模型
│   │   │   └── response/         # 响应模型
│   │   ├── router/               # 路由配置
│   │   │   ├── apikey.go         # API Key 路由
│   │   │   ├── auth.go           # 认证路由
│   │   │   ├── document.go       # 文档路由
│   │   │   ├── library.go        # 库路由
│   │   │   ├── mcp.go            # MCP 路由
│   │   │   └── search.go         # 搜索路由
│   │   └── service/              # 业务逻辑
│   │       ├── apikey.go         # API Key 服务
│   │       ├── document.go       # 文档服务
│   │       ├── github_import.go  # GitHub 导入服务
│   │       ├── library.go        # 库服务
│   │       ├── mcp.go            # MCP 服务
│   │       ├── processor.go      # 文档处理器
│   │       └── search.go         # 搜索服务
│   ├── pkg/                      # 公共包
│   │   ├── cache/                # 缓存接口（Redis）
│   │   ├── chunker/              # 文档分块（TokenBased）
│   │   ├── config/               # 配置管理
│   │   ├── core/                 # 核心组件（Zap、Server）
│   │   ├── embedding/            # Embedding 服务（OpenAI）
│   │   ├── github/               # GitHub API 客户端
│   │   ├── global/               # 全局变量
│   │   ├── parser/               # 文档解析（Markdown）
│   │   ├── storage/              # 存储服务（七牛云/本地）
│   │   ├── utils/                # 工具函数
│   │   └── vectorstore/          # 向量存储（pgvector）
│   ├── scripts/                  # 脚本工具
│   ├── uploads/                  # 上传文件目录
│   ├── Dockerfile
│   └── main.go
│
├── web-mcp/                      # 前端管理界面
│   ├── src/
│   │   ├── api/                  # API 接口
│   │   │   ├── apikey.ts         # API Key 接口
│   │   │   ├── document.ts       # 文档接口
│   │   │   ├── library.ts        # 库接口（含 GitHub 导入）
│   │   │   └── search.ts         # 搜索接口
│   │   ├── components/           # Vue 组件
│   │   │   ├── AddVersionModal.vue  # 版本添加弹窗（支持 Local/GitHub）
│   │   │   ├── AppHeader.vue     # 顶部导航
│   │   │   ├── AppFooter.vue     # 底部栏
│   │   │   └── PersonalDropdown.vue  # 用户下拉菜单
│   │   ├── router/               # 路由配置
│   │   ├── stores/               # Pinia 状态管理
│   │   │   └── user.ts           # 用户状态
│   │   ├── utils/                # 工具函数
│   │   │   ├── request.ts        # Axios 封装
│   │   │   ├── token.ts          # Token 管理
│   │   │   ├── sse.ts            # SSE 流式处理
│   │   │   └── deviceId.ts       # 设备 ID
│   │   └── views/                # 页面视图
│   │       ├── dashboard/        # Dashboard 页面
│   │       ├── home/             # 首页
│   │       ├── library/          # 库管理页面
│   │       │   ├── index.vue     # 库列表
│   │       │   ├── detail.vue    # 库详情（搜索测试）
│   │       │   └── admin.vue     # 库管理（文档上传）
│   │       ├── search/           # 搜索页面
│   │       ├── layout/           # 布局组件
│   │       └── SSOCallback.vue   # SSO 回调
│   └── package.json
│
├── docker-compose.yml            # Docker 编排
├── docker-compose.prod.yml       # 生产环境编排
└── README.md
```

## 📋 开发计划

### MVP（第 1-2 周）✅
- [x] 项目骨架
- [x] 配置管理
- [x] 数据库模型
- [x] API 路由
- [x] MCP 端点
- [x] SSO JWT 认证
- [x] API Key 管理
- [x] 文档解析（Markdown）
- [x] Embedding 生成（OpenAI）
- [x] 向量搜索（pgvector）
- [x] 文档上传与处理
- [x] 前端库管理界面

### 第二阶段（第 3-4 周）✅
- [ ] PDF/DOCX 解析
- [x] 混合搜索（向量 + BM25）
- [x] 重排序算法（3 指标）
- [x] Redis 缓存优化（Embedding 缓存 + 搜索结果缓存 + GetOrSet 模式）
- [x] 前端搜索结果展示
- [x] GitHub 仓库导入功能
- [ ] MCP IDE 集成测试

## 📝 开发日志

### 2025-12-23

#### Added
- **活动日志系统 (Activity Log)**
  - 新增 `pkg/actlog` 包：异步批量活动日志记录
    - `Buffer`：缓冲区实现，支持批量写入（默认 50 条/批，2 秒刷新）
    - `TaskLogger`：任务级别日志器，预填充 libraryID、taskID、version 等公共字段
    - 支持 `WithActor`、`WithTarget`、`WithTaskID`、`WithVersion` 等选项
  - 新增 `ActivityLog` 数据库模型：记录库操作事件
  - 新增 `GET /api/v1/logs` API：获取库的最新任务日志
  - 前端 `detail.vue` 新增 Logs Tab：终端风格日志面板，支持自动轮询

- **GitHub 快速导入功能**
  - 新增 `POST /api/v1/libraries/github/init-import` API
    - 输入 GitHub URL → 自动解析仓库 → 验证连通性 → 检查重复 → 创建库 → 异步导入
    - 返回 `library_id` 和 `version`，前端跳转到 logs tab 查看进度
  - 新增 `AddDocsModal.vue` 组件：支持 GitHub 和 Local 两种导入方式
  - 新增 `pkg/utils/github.go`：`ParseGitHubURL`、`ExtractRepoName` 工具函数
  - 新增 `pkg/utils/task_id.go`：`GenerateTaskID` 生成 ULID 格式任务 ID

- **版本添加弹窗重构**
  - 新增 `AddVersionModal.vue` 组件：统一 Local 和 GitHub 两种模式
    - Local 模式：输入版本名创建空版本
    - GitHub 模式：选择 tag 自动导入
  - 版本创建成功后跳转到 logs tab 查看进度

#### Changed
- **GitHub 导入路由统一**
  - `POST /libraries/:id/import-github` → `POST /libraries/github/import?id=xxx`
  - `POST /libraries/:id/import-github-sse` → `POST /libraries/github/import-sse?id=xxx`
  - 新增 `POST /libraries/github/init-import`（快速导入）

- **活动日志集成**
  - `ImportFromGitHub`：记录 `github.import.start`、`github.import.download`、`github.import.complete` 等事件
  - `RefreshVersion`：记录 `version.refresh` 事件
  - `InitImportFromGitHub`：记录 `library.create` 和 `github.import.start` 事件
  - 所有日志包含 `actor_id`、`task_id`、`version`、`target_type`、`target_id` 等字段

- **API 层同步写入开始日志**
  - 在 goroutine 启动前同步写入"开始"日志，确保 API 返回前日志已入库
  - 解决前端跳转后日志显示 `status: complete` 的问题

- **前端 Tab 切换优化**
  - `onMounted` 根据当前 tab 加载对应数据，避免不必要的请求
  - 版本变化时只加载当前 tab 的数据
  - 切换到 context tab 时，如果没有搜索结果则自动加载

- **LibraryCreate 支持 DefaultVersion**
  - `LibraryCreate` 请求新增 `default_version` 字段
  - GitHub 导入时默认版本设为 `latest`

#### Fixed
- **版本重复检查**
  - `ImportFromGitHub` API 在启动 goroutine 前检查版本是否已存在
  - 避免重复导入同一版本

- **TaskID 统一**
  - API 层生成 taskID 并传递给服务方法，避免同一任务出现多个 taskID

---

### 2025-12-21

#### Added
- **GitHub 仓库导入功能**
  - 新增 `GitHubImportService`：从 GitHub 仓库直接导入 Markdown 文档
  - 支持指定分支（branch）或标签（tag）导入
  - 支持路径过滤（`path_filter`）和排除模式（`excludes`）
  - 动态下载策略：小仓库使用多 API 并行下载，大仓库（>100MB）使用 tarball 流式下载
  - SSE 实时进度推送：fetching_tree → downloading → processing → completed
  - 自动创建版本：仅在有成功导入文件时才创建版本，避免孤立版本

- **GitHub 版本列表 API**
  - 新增 `GET /api/v1/libraries/github/releases?repo=owner/repo`
  - 返回仓库信息（default_branch、description）和每个大版本的最新 tag

- **GitHub 客户端**
  - 新增 `pkg/github/client.go`：封装 GitHub API 调用
  - 支持 Token 认证和代理配置
  - 实现 `GetRepoInfo`、`GetTree`、`FilterTree`、`GetMajorVersions` 等方法
  - 支持 tarball 流式下载（`DownloadTarballFiles`）

- **LLM 富化并发优化**
  - `enrichChunks` 改用 5 个 worker 并发处理，性能提升约 5 倍
  - Worker Pool 模式：所有任务通过 channel 分发，固定 worker 数量

#### Changed
- **配置新增 GitHub 字段**
  - `config.yaml` 新增 `github.token` 和 `github.proxy` 配置项
  - 支持企业内网代理访问 GitHub API

- **七牛云存储上传优化**
  - 使用 `putExtra.MimeType` 设置 MIME 类型，替代 `putPolicy.MimeLimit`

---

### 2025-12-19

#### Added
- **无感知更新（Transactional Document Refresh）**
  - `DocumentChunk` 新增 `BatchVersion` 字段，支持版本化原子切换
  - 新增 `ProcessDocumentForRefresh()` 方法，返回 chunks 而非直接写库
  - 重写 `RefreshVersionWithCallback()`：先生成 pending chunks → 原子切换 → 软删除旧数据
  - 刷新过程中检索不受影响，用户无感知

- **版本刷新 SSE 实时进度推送**
  - 新增 `RefreshVersionSSE` API 端点 (`POST /libraries/:id/versions/:version/refresh-sse`)
  - 新增 `library_refresh_sse.go` 定义 `RefreshStatus` 结构和 SSE 写入器
  - 前端 `admin.vue` 新增刷新进度弹窗：进度条 + 文档状态列表

- **七牛云存储 Download 方法实现**
  - `qiniu.go` 实现 `Download()` 方法，通过 HTTP 获取文件内容
  - 修复文档刷新时从本地读取改为云存储下载

#### Changed
- **Processor 重构**
  - 提取 `processDocumentCore()` 公共方法，`ProcessDocument` 和 `ProcessDocumentForRefresh` 复用
  - 避免代码重复

- **GetVersions 统计修复**
  - `TokenCount` 和 `ChunkCount` 从硬编码 0 改为数据库聚合计算 (`SUM`)

- **Document List 接口优化**
  - 不传 `version` 时自动使用 `library.DefaultVersion`
  - 修复 GORM 链问题：使用 `Session()` 克隆避免 `Count()` 影响 `Find()`

---

### 2025-12-18

#### Added
- **多 Topic 搜索 + RRF 合并**
  - 支持逗号/空格分隔的多 topic 查询：`routing, middleware, binding`
  - 每个 topic 独立搜索，使用 Reciprocal Rank Fusion (RRF) 算法合并结果
  - RRF 公式：`score(d) = Σ 1/(k + rank)`，k=60（Elasticsearch 默认值）
  - 并行搜索：多个 topic 并发执行，提升响应速度

- **搜索结果缓存**
  - 每个子 topic 的搜索结果独立缓存，支持跨查询复用
  - 缓存 Key 格式（递进关系）：`search:topic:{library_id}:{version}:{mode}:{topic_hash}`
  - TTL：24 小时
  - 性能提升：多 topic 热启动快 20 倍（0.82s → 0.04s）

- **通用缓存工具 `GetOrSet[T]`**
  - 实现 Cache-Aside Pattern（旁路缓存模式）
  - 泛型支持，自动处理缓存命中/未命中逻辑
  - 位置：`pkg/cache/cache.go`

- **Redis 升级到 v9**
  - 统一使用 `github.com/redis/go-redis/v9`
  - 支持 Context 参数
  - `NewRedisCacheWithClient()` 复用全局 Redis 客户端

#### Changed
- **全局变量新增 `global.Cache`**
  - 通用缓存接口，用于搜索结果缓存等场景
  - 初始化顺序：Redis → Cache → Embedding

---

### 2025-12-17

#### Added
- **文档处理流程重构（参考 Context7 和业界最佳实践）**
  - 新增 Pre-Chunking 预处理：移除徽章、HTML 标签、空白行等无效内容
  - 新增 LLM Enrich 阶段：使用 LLM 为每个块生成 Title 和 Description
  - 处理流程：Parse → Pre-Process → Chunk → Enrich → Embed → Store

- **Markdown 分块逻辑优化**
  - 修复空标题问题：从标题行之后开始提取内容，跳过只有标题没有内容的 section
  - 简化 ChunkType：只保留 `code` 和 `info` 两种类型（有代码块 → code，无 → info）
  - 标题层级传递：空标题的 headers 会传递给下一个有内容的 section

- **LLM Service 更新**
  - 简化 `EnrichInput`：Content、Headers、Language、Source
  - 简化 `EnrichOutput`：只返回 Title 和 Description
  - 优化提示词：中文输出，简洁明了

#### Changed
- **处理流程进度调整**
  - parsing: 5% → preprocessing: 10% → chunking: 20% → enriching: 35% → embedding: 60% → saving: 85% → completed: 100%

---

### 2025-12-16 (续)

#### Added
- **文档块获取 API 重构**
  - 合并两个 GetChunks 路由为统一端点：`GET /documents/chunks/:mode/:libid/*version`
  - 支持 `mode` 参数（code/info）按类型筛选文档块
  - 版本参数可选，未指定时默认使用库的 `DefaultVersion`
  - 后端 `GetChunks()` 方法支持 mode 和 version 过滤

- **前端文档块格式化**
  - 新增 `formatCodeChunk()` 和 `formatInfoChunk()` 辅助函数
  - Code 模式：标题 → 来源 → 描述 → 代码块（带语言标记）
  - Info 模式：标题 → 来源 → 描述 → 正文内容
  - 块之间使用分隔符 `\n\n--------------------------------\n\n` 分隔

- **前端 Code/Info 标签页切换**
  - 导入 `getLatestInfo()` 函数
  - `fetchDocument()` 根据 `searchMode` 调用对应 API
  - 添加 `watch(searchMode)` 监听，切换时自动加载内容

#### Changed
- **路由路径调整**
  - `/documents` → `/documents/list`（文档列表）
  - `/documents/:id` → `/documents/detail/:id`（文档详情）
  - 新增 `/documents/chunks/:mode/:libid/*version`（文档块）

- **API 响应格式**
  - `getChunks()` 返回 `ChunksResponse` 包含 chunks 数组
  - `getLatestCode()`、`getLatestInfo()` 返回合并后的 `DocumentContent`

#### Fixed
- **路由冲突修复**
  - 解决 `:id` 和 `:mode` 参数冲突，将 chunks 路由独立为 `/documents/chunks/...`
  
- **SSO 设备查询修复**
  - `auth_service.go` 设备查询添加 `app_id` 条件
  - 修复不同应用设备记录互相干扰的问题

- **Code 模式显示修复**
  - `formatCodeChunk()` 优先使用 `code` 字段，否则使用 `chunk_text`
  - 确保代码块正常显示

---

### 2025-12-16

#### Added
- **版本导航功能**
  - `admin.vue` 版本名称可点击，跳转到对应版本页面
  - 默认版本跳转到 `/libraries/{id}`，其他版本跳转到 `/libraries/{id}/{version}`
  - 路由配置：`/libraries/:id/:version/:title` 用于显示特定版本的文档

- **文档列表版本过滤**
  - `detail.vue` 添加 Documents tab，显示当前版本的文档列表
  - 支持分页显示文档（标题、Tokens、Snippets、更新时间）
  - `getDocuments()` API 支持可选的 `version` 参数

- **用户信息全局缓存**
  - `stores/user.ts` 实现 Promise 缓存机制
  - `initUserState()` 只执行一次，避免重复请求 `/api/v1/user/info`
  - `logout()` 时重置缓存标记

- **日志系统优化**
  - 使用 `gin.Default()` 替代 `gin.New()` + 自定义中间件
  - 控制台显示 Gin 原生彩色日志（DebugMode）
  - 文件日志保持普通格式（无颜色代码）

#### Changed
- **后端 API 参数精简**
  - `DocumentList` 请求模型：移除 `title`、`file_type`、`status` 参数
  - 只保留 `library_id`（必需）、`version`（可选）、分页参数
  - `DocumentService.List()` 简化过滤逻辑

- **页面结构重组**
  - `admin.vue` 移除 Documents tab（只保留 Configuration 和 Versions）
  - `detail.vue` 添加 Documents tab（显示版本相关文档）
  - 明确页面职责：admin 管理全局，detail 展示版本内容

- **路由调整**
  - `/libraries/:id` 不区分版本（显示默认版本）
  - `/libraries/:id/:version/:title` 显示特定版本的文档内容
  - 移除了错误的 `/libraries/:id/:version` 路由

#### Fixed
- 控制台日志无彩色输出问题
- 日志文件包含 ANSI 颜色代码问题
- 版本参数在文档内容 API 中未生效问题

---

### 2025-12-15

#### Added
- **版本管理系统**
  - 库创建时初始化 `default_version = "default"`，`versions = []`（versions 不包含 default）
  - 实现 `GetVersions` API，直接从 Library 表读取，返回 default_version 在前，versions 倒序在后
  - 前端 `library.ts` 添加 `getVersions()` 接口

- **后端初始化**
  - `main.go` 添加 `InitStorage()` 和 `InitLLM()` 初始化
  - 添加时区同步：`time.Local = time.FixedZone("CST", 8*3600)`

- **前端 API 统一**
  - 重构 `library.ts` 接口定义，分离 `LibraryListItem` 和 `Library` 类型
  - 移除 `admin.vue` 原生 fetch 调用，统一使用 API 接口
  - `admin.vue` 延迟加载版本列表（只在切换标签页时加载）

- **时间显示优化**
  - 更新 `home/index.vue` 和 `detail.vue` 的 `formatDate()` 函数
  - 支持分钟级精度：`just now`、`5 minutes`、`1 hour`、`2 days`、`1 week` 等
  - 处理未来时间和无效时间戳（显示 "now"）

#### Changed
- Library 数据模型：`Versions` 字段分离，不包含 "default"
- 库列表响应：使用 `LibraryListItem` 精简字段
- 路由参数格式：`:id` → `/:id`（Gin 标准格式）

#### Fixed
- 时间戳显示 "-1 days" 问题
- 版本列表 API 404 问题

---

### 2025-12-07

#### Added
- **生产环境部署**
  - 完成 `mcp.hsk423.cn` 域名部署
  - Nginx 反向代理配置（HTTP→HTTPS 重定向、SSL、SSE 支持）
  - Docker 镜像构建与远程部署脚本 (`deploy.sh`)
  - 服务集成到 `blog-network` 网络，与 `nginx-proxy` 互通

#### Changed
- **部署脚本优化**
  - 单服务部署添加 `stop` 步骤，确保配置更新生效
  - 部署后自动清理 `.tar` 镜像文件
  - 本地构建后清理悬空镜像 (`docker image prune`)
  - 创建 `deploy.example.sh` 模板，隐藏敏感信息

- **配置文件安全**
  - `deploy.sh` 加入 `.gitignore`（包含服务器 IP）

#### Fixed
- 前端 `VITE_BASE_API` 环境变量未注入 Docker 构建
- Nginx 配置中 `mcp.hsk423.cn` 重定向问题（浏览器缓存）

---

### 2025-12-06

#### Added
- **搜索结果展示优化**
  - 修改 `SearchResultItem` 结构，添加 `Title`、`Source`、`Tokens`、`Relevance` 字段
  - 实现 `extractDeepestTitle()` 从 Metadata 提取最深层级标题
  - 向量搜索和 BM25 搜索 JOIN documents 表获取文档标题
  - 前端搜索结果改为卡片列表展示（标题、来源、tokens、相关性分数）

#### Changed
- **搜索模式优化**
  - `code` 模式搜索 `code + mixed` 类型
  - `info` 模式搜索 `info + mixed` 类型

#### Fixed
- 库详情页 Tokens 显示错误（`chunk_count` → `token_count`）

### 2025-12-05

#### Added
- **文档上传与处理**
  - 实现 Markdown 文档上传 (`document.go`)
  - SSE 流式进度反馈 (`sse.ts`)
  - 文档处理器：分块、Embedding 生成、向量存储 (`processor.go`)
  - 文档管理页面 (`admin.vue`)

- **库管理完善**
  - 库列表页面 (`library/index.vue`)
  - 库详情页面 (`library/detail.vue`)
  - 库统计信息（token_count、document_count、chunk_count）
  - 获取最新文档内容 API (`getLatestCode`)

- **Token 刷新机制**
  - Axios 拦截器自动刷新过期 Token (`request.ts`)
  - Token 管理工具 (`token.ts`)

#### Changed
- 优化 OpenAI Embedding 代理配置 (`openai_proxy.go`)
- 完善 Zap 日志配置 (`zap.go`)

### 2025-12-04

#### Added
- **SSO 认证集成**
  - 实现 SSO JWT 中间件 (`sso_jwt.go`)
  - 创建认证 API (`auth.go`)、用户 API (`user.go`)
  - 前端 SSO 回调处理 (`SSOCallback.vue`)
  - 用户状态管理 (`stores/user.ts`)

- **API Key 管理系统**
  - 新增 API Key 数据模型 (`api_keys.go`)
  - 实现 API Key CRUD 接口 (`apikey.go`)
  - 创建 API Key 认证中间件，支持 `MCP_API_KEY` Header
  - 前端 API Key 管理 (`apikey.ts`)

- **前端 Dashboard 页面**
  - 实现 Dashboard 主页 (`dashboard/index.vue`)
  - MCP 配置展示卡片
  - API Keys 管理表格（参考 Context7 设计）
  - 组件：`AppHeader.vue`、`AppFooter.vue`、`PersonalDropdown.vue`

- **路由完善**
  - 拆分路由模块：`library.go`、`document.go`、`search.go`、`auth.go`、`user.go`、`apikey.go`
  - 统一路由初始化 (`router/enter.go`)

- **核心包实现**
  - Embedding 服务：OpenAI 集成 (`openai.go`、`openai_proxy.go`)
  - 文档分块：TokenBased 分块器 (`token_based.go`)
  - 缓存：Redis 缓存 (`redis.go`)
  - 文档解析：Markdown 解析器 (`markdown.go`)
  - 向量存储：pgvector 存储 (`store.go`)

- **业务逻辑层**
  - 搜索服务：向量搜索 + BM25 + 重排序 (`search.go`)
  - 库管理服务 (`library.go`)
  - MCP 服务：search-libraries、get-library-docs (`mcp.go`)
  - 文档处理器 (`processor.go`)
  - 文档服务 (`document.go`)

- **前端项目初始化**
  - Vue3 + TypeScript + Vite + TailwindCSS
  - 页面：`layout`、`library`、`search`、`dashboard`
  - API 接口：`search.ts`、`document.ts`、`library.ts`、`apikey.ts`

### 2025-12-02

#### Added
- **后端基础架构**
  - 核心包骨架：`parser`、`embedding`、`vectorstore`、`cache`、`chunker`
  - 配置管理：`conf_system`、`conf_postgres`、`conf_redis`、`conf_jwt`、`conf_zap`、`conf_embedding`、`conf_sso`
  - 数据模型：`library`、`document`、`document_chunk`、`api_keys`、`search_cache`、`statistics`
  - API 骨架：`library.go`、`document.go`、`search.go`、`mcp.go`
  - 路由：`base.go`、`mcp.go`
  - Docker 配置：`docker-compose.yml`、`docker-compose.prod.yml`、`Dockerfile`

---

**日志说明**：
- 日志按时间倒序排列（最新在上）
- 日期基于文件修改时间统计
- 使用语义化版本分类：Added（新增）、Changed（变更）、Fixed（修复）、Removed（移除）

## 📄 License

MIT
