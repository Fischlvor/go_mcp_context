package request

// APIKeyCreate 创建 API Key 请求
type APIKeyCreate struct {
	Name string `json:"name" binding:"required,max=100"`
}
