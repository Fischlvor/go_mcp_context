package request

// GitHubImportRequest GitHub 导入请求
type GitHubImportRequest struct {
	Repo       string   `json:"repo" binding:"required"` // owner/repo
	Branch     string   `json:"branch"`                  // 分支名（与 Tag 二选一）
	Tag        string   `json:"tag"`                     // 特定 tag
	Version    string   `json:"version"`                 // 存储为的版本名
	PathFilter string   `json:"path_filter"`             // 只导入指定路径（如 docs/）
	Excludes   []string `json:"excludes"`                // 排除模式
	TaskID     string   `json:"-"`                       // 任务ID（内部使用，不从 JSON 解析）
}

// GitHubReleasesQuery 获取 GitHub 版本列表请求
type GitHubReleasesQuery struct {
	Repo     string `form:"repo" binding:"required"` // owner/repo
	MaxCount int    `form:"max_count"`               // 最多返回几个版本，默认 20
}

// GitHubInitImportRequest 从 GitHub 初始化导入请求（创建库+导入）
type GitHubInitImportRequest struct {
	GitHubURL string `json:"github_url" binding:"required"` // 完整 GitHub URL，如 https://github.com/gin-gonic/gin
}
