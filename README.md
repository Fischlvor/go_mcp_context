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
│   │   ├── api/              # HTTP 处理器
│   │   ├── initialize/       # 初始化模块
│   │   ├── middleware/       # 中间件（JWT、API Key）
│   │   ├── model/            # 数据模型
│   │   ├── router/           # 路由配置
│   │   └── service/          # 业务逻辑
│   ├── pkg/
│   │   ├── cache/            # 缓存接口
│   │   ├── chunker/          # 文档分块
│   │   ├── config/           # 配置管理
│   │   ├── embedding/        # Embedding 服务
│   │   ├── global/           # 全局变量
│   │   ├── parser/           # 文档解析
│   │   └── vectorstore/      # 向量存储
│   ├── Dockerfile
│   └── main.go
├── web-mcp/                  # 前端管理界面
│   ├── src/
│   │   ├── api/              # API 接口
│   │   ├── components/       # Vue 组件
│   │   ├── router/           # 路由配置
│   │   ├── stores/           # Pinia 状态管理
│   │   ├── utils/            # 工具函数
│   │   └── views/            # 页面视图
│   └── package.json
├── docker-compose.yml
└── README.md
```

## 📋 开发计划

### MVP（第 1-2 周）
- [x] 项目骨架
- [x] 配置管理
- [x] 数据库模型
- [x] API 路由
- [x] MCP 端点
- [x] SSO JWT 认证
- [x] API Key 管理
- [ ] 文档解析（Markdown）
- [ ] Embedding 生成
- [ ] 向量搜索

### 第二阶段（第 3-4 周）
- [ ] PDF/DOCX 解析
- [ ] 混合搜索
- [ ] 重排序算法
- [ ] Redis 缓存
- [ ] 前端界面完善

## 📝 开发日志

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

### 2025-12-03

#### Added
- **Embedding 服务实现**
  - OpenAI Embedding 集成 (`openai.go`、`openai_proxy.go`)
  - Embedding 配置 (`conf_embedding.go`)
  - Embedding 初始化 (`initialize/embedding.go`)

- **业务逻辑层**
  - 搜索服务 (`service/search.go`)
  - 库管理服务 (`service/library.go`)
  - MCP 服务 (`service/mcp.go`)
  - 文档处理器 (`service/processor.go`)
  - 文档服务 (`service/document.go`)

- **前端项目初始化**
  - Vue3 + TypeScript + Vite 项目搭建
  - 页面：`layout`、`library`、`search`、`document`
  - API 接口：`search.ts`、`document.ts`、`library.ts`

### 2025-12-02

#### Added
- **后端基础架构**
  - 核心包：`parser`、`embedding`、`vectorstore`、`cache`、`chunker`
  - 配置管理：`conf_system`、`conf_postgres`、`conf_redis`、`conf_jwt`、`conf_zap`
  - 数据模型：`library`、`document`、`document_chunk`、`search_cache`、`statistics`
  - API 骨架：`library.go`、`document.go`、`search.go`、`mcp.go`
  - 路由：`base.go`、`mcp.go`
  - Docker 配置：`docker-compose.prod.yml`

---

**日志说明**：
- 日志按时间倒序排列（最新在上）
- 日期基于文件创建时间统计
- 使用语义化版本分类：Added（新增）、Changed（变更）、Fixed（修复）、Removed（移除）

## 📄 License

MIT
