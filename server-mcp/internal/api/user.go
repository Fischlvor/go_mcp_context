package api

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"go-mcp-context/internal/model/response"
	"go-mcp-context/pkg/global"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type UserApi struct{}

// GetUserInfo 获取当前用户信息（从 SSO 获取）
// GET /api/v1/user/info
func (a *UserApi) GetUserInfo(c *gin.Context) {
	// 从请求头获取 access_token
	authHeader := c.GetHeader("Authorization")
	if authHeader == "" || len(authHeader) < 8 {
		response.NoAuth("未提供认证 token", c)
		return
	}
	accessToken := authHeader[7:] // 去掉 "Bearer " 前缀

	// 调用 SSO 获取用户信息
	userInfo, err := getUserInfoFromSSO(accessToken)
	if err != nil {
		global.Log.Error("获取用户信息失败", zap.Error(err))
		response.FailWithMessage("获取用户信息失败: "+err.Error(), c)
		return
	}

	response.OkWithData(userInfo, c)
}

// getUserInfoFromSSO 从 SSO 获取用户信息
func getUserInfoFromSSO(accessToken string) (map[string]interface{}, error) {
	ssoURL := fmt.Sprintf("%s/api/user/info", global.Config.SSO.ServiceURL)

	req, err := http.NewRequest("GET", ssoURL, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Authorization", "Bearer "+accessToken)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode != 200 {
		return nil, fmt.Errorf("SSO 返回错误: %s", string(body))
	}

	var result struct {
		Code    int                    `json:"code"`
		Message string                 `json:"message"`
		Data    map[string]interface{} `json:"data"`
	}

	if err := json.Unmarshal(body, &result); err != nil {
		return nil, err
	}

	if result.Code != 0 {
		return nil, fmt.Errorf("SSO 错误: %s", result.Message)
	}

	return result.Data, nil
}
