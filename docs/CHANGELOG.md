# 📝 开发日志

日志按时间倒序排列（最新在上），使用语义化版本分类：Added（新增）、Changed（变更）、Fixed（修复）、Removed（移除）

---

## 2025-12-26

### Added

- **混合搜索RRF算法**  
  - 新增 `normalizeScoresMinMax()` 方法实现Min-Max得分归一化，将向量搜索和BM25搜索得分统一到[0,1]范围
  - 新增 `hybridRRF()` 方法实现基于排名的RRF（Reciprocal Rank Fusion）算法，解决不同搜索方式得分冲突问题
- **核心架构文档**  
  - 新增 `docs/SEARCH.md` 文档，详细记录混合搜索架构、RRF算法实现、性能指标和调优建议
  - 新增 `docs/CACHE.md` 文档，完整描述多层缓存架构、Redis集成、缓存策略和性能优化

### Changed

- **搜索算法优化**  
  - 重构 `mergeAndRerank()` 方法，从简单线性加权改为RRF算法融合，提高搜索结果质量
  - 调整权重配置：向量搜索权重提升至0.7，BM25权重调整为0.3，更好平衡语义搜索和关键词匹配
  - 统一单查询和多查询流程，都使用RRF算法进行结果融合，保证搜索一致性

### Performance

- **搜索质量提升**  
  - RRF算法基于排名而非绝对得分进行融合，避免某种搜索方式占主导地位
  - Min-Max归一化确保不同搜索方式得分可比较，提升融合效果
  - 业界验证的算法（Elasticsearch 8.8+、OpenSearch 2.11+默认使用），稳定性和效果有保障

---

## 2025-12-25

### Added

- **核心文档合集**  
  - 新增 `docs/API.md`、`docs/ARCHITECTURE.md`、`docs/DEPLOYMENT.md`、`docs/MCP.md`，覆盖 API 规范、架构设计、部署指南与 MCP 集成说明。
- **数据库日志与索引**  
  - `configs/config.example.yaml`、`pkg/config/conf_postgres.go` 与 `initialize/gorm.go` 支持将 GORM 日志写入可轮转文件，并为 `document_chunks` 预计算 `chunk_tsvector_simple` 索引。

### Changed

- **库管理**  
  - Local 类型的库仅需名称和描述，`LibraryCreate/Update`、`LibraryService`、`library.ts`、`AddDocsModal.vue`、`home/index.vue` 等前端表单均去掉 `source_type/source_url` 字段。  
  - GitHub 快速导入会调用 LLM 自动生成更友好的库名，并将 repo 默认分支等元信息写入活动日志。
- **MCP 工具 & 搜索**  
  - `get-library-docs` 允许 `libraryId` 可选，支持跨库搜索与多 topic（逗号分隔），并在响应中返回 `mode`、`version`、`language`、`code/content` 等字段。  
  - `vectorSearch` / `bm25Search` 不再每次 `to_tsvector(chunk_text)` 也不再放任 50+ 条结果，而是依赖预计算 `chunk_tsvector_simple` 与模型级过滤（status、deleted_at、library_id、version、chunk_type），命中新建索引后只返回 10 条排序结果，响应更快更稳定。
  - `search.go`、`search_service` 统一要求 `version`、`library_id`，将分页 `limit` 固定为 10，并只在 info 模式回传 `content`。  
  - `document.go` 获取文档块时同样限制 10 条并沿用搜索逻辑。
- **GitHub 导入与处理流水**  
  - Tarball 下载切换到 `codeload.github.com`，增加命中/扫描统计日志；文档处理器新增 `isFinalTask` 标记，区分独立上传与批量导入的日志级别。

### Performance

- **搜索时间**  
  - `vectorSearch`：优化前单次查询 ~900ms，启用模型查询 + 联合索引 + 预计算索引后常规请求降至 ~10ms。  
  - `bm25Search`：优化前 ~200ms，改用 `chunk_tsvector_simple` GIN 索引 + 联合索引后稳定在 ~30ms。

### Fixed

