# MCP 接口文档

## 概述

MCP (Model Context Protocol) 接口使用 JSON-RPC 2.0 协议，供 AI IDE（如 VS Code、Cursor、Windsurf）调用获取文档上下文。

- **Base URL**: `http://localhost:8090` 或 `https://mcp.hsk423.cn`
- **认证方式**: `MCP_API_KEY: <API_KEY>` (HTTP Header)
- **协议**: JSON-RPC 2.0

---

## 健康检查

```http
GET /mcp/health
```

**响应：**

```json
{
  "status": "ok",
  "version": "1.0.0"
}
```

---

## 获取工具列表

```http
GET /mcp/tools
```

**响应：**

```json
{
  "tools": [
    {
      "name": "search-libraries",
      "description": "Search for documentation libraries by name",
      "inputSchema": { ... }
    },
    {
      "name": "get-library-docs",
      "description": "Get documentation from a specific library",
      "inputSchema": { ... }
    }
  ]
}
```

---

## MCP 方法列表

MCP 服务器支持以下 7 个方法：

### 1. initialize - 初始化连接

**请求：**
```json
{
  "jsonrpc": "2.0",
  "id": 1,
  "method": "initialize",
  "params": {
    "protocolVersion": "2025-11-25",
    "capabilities": {},
    "clientInfo": {
      "name": "CoStrict",
      "version": "2.1.4"
    }
  }
}
```

**响应：**
```json
{
  "jsonrpc": "2.0",
  "id": 1,
  "result": {
    "protocolVersion": "2025-11-25",
    "capabilities": {
      "tools": { "listChanged": true },
      "resources": { "subscribe": true, "listChanged": true },
      "logging": {}
    },
    "serverInfo": {
      "name": "go-mcp-context",
      "version": "1.0.0"
    }
  }
}
```

---

### 2. notifications/initialized - 初始化完成通知

**请求：**
```json
{
  "jsonrpc": "2.0",
  "method": "notifications/initialized"
}
```

**说明：** 客户端在初始化完成后发送此通知，服务器不需要返回响应。

---

### 3. tools/list - 获取工具列表

**请求：**
```json
{
  "jsonrpc": "2.0",
  "id": 1,
  "method": "tools/list"
}
```

**响应：**
```json
{
  "jsonrpc": "2.0",
  "id": 1,
  "result": {
    "tools": [
      {
        "name": "search-libraries",
        "description": "Search for documentation libraries by name. Returns matching libraries with metadata including available versions. Use this method to discover libraries and get their version information (versions array and defaultVersion) before calling get-library-docs.",
        "inputSchema": { ... }
      },
      {
        "name": "get-library-docs",
        "description": "Get documentation for a specific library. Requires libraryId, topic, and version. Supports comma-separated topics for multi-topic search.",
        "inputSchema": { ... }
      }
    ]
  }
}
```

---

### 4. tools/call - 调用工具

**请求：**
```json
{
  "jsonrpc": "2.0",
  "id": 1,
  "method": "tools/call",
  "params": {
    "name": "search-libraries",
    "arguments": {
      "libraryName": "gin"
    }
  }
}
```

**说明：** 通过此方法调用 `search-libraries` 和 `get-library-docs` 两个工具。详见下文的工具调用部分。

---

### 5. resources/list - 获取资源列表

**请求：**
```json
{
  "jsonrpc": "2.0",
  "id": 1,
  "method": "resources/list"
}
```

**响应：**
```json
{
  "jsonrpc": "2.0",
  "id": 1,
  "result": {
    "resources": [
      {
        "uri": "go-mcp-context:///library/1",
        "name": "Gin",
        "description": "Gin is a HTTP web framework written in Go",
        "mimeType": "application/json"
      },
      {
        "uri": "go-mcp-context:///library/2",
        "name": "Go Standard Library",
        "description": "The Go standard library documentation",
        "mimeType": "application/json"
      }
    ]
  }
}
```

**说明：** 返回所有可用的库资源列表，每个库对应一个资源。

---

### 6. resources/templates/list - 获取资源模板列表

**请求：**
```json
{
  "jsonrpc": "2.0",
  "id": 1,
  "method": "resources/templates/list"
}
```

