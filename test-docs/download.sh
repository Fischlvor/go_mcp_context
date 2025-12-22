#!/bin/bash

# 测试文档下载脚本
# 使用 GitHub API + 多线程下载
# 支持指定 tag 版本，保留原始目录结构
# 代理: 10.21.71.52:7890

# 代理设置
export http_proxy="http://10.21.71.52:7890"
export https_proxy="http://10.21.71.52:7890"

# GitHub Token (提高 API 速率限制: 60/h -> 5000/h)
GITHUB_TOKEN="***"
GITHUB_AUTH_HEADER="Authorization: token $GITHUB_TOKEN"

BASE_DIR="$(cd "$(dirname "$0")" && pwd)"
PARALLEL_JOBS=10  # 并行下载数

echo "下载目录: $BASE_DIR"
echo "并行数: $PARALLEL_JOBS"

# 创建目录
create_dir() {
    mkdir -p "$BASE_DIR/$1/$2"
}

# 下载单个文件并创建目录（供 xargs 调用）
download_with_mkdir() {
    local raw_url=$1
    local dest=$2
    mkdir -p "$(dirname "$dest")"
    curl -sL "$raw_url" -o "$dest" 2>/dev/null
}
export -f download_with_mkdir

# ==========================================
# 获取仓库的所有 tags（版本号）
# 参数: repo, limit (可选，默认 30)
# 返回: 版本号列表（每行一个）
# 注意: 包含所有 tag，包括 alpha、beta、rc 等预发布版本
# ==========================================
get_repo_tags() {
    local repo=$1
    local limit=${2:-30}
    
    # GitHub Tags API，支持分页
    local tags_url="https://api.github.com/repos/$repo/tags?per_page=$limit"
    local tags_json=$(curl -sL -H "$GITHUB_AUTH_HEADER" "$tags_url")
    
    # 提取 tag 名称
    echo "$tags_json" | grep -o '"name": "[^"]*"' | cut -d'"' -f4
}

# ==========================================
# 获取仓库的正式发布版本（并行分页，速度更快）
# 参数: repo
# 返回: 版本号列表（每行一个，只包含正式版）
# ==========================================
get_repo_releases() {
    local repo=$1
    local base_url="https://api.github.com/repos/$repo/releases?per_page=100"
    
    # 先请求一次获取总页数（从 Link 响应头）
    local header=$(curl -sI -H "$GITHUB_AUTH_HEADER" "$base_url")
    local last_page=$(echo "$header" | grep -i '^link:' | grep -oE 'page=[0-9]+>; rel="last"' | grep -oE '[0-9]+')
    
    # 如果没有 Link 头，说明只有 1 页
    if [ -z "$last_page" ]; then
        last_page=1
    fi
    
    # 并行请求所有页面
    local all_versions=""
    if [ "$last_page" -gt 1 ]; then
        # 使用 xargs 并行请求（最多 5 个并发）
        all_versions=$(seq 1 $last_page | xargs -P 5 -I {} curl -sL -H "$GITHUB_AUTH_HEADER" "${base_url}&page={}" | \
            grep -o '"tag_name": "[^"]*"' | \
            cut -d'"' -f4 | \
            grep -E '^v?[0-9]+\.[0-9]+(\.[0-9]+)?$')
    else
        # 只有 1 页，直接请求
        all_versions=$(curl -sL -H "$GITHUB_AUTH_HEADER" "$base_url" | \
            grep -o '"tag_name": "[^"]*"' | \
            cut -d'"' -f4 | \
            grep -E '^v?[0-9]+\.[0-9]+(\.[0-9]+)?$')
    fi
    
    echo "$all_versions" | grep -v '^$' | sort -V -r | uniq
}

