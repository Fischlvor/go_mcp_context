#!/bin/bash

# ============================================
# MCP 文档检索性能测试
# ============================================

set -e

# 检查参数
if [ $# -lt 1 ]; then
    echo "用法: $0 <MCP_API_KEY> [LIBRARY_ID] [VERSION] [MODE]"
    echo ""
    echo "参数:"
    echo "  MODE: remote (默认，通过 HTTPS) 或 local (服务器本地测试)"
    echo ""
    echo "示例:"
    echo "  $0 mcp_xxxxxxxx                    # 远程测试"
    echo "  $0 mcp_xxxxxxxx 6 v1.11.0          # 指定库和版本"
    echo "  $0 mcp_xxxxxxxx 6 v1.11.0 local    # 本地测试"
    exit 1
fi

# 配置
API_KEY="$1"
LIBRARY_ID="${2:-6}"
VERSION="${3:-v1.11.0}"
MODE="${4:-remote}"

# 如果是 local 模式，通过 SSH 在服务器上执行
if [ "$MODE" = "local" ]; then
    echo "本地模式：通过 SSH 在服务器上执行测试..."
    echo ""
    
    REMOTE_HOST="root@8.148.64.96"
    REMOTE_SCRIPT="/tmp/benchmark_local_$$.sh"
    REMOTE_RESULT="/tmp/benchmark_result_local_$(date +%Y%m%d_%H%M%S).txt"
    
    # 获取本地脚本目录
    LOCAL_SCRIPT_DIR="$(cd "$(dirname "$0")" && pwd)"
    LOCAL_RESULT="$LOCAL_SCRIPT_DIR/benchmark_result_local_$(date +%Y%m%d_%H%M%S).txt"
    
    # 复制脚本到服务器
    scp "$0" "$REMOTE_HOST:$REMOTE_SCRIPT" > /dev/null 2>&1
    
    # 在服务器上执行
    # 如果设置了 BATCH_MODE 环境变量，只输出到标准输出
    if [ -n "$BATCH_MODE" ]; then
        ssh "$REMOTE_HOST" "bash $REMOTE_SCRIPT '$API_KEY' '$LIBRARY_ID' '$VERSION' 'local_exec' 2>&1 | sed 's/\x1b\[[0-9;]*m//g'"
    else
        # 单独执行时，保存到本地文件
        ssh "$REMOTE_HOST" "bash $REMOTE_SCRIPT '$API_KEY' '$LIBRARY_ID' '$VERSION' 'local_exec' 2>&1 | sed 's/\x1b\[[0-9;]*m//g' | tee $REMOTE_RESULT"
        scp "$REMOTE_HOST:$REMOTE_RESULT" "$LOCAL_RESULT" > /dev/null 2>&1
        echo ""
        echo "结果已保存到: $LOCAL_RESULT"
    fi
    
    # 清理服务器上的临时文件
    ssh "$REMOTE_HOST" "rm -f $REMOTE_SCRIPT $REMOTE_RESULT" > /dev/null 2>&1
    
    exit 0
fi

# 根据模式设置 URL
if [ "$MODE" = "local_exec" ]; then
    MCP_URL="http://192.168.16.7:8090/mcp"
    TEST_MODE="本地测试 (容器内网)"
else
    MCP_URL="https://mcp.hsk423.cn/mcp"
    TEST_MODE="远程测试 (HTTPS)"
fi

# 颜色输出
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m'

echo -e "${BLUE}========================================${NC}"
echo -e "${BLUE}MCP 文档检索性能测试${NC}"
echo -e "${BLUE}========================================${NC}"
echo ""
echo "配置:"
echo "  测试模式:   $TEST_MODE"
echo "  MCP URL:    $MCP_URL"
echo "  Library ID: $LIBRARY_ID"
echo "  Version:    $VERSION"
echo ""

# 获取脚本所在目录
SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"

# 创建临时目录（相对于脚本目录）
TEMP_DIR="$SCRIPT_DIR/benchmark_temp"
mkdir -p "$TEMP_DIR"

# 创建结果文件
RESULT_FILE="$SCRIPT_DIR/benchmark_result_$(date +%Y%m%d_%H%M%S).txt"

# 重定向输出到文件和终端（去除颜色代码）
exec > >(tee -a >(sed 's/\x1b\[[0-9;]*m//g' > "$RESULT_FILE"))
exec 2>&1

# ============================================
# 测试 1：冷启动 (无缓存)
# ============================================
echo -e "${YELLOW}[1] 冷启动测试 (无缓存)${NC}"
echo "说明: 每次请求使用不同的 topic，避免缓存命中"
echo ""

COLD_TIME=0
for i in {1..10}; do
    # 生成唯一的 topic
    TOPIC="cold_start_test_$(date +%s%N)_$i"
    
    cat > "$TEMP_DIR/cold_$i.json" << EOF
{
  "jsonrpc": "2.0",
  "id": $i,
  "method": "tools/call",
  "params": {
    "name": "get-library-docs",
    "arguments": {
      "libraryId": $LIBRARY_ID,
      "version": "$VERSION",
      "topic": "$TOPIC",
      "mode": "code",
      "page": 1
    }
  }
}
EOF
    
    START=$(date +%s%N)
    RESPONSE=$(curl -s -X POST "$MCP_URL" \
      -H "Content-Type: application/json" \
      -H "MCP_API_KEY: $API_KEY" \
      -d @"$TEMP_DIR/cold_$i.json")
    END=$(date +%s%N)
    
    # 检查是否有错误
    if echo "$RESPONSE" | grep -q '"error"'; then
        echo -e "  请求 $i: ${RED}失败${NC}"
        echo "  错误: $(echo $RESPONSE | jq -r '.error.message' 2>/dev/null || echo $RESPONSE)"
        continue
    fi
    
    ELAPSED=$((($END - $START) / 1000000))
    COLD_TIME=$(($COLD_TIME + $ELAPSED))
    echo "  请求 $i: ${ELAPSED}ms"
done

AVG_COLD=$(($COLD_TIME / 10))
echo -e "${GREEN}✓ 冷启动平均延迟: ${AVG_COLD}ms${NC}"
echo ""

# ============================================
# 测试 2：热启动 (缓存命中)
# ============================================
echo -e "${YELLOW}[2] 热启动测试 (缓存命中)${NC}"
echo "说明: 重复相同的请求，测试缓存效果"
echo ""

cat > "$TEMP_DIR/hot.json" << EOF
{
  "jsonrpc": "2.0",
  "id": 1,
  "method": "tools/call",
  "params": {
    "name": "get-library-docs",
    "arguments": {
      "libraryId": $LIBRARY_ID,
      "version": "$VERSION",
      "topic": "routing",
      "mode": "code",
      "page": 1
    }
  }
}
EOF

# 预热缓存
echo "预热缓存..."
curl -s -X POST "$MCP_URL" \
  -H "Content-Type: application/json" \
  -H "MCP_API_KEY: $API_KEY" \
  -d @"$TEMP_DIR/hot.json" > /dev/null

# 测试缓存命中 (100 请求, 10 并发)
echo "执行 100 个相同查询 (10 并发)..."
ab -n 100 -c 10 \
   -p "$TEMP_DIR/hot.json" \
   -T "application/json" \
   -H "MCP_API_KEY: $API_KEY" \
   "$MCP_URL" > "$SCRIPT_DIR/hot_result.txt" 2>&1

# 提取关键指标
HOT_QPS=$(grep "Requests per second" "$SCRIPT_DIR/hot_result.txt" | awk '{print $4}')
HOT_P50=$(grep "50%" "$SCRIPT_DIR/hot_result.txt" | awk '{print $2}')
HOT_P95=$(grep "95%" "$SCRIPT_DIR/hot_result.txt" | awk '{print $2}')
HOT_P99=$(grep "99%" "$SCRIPT_DIR/hot_result.txt" | awk '{print $2}')

echo -e "${GREEN}✓ 热启动 QPS:  ${HOT_QPS}${NC}"
echo -e "${GREEN}✓ 热启动 P50:  ${HOT_P50}ms${NC}"
echo -e "${GREEN}✓ 热启动 P95:  ${HOT_P95}ms${NC}"
echo -e "${GREEN}✓ 热启动 P99:  ${HOT_P99}ms${NC}"
echo ""

# ============================================
# 测试 3：混合场景
# ============================================
echo -e "${YELLOW}[3] 混合场景测试 (80% 热门 + 20% 新查询)${NC}"
echo "说明: 模拟真实场景"
echo ""

# 生成 5 个热门查询
for i in {1..5}; do
    cat > "$TEMP_DIR/popular_$i.json" << EOF
{
  "jsonrpc": "2.0",
  "id": $i,
  "method": "tools/call",
  "params": {
    "name": "get-library-docs",
    "arguments": {
      "libraryId": $LIBRARY_ID,
      "version": "$VERSION",
      "topic": "popular_topic_$i",
      "mode": "code",
      "page": 1
    }
  }
}
EOF
done

# 预热热门查询
echo "预热热门查询缓存..."
for i in {1..5}; do
    curl -s -X POST "$MCP_URL" \
      -H "Content-Type: application/json" \
      -H "MCP_API_KEY: $API_KEY" \
      -d @"$TEMP_DIR/popular_$i.json" > /dev/null
done

# 执行混合测试
echo "执行混合场景测试 (100 个请求)..."
MIXED_TIME=0
for i in {1..100}; do
    RAND=$((RANDOM % 100))
    if [ $RAND -lt 80 ]; then
        # 80% 热门查询
        POPULAR_ID=$(((RANDOM % 5) + 1))
        REQUEST_FILE="$TEMP_DIR/popular_$POPULAR_ID.json"
    else
        # 20% 新查询
        TOPIC="new_topic_$i_$RANDOM"
        cat > "$TEMP_DIR/new_$i.json" << EOF
{
  "jsonrpc": "2.0",
  "id": $i,
  "method": "tools/call",
  "params": {
    "name": "get-library-docs",
    "arguments": {
      "libraryId": $LIBRARY_ID,
      "version": "$VERSION",
      "topic": "$TOPIC",
      "mode": "code",
      "page": 1
    }
  }
}
EOF
        REQUEST_FILE="$TEMP_DIR/new_$i.json"
    fi
    
    START=$(date +%s%N)
    curl -s -X POST "$MCP_URL" \
      -H "Content-Type: application/json" \
      -H "MCP_API_KEY: $API_KEY" \
      -d @"$REQUEST_FILE" > /dev/null
    END=$(date +%s%N)
    
    ELAPSED=$((($END - $START) / 1000000))
    MIXED_TIME=$(($MIXED_TIME + $ELAPSED))
    
    if [ $(($i % 20)) -eq 0 ]; then
        echo "  已完成 $i/100 个请求"
    fi
done

AVG_MIXED=$(($MIXED_TIME / 100))
echo -e "${GREEN}✓ 混合场景平均延迟: ${AVG_MIXED}ms${NC}"
echo ""

# ============================================
# 测试总结
# ============================================
echo -e "${BLUE}========================================${NC}"
echo -e "${BLUE}测试总结${NC}"
echo -e "${BLUE}========================================${NC}"
echo ""
echo "1. 冷启动 (无缓存):"
echo "   - 平均延迟: ${AVG_COLD}ms"
echo ""
echo "2. 热启动 (缓存命中):"
echo "   - QPS:      ${HOT_QPS}"
echo "   - P50:      ${HOT_P50}ms"
echo "   - P95:      ${HOT_P95}ms"
echo "   - P99:      ${HOT_P99}ms"
echo ""
echo "3. 混合场景 (80% 热门 + 20% 新查询):"
echo "   - 平均延迟: ${AVG_MIXED}ms"
echo ""

# 计算缓存加速比
if [ -n "$HOT_P50" ] && [ "$HOT_P50" != "0" ]; then
    SPEEDUP=$(echo "scale=2; $AVG_COLD / $HOT_P50" | bc)
    echo -e "${YELLOW}缓存加速比: ${SPEEDUP}x${NC}"
    echo ""
fi

# 性能评估
echo "性能评估:"
if [ $AVG_COLD -lt 300 ]; then
    echo -e "  冷启动: ${GREEN}✓ 优秀 (<300ms)${NC}"
elif [ $AVG_COLD -lt 500 ]; then
    echo -e "  冷启动: ${YELLOW}○ 良好 (<500ms)${NC}"
else
    echo -e "  冷启动: ${RED}✗ 需要优化 (>500ms)${NC}"
fi

if [ -n "$HOT_P50" ]; then
    if [ $HOT_P50 -lt 10 ]; then
        echo -e "  热启动: ${GREEN}✓ 优秀 (<10ms)${NC}"
    elif [ $HOT_P50 -lt 20 ]; then
        echo -e "  热启动: ${YELLOW}○ 良好 (<20ms)${NC}"
    else
        echo -e "  热启动: ${RED}✗ 需要优化 (>20ms)${NC}"
    fi
fi

if [ -n "$HOT_QPS" ]; then
    QPS_INT=$(echo "$HOT_QPS" | cut -d. -f1)
    if [ $QPS_INT -gt 100 ]; then
        echo -e "  QPS:    ${GREEN}✓ 优秀 (>100)${NC}"
    elif [ $QPS_INT -gt 50 ]; then
        echo -e "  QPS:    ${YELLOW}○ 良好 (>50)${NC}"
    else
        echo -e "  QPS:    ${RED}✗ 需要优化 (<50)${NC}"
    fi
fi

# 清理临时文件
echo ""
echo "清理临时文件..."
rm -rf "$TEMP_DIR"

echo ""
echo -e "${GREEN}✓ 测试完成！${NC}"
