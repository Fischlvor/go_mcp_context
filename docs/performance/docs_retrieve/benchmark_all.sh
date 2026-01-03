#!/bin/bash

# ============================================
# 批量性能测试脚本
# 测试所有库的所有版本
# ============================================

set -e

if [ $# -lt 1 ]; then
    echo "用法: $0 <MCP_API_KEY> [MODE]"
    echo ""
    echo "参数:"
    echo "  MODE: remote (默认) 或 local (服务器本地测试)"
    echo ""
    echo "示例:"
    echo "  $0 mcp_xxxxxxxx"
    echo "  $0 mcp_xxxxxxxx local"
    exit 1
fi

API_KEY="$1"
MODE="${2:-remote}"

# 获取脚本所在目录
SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"

# 创建结果目录
RESULTS_DIR="$SCRIPT_DIR/results_$(date +%Y%m%d_%H%M%S)"
mkdir -p "$RESULTS_DIR"

echo "========================================"
echo "批量性能测试"
echo "========================================"
echo ""
echo "结果目录: $RESULTS_DIR"
echo "测试模式: $MODE"
echo ""

# 定义测试库和版本
declare -A LIBRARIES

# Gin - 9 个版本
LIBRARIES[6]="v1.11.0" # v1.10.1 v1.9.1 v1.8.2 v1.7.7 v1.6.3 v1.5.0 v1.4.0 v1.0.0"

# Chi - 6 个版本
LIBRARIES[9]="v5.2.3" # v4.1.2 v3.3.4 v2.1.0 v1.5.4 v0.9.0"

# 其他库 - latest 版本
LIBRARIES[10]="latest"  # go-mcp-context
LIBRARIES[12]="latest"  # Gin Examples
LIBRARIES[17]="latest"  # GORM
LIBRARIES[18]="latest"  # Vite
LIBRARIES[20]="latest"  # Echo
LIBRARIES[21]="latest"  # Tailwind CSS
LIBRARIES[22]="latest"  # Go MCP Context
LIBRARIES[23]="latest"  # JWT
LIBRARIES[24]="latest"  # UUID
LIBRARIES[25]="latest"  # Tiktoken
LIBRARIES[26]="latest"  # Gorm

# 库名映射
declare -A LIB_NAMES
LIB_NAMES[6]="Gin"
LIB_NAMES[9]="Chi"
LIB_NAMES[10]="go-mcp-context"
LIB_NAMES[12]="Gin Examples"
LIB_NAMES[17]="GORM"
LIB_NAMES[18]="Vite"
LIB_NAMES[20]="Echo"
LIB_NAMES[21]="Tailwind CSS"
LIB_NAMES[22]="Go MCP Context"
LIB_NAMES[23]="JWT"
LIB_NAMES[24]="UUID"
LIB_NAMES[25]="Tiktoken"
LIB_NAMES[26]="Gorm"

# 统计信息
TOTAL_TESTS=0
COMPLETED_TESTS=0
FAILED_TESTS=0

# 计算总测试数
for lib_id in "${!LIBRARIES[@]}"; do
    versions="${LIBRARIES[$lib_id]}"
    for version in $versions; do
        TOTAL_TESTS=$((TOTAL_TESTS + 1))
    done
done

echo "总共需要测试: $TOTAL_TESTS 个库版本"
echo ""

# 创建汇总文件
SUMMARY_FILE="$RESULTS_DIR/summary.txt"
echo "========================================" > "$SUMMARY_FILE"
echo "批量性能测试汇总" >> "$SUMMARY_FILE"
echo "========================================" >> "$SUMMARY_FILE"
echo "" >> "$SUMMARY_FILE"
echo "测试时间: $(date)" >> "$SUMMARY_FILE"
echo "测试模式: $MODE" >> "$SUMMARY_FILE"
echo "总测试数: $TOTAL_TESTS" >> "$SUMMARY_FILE"
echo "" >> "$SUMMARY_FILE"

# 遍历所有库和版本
for lib_id in $(echo "${!LIBRARIES[@]}" | tr ' ' '\n' | sort -n); do
    lib_name="${LIB_NAMES[$lib_id]}"
    versions="${LIBRARIES[$lib_id]}"
    
    echo "----------------------------------------"
    echo "测试库: $lib_name (ID: $lib_id)"
    echo "----------------------------------------"
    
    for version in $versions; do
        COMPLETED_TESTS=$((COMPLETED_TESTS + 1))
        
        echo ""
        echo "[$COMPLETED_TESTS/$TOTAL_TESTS] 测试: $lib_name $version"
        echo ""
        
        # 运行测试（设置 BATCH_MODE 环境变量）
        RESULT_FILE="$RESULTS_DIR/${lib_id}_${version}.txt"
        
        if BATCH_MODE=1 bash "$SCRIPT_DIR/benchmark.sh" "$API_KEY" "$lib_id" "$version" "$MODE" > "$RESULT_FILE" 2>&1; then
            echo "✓ 测试成功"
            
            # 提取关键指标（修复正则表达式）
            COLD_AVG=$(grep "冷启动平均延迟:" "$RESULT_FILE" | grep -oP '\d+(?=ms)' | head -1)
            HOT_QPS=$(grep "热启动 QPS:" "$RESULT_FILE" | grep -oP '\d+\.\d+' | head -1)
            HOT_P50=$(grep "热启动 P50:" "$RESULT_FILE" | grep -oP '\d+(?=ms)' | head -1)
            MIXED_AVG=$(grep "混合场景平均延迟:" "$RESULT_FILE" | grep -oP '\d+(?=ms)' | head -1)
            
            # 写入汇总
            echo "[$lib_name $version]" >> "$SUMMARY_FILE"
            echo "  冷启动: ${COLD_AVG:-N/A}ms" >> "$SUMMARY_FILE"
            echo "  热启动 QPS: ${HOT_QPS:-N/A}" >> "$SUMMARY_FILE"
            echo "  热启动 P50: ${HOT_P50:-N/A}ms" >> "$SUMMARY_FILE"
            echo "  混合场景: ${MIXED_AVG:-N/A}ms" >> "$SUMMARY_FILE"
            echo "" >> "$SUMMARY_FILE"
        else
            echo "✗ 测试失败"
            FAILED_TESTS=$((FAILED_TESTS + 1))
            
            echo "[$lib_name $version] - 失败" >> "$SUMMARY_FILE"
            echo "" >> "$SUMMARY_FILE"
        fi
        
        # 短暂延迟，避免过载
        sleep 2
    done
    
    echo ""
done

# 最终统计
echo "========================================"
echo "测试完成"
echo "========================================"
echo ""
echo "总测试数: $TOTAL_TESTS"
echo "成功: $((TOTAL_TESTS - FAILED_TESTS))"
echo "失败: $FAILED_TESTS"
echo ""
echo "结果目录: $RESULTS_DIR"
echo "汇总文件: $SUMMARY_FILE"
echo ""

# 写入最终统计
echo "========================================" >> "$SUMMARY_FILE"
echo "测试统计" >> "$SUMMARY_FILE"
echo "========================================" >> "$SUMMARY_FILE"
echo "" >> "$SUMMARY_FILE"
echo "总测试数: $TOTAL_TESTS" >> "$SUMMARY_FILE"
echo "成功: $((TOTAL_TESTS - FAILED_TESTS))" >> "$SUMMARY_FILE"
echo "失败: $FAILED_TESTS" >> "$SUMMARY_FILE"
echo "" >> "$SUMMARY_FILE"

# 显示汇总
cat "$SUMMARY_FILE"
