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
│   │   ├── service/          # 业务逻辑
│   │   └── transport/        # 传输层（协议抽象）
│   │       ├── interface.go  # ResponseWriter 接口
│   │       ├── types.go      # 数据结构定义
│   │       ├── detector.go   # 协议检测
│   │       ├── factory.go    # 响应写入器工厂
│   │       ├── http/         # HTTP 协议实现
│   │       ├── streamable/   # Streamable HTTP 实现
│   │       └── sse/          # SSE 协议实现
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
│   ├── test/                 # 测试（覆盖率 81.0%）
│   │   ├── unit/             # 单元测试
│   │   ├── integration/      # 集成测试
│   │   ├── README.md         # 测试文档
│   │   ├── COVERAGE_LIMITATIONS.md  # 覆盖率限制说明
│   │   └── Makefile          # 测试命令
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

---

## MCP 处理架构

### 分层设计

```
┌─────────────────────────────────────────┐
│         HTTP 请求入口                    │
│    (internal/api/mcp.go)                │
└──────────────┬──────────────────────────┘
               │
               ▼
┌─────────────────────────────────────────┐
│      传输层 (Transport Layer)            │
│  - 协议检测 (HTTP/SSE/Streamable)      │
│  - 响应写入器工厂                       │
│  (internal/transport/)                  │
└──────────────┬──────────────────────────┘
               │
               ▼
┌─────────────────────────────────────────┐
│      处理器层 (Handler Layer)            │
│  - MCP 请求分发                         │
│  - 方法路由                             │
│  (internal/service/mcp_handler.go)     │
└──────────────┬──────────────────────────┘
               │
               ▼
┌─────────────────────────────────────────┐
│      业务逻辑层 (Service Layer)          │
│  - SearchLibraries                      │
│  - GetLibraryDocs                       │
│  - GetAllLibraries                      │
│  (internal/service/mcp.go)             │
└─────────────────────────────────────────┘
```

### 核心特性

#### 1. 协议无关的处理架构
- **业务逻辑与传输协议完全解耦**
- 支持多种传输协议（HTTP、SSE、Streamable HTTP）
- 新增协议只需实现 `ResponseWriter` 接口

#### 2. 统一的 MCP 请求处理
- `MCPHandler.ProcessRequest()` 方法分发所有 MCP 请求
- 支持的方法：
  - `initialize` / `notifications/initialized` - 初始化
  - `tools/list` / `tools/call` - 工具管理
  - `resources/list` / `resources/templates/list` / `resources/read` - 资源管理

#### 3. MCP 规范响应格式
- 所有 `tools/call` 响应遵循规范：
  ```json
  {
    "content": [
      {
        "type": "text",
        "text": "<JSON string>"
      }
    ]
  }
  ```
- 实际数据在 `text` 字段中以 JSON 字符串形式存储
- 与 Costrict 客户端的 `McpToolCallResponse` 类型完全兼容

### 传输层接口

```go
// ResponseWriter 接口 - 所有传输协议必须实现
type ResponseWriter interface {
    WriteResponse(resp *MCPResponse) error
    WriteError(err *MCPError, id interface{}) error
    WriteEvent(event interface{}) error
}

// RequestContext - 请求上下文
type RequestContext struct {
    Transport TransportType
    Method    string
    Params    interface{}
    ID        interface{}
    GinCtx    *gin.Context
}

// TransportType - 传输协议类型
type TransportType string

const (
    TransportHTTP       TransportType = "http"
    TransportSSE        TransportType = "sse"
    TransportStreamable TransportType = "streamable"
)
```

### 工作流程示例

**请求：** `POST /mcp` with `tools/call` method

```
1. API 层 (HandleRequest)
   ├─ 解析 JSON-RPC 请求
   ├─ 检测传输协议类型
   └─ 创建响应写入器

2. 传输层 (Transport)
   ├─ 自动检测协议（HTTP/SSE/Streamable）
   └─ 创建对应的 ResponseWriter

3. 处理器层 (MCPHandler)
   ├─ 根据 method 分发请求
   ├─ 调用 handleToolsCall()
   └─ 获取业务逻辑结果

4. 业务逻辑层 (MCPService)
   ├─ SearchLibraries() / GetLibraryDocs()
   └─ 返回原始数据

5. 处理器层 (MCPHandler)
   ├─ 将原始数据转换为 MCP 规范格式
   └─ 调用 ResponseWriter.WriteResponse()

6. 传输层 (ResponseWriter)
   └─ 根据协议类型发送响应

7. 客户端接收
   └─ 解析 McpToolCallResponse 格式
```

### 版本管理设计

#### search-libraries 响应
- `defaultVersion`：默认版本（通常为 `latest`）
- `versions` 数组：所有可用版本列表
- **关键设计**：`defaultVersion` 总是包含在 `versions` 数组中

#### LLM 工作流指导
1. 调用 `search-libraries` 获取库的版本信息
2. 从响应的 `versions` 数组中选择一个版本
3. 调用 `get-library-docs` 时使用选中的版本号

这样设计的优势：
- LLM 可直接使用 `versions` 数组中的任意版本
- 避免硬编码版本号
- 支持库的多个版本管理

---

## 可扩展性设计

### 新增传输协议
只需在 `internal/transport/` 中添加新的 Writer 实现：

```go
// 例如：添加 WebSocket 支持
type WebSocketResponseWriter struct {
    conn *websocket.Conn
}

func (w *WebSocketResponseWriter) WriteResponse(resp *MCPResponse) error {
    return w.conn.WriteJSON(resp)
}
```

### 新增 MCP 方法
只需在 `MCPHandler` 中添加新的 handle 函数：

```go
func (h *MCPHandler) handleNewMethod(req *RequestContext, writer ResponseWriter) error {
    // 业务逻辑
    result := h.mcpService.NewMethod(...)
    
    // 转换为 MCP 规范格式
    mcpResult := convertToMCPFormat(result)
    
    // 发送响应
    return writer.WriteResponse(&MCPResponse{
        JSONRPC: "2.0",
        ID:      req.ID,
        Result:  mcpResult,
    })
}
```

### 修改响应格式
只需修改 ResponseWriter 的实现，无需改动业务逻辑