- 搜索接口在 `limit > 50` 时可能返回超量结果、或缺省版本导致不同库混查的问题。  
- `document_chunks` 全文索引改为预先维护的 `chunk_tsvector_simple`，避免频繁运行 `to_tsvector` 带来的性能开销。  
- GitHub tarball 下载在遇到分支/标签 URL、或个别文件读取失败时会立即报错的问题（现改为日志+跳过）。

### Database

- 新增 `idx_chunks_text_simple_active` / `idx_chunks_text_simple_active_code` / `idx_chunks_text_simple_active_info` / `idx_chunks_library_version_type` 等索引，配合 `chunk_tsvector_simple` 字段，提升 `bm25Search` 与版本过滤查询性能。

---

## 2025-12-24

### Added

- **MCP 调用日志系统**
  - 新增 `mcp_call_logs` 表：记录 MCP 接口调用详情
  - 新增 `pkg/bufferedwriter/mcplog` 包：异步批量写入 MCP 调用日志
  - 记录字段：`actor_id`、`func_name`、`library_id`、`params`（JSON）、`result_count`、`latency_ms`、`status`、`error_msg`
  - 在 `handleToolCall` 统一记录日志，新增工具只需实现 `doXxx(args) MCPToolResult`

- **API Key 使用次数统计**
  - `api_keys` 表新增 `usage_count` 字段
  - 每次 MCP 调用自动 `usage_count + 1`（数据库原子操作）

- **actlog 同步写入方法**
  - 新增 `StatusStart = "start"` 状态常量：标识任务开始
  - 新增 `LogSync()` / `InfoSync()` / `InfoStartSync()`：同步写入日志（绕过缓冲区）
  - 新增 `TaskLogger.InfoStartSync()`：任务开始日志同步写入
  - 确保 API 返回前开始日志已入库，前端跳转后能立即查到

### Changed

- **MCP 接口重构**
  - `search-libraries` 返回：`libraryId`（uint）+ `versions`（数组）+ `defaultVersion`
  - `get-library-docs` 参数：`libraryId`（uint）+ `version`（可选）+ `topic` + `mode` + `page`
  - 移除旧的 `id: "name/version"` 格式

- **缓冲写入器统一管理**
  - 新增 `initialize/buffered_writers.go`
  - `InitBufferedWriters()`：统一初始化 actlog、stats、mcplog
  - `CloseBufferedWriters()`：统一关闭并刷新缓冲区
  - `main.go` 简化为两行调用

- **文档上传接口重构**
  - 前端改用普通 API 上传（`uploadDocument`），不再使用 SSE
  - 上传完成后跳转到 logs tab 查看处理进度
  - SSE 上传代码保留为注释备用
  - 上传时显示简单 loading 状态，不再显示进度条

- **actlog 日志规范化**
  - 开始日志：`status = "start"`（通过 `InfoStartSync`）
  - 过程日志：`status = "info"`（通过 `Info`）
  - 结束日志：`status = "success"`（通过 `Success`）
  - 错误日志：`status = "error"`（通过 `Error`）
  - 日志查询接口：只有 `success` 返回 `complete`，其他都是 `processing`

- **actlog 使用统一化**
  - `InitImportFromGitHub`：使用根 `TaskLogger`，各阶段通过 `WithTarget` 派生
  - `ImportFromGitHub`：开始日志在 API 层同步写入
  - `RefreshVersion`：改用 `actLogger.InfoStartSync()` 替代直接写 DB
  - `Upload`：开始日志在 service 层同步写入，完成日志在调用方记录

- **前端日志渲染优化**
  - 优先根据 `status` 渲染颜色（start=紫色, success=绿色, error=红色, warning=黄色）
  - `status === 'info'` 时再根据 `event` 渲染

- **前端 UI 调整**
  - 注释掉 Header 的 Plans, Learn, Try Live, Install 导航项
  - "Chat with Docs" 按钮置灰

---

## 2025-12-23

### Added

