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
// @Summary 创建 API Key
// @Description 为当前用户创建新的 API Key，用于 MCP 协议调用（需要认证）
// @Tags API Keys
// @Accept json
// @Produce json
// @Security JWTAuth
// @Param data body request.APIKeyCreate true "API Key 信息"
// @Success 200 {object} response.Response{data=response.APIKeyCreateResponse}
// @Failure 400 {object} response.Response
// @Failure 401 {object} response.Response
// @Router /api/v1/api-keys/create [post]
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
// @Summary 获取 API Key 列表
// @Description 获取当前用户的所有 API Key（需要认证）
// @Tags API Keys
// @Accept json
// @Produce json
// @Security JWTAuth
// @Success 200 {object} response.Response{data=[]response.APIKeyListItem}
// @Failure 401 {object} response.Response
// @Router /api/v1/api-keys/list [get]
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
// @Summary 删除 API Key
// @Description 删除指定的 API Key（需要认证）
// @Tags API Keys
// @Accept json
// @Produce json
// @Security JWTAuth
// @Param id path int true "API Key ID"
// @Success 200 {object} response.Response{data=nil}
// @Failure 401 {object} response.Response
// @Failure 404 {object} response.Response
// @Router /api/v1/api-keys/:id [delete]
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
