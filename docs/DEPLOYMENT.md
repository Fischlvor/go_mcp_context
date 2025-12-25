# 部署指南

## 环境要求

- Go 1.23+
- Node.js 22.15.0+ (推荐使用 nvm 管理)
- Docker & Docker Compose
- OpenAI API Key
- PostgreSQL 15+ (with pgvector)
- Redis 6+

---

## 本地开发

```bash
# 克隆项目
git clone https://github.com/Fischlvor/go_mcp_context.git
cd go-mcp-context

# 启动依赖服务
docker-compose up -d postgres redis

# 运行后端（默认使用 configs/config.yaml）
cd server-mcp
go run main.go
# 生产环境使用 APP_ENV=prod go run main.go（加载 configs/config.prod.yaml）

# 运行前端
cd web-mcp
nvm use 22.15.0
npm install
npm run dev
```

---

## Docker 部署

### Docker 部署（开发环境）

```bash
# 启动所有服务
docker-compose up -d
```

### Docker 部署（生产环境）

**基础设施服务（PostgreSQL + Redis）：**

在 `go_blog` 项目中启动基础设施：

```yaml
# go_blog/docker-compose.base.yml（摘录）
services:
  # PostgreSQL + pgvector
  postgres:
    image: pgvector/pgvector:pg15
    container_name: agent_postgres
    restart: always
    environment:
      POSTGRES_USER: postgres
      POSTGRES_PASSWORD: root
      POSTGRES_DB: agent_db
    ports:
      - "15432:5432"
    volumes:
      - ./docker-data/postgres:/var/lib/postgresql/data
    networks:
      - infrastructure-network
    healthcheck:
      test: ["CMD-SHELL", "pg_isready -U postgres"]
      interval: 10s
      timeout: 5s
      retries: 5

  # Redis
  redis:
    image: redis:6.2
    container_name: redis6
    restart: always
    ports:
      - "16379:6379"
    volumes:
      - ./docker-data/redis/redis.conf:/usr/local/etc/redis/redis.conf
      - ./docker-data/redis:/data
    command: redis-server /usr/local/etc/redis/redis.conf
    networks:
      - infrastructure-network

networks:
  infrastructure-network:
    driver: bridge
    name: infrastructure-network
```

**MCP 服务（后端 + 前端）：**

```yaml
# docker-compose.prod.yml
version: '3.8'

services:
  # MCP Context 后端服务
  server-mcp:
    image: server-mcp:latest
    container_name: server-mcp
    restart: always
    environment:
      - TZ=Asia/Shanghai
      - APP_ENV=prod
    volumes:
      - ./deploy/server-mcp/configs:/app/configs
      - ./deploy/server-mcp/uploads:/app/uploads
      - ./deploy/server-mcp/log:/app/log
      - ./deploy/server-mcp/keys:/app/keys
    networks:
      - mcp-network
      - infrastructure-network
      - blog-network
    healthcheck:
      test: ["CMD", "wget", "--no-verbose", "--tries=1", "-O", "/dev/null", "http://localhost:8090/api/base/health"]
      interval: 30s
      timeout: 10s
      retries: 3
      start_period: 40s

  # MCP Context 前端服务
  web-mcp:
    image: web-mcp:latest
    container_name: web-mcp
    restart: always
    environment:
      - TZ=Asia/Shanghai
    networks:
      - mcp-network
      - blog-network
    depends_on:
      - server-mcp
    healthcheck:
      test: ["CMD", "wget", "--no-verbose", "--tries=1", "-O", "/dev/null", "http://localhost/"]
      interval: 30s
      timeout: 10s
      retries: 3
      start_period: 10s

networks:
  mcp-network:
    driver: bridge
    name: mcp-network
  infrastructure-network:
    external: true
    name: infrastructure-network
  blog-network:
    external: true
    name: blog-network
```

---

## 部署脚本使用

项目提供 `deploy.sh` 脚本，支持本地构建 Docker 镜像并部署到远程服务器。

### 使用方法

```bash
# 完整部署（构建 + 上传 + 部署）
./deploy.sh all

# 分步执行
./deploy.sh build    # 或 1，构建 Docker 镜像
./deploy.sh upload   # 或 2，保存镜像并上传到服务器
./deploy.sh deploy   # 或 3，远程服务器部署

# 单服务部署
./deploy.sh single server-mcp all      # 完整部署 server-mcp
./deploy.sh single web-mcp all         # 完整部署 web-mcp
./deploy.sh single server-mcp build    # 只构建 server-mcp
./deploy.sh single server-mcp upload   # 只上传 server-mcp
./deploy.sh single server-mcp deploy   # 只部署 server-mcp
```

### 部署流程

1. **build** - 本地构建 `server-mcp` 和 `web-mcp` Docker 镜像
2. **upload** - 将镜像导出为 tar 文件，通过 SCP 上传到服务器
3. **deploy** - SSH 到服务器，加载镜像并启动容器

### 配置说明

部署前需要修改 `deploy.sh` 中的配置：

```bash
REMOTE_HOST="your-server-ip"      # 服务器 IP
REMOTE_USER="root"                # SSH 用户名
REMOTE_PORT="22"                  # SSH 端口
REMOTE_BASE_DIR="/path/to/deploy" # 远程部署目录
```

### 目录结构

部署后服务器目录结构：

```text
/path/to/deploy/
├── docker-compose.prod.yml
├── docker_images/           # Docker 镜像 tar 文件
├── deploy/
│   └── server-mcp/
│       ├── configs/         # 配置文件
│       ├── uploads/         # 上传文件
│       ├── log/             # 日志文件
│       └── keys/            # SSO 公钥等
```

