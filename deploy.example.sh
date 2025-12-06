#!/bin/bash

# CI/CD 部署脚本
# 功能：本地构建Docker镜像，上传到远程服务器并部署
# 支持分步执行：./deploy.sh [step]
# 步骤：1=构建镜像, 2=保存上传, 3=远程部署

set -e  # 遇到错误立即退出

# ==================== 配置区域 ====================
# 远程服务器配置
REMOTE_HOST="your_server_ip"     # 服务器IP
REMOTE_USER="root"            # SSH用户名
REMOTE_PORT="22"              # SSH端口，默认22

# 远程服务器路径
REMOTE_BASE_DIR="/path/to/your/deploy/dir"
REMOTE_IMAGE_DIR="${REMOTE_BASE_DIR}/docker_images"
REMOTE_COMPOSE_DIR="${REMOTE_BASE_DIR}"
REMOTE_DEPLOY_DIR="${REMOTE_BASE_DIR}/deploy"

# 本地项目根目录
LOCAL_PROJECT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"

# Docker镜像名称和标签
IMAGE_TAG="latest"
SERVER_MCP_IMAGE="server-mcp:${IMAGE_TAG}"
WEB_MCP_IMAGE="web-mcp:${IMAGE_TAG}"

# 颜色输出
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

# ==================== 辅助函数 ====================
log_info() {
    echo -e "${GREEN}[INFO]${NC} $1"
}

log_warn() {
    echo -e "${YELLOW}[WARN]${NC} $1"
}

log_error() {
    echo -e "${RED}[ERROR]${NC} $1"
}

# 检查命令是否存在
check_command() {
    if ! command -v $1 &> /dev/null; then
        log_error "$1 未安装，请先安装"
        exit 1
    fi
}

# ==================== 步骤函数 ====================

# 步骤1：构建Docker镜像
step_build_images() {
    log_info "=========================================="
    log_info "步骤1: 构建Docker镜像"
    log_info "=========================================="
    
    # 检查必要的命令
    log_info "检查必要的命令..."
    check_command docker
    
    log_info "开始构建Docker镜像..."
    
    log_info "构建 server-mcp 镜像..."
    cd "${LOCAL_PROJECT_DIR}/server-mcp"
    docker build -t ${SERVER_MCP_IMAGE} .
    
    log_info "构建 web-mcp 镜像..."
    cd "${LOCAL_PROJECT_DIR}/web-mcp"
    docker build -t ${WEB_MCP_IMAGE} .
    
    log_info "所有镜像构建完成！"
    
    # 清理多阶段构建的中间镜像和悬空镜像
    log_info "清理多阶段构建的中间镜像..."
    docker image prune -f 2>/dev/null || true
    docker builder prune -f 2>/dev/null || true
    log_info "镜像清理完成！"
}

# 按服务名构建单个镜像
step_build_single_image() {
    local service="$1"
    check_command docker
    case "$service" in
        server-mcp)
            log_info "构建 server-mcp 镜像..."
            cd "${LOCAL_PROJECT_DIR}/server-mcp"
            docker build -t ${SERVER_MCP_IMAGE} .
            ;;
        web-mcp)
            log_info "构建 web-mcp 镜像..."
            cd "${LOCAL_PROJECT_DIR}/web-mcp"
            docker build -t ${WEB_MCP_IMAGE} .
            ;;
        *)
            log_error "未知服务: $service"
            exit 1
            ;;
    esac
    log_info "单个服务镜像构建完成：$service"
}