# ==========================================
# 下载每个大版本的最新版本（Context7 策略）
# 如果只有一个大版本，则下载每个次版本(minor)的最新
# 参数: repo, lib, path_filter, max_count (最多下载几个版本)
# ==========================================
download_major_versions() {
    local repo=$1
    local lib=$2
    local path_filter=$3
    local max_count=${4:-10}
    
    echo ""
    echo "=========================================="
    echo "🔍 获取 $repo 每个大版本的最新版本..."
    echo "=========================================="
    
    # 获取正式发布版本
    local versions=$(get_repo_releases "$repo" 100)
    
    if [ -z "$versions" ]; then
        echo "  [警告] 无 releases，尝试使用 tags..."
        versions=$(get_repo_tags "$repo" 100 | grep -vE '(alpha|beta|rc|dev|pre|canary|nightly)')
    fi
    
    if [ -z "$versions" ]; then
        echo "  [错误] 无法获取版本列表"
        return
    fi
    
    # 先统计有多少个大版本
    local unique_majors=$(echo "$versions" | sed -E 's/^v?([0-9]+)\..*/\1/' | sort -u | wc -l)
    
    local selected_versions=""
    local seen_keys=""
    
    if [ "$unique_majors" -eq 1 ]; then
        # 只有一个大版本，按 minor 版本分组
        echo "  [策略] 只有 1 个大版本，改为按次版本(minor)分组"
        
        for ver in $versions; do
            # 提取 major.minor (v1.10.x -> 1.10)
            local minor_key=$(echo "$ver" | sed -E 's/^v?([0-9]+\.[0-9]+)\..*/\1/')
            
            if ! echo "$seen_keys" | grep -q "^${minor_key}$"; then
                seen_keys="$seen_keys
$minor_key"
                selected_versions="$selected_versions $ver"
                
                local count=$(echo "$selected_versions" | wc -w)
                if [ "$count" -ge "$max_count" ]; then
                    break
                fi
            fi
        done
        
        echo "  找到次版本 (每个取最新):"
    else
        # 多个大版本，按 major 版本分组
        for ver in $versions; do
            # 提取大版本号 (v1.x.x -> 1)
            local major=$(echo "$ver" | sed -E 's/^v?([0-9]+)\..*/\1/')
            
            if ! echo "$seen_keys" | grep -q "^${major}$"; then
                seen_keys="$seen_keys
$major"
                selected_versions="$selected_versions $ver"
                
                local count=$(echo "$selected_versions" | wc -w)
                if [ "$count" -ge "$max_count" ]; then
                    break
                fi
            fi
        done
        
        echo "  找到 $unique_majors 个大版本 (每个取最新):"
    fi
    
    for ver in $selected_versions; do
        echo "    $ver"
    done
    
    # 逐个下载
    for version in $selected_versions; do
        download_repo_docs "$repo" "$version" "$lib" "$version" "$path_filter"
    done
}

# ==========================================
# 下载无 tag 的文档仓库（使用分支）
# 参数: repo, branch, lib, version, path_filter
# version 建议: "latest", "v3", "v19" 等有意义的名称
# ==========================================
download_docs_repo() {
    local repo=$1
    local branch=${2:-"main"}
    local lib=$3
    local version=${4:-"latest"}
    local path_filter=${5:-""}
    
    echo ""
    echo "=========================================="
    echo "📚 $lib/$version (from $branch branch)"
    echo "   仓库: $repo"
    echo "   分支: $branch"
    [ -n "$path_filter" ] && echo "   路径: $path_filter"
    echo "=========================================="
    
    download_repo_docs "$repo" "$branch" "$lib" "$version" "$path_filter"
}

