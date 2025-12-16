# 第一阶段：构建阶段
FROM golang:1.23.0-alpine AS builder

# 设置工作目录
WORKDIR /app

# 使用国内镜像源加速,并安装必要的系统依赖
RUN sed -i 's/dl-cdn.alpinelinux.org/mirrors.aliyun.com/g' /etc/apk/repositories && \
    apk add --no-cache ca-certificates tzdata git

# 设置Go代理（加速下载）
ENV GOPROXY=https://goproxy.cn,direct
ENV GOSUMDB=sum.golang.org

# 复制go mod文件
COPY go.mod go.sum ./

# 下载依赖
RUN go mod download

# 复制源代码
COPY . .

# 设置环境变量
ENV CGO_ENABLED=0
ENV GOOS=linux
ENV GOARCH=amd64

# 构建应用
RUN go build -ldflags="-s -w" -o main .

# 第二阶段：运行阶段
FROM alpine:latest

# 使用国内镜像源加速,安装运行时依赖并清理,同时创建必要的目录
RUN sed -i 's/dl-cdn.alpinelinux.org/mirrors.aliyun.com/g' /etc/apk/repositories && \
    apk --no-cache add ca-certificates tzdata wget && \
    rm -rf /var/cache/apk/* /tmp/* /var/tmp/* && \
    mkdir -p /app/configs /app/uploads /app/logs /app/keys

# 设置工作目录
WORKDIR /app

# 从构建阶段复制二进制文件和配置
COPY --from=builder /app/main .
COPY --from=builder /app/configs ./configs

# 暴露端口
EXPOSE 8090

# 启动应用 (config file is selected by APP_ENV environment variable)
CMD ["./main"]