**响应：**
```json
{
  "jsonrpc": "2.0",
  "id": 1,
  "result": {
    "resourceTemplates": [
      {
        "uriTemplate": "go-mcp-context:///library/{libraryId}",
        "name": "Library by ID",
        "description": "Get library information by ID",
        "mimeType": "application/json"
      },
      {
        "uriTemplate": "go-mcp-context:///docs/chunk/{libraryId}/{version}/{topic}",
        "name": "Documentation Chunk",
        "description": "Get documentation chunk for a specific library version and topic (supports comma-separated topics like overview,api,examples)",
        "mimeType": "text/markdown"
      }
    ]
  }
}
```

**说明：** 返回资源的 URI 模板，客户端可以根据模板生成具体的资源 URI。

**模板说明：**

| 模板 | 说明 | 参数 |
|------|------|------|
| `go-mcp-context:///library/{libraryId}` | 获取库的基本信息 | `libraryId` - 库ID |
| `go-mcp-context:///docs/chunk/{libraryId}/{version}/{topic}` | 获取库的文档块 | `libraryId` - 库ID<br/>`version` - 版本号<br/>`topic` - 主题（支持逗号分隔） |

---

### 7. resources/read - 读取资源内容

**请求：**
```json
{
  "jsonrpc": "2.0",
  "id": 1,
  "method": "resources/read",
  "params": {
    "uri": "go-mcp-context:///library/1"
  }
}
```

**响应：**
```json
{
  "jsonrpc": "2.0",
  "id": 1,
  "result": {
    "contents": [
      {
        "uri": "go-mcp-context:///library/1",
        "mimeType": "application/json",
        "text": "{\"id\": 1, \"name\": \"Gin\", \"description\": \"Gin is a HTTP web framework\", \"versions\": [\"latest\", \"v1.9.0\", \"v1.8.0\"], \"defaultVersion\": \"latest\", \"snippets\": 150}"
      }
    ]
  }
}
```

**说明：** 读取指定资源的内容，返回库的元数据信息。

---

## 工具调用详解

```http
POST /mcp
MCP_API_KEY: <API_KEY>
```

### search-libraries

搜索文档库（支持语义向量搜索 + 模糊匹配降级）。

**搜索策略：**
1. **优先使用向量搜索**：基于语义相似度（cosine distance）进行搜索
2. **降级到模糊匹配**：当向量搜索失败或无结果时，使用 SQL ILIKE 模糊匹配

**请求：**

```json
{
  "jsonrpc": "2.0",
  "id": 1,
  "method": "tools/call",
  "params": {
    "name": "search-libraries",
    "arguments": {
      "libraryName": "gin"
    }
  }
}
```

**响应：**

```json
{
  "jsonrpc": "2.0",
  "id": 1,
  "result": {
    "content": [
      {
        "type": "text",
        "text": "{\"libraries\": [{\"libraryId\": 1, \"name\": \"gin\", \"versions\": [\"latest\", \"v1.9.0\", \"v1.8.0\"], \"defaultVersion\": \"latest\", \"description\": \"Gin is a HTTP web framework\", \"snippets\": 150, \"score\": 0.95}]}"
      }
    ]
  }
}
```

**响应说明：**

响应遵循 MCP 规范，将实际数据包装在 `content` 数组中，`text` 字段包含 JSON 字符串。

**数据字段说明（text 字段中的 JSON）：**

| 字段 | 类型 | 说明 |
|------|------|------|
| libraryId | uint | 库 ID（用于 get-library-docs） |
| name | string | 库名称 |
| versions | string[] | 额外版本列表（不含 defaultVersion） |
| defaultVersion | string | 默认版本（通常为 `latest`） |
| description | string | 库描述 |
| snippets | int | 文档片段数量 |
| score | float | 匹配分数（0-1） |

**LLM 工作流指导：**

1. 调用 `search-libraries` 获取库的版本信息
2. 从响应的 `versions` 数组中选择一个版本
3. 调用 `get-library-docs` 时使用选中的版本号

---

### get-library-docs

获取库文档。

**请求：**

```json
{
  "jsonrpc": "2.0",
  "id": 2,
  "method": "tools/call",
  "params": {
    "name": "get-library-docs",
    "arguments": {
      "libraryId": 1,
      "version": "latest",
      "topic": "middleware",
      "mode": "code",
      "page": 1
    }
  }
}
```

**参数说明：**

| 参数 | 类型 | 必填 | 说明 |
|------|------|------|------|
| topic | string | 是 | 搜索主题，支持逗号分隔多个主题（如 `routing, middleware`） |
| libraryId | uint | 否 | 库 ID（从 search-libraries 获取）。不传则全局搜索所有库 |
| version | string | 否 | 版本号。不传则搜索所有版本 |
| mode | string | 否 | `code`（代码示例）或 `info`（文档说明）。不传则搜索所有类型 |
| page | int | 否 | 分页 1-10，默认 1 |

