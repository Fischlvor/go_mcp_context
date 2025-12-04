package api

import (
	"strconv"

	"go-mcp-context/internal/model/request"
	"go-mcp-context/internal/model/response"

	"github.com/gin-gonic/gin"
	"github.com/gofrs/uuid"
)

type ApiKeyApi struct{}

// getUserUUIDFromContext 从上下文获取用户 UUID（避免循环导入）
func getUserUUIDFromContext(c *gin.Context) uuid.UUID {
	if val, exists := c.Get("user_uuid"); exists {
		if userUUID, ok := val.(uuid.UUID); ok {
			return userUUID
		}
	}
	return uuid.Nil
}

// Create 创建 API Key
// POST /api/v1/api-keys/create
func (a *ApiKeyApi) Create(c *gin.Context) {
	var req request.APIKeyCreate
	if err := c.ShouldBindJSON(&req); err != nil {
		response.FailWithMessage("参数错误: "+err.Error(), c)
		return
	}

	userUUID := getUserUUIDFromContext(c)
	if userUUID.IsNil() {
		response.NoAuth("未登录", c)
		return
	}

	result, err := apiKeyService.Create(userUUID.String(), &req)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}

	response.OkWithData(result, c)
}

// List 获取 API Key 列表
// GET /api/v1/api-keys/list
func (a *ApiKeyApi) List(c *gin.Context) {
	userUUID := getUserUUIDFromContext(c)
	if userUUID.IsNil() {
		response.NoAuth("未登录", c)
		return
	}

	list, err := apiKeyService.List(userUUID.String())
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}

	response.OkWithData(list, c)
}

// Delete 删除 API Key
// DELETE /api/v1/api-keys/:id
func (a *ApiKeyApi) Delete(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		response.FailWithMessage("无效的 ID", c)
		return
	}

	userUUID := getUserUUIDFromContext(c)
	if userUUID.IsNil() {
		response.NoAuth("未登录", c)
		return
	}

	if err := apiKeyService.Delete(userUUID.String(), uint(id)); err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}

	response.OkWithMessage("删除成功", c)
}