- **MCP 调用统计系统**
  - 新增 `pkg/bufferedwriter` 公共包：提取通用异步批量写入逻辑
    - `Buffer[T]`：泛型缓冲区，支持批量写入 + 定时刷新
    - `Writer[T]`：写入器接口（`WriteBatch` + `Close`）
  - 新增 `pkg/bufferedwriter/stats` 包：MCP 调用统计
    - `stats.Increment(metricName, delta)`：全局统计
    - `stats.IncrementWithLibrary(libraryID, metricName, delta)`：库级统计
    - 支持 PostgreSQL upsert（`ON CONFLICT` 累加）
  - 重构 `pkg/bufferedwriter/actlog`：复用公共 Buffer
  - 新增 `Statistics` 模型常量：`MetricMCPGetLibraryDocs`、`MetricMCPSearchLibraries`
  - MCP 接口 `search-libraries` 和 `get-library-docs` 自动记录调用次数

- **个人统计 API**
  - 新增 `GET /api/v1/stats/my`：获取当前用户统计（库数、文档数、Token 数、MCP 调用数）
  - 新增 `StatsService.GetUserStats()`：带 5 分钟缓存
  - 前端 Dashboard 展示个人统计卡片

- **库列表热门排序**
  - `GET /api/v1/libraries` 支持 `sort=popular` 参数
  - 按 MCP 调用次数（`mcp.func.get_library_docs`）降序排序
  - 前端首页 Popular Tab 使用热门排序

- **MCP 搜索版本过滤**
  - `get-library-docs` 现在正确按版本过滤文档
  - `SearchDocuments` 和 `GetChunksByLibrary` 支持 `version` 参数

### Changed

- **bufferedwriter 重构**
  - `pkg/actlog` → `pkg/bufferedwriter/actlog`
  - `pkg/stats` → `pkg/bufferedwriter/stats`
  - 两者共用 `bufferedwriter.Buffer[T]` 和 `bufferedwriter.Writer[T]` 接口

- **Library 模型新增 CreatedBy**
  - `Library` 新增 `created_by` 字段，记录创建者 UUID
  - `InitFromGitHub` 传递 `createdBy` 参数

---

## 2025-12-22

### Added

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

### Changed

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

### Fixed

- **版本重复检查**
  - `ImportFromGitHub` API 在启动 goroutine 前检查版本是否已存在
  - 避免重复导入同一版本

- **TaskID 统一**
  - API 层生成 taskID 并传递给服务方法，避免同一任务出现多个 taskID

---

## 2025-12-21

### Added

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

### Changed

- **配置新增 GitHub 字段**
  - `config.yaml` 新增 `github.token` 和 `github.proxy` 配置项
  - 支持企业内网代理访问 GitHub API

- **七牛云存储上传优化**
  - 使用 `putExtra.MimeType` 设置 MIME 类型，替代 `putPolicy.MimeLimit`

---

## 2025-12-19

### Added

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

### Changed

- **Processor 重构**
  - 提取 `processDocumentCore()` 公共方法，`ProcessDocument` 和 `ProcessDocumentForRefresh` 复用
  - 避免代码重复

- **GetVersions 统计修复**
  - `TokenCount` 和 `ChunkCount` 从硬编码 0 改为数据库聚合计算 (`SUM`)

- **Document List 接口优化**
  - 不传 `version` 时自动使用 `library.DefaultVersion`
  - 修复 GORM 链问题：使用 `Session()` 克隆避免 `Count()` 影响 `Find()`

---

## 2025-12-18

### Added

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

### Changed

- **全局变量新增 `global.Cache`**
  - 通用缓存接口，用于搜索结果缓存等场景
  - 初始化顺序：Redis → Cache → Embedding

---

## 2025-12-17

### Added

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

### Changed

- **处理流程进度调整**
  - parsing: 5% → preprocessing: 10% → chunking: 20% → enriching: 35% → embedding: 60% → saving: 85% → completed: 100%

---

## 2025-12-16

### Added

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

### Changed

