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

- 🔌 **MCP 协议规范化** - 完全遵循 MCP 规范，与 Costrict 等客户端无缝兼容
  - 协议无关的处理架构，支持 HTTP、SSE、Streamable HTTP 多种传输协议
  - 统一的 MCP 请求处理器，支持所有 MCP 方法
  - 规范的响应格式包装，确保客户端正确解析
- 🔍 **向量检索** - 基于 PostgreSQL + pgvector 的高性能向量搜索
- 📄 **多格式文档** - 支持 Markdown、PDF、DOCX、Swagger 等格式
- 🔀 **混合搜索** - 向量相似度 + BM25 关键词搜索
- 📊 **智能重排序** - 多指标评分优化搜索结果
- 🔐 **双重认证** - SSO JWT 管理 + API Key MCP 调用
- 🎨 **现代化 UI** - Vue3 + TypeScript + TailwindCSS
- 📦 **版本管理** - 支持库的多个版本，LLM 工作流指导

---

## 🚀 快速开始

详情见：[部署指南](docs/DEPLOYMENT.md)

---

## 🔧 IDE 配置

详情见：[https://mcp.hsk423.cn/dashboard](https://mcp.hsk423.cn/dashboard)

---

## 📚 文档

- [API 文档](docs/API.md) - REST API 和 MCP 工具说明
- [架构文档](docs/ARCHITECTURE.md) - 技术栈、项目结构、数据模型
- [部署指南](docs/DEPLOYMENT.md) - 环境配置、Docker 部署、Nginx 配置
- [开发日志](docs/CHANGELOG.md) - 版本更新记录

---

## 📋 开发计划

### MVP ✅

- [x] MCP 协议端点
- [x] SSO JWT + API Key 认证
- [x] 文档解析（Markdown）
- [x] Embedding 生成（OpenAI）
- [x] 向量搜索（pgvector）
- [x] 前端管理界面

### 进行中

- [ ] PDF/DOCX 解析
- [x] MCP IDE 集成测试

---

## 📄 License

MIT
