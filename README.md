# go-mcp-context

<div align="center">

![Go Version](https://img.shields.io/badge/Go-1.23-blue)
![Vue Version](https://img.shields.io/badge/Vue-3.5-green)
![License](https://img.shields.io/badge/license-MIT-orange)

私有化的 Context7 替代方案，为企业内网的 AI IDE 提供实时、准确的技术文档和代码示例

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
GET    /api/v1/libraries      # 获取库列表
POST   /api/v1/libraries      # 创建库
GET    /api/v1/libraries/:id  # 获取库详情
PUT    /api/v1/libraries/:id  # 更新库
DELETE /api/v1/libraries/:id  # 删除库
```

### 文档管理

```
POST   /api/v1/documents/upload  # 上传文档
GET    /api/v1/documents/:id     # 获取文档
DELETE /api/v1/documents/:id     # 删除文档
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
│   │   ├── global/               # 全局变量
│   │   ├── parser/               # 文档解析（Markdown）
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
│   │   │   ├── library.ts        # 库接口
│   │   │   └── search.ts         # 搜索接口
│   │   ├── components/           # Vue 组件
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

### 第二阶段（第 3-4 周）🚧
- [ ] PDF/DOCX 解析
- [x] 混合搜索（向量 + BM25）
- [x] 重排序算法（3 指标）
- [ ] Redis 缓存优化
- [x] 前端搜索结果展示
- [ ] MCP IDE 集成测试

## 📝 开发日志

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

### 2025-12-02 ~ 2025-12-03

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
