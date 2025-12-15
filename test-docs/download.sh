#!/bin/bash

# 测试文档下载脚本
# 使用 GitHub API + 多线程下载
# 只下载 docs/ 或 src/ 目录下的文档
# 代理: 10.21.71.52:7890

# 代理设置
export http_proxy="http://10.21.71.52:7890"
export https_proxy="http://10.21.71.52:7890"

BASE_DIR="$(cd "$(dirname "$0")" && pwd)"
PARALLEL_JOBS=10  # 并行下载数

echo "下载目录: $BASE_DIR"
echo "并行数: $PARALLEL_JOBS"

# 创建目录
create_dir() {
    mkdir -p "$BASE_DIR/$1/$2"
}

# 下载单个文件（供 xargs 调用）
download_single() {
    local raw_url=$1
    local dest=$2
    curl -sL "$raw_url" -o "$dest" 2>/dev/null
}
export -f download_single

# 使用 GitHub API + 多线程下载文档
# 参数: repo, branch, lib, version, path_filter (只下载包含此路径的文件)
download_repo_docs() {
    local repo=$1
    local branch=$2
    local lib=$3
    local version=$4
    local path_filter=${5:-""}  # 如 "docs" "src/content" 等
    
    echo ""
    echo "=========================================="
    echo "📚 $lib/$version"
    echo "   仓库: $repo"
    echo "   分支: $branch"
    [ -n "$path_filter" ] && echo "   路径: $path_filter"
    echo "=========================================="
    
    create_dir "$lib" "$version"
    
    # 使用 GitHub Tree API 获取完整目录树
    local tree_url="https://api.github.com/repos/$repo/git/trees/$branch?recursive=1"
    echo "  获取目录树..."
    
    local tree_json=$(curl -sL "$tree_url")
    
    # 检查是否成功
    if echo "$tree_json" | grep -q '"message"'; then
        echo "  [错误] 无法获取目录树"
        echo "$tree_json" | head -5
        return
    fi
    
    # 提取 md/mdx 文件，过滤路径
    local files
    if [ -n "$path_filter" ]; then
        files=$(echo "$tree_json" | \
            grep -o '"path": "[^"]*"' | \
            cut -d'"' -f4 | \
            grep -E '\.(md|mdx)$' | \
            grep "^$path_filter")
    else
        files=$(echo "$tree_json" | \
            grep -o '"path": "[^"]*"' | \
            cut -d'"' -f4 | \
            grep -E '\.(md|mdx)$')
    fi
    
    # 统计文件数
    local total=$(echo "$files" | grep -c .)
    echo "  找到 $total 个文档文件"
    
    # 生成下载任务列表
    local task_file=$(mktemp)
    for file_path in $files; do
        local safe_name=$(echo "$file_path" | tr '/' '_')
        local raw_url="https://raw.githubusercontent.com/$repo/$branch/$file_path"
        echo "$raw_url $BASE_DIR/$lib/$version/$safe_name" >> "$task_file"
    done
    
    # 多线程下载
    echo "  开始多线程下载..."
    cat "$task_file" | xargs -P $PARALLEL_JOBS -L 1 bash -c 'curl -sL "$0" -o "$1" 2>/dev/null && echo -n "."'
    echo ""
    
    rm -f "$task_file"
    
    local downloaded=$(find "$BASE_DIR/$lib/$version" -type f | wc -l)
    echo "  ✓ 完成: $downloaded 个文件"
}

# 直接下载 README 等单个文件
download_readme() {
    local repo=$1
    local branch=$2
    local lib=$3
    local version=$4
    
    echo ""
    echo "=========================================="
    echo "📚 $lib/$version (README)"
    echo "=========================================="
    
    create_dir "$lib" "$version"
    
    local url="https://raw.githubusercontent.com/$repo/$branch/README.md"
    echo "  下载: README.md"
    curl -sL "$url" -o "$BASE_DIR/$lib/$version/README.md" 2>/dev/null
    
    # 尝试下载 CHANGELOG
    url="https://raw.githubusercontent.com/$repo/$branch/CHANGELOG.md"
    curl -sL "$url" -o "$BASE_DIR/$lib/$version/CHANGELOG.md" 2>/dev/null && echo "  下载: CHANGELOG.md"
    
    echo "  ✓ 完成"
}

echo ""
echo "============================================"
echo "  开始下载测试文档"
echo "  代理: $http_proxy"
echo "============================================"

# ==========================================
# Vue.js 文档 - 只下载 src/ 目录
# ==========================================

download_repo_docs "vuejs/docs" "main" "vue" "3.4" "src/"

# ==========================================
# React 文档 - 只下载 src/content/ 目录
# ==========================================

download_repo_docs "reactjs/react.dev" "main" "react" "18" "src/content/"

# ==========================================
# Next.js 文档 - 只下载 docs/ 目录
# ==========================================

download_repo_docs "vercel/next.js" "canary" "nextjs" "14" "docs/"

# ==========================================
# Tailwind CSS
# ==========================================

download_readme "tailwindlabs/tailwindcss" "master" "tailwindcss" "3.4"

# ==========================================
# TypeScript
# ==========================================

download_readme "microsoft/TypeScript" "main" "typescript" "5.0"

# ==========================================
# Gin (Go Web Framework)
# ==========================================

download_readme "gin-gonic/gin" "master" "gin" "1.9"

# ==========================================
# GORM
# ==========================================

# GORM 主仓库 README
download_readme "go-gorm/gorm" "master" "gorm" "2.0"

# GORM 文档仓库 - 只下载 pages/docs/ 目录
download_repo_docs "go-gorm/gorm.io" "master" "gorm" "2.0-docs" "pages/docs/"

# ==========================================
# Echo (Go Web Framework)
# ==========================================

download_readme "labstack/echo" "master" "echo" "4.0"

# ==========================================
# Vite - 只下载 docs/ 目录
# ==========================================

download_repo_docs "vitejs/vite" "main" "vite" "5.0" "docs/"

echo ""
echo "============================================"
echo "  下载完成！"
echo "============================================"
echo ""
echo "目录结构："
find "$BASE_DIR" -type d | grep -v "^\.$" | sort
echo ""
echo "文件统计："
for lib in "$BASE_DIR"/*/; do
    if [ -d "$lib" ]; then
        lib_name=$(basename "$lib")
        for ver in "$lib"*/; do
            if [ -d "$ver" ]; then
                ver_name=$(basename "$ver")
                count=$(find "$ver" -maxdepth 1 -type f \( -name "*.md" -o -name "*.mdx" \) 2>/dev/null | wc -l)
                if [ "$count" -gt 0 ]; then
                    echo "  $lib_name/$ver_name: $count 个文件"
                fi
            fi
        done
    fi
done

echo ""
echo "总文件数："
find "$BASE_DIR" -type f \( -name "*.md" -o -name "*.mdx" \) | wc -l
