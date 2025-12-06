package response

// PageResult 分页结果
type PageResult struct {
	List     interface{} `json:"list"`
	Total    int64       `json:"total"`
	Page     int         `json:"page"`
	PageSize int         `json:"page_size"`
}

// ProcessStatus 文档处理状态（用于 SSE 推送）
type ProcessStatus struct {
	Stage    string `json:"stage"`    // uploaded, parsing, chunking, embedding, saving, completed, failed
	Progress int    `json:"progress"` // 0-100
	Message  string `json:"message"`
	Status   string `json:"status"` // processing, active, failed
}