# 步骤2：保存镜像并上传
step_save_and_upload() {
    log_info "=========================================="
    log_info "步骤2: 保存镜像并上传到远程服务器"
    log_info "=========================================="
    
    check_command ssh
    check_command scp
    
    # 保存镜像为tar文件
    log_info "保存镜像为tar文件..."
    TEMP_DIR=$(mktemp -d)
    cd ${TEMP_DIR}
    
    docker save ${SERVER_MCP_IMAGE} -o server-mcp.tar
    docker save ${WEB_MCP_IMAGE} -o web-mcp.tar
    
    log_info "镜像保存完成！"
    
    # 创建远程目录
    log_info "创建远程服务器目录..."
    ssh -T -q -p ${REMOTE_PORT} ${REMOTE_USER}@${REMOTE_HOST} << EOF
        mkdir -p ${REMOTE_IMAGE_DIR}
        mkdir -p ${REMOTE_COMPOSE_DIR}
        mkdir -p ${REMOTE_DEPLOY_DIR}/server-mcp/configs
        mkdir -p ${REMOTE_DEPLOY_DIR}/server-mcp/uploads
        mkdir -p ${REMOTE_DEPLOY_DIR}/server-mcp/log
        mkdir -p ${REMOTE_DEPLOY_DIR}/server-mcp/keys
EOF
    
    # 上传镜像文件
    log_info "上传镜像文件到远程服务器..."
    scp -P ${REMOTE_PORT} ${TEMP_DIR}/*.tar ${REMOTE_USER}@${REMOTE_HOST}:${REMOTE_IMAGE_DIR}/
    
    # 上传docker-compose文件
    log_info "上传docker-compose文件..."
    scp -P ${REMOTE_PORT} ${LOCAL_PROJECT_DIR}/docker-compose.prod.yml ${REMOTE_USER}@${REMOTE_HOST}:${REMOTE_COMPOSE_DIR}/docker-compose.prod.yml
    
    # 上传配置文件（正式环境）
    log_info "检查并上传正式环境配置文件..."
    
    # 复制 server-mcp 正式环境配置文件到挂载位置
    if [ -f "${LOCAL_PROJECT_DIR}/server-mcp/configs/config.prod.yaml" ]; then
        log_info "上传 server-mcp 正式环境配置文件..."
        scp -P ${REMOTE_PORT} ${LOCAL_PROJECT_DIR}/server-mcp/configs/config.prod.yaml ${REMOTE_USER}@${REMOTE_HOST}:${REMOTE_DEPLOY_DIR}/server-mcp/configs/ || log_warn "上传正式环境配置文件失败"
    else
        log_warn "server-mcp 正式环境配置文件不存在于本地: ${LOCAL_PROJECT_DIR}/server-mcp/configs/config.prod.yaml"
    fi
    
    # 上传 server-mcp 密钥文件
    if [ -d "${LOCAL_PROJECT_DIR}/server-mcp/keys" ] && [ "$(ls -A ${LOCAL_PROJECT_DIR}/server-mcp/keys 2>/dev/null)" ]; then
        log_info "上传 server-mcp 密钥文件..."
        scp -P ${REMOTE_PORT} -r ${LOCAL_PROJECT_DIR}/server-mcp/keys/* ${REMOTE_USER}@${REMOTE_HOST}:${REMOTE_DEPLOY_DIR}/server-mcp/keys/ 2>/dev/null || log_warn "密钥文件已存在，跳过上传"
    fi
    
    # 清理本地临时文件
    log_info "清理本地临时文件..."
    rm -rf ${TEMP_DIR}
    
    log_info "步骤2完成！"
}

# 保存并上传单个服务镜像
step_save_and_upload_single() {
    local service="$1"
    log_info "=========================================="
    log_info "步骤2(单服务): 保存并上传 ${service} 镜像"
    log_info "=========================================="
    check_command ssh
    check_command scp
    local tar_name=""
    local image_name=""
    case "$service" in
        server-mcp) tar_name="server-mcp.tar"; image_name="${SERVER_MCP_IMAGE}" ;;
        web-mcp)    tar_name="web-mcp.tar";    image_name="${WEB_MCP_IMAGE}" ;;
        *)
            log_error "未知服务: $service"
            exit 1
            ;;
    esac
    local TEMP_DIR
    TEMP_DIR=$(mktemp -d)
    cd ${TEMP_DIR}
    log_info "保存镜像为 tar：${tar_name}"
    docker save ${image_name} -o "${tar_name}"
    log_info "创建远程目录..."
    ssh -T -q -p ${REMOTE_PORT} ${REMOTE_USER}@${REMOTE_HOST} << EOF
        mkdir -p ${REMOTE_IMAGE_DIR}
        mkdir -p ${REMOTE_COMPOSE_DIR}
        mkdir -p ${REMOTE_DEPLOY_DIR}/server-mcp/configs
        mkdir -p ${REMOTE_DEPLOY_DIR}/server-mcp/uploads
        mkdir -p ${REMOTE_DEPLOY_DIR}/server-mcp/log
        mkdir -p ${REMOTE_DEPLOY_DIR}/server-mcp/keys
EOF
    log_info "上传 tar 到远程..."
    scp -P ${REMOTE_PORT} "${tar_name}" ${REMOTE_USER}@${REMOTE_HOST}:${REMOTE_IMAGE_DIR}/
    log_info "上传 docker-compose 文件..."
    scp -P ${REMOTE_PORT} ${LOCAL_PROJECT_DIR}/docker-compose.prod.yml ${REMOTE_USER}@${REMOTE_HOST}:${REMOTE_COMPOSE_DIR}/docker-compose.prod.yml
    # 上传必要配置与密钥（尽量幂等）
    if [ -f "${LOCAL_PROJECT_DIR}/server-mcp/configs/config.prod.yaml" ]; then
        scp -P ${REMOTE_PORT} ${LOCAL_PROJECT_DIR}/server-mcp/configs/config.prod.yaml ${REMOTE_USER}@${REMOTE_HOST}:${REMOTE_DEPLOY_DIR}/server-mcp/configs/ 2>/dev/null || true
    fi
    if [ -d "${LOCAL_PROJECT_DIR}/server-mcp/keys" ] && [ "$(ls -A ${LOCAL_PROJECT_DIR}/server-mcp/keys 2>/dev/null)" ]; then
        scp -P ${REMOTE_PORT} -r ${LOCAL_PROJECT_DIR}/server-mcp/keys/* ${REMOTE_USER}@${REMOTE_HOST}:${REMOTE_DEPLOY_DIR}/server-mcp/keys/ 2>/dev/null || true
    fi
    rm -rf ${TEMP_DIR}
    log_info "步骤2(单服务) 完成！"
}

# 步骤3：远程部署
step_deploy() {
    log_info "=========================================="
    log_info "步骤3: 在远程服务器部署"
    log_info "=========================================="
    
    # 在远程服务器加载镜像并部署
    log_info "在远程服务器加载镜像并部署..."
    ssh -p ${REMOTE_PORT} ${REMOTE_USER}@${REMOTE_HOST} << EOF
        set -e
        
        echo "加载Docker镜像..."
        cd ${REMOTE_IMAGE_DIR}
        docker load -i server-mcp.tar
        docker load -i web-mcp.tar
        
        # 停止旧业务服务容器（如果存在）
        cd ${REMOTE_COMPOSE_DIR}
        echo "停止旧业务服务容器..."
        docker-compose -f docker-compose.prod.yml down 2>/dev/null || true
        
        # 启动业务服务
        # 注意：基础服务（postgres, redis）需要单独管理，不在此脚本中处理
        echo "启动业务服务..."
        docker-compose -f docker-compose.prod.yml up -d
        
        # 清理旧的镜像 tar 文件
        echo "清理临时文件..."
        rm -f ${REMOTE_IMAGE_DIR}/*.tar
        
        echo "部署完成！"
EOF
    
    log_info "=========================================="
    log_info "部署完成！"
    log_info "=========================================="
    log_info "远程服务器: ${REMOTE_HOST}"
    log_info "部署目录: ${REMOTE_COMPOSE_DIR}"
    log_info ""
    log_info "查看服务状态:"
    log_info "  ssh ${REMOTE_USER}@${REMOTE_HOST} 'cd ${REMOTE_COMPOSE_DIR} && docker-compose -f docker-compose.prod.yml ps'"
    log_info ""
    log_info "查看日志:"
    log_info "  ssh ${REMOTE_USER}@${REMOTE_HOST} 'cd ${REMOTE_COMPOSE_DIR} && docker-compose -f docker-compose.prod.yml logs -f'"
}

# 远程部署单个服务
step_deploy_single() {
    local service="$1"
    log_info "=========================================="
    log_info "步骤3(单服务): 在远程服务器部署 ${service}"
    log_info "=========================================="
    local tar_name=""
    case "$service" in
        server-mcp) tar_name="server-mcp.tar" ;;
        web-mcp)    tar_name="web-mcp.tar" ;;
        *)
            log_error "未知服务: $service"
            exit 1
            ;;
    esac
    ssh -p ${REMOTE_PORT} ${REMOTE_USER}@${REMOTE_HOST} << EOF
        set -e
        echo "加载 ${service} Docker 镜像..."
        cd ${REMOTE_IMAGE_DIR}
        if [ -f "${tar_name}" ]; then
            docker load -i ${tar_name}
        else
            echo "未找到镜像 tar：${tar_name}，请先执行 upload。"
            exit 1
        fi
        echo "停止旧容器..."
        cd ${REMOTE_COMPOSE_DIR}
        docker-compose -f docker-compose.prod.yml stop ${service} 2>/dev/null || true
        
        echo "启动 ${service}..."
        docker-compose -f docker-compose.prod.yml up -d ${service}
        
        # 清理镜像 tar 文件
        echo "清理临时文件..."
        rm -f ${REMOTE_IMAGE_DIR}/${tar_name}
        
        echo "单服务部署完成！"
EOF
    log_info "单服务部署完成：${service}"
}

# ==================== 主流程 ====================

# 检查参数
if [ $# -eq 0 ]; then
    log_error "必须指定执行步骤！"
    echo ""
    echo "用法: $0 <步骤>"
    echo ""
    echo "步骤选项:"
    echo "  1 或 build  - 构建Docker镜像"
    echo "  2 或 upload - 保存镜像并上传到远程服务器"
    echo "  3 或 deploy - 在远程服务器部署"
    echo "  all 或 ALL  - 执行所有步骤（完整部署流程）"
    echo ""
    echo "示例:"
    echo "  $0 1        # 构建镜像"
    echo "  $0 build    # 构建镜像（别名）"
    echo "  $0 2        # 上传文件"
    echo "  $0 upload   # 上传文件（别名）"
    echo "  $0 3        # 远程部署"
    echo "  $0 deploy   # 远程部署（别名）"
    echo "  $0 all      # 执行完整部署流程"
    echo ""
    echo "单服务全流程/分步:"
    echo "  $0 single <service> [build|upload|deploy|all]"
    echo "支持的 <service>：server-mcp | web-mcp"
    echo "示例："
    echo "  $0 single server-mcp           # 单服务完整流程"
    echo "  $0 single web-mcp build        # 仅构建 web-mcp"
    echo "  $0 single server-mcp upload    # 仅上传 server-mcp"
    echo "  $0 single web-mcp deploy       # 仅部署 web-mcp"
    exit 1
fi

# 解析参数
STEP=$1

case ${STEP} in
    1|build)
        step_build_images
        ;;
    2|upload)
        step_save_and_upload
        ;;
    3|deploy)
        step_deploy
        ;;
    all|ALL)
        log_info "=========================================="
        log_info "执行完整部署流程（所有步骤）"
        log_info "=========================================="
        step_build_images
        step_save_and_upload
        step_deploy
        log_info "=========================================="
        log_info "完整部署流程执行完成！"
        log_info "=========================================="
        ;;
    single)
        SERVICE="$2"
        ACTION="$3"
        if [ -z "${SERVICE}" ]; then
            log_error "single 模式需要指定 <service>"
            echo "用法：$0 single <service> [build|upload|deploy|all]"
            exit 1
        fi
        if [ -z "${ACTION}" ]; then
            ACTION="all"
        fi
        case "${ACTION}" in
            build)
                step_build_single_image "${SERVICE}"
                ;;
            upload)
                step_save_and_upload_single "${SERVICE}"
                ;;
            deploy)
                step_deploy_single "${SERVICE}"
                ;;
            all|ALL)
                step_build_single_image "${SERVICE}"
                step_save_and_upload_single "${SERVICE}"
                step_deploy_single "${SERVICE}"
                ;;
            *)
                log_error "未知动作: ${ACTION}，可选：build|upload|deploy|all"
                exit 1
                ;;
        esac
        ;;
    *)
        log_error "未知步骤: ${STEP}"
        echo ""
        echo "用法: $0 <步骤>"
        echo ""
        echo "步骤选项:"
        echo "  1 或 build  - 构建Docker镜像"
        echo "  2 或 upload - 保存镜像并上传到远程服务器"
        echo "  3 或 deploy - 在远程服务器部署"
        echo "  all 或 ALL  - 执行所有步骤（完整部署流程）"
        exit 1
        ;;
esac