# 使用 GitHub API + 多线程下载文档
# 参数: repo, ref (branch 或 tag), lib, version, path_filter (只下载包含此路径的文件)
# ref 可以是: main, master, v1.10.0, v2.0.0 等
download_repo_docs() {
    local repo=$1
    local ref=$2          # 分支名或 tag 名（如 v1.10.0）
    local lib=$3
    local version=$4
    local path_filter=${5:-""}  # 如 "docs" "src/content" 等
    
    echo ""
    echo "=========================================="
    echo "📚 $lib/$version"
    echo "   仓库: $repo"
    echo "   Ref:  $ref"
    [ -n "$path_filter" ] && echo "   路径: $path_filter"
    echo "=========================================="
    
    create_dir "$lib" "$version"
    
    # 使用 GitHub Tree API 获取完整目录树
    # ref 可以是分支名或 tag 名，API 都支持
    local tree_url="https://api.github.com/repos/$repo/git/trees/$ref?recursive=1"
    echo "  获取目录树..."
    
    local tree_json=$(curl -sL -H "$GITHUB_AUTH_HEADER" "$tree_url")
    
    # 检查是否成功
    if echo "$tree_json" | grep -q '"message"'; then
        echo "  [错误] 无法获取目录树，尝试获取 tag SHA..."
        
        # 如果直接用 tag 失败，尝试先获取 tag 对应的 commit SHA
        local tag_url="https://api.github.com/repos/$repo/git/ref/tags/$ref"
        local tag_json=$(curl -sL -H "$GITHUB_AUTH_HEADER" "$tag_url")
        local sha=$(echo "$tag_json" | grep -o '"sha": "[^"]*"' | head -1 | cut -d'"' -f4)
        
        if [ -n "$sha" ]; then
            echo "  找到 SHA: $sha"
            tree_url="https://api.github.com/repos/$repo/git/trees/$sha?recursive=1"
            tree_json=$(curl -sL -H "$GITHUB_AUTH_HEADER" "$tree_url")
        else
            echo "  [错误] 无法获取 tag 信息"
            echo "$tag_json" | head -5
            return
        fi
    fi
    
    # 提取 md/mdx 文件，过滤路径
    # 参考 Context7 的排除规则：
    # 排除目录: .github, test(s), dist, node_modules, vendor, fixtures, bench
    #          archive/archived/deprecated/legacy/old/outdated
    #          i18n 非英语目录, zh-cn/zh-tw 等
    # 排除文件: CHANGELOG, LICENSE, CODE_OF_CONDUCT
    # 注意：保留 examples（代码示例有用）、CONTRIBUTING（贡献指南可能有用）
    local files
    local exclude_dirs='(^\.github/|/\.github/|^test/|/test/|^tests/|/tests/|/__tests__/|^dist/|/dist/|/node_modules/|/vendor/|^vendor/|/fixtures/|^fixtures/|/bench/|^bench/|archive|deprecated|legacy|/old/|^old/|outdated|/i18n/|^i18n/|/zh-cn/|/zh-tw/|/zh-hk/)'
    local exclude_files='^(CHANGELOG|LICENSE|CODE_OF_CONDUCT)\.(md|mdx)$'
    
    if [ -n "$path_filter" ]; then
        files=$(echo "$tree_json" | \
            grep -o '"path": "[^"]*"' | \
            cut -d'"' -f4 | \
            grep -E '\.(md|mdx)$' | \
            grep "^$path_filter" | \
            grep -vE "$exclude_dirs" | \
            grep -viE "$exclude_files")
    else
        files=$(echo "$tree_json" | \
            grep -o '"path": "[^"]*"' | \
            cut -d'"' -f4 | \
            grep -E '\.(md|mdx)$' | \
            grep -vE "$exclude_dirs" | \
            grep -viE "$exclude_files")
    fi
    
    # 统计文件数
    local total=$(echo "$files" | grep -c .)
    echo "  找到 $total 个文档文件"
    
    # 生成下载任务列表（保留目录结构）
    local task_file=$(mktemp)
    for file_path in $files; do
        # 保留原始目录结构
        local raw_url="https://raw.githubusercontent.com/$repo/$ref/$file_path"
        local dest_path="$BASE_DIR/$lib/$version/$file_path"
        echo "$raw_url $dest_path" >> "$task_file"
    done
    
    # 多线程下载（会自动创建子目录）
    echo "  开始多线程下载..."
    cat "$task_file" | xargs -P $PARALLEL_JOBS -L 1 bash -c 'download_with_mkdir "$0" "$1" && echo -n "."'
    echo ""
    
    rm -f "$task_file"
    
    local downloaded=$(find "$BASE_DIR/$lib/$version" -type f | wc -l)
    echo "  ✓ 完成: $downloaded 个文件"
}

echo ""
echo "============================================"
echo "  开始下载测试文档"
echo "  代理: $http_proxy"
echo "============================================"

# ==========================================
# Gin (Go Web Framework) - 每个大版本取最新
# ==========================================

download_major_versions "gin-gonic/gin" "gin" "" 20

# Gin 示例仓库（无 tag）
download_docs_repo "gin-gonic/examples" "master" "gin-examples" "latest" ""

# ==========================================
# GORM 文档（文档在 gorm.io 仓库，无 tag）
# ==========================================

download_docs_repo "go-gorm/gorm.io" "master" "gorm" "latest" "pages/docs/"

# ==========================================
# Echo (Go Web Framework) - 每个大版本取最新
# ==========================================

download_major_versions "labstack/echo" "echo" "" 20

# ==========================================
# Next.js 文档 - 每个大版本取最新
# ==========================================

download_major_versions "vercel/next.js" "nextjs" "docs/" 20

# ==========================================
# Vite - 每个大版本取最新
# ==========================================

download_major_versions "vitejs/vite" "vite" "docs/" 20

# ==========================================
# Vue.js 文档（无 tag，使用 main）
# ==========================================

download_docs_repo "vuejs/docs" "main" "vue" "latest" "src/"

# ==========================================
# React 文档（无 tag，使用 main）
# ==========================================

download_docs_repo "reactjs/react.dev" "main" "react" "latest" "src/content/"

# ==========================================
# Tailwind CSS - 每个大版本取最新
# ==========================================

download_major_versions "tailwindlabs/tailwindcss" "tailwindcss" "" 20

# ==========================================
# TypeScript - 每个大版本取最新
# ==========================================

download_major_versions "microsoft/TypeScript" "typescript" "" 20

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
