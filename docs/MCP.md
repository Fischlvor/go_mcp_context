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

## 调用工具

```http
POST /mcp
MCP_API_KEY: <API_KEY>
```

### search-libraries

搜索文档库。

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
    "libraries": [
      {
        "libraryId": 1,
        "name": "gin",
        "versions": ["v1.9.0", "v1.8.0"],
        "defaultVersion": "latest",
        "description": "Gin is a HTTP web framework",
        "snippets": 150,
        "score": 0.95
      }
    ]
  }
}
```

**响应字段说明：**

| 字段 | 类型 | 说明 |
|------|------|------|
| libraryId | uint | 库 ID（用于 get-library-docs） |
| name | string | 库名称 |
| versions | string[] | 额外版本列表（不含 defaultVersion） |
| defaultVersion | string | 默认版本（通常为 `latest`） |
| description | string | 库描述 |
| snippets | int | 文档片段数量 |
| score | float | 匹配分数（0-1） |

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
    "libraryId": 6,
    "documents": [
      {
        "title": "Defining Routes with Different HTTP Methods in Gin",
        "description": "This code snippet demonstrates how to define routes for various HTTP methods using the Gin framework.",
        "source": "mcp/docs/gin/v1.9.1/docs/doc.md",
        "language": "go",
        "code": "func main() {\n  router := gin.Default()\n  router.GET(\"/someGet\", getting)\n  router.POST(\"/somePost\", posting)\n  router.Run()\n}",
        "tokens": 319,
        "relevance": 0.134
      }
    ],
    "page": 1,
    "hasMore": true
  }
}
```

**响应（info 模式）：**

```json
{
  "jsonrpc": "2.0",
  "id": 2,
  "result": {
    "libraryId": 6,
    "documents": [
      {
        "title": "Gin Web Framework > Getting started > Installation",
        "source": "mcp/docs/gin/v1.9.1/README.md",
        "content": "To install Gin package, you need to install Go and set your Go workspace first....",
        "tokens": 105,
        "relevance": 0.096
      }
    ],
    "page": 1,
    "hasMore": true
  }
}
```

**响应字段说明（code 模式）：**

| 字段 | 类型 | 说明 |
|------|------|------|
| title | string | 代码标题（LLM 生成） |
| description | string | 代码描述（LLM 生成） |
| source | string | 来源文件路径 |
| language | string | 代码语言 |
| code | string | 代码内容 |
| tokens | int | Token 数量 |
| relevance | float | 相关性分数（0-1） |

**响应字段说明（info 模式）：**

| 字段 | 类型 | 说明 |
|------|------|------|
| title | string | 文档标题层级 |
| source | string | 来源文件路径 |
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
