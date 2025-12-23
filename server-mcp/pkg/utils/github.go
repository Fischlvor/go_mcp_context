package utils

import (
	"errors"
	"strings"
)

// ParseGitHubURL 解析 GitHub URL，返回 owner/repo
// 支持格式：
// - https://github.com/owner/repo
// - https://github.com/owner/repo.git
// - github.com/owner/repo
func ParseGitHubURL(urlStr string) (string, error) {
	urlStr = strings.TrimSpace(urlStr)
	urlStr = strings.TrimSuffix(urlStr, ".git")
	urlStr = strings.TrimSuffix(urlStr, "/")

	// 移除协议前缀
	urlStr = strings.TrimPrefix(urlStr, "https://")
	urlStr = strings.TrimPrefix(urlStr, "http://")

	// 检查是否是 github.com
	if !strings.HasPrefix(urlStr, "github.com/") {
		return "", errors.New("仅支持 github.com 仓库")
	}

	// 提取 owner/repo
	parts := strings.Split(strings.TrimPrefix(urlStr, "github.com/"), "/")
	if len(parts) < 2 || parts[0] == "" || parts[1] == "" {
		return "", errors.New("URL 格式错误，应为 github.com/owner/repo")
	}

	return parts[0] + "/" + parts[1], nil
}

// ExtractRepoName 从 owner/repo 中提取 repo 名称
func ExtractRepoName(repo string) string {
	parts := strings.Split(repo, "/")
	if len(parts) >= 2 {
		return parts[1]
	}
	return repo
}
