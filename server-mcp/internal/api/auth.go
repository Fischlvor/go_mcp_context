package api

import (
	"bytes"
	"encoding/json"
	"fmt"
	"go-mcp-context/internal/model/response"
	"go-mcp-context/pkg/global"
	"io"
	"net/http"
	"net/url"
	"time"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type AuthApi struct{}

// TokenResponse SSO token 响应
type TokenResponse struct {
	AccessToken  string      `json:"access_token"`
	RefreshToken string      `json:"refresh_token"`
	TokenType    string      `json:"token_type"`
	ExpiresIn    int         `json:"expires_in"`
	UserInfo     interface{} `json:"user_info,omitempty"`
}

// GetSSOLoginURL 获取 SSO 授权地址
// GET /api/auth/sso_login_url?redirect_uri=xxx&return_url=xxx
func (a *AuthApi) GetSSOLoginURL(c *gin.Context) {
	// 获取回调地址
	redirectURI := c.Query("redirect_uri")
	if redirectURI == "" {
		redirectURI = global.Config.SSO.CallbackURL
	}

	// 获取 return_url 参数（用户想访问的页面）
	returnURL := c.Query("return_url")
	if returnURL == "" {
		returnURL = "/"
	}

	// 构建 state 参数
	state := fmt.Sprintf(`{"return_url":"%s"}`, returnURL)

	// 构建 SSO 授权 URL
	ssoLoginURL := fmt.Sprintf("%s/api/oauth/authorize?app_id=%s&redirect_uri=%s&state=%s",
		global.Config.SSO.WebURL,
		global.Config.SSO.ClientID,
		url.QueryEscape(redirectURI),
		url.QueryEscape(state),
	)

	response.OkWithData(gin.H{
		"sso_login_url": ssoLoginURL,
	}, c)
}

// SSOCallback SSO 回调接口（后端用 code 换 token）
// GET /api/auth/callback?code=xxx&redirect_uri=xxx&state=xxx
func (a *AuthApi) SSOCallback(c *gin.Context) {
	code := c.Query("code")
	state := c.Query("state")
	redirectURI := c.Query("redirect_uri")

	if code == "" {
		response.FailWithMessage("缺少授权码", c)
		return
	}

	// 用 code 向 SSO 换取 token
	tokenResp, err := exchangeCodeForToken(code, redirectURI)
	if err != nil {
		global.Log.Error("换取 token 失败", zap.Error(err))
		response.FailWithMessage("换取 token 失败: "+err.Error(), c)
		return
	}

	// refresh_token 存入后端 Session（不返回给前端）
	session := sessions.Default(c)
	session.Set("refresh_token", tokenResp.RefreshToken)
	session.Set("refresh_token_expires_at", time.Now().Add(7*24*time.Hour).Unix())
	if err := session.Save(); err != nil {
		global.Log.Error("保存 session 失败", zap.Error(err))
		response.FailWithMessage("保存会话失败", c)
		return
	}

	// 返回 access_token 给前端
	response.OkWithData(gin.H{
		"access_token": tokenResp.AccessToken,
		"token_type":   tokenResp.TokenType,
		"expires_in":   tokenResp.ExpiresIn,
		"state":        state,
	}, c)
}

// Logout 登出
// POST /api/auth/logout
func (a *AuthApi) Logout(c *gin.Context) {
	// 清除 Session
	session := sessions.Default(c)
	session.Clear()
	if err := session.Save(); err != nil {
		global.Log.Error("清除 session 失败", zap.Error(err))
	}

	response.OkWithMessage("登出成功", c)
}

// RefreshAccessTokenFromSSO 后端用 refresh_token 向 SSO 刷新（供中间件调用）
func RefreshAccessTokenFromSSO(refreshToken string) (*TokenResponse, error) {
	return refreshAccessToken(refreshToken)
}

// exchangeCodeForToken 用 code 向 SSO 换取 token
func exchangeCodeForToken(code, redirectURI string) (*TokenResponse, error) {
	ssoURL := fmt.Sprintf("%s/api/auth/token", global.Config.SSO.ServiceURL)

	requestBody := map[string]string{
		"grant_type":    "authorization_code",
		"code":          code,
		"client_id":     global.Config.SSO.ClientID,
		"client_secret": global.Config.SSO.ClientSecret,
		"redirect_uri":  redirectURI,
	}

	jsonData, err := json.Marshal(requestBody)
	if err != nil {
		return nil, err
	}

	resp, err := http.Post(ssoURL, "application/json", bytes.NewBuffer(jsonData))
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
		Code    int            `json:"code"`
		Message string         `json:"message"`
		Data    *TokenResponse `json:"data"`
	}

	if err := json.Unmarshal(body, &result); err != nil {
		return nil, err
	}

	if result.Code != 0 {
		return nil, fmt.Errorf("SSO 错误: %s", result.Message)
	}

	return result.Data, nil
}

// refreshAccessToken 用 refresh_token 刷新 access_token
func refreshAccessToken(refreshToken string) (*TokenResponse, error) {
	ssoURL := fmt.Sprintf("%s/api/auth/token", global.Config.SSO.ServiceURL)

	requestBody := map[string]string{
		"grant_type":    "refresh_token",
		"refresh_token": refreshToken,
		"client_id":     global.Config.SSO.ClientID,
		"client_secret": global.Config.SSO.ClientSecret,
	}

	jsonData, err := json.Marshal(requestBody)
	if err != nil {
		return nil, err
	}

	resp, err := http.Post(ssoURL, "application/json", bytes.NewBuffer(jsonData))
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
		Code    int            `json:"code"`
		Message string         `json:"message"`
		Data    *TokenResponse `json:"data"`
	}

	if err := json.Unmarshal(body, &result); err != nil {
		return nil, err
	}

	if result.Code != 0 {
		return nil, fmt.Errorf("SSO 错误: %s", result.Message)
	}

	return result.Data, nil
}
