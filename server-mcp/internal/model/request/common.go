package request

// PageInfo 分页参数
type PageInfo struct {
	Page     int `json:"page" form:"page"`           // 页码
	PageSize int `json:"page_size" form:"page_size"` // 每页大小
}

// IDRequest ID 请求参数
type IDRequest struct {
	ID uint `json:"id" form:"id" uri:"id" binding:"required"`
}
