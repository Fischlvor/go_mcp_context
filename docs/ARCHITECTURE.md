# 项目架构

## 技术栈

### 后端

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

### 前端

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

---

## 项目结构

```text
go-mcp-context/
├── server-mcp/               # MCP 后端服务
│   ├── cmd/                  # 主程序入口
│   ├── configs/              # 配置文件
│   ├── internal/
│   │   ├── api/              # HTTP 处理器
│   │   ├── initialize/       # 初始化模块
│   │   ├── middleware/       # 中间件
│   │   ├── model/            # 数据模型
│   │   │   ├── database/     # 数据库模型
│   │   │   ├── request/      # 请求模型
│   │   │   └── response/     # 响应模型
│   │   ├── router/           # 路由配置
│   │   └── service/          # 业务逻辑
│   ├── pkg/                  # 公共包
│   │   ├── bufferedwriter/   # 异步批量写入框架
│   │   ├── cache/            # 缓存接口
│   │   ├── chunker/          # 文档分块
│   │   ├── config/           # 配置管理
│   │   ├── core/             # 核心组件
│   │   ├── embedding/        # Embedding 服务
│   │   ├── github/           # GitHub API 客户端
│   │   ├── global/           # 全局变量
│   │   ├── parser/           # 文档解析
│   │   ├── storage/          # 存储服务
│   │   ├── utils/            # 工具函数
│   │   └── vectorstore/      # 向量存储
│   ├── scripts/              # 脚本工具
│   ├── uploads/              # 上传文件目录
│   ├── Dockerfile
│   └── main.go
│
├── web-mcp/                  # 前端管理界面
│   ├── src/
│   │   ├── api/              # API 接口
│   │   ├── components/       # Vue 组件
│   │   ├── router/           # 路由配置
│   │   ├── stores/           # Pinia 状态管理
│   │   ├── utils/            # 工具函数
│   │   └── views/            # 页面视图
│   └── package.json
│
├── docs/                     # 文档
│   ├── API.md                # API 文档
│   ├── ARCHITECTURE.md       # 架构文档
│   ├── CHANGELOG.md          # 开发日志
│   └── DEPLOYMENT.md         # 部署指南
│
├── docker-compose.yml        # Docker 编排
├── docker-compose.prod.yml   # 生产环境编排
└── README.md
```

---

## 数据模型

### Library（文档库）

| 字段 | 类型 | 说明 |
|------|------|------|
| id | uint | 主键 |
| name | string | 库名称 |
| description | string | 描述 |
| default_version | string | 默认版本 |
| versions | []string | 版本列表 |
| created_by | string | 创建者 UUID |

### Document（文档）

| 字段 | 类型 | 说明 |
|------|------|------|
| id | uint | 主键 |
| library_id | uint | 所属库 |
| version | string | 版本 |
| title | string | 标题 |
| file_type | string | 文件类型 |
| token_count | int | Token 数 |
| chunk_count | int | 分块数 |

### DocumentChunk（文档块）

| 字段 | 类型 | 说明 |
|------|------|------|
| id | uint | 主键 |
| document_id | uint | 所属文档 |
| chunk_type | string | code / info |
| chunk_text | string | 文本内容 |
| code | string | 代码内容 |
| embedding | vector | 向量 |
| title | string | LLM 生成标题 |
| description | string | LLM 生成描述 |

### APIKey（API 密钥）

| 字段 | 类型 | 说明 |
|------|------|------|
| id | uint | 主键 |
| user_uuid | string | 用户 UUID |
| token_hash | string | Token 哈希 |
| suffix | string | Token 后缀（显示用） |
| name | string | 名称 |
| usage_count | int64 | 使用次数 |
| last_used_at | time | 最后使用时间 |

### MCPCallLog（MCP 调用日志）

| 字段 | 类型 | 说明 |
|------|------|------|
| id | uint | 主键 |
| actor_id | string | 调用者 UUID |
| func_name | string | 函数名 |
| library_id | uint | 库 ID |
| params | jsonb | 请求参数 |
| result_count | int | 结果数量 |
| latency_ms | int64 | 延迟（毫秒） |
| status | string | success / error |
| error_msg | string | 错误信息 |