- **路由路径调整**
  - `/documents` → `/documents/list`（文档列表）
  - `/documents/:id` → `/documents/detail/:id`（文档详情）
  - 新增 `/documents/chunks/:mode/:libid/*version`（文档块）

- **API 响应格式**
  - `getChunks()` 返回 `ChunksResponse` 包含 chunks 数组
  - `getLatestCode()`、`getLatestInfo()` 返回合并后的 `DocumentContent`

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

### Fixed

- **路由冲突修复**
  - 解决 `:id` 和 `:mode` 参数冲突，将 chunks 路由独立为 `/documents/chunks/...`
  
- **SSO 设备查询修复**
  - `auth_service.go` 设备查询添加 `app_id` 条件
  - 修复不同应用设备记录互相干扰的问题

- **Code 模式显示修复**
  - `formatCodeChunk()` 优先使用 `code` 字段，否则使用 `chunk_text`
  - 确保代码块正常显示

- 控制台日志无彩色输出问题
- 日志文件包含 ANSI 颜色代码问题
- 版本参数在文档内容 API 中未生效问题

---

## 2025-12-15

### Added

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

### Changed

- Library 数据模型：`Versions` 字段分离，不包含 "default"
- 库列表响应：使用 `LibraryListItem` 精简字段
- 路由参数格式：`:id` → `/:id`（Gin 标准格式）

### Fixed

- 时间戳显示 "-1 days" 问题
- 版本列表 API 404 问题

---

## 2025-12-07

### Added

- **生产环境部署**
  - 完成 `mcp.hsk423.cn` 域名部署
  - Nginx 反向代理配置（HTTP→HTTPS 重定向、SSL、SSE 支持）
  - Docker 镜像构建与远程部署脚本 (`deploy.sh`)
  - 服务集成到 `blog-network` 网络，与 `nginx-proxy` 互通

### Changed

- **部署脚本优化**
  - 单服务部署添加 `stop` 步骤，确保配置更新生效
  - 部署后自动清理 `.tar` 镜像文件
  - 本地构建后清理悬空镜像 (`docker image prune`)
  - 创建 `deploy.example.sh` 模板，隐藏敏感信息

- **配置文件安全**
  - `deploy.sh` 加入 `.gitignore`（包含服务器 IP）

### Fixed

- 前端 `VITE_BASE_API` 环境变量未注入 Docker 构建
- Nginx 配置中 `mcp.hsk423.cn` 重定向问题（浏览器缓存）

---

## 2025-12-06

### Added

- **搜索结果展示优化**
  - 修改 `SearchResultItem` 结构，添加 `Title`、`Source`、`Tokens`、`Relevance` 字段
  - 实现 `extractDeepestTitle()` 从 Metadata 提取最深层级标题
  - 向量搜索和 BM25 搜索 JOIN documents 表获取文档标题
  - 前端搜索结果改为卡片列表展示（标题、来源、tokens、相关性分数）

### Changed

- **搜索模式优化**
  - `code` 模式搜索 `code + mixed` 类型
  - `info` 模式搜索 `info + mixed` 类型

### Fixed

- 库详情页 Tokens 显示错误（`chunk_count` → `token_count`）

---

## 2025-12-05

### Added

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

### Changed

- 优化 OpenAI Embedding 代理配置 (`openai_proxy.go`)
- 完善 Zap 日志配置 (`zap.go`)

---

## 2025-12-04

### Added

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

---

## 2025-12-02

### Added

- **后端基础架构**
  - 核心包骨架：`parser`、`embedding`、`vectorstore`、`cache`、`chunker`
  - 配置管理：`conf_system`、`conf_postgres`、`conf_redis`、`conf_jwt`、`conf_zap`、`conf_embedding`、`conf_sso`
  - 数据模型：`library`、`document`、`document_chunk`、`api_keys`、`search_cache`、`statistics`
  - API 骨架：`library.go`、`document.go`、`search.go`、`mcp.go`
  - 路由：`base.go`、`mcp.go`
  - Docker 配置：`docker-compose.yml`、`docker-compose.prod.yml`、`Dockerfile`
