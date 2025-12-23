package github

// ==========================================
// GitHub API 响应数据结构
// ==========================================

// Repo 仓库信息
type Repo struct {
	Name          string `json:"name"`
	FullName      string `json:"full_name"`
	Description   string `json:"description"`
	DefaultBranch string `json:"default_branch"`
	Size          int    `json:"size"` // 仓库大小（KB）
}

// Release 版本信息
type Release struct {
	TagName    string `json:"tag_name"`
	Name       string `json:"name"`
	Prerelease bool   `json:"prerelease"`
	Draft      bool   `json:"draft"`
}

// TreeItem 目录树项
type TreeItem struct {
	Path string `json:"path"`
	Mode string `json:"mode"`
	Type string `json:"type"` // blob, tree
	Size int    `json:"size"`
	SHA  string `json:"sha"`
	URL  string `json:"url"`
}

// Tree 目录树响应
type Tree struct {
	SHA       string     `json:"sha"`
	URL       string     `json:"url"`
	Tree      []TreeItem `json:"tree"`
	Truncated bool       `json:"truncated"`
}