**响应（code 模式）：**

```json
{
  "jsonrpc": "2.0",
  "id": 2,
  "result": {
    "content": [
      {
        "type": "text",
        "text": "{\"libraryId\": 6, \"documents\": [{\"title\": \"Defining Routes with Different HTTP Methods in Gin\", \"description\": \"This code snippet demonstrates how to define routes for various HTTP methods using the Gin framework.\", \"source\": \"mcp/docs/gin/v1.9.1/docs/doc.md\", \"version\": \"v1.9.1\", \"mode\": \"code\", \"language\": \"go\", \"code\": \"func main() {...}\", \"tokens\": 319, \"relevance\": 0.134}], \"page\": 1, \"hasMore\": true}"
      }
    ]
  }
}
```

**响应（info 模式）：**

```json
{
  "jsonrpc": "2.0",
  "id": 2,
  "result": {
    "content": [
      {
        "type": "text",
        "text": "{\"libraryId\": 6, \"documents\": [{\"title\": \"Gin Web Framework > Getting started > Installation\", \"source\": \"mcp/docs/gin/v1.9.1/README.md\", \"version\": \"v1.9.1\", \"mode\": \"info\", \"content\": \"To install Gin package, you need to install Go and set your Go workspace first....\", \"tokens\": 105, \"relevance\": 0.096}], \"page\": 1, \"hasMore\": true}"
      }
    ]
  }
}
```

**响应说明：**

响应遵循 MCP 规范，将实际数据包装在 `content` 数组中，`text` 字段包含 JSON 字符串。

**数据字段说明（text 字段中的 JSON）：**

**code 模式：**

| 字段 | 类型 | 说明 |
|------|------|------|
| title | string | 代码标题（LLM 生成） |
| description | string | 代码描述（LLM 生成） |
| source | string | 来源文件路径 |
| version | string | 文档版本 |
| mode | string | 返回模式（`code`） |
| language | string | 代码语言 |
| code | string | 代码内容 |
| tokens | int | Token 数量 |
| relevance | float | 相关性分数（0-1） |

**info 模式：**

| 字段 | 类型 | 说明 |
|------|------|------|
| title | string | 文档标题层级 |
| source | string | 来源文件路径 |
| version | string | 文档版本 |
| mode | string | 返回模式（`info`） |
| content | string | 文档内容（Markdown） |
| tokens | int | Token 数量 |
| relevance | float | 相关性分数（0-1） |

---

## IDE 配置

### Cursor

```json
{
  "mcpServers": {
    "go-mcp-context": {
      "url": "https://mcp.hsk423.cn/mcp",
      "headers": {
        "MCP_API_KEY": "YOUR_API_KEY"
      }
    }
  }
}
```

### Claude Code

```bash
claude mcp add --transport http go-mcp-context https://mcp.hsk423.cn/mcp \
  --header "MCP_API_KEY: YOUR_API_KEY"
```

### VS Code

在 `settings.json` 中添加：

```json
"mcp": {
  "servers": {
    "go-mcp-context": {
      "type": "http",
      "url": "https://mcp.hsk423.cn/mcp",
      "headers": {
        "MCP_API_KEY": "YOUR_API_KEY"
      }
    }
  }
}
```

### Windsurf

```json
{
  "mcpServers": {
    "go-mcp-context": {
      "serverUrl": "https://mcp.hsk423.cn/mcp",
      "headers": {
        "MCP_API_KEY": "YOUR_API_KEY"
      }
    }
  }
}
```

### Codex

在 `codex.toml` 中添加：

```toml
[mcp_servers.go-mcp-context]
url = "https://mcp.hsk423.cn/mcp"
http_headers = { "MCP_API_KEY" = "YOUR_API_KEY" }
```

### Gemini CLI

```json
{
  "mcpServers": {
    "go-mcp-context": {
      "httpUrl": "https://mcp.hsk423.cn/mcp",
      "headers": {
        "MCP_API_KEY": "YOUR_API_KEY",
        "Accept": "application/json, text/event-stream"
      }
    }
  }
}
```

### 本地开发

将 URL 替换为 `http://localhost:8090/mcp`。

---

## API Key 获取

1. 登录 Web 管理界面
2. 进入「设置」→「API Keys」
3. 点击「创建 API Key」
4. 复制生成的 Key（仅显示一次）

详见 [API 文档 - API Key 管理接口](./API.md#api-key-管理接口)