---

## 环境变量

| 变量 | 必填 | 说明 |
|------|------|------|
| OPENAI_API_KEY | 是 | OpenAI API 密钥 |
| JWT_SECRET | 是 | JWT 签名密钥 |
| DB_HOST | 否 | PostgreSQL 主机，默认 localhost |
| DB_PORT | 否 | PostgreSQL 端口，默认 5432 |
| DB_USER | 否 | PostgreSQL 用户，默认 postgres |
| DB_PASSWORD | 否 | PostgreSQL 密码 |
| DB_NAME | 否 | 数据库名，默认 mcp_context |
| REDIS_HOST | 否 | Redis 主机，默认 localhost |
| REDIS_PORT | 否 | Redis 端口，默认 6379 |
| GITHUB_TOKEN | 否 | GitHub API Token（导入功能） |
| GITHUB_PROXY | 否 | GitHub API 代理地址 |

---

## Nginx 配置示例

以下是生产环境使用的 Nginx 配置（Docker 容器间通信）：

```nginx
# ==================== 限流配置 ====================
limit_req_zone $binary_remote_addr zone=global_limit:10m rate=10000r/s;
limit_req_log_level warn;
limit_req_status 429;

# HTTP 重定向到 HTTPS
server {
    listen 80;
    server_name mcp.hsk423.cn;
    return 301 https://mcp.hsk423.cn$request_uri;
}

# HTTPS 主配置
server {
    # Gzip 压缩
    gzip on;
    gzip_vary on;
    gzip_disable "MSIE [1-6]\.";
    gzip_static on;
    gzip_min_length 256;
    gzip_buffers 32 8k;
    gzip_http_version 1.1;
    gzip_comp_level 5;
    gzip_proxied any;
    gzip_types text/plain text/css text/xml application/javascript application/x-javascript application/xml application/xml+rss application/emacscript application/json image/svg+xml;

    listen 443 ssl;
    server_name mcp.hsk423.cn;
    
    # SSL 证书
    ssl_certificate /etc/nginx/ssl/certificate.pem;
    ssl_certificate_key /etc/nginx/ssl/private.key;
    ssl_session_timeout 5m;
    ssl_session_cache shared:MozSSL:10m;
    ssl_ciphers ECDHE-ECDSA-AES128-GCM-SHA256:ECDHE-RSA-AES128-GCM-SHA256:ECDHE-ECDSA-AES256-GCM-SHA384:ECDHE-RSA-AES256-GCM-SHA384:ECDHE-ECDSA-CHACHA20-POLY1305:ECDHE-RSA-CHACHA20-POLY1305:DHE-RSA-AES128-GCM-SHA256:DHE-RSA-AES256-GCM-SHA384;
    ssl_protocols TLSv1.2 TLSv1.3;
    ssl_prefer_server_ciphers on;

    # 安全头
    add_header Strict-Transport-Security "max-age=63072000; includeSubDomains; preload" always;
    add_header X-Frame-Options DENY;
    add_header X-Content-Type-Options nosniff;
    add_header X-XSS-Protection "1; mode=block";

    # SSE 流式响应特殊配置（文档上传进度）
    location /api/v1/document/upload {
        proxy_buffering off;
        proxy_request_buffering off;
        proxy_connect_timeout 300s;
        proxy_send_timeout 300s;
        proxy_read_timeout 300s;
        proxy_http_version 1.1;
        proxy_set_header Connection "";
        proxy_set_header Host $host;
        proxy_set_header X-Real-IP $remote_addr;
        proxy_set_header X-Forwarded-For $proxy_add_x_forwarded_for;
        proxy_set_header X-Forwarded-Proto $scheme;
        proxy_set_header Accept-Encoding "";
        proxy_pass http://server-mcp:8090/api/v1/document/upload;
    }

    # API 代理到后端服务
    location /api/ {
        limit_req zone=global_limit burst=1000 nodelay;
        
        proxy_set_header Host $host;
        proxy_set_header X-Real-IP $remote_addr;
        proxy_set_header REMOTE-HOST $remote_addr;
        proxy_set_header X-Forwarded-For $proxy_add_x_forwarded_for;
        proxy_set_header X-Forwarded-Proto $scheme;
        proxy_pass http://server-mcp:8090/api/;
    }

    # MCP 端点代理
    location /mcp {
        proxy_set_header Host $host;
        proxy_set_header X-Real-IP $remote_addr;
        proxy_set_header X-Forwarded-For $proxy_add_x_forwarded_for;
        proxy_set_header X-Forwarded-Proto $scheme;
        proxy_pass http://server-mcp:8090/mcp;
    }

    # 静态文件代理到前端容器
    location / {
        proxy_set_header Host $host;
        proxy_set_header X-Real-IP $remote_addr;
        proxy_set_header X-Forwarded-For $proxy_add_x_forwarded_for;
        proxy_set_header X-Forwarded-Proto $scheme;
        proxy_pass http://web-mcp:80;
    }
}
```

**说明：**

- `server-mcp:8090` - 后端服务容器名和端口
- `web-mcp:80` - 前端服务容器名和端口
- SSE 接口需要禁用 `proxy_buffering` 以支持流式响应
- 使用 `limit_req` 进行请求限流保护

---

## 数据库初始化

PostgreSQL 需要安装 pgvector 扩展：

```sql
CREATE EXTENSION IF NOT EXISTS vector;
```

GORM 会自动迁移表结构，无需手动创建表。
