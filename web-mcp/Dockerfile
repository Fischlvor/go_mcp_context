# 第一阶段：构建阶段
FROM node:22-alpine AS builder

# 写入前端构建所需环境变量（可被 --build-arg 覆盖）
ARG VITE_BASE_API=/api/v1
ENV VITE_BASE_API=$VITE_BASE_API

WORKDIR /app

# 使用国内镜像源加速
RUN sed -i 's/dl-cdn.alpinelinux.org/mirrors.aliyun.com/g' /etc/apk/repositories

# 设置 npm 镜像
RUN npm config set registry https://registry.npmmirror.com

# 复制 package 文件
COPY package*.json ./

# 安装依赖
RUN npm ci

# 复制源代码
COPY . .

# 构建
RUN npm run build-only

# 第二阶段：运行阶段
FROM nginx:alpine

# 使用国内镜像源加速,安装 wget 用于 healthcheck
RUN sed -i 's/dl-cdn.alpinelinux.org/mirrors.aliyun.com/g' /etc/apk/repositories && \
    apk --no-cache add wget

# 复制构建产物到nginx目录
COPY --from=builder /app/dist /usr/share/nginx/html

# 复制nginx配置
COPY nginx.conf /etc/nginx/conf.d/default.conf

# 暴露端口
EXPOSE 80

CMD ["nginx", "-g", "daemon off;"]
