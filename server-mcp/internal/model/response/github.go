package response

// GitHubImportProgress GitHub 导入进度（SSE 推送）
type GitHubImportProgress struct {
	Stage    string `json:"stage"`    // fetching_tree, downloading, processing, completed, failed
	Current  int    `json:"current"`  // 当前进度
	Total    int    `json:"total"`    // 总数
	Message  string `json:"message"`  // 状态消息
	FileName string `json:"filename"` // 当前文件名
}

// GitHubRepoInfo GitHub 仓库信息
type GitHubRepoInfo struct {
	Repo          string   `json:"repo"`
	DefaultBranch string   `json:"default_branch"`
	Description   string   `json:"description"`
	Versions      []string `json:"versions"`
}

// GitHubInitImportResponse 从 GitHub 初始化导入响应
type GitHubInitImportResponse struct {
	LibraryID uint   `json:"library_id"`
	Version   string `json:"version"`
}
