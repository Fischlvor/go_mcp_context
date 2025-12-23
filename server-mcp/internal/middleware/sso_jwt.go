package middleware

import (
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"errors"
	"go-mcp-context/internal/api"
	"go-mcp-context/internal/model/response"
	"go-mcp-context/pkg/global"
	"os"
	"strconv"
	"strings"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"github.com/gofrs/uuid"
	"github.com/golang-jwt/jwt/v4"
	"go.uber.org/zap"
)

var (
	TokenExpired     = errors.New("token is expired")
	TokenNotValidYet = errors.New("token not active yet")
	TokenMalformed   = errors.New("that's not even a token")
	TokenInvalid     = errors.New("couldn't handle this token")

	// SSOPublicKey SSO 公钥
	SSOPublicKey *rsa.PublicKey
)

// SSOClaims SSO 颁发的 JWT Claims
type SSOClaims struct {
	UserUUID  uuid.UUID `json:"user_uuid"`
	AppID     string    `json:"app_id"`
	DeviceID  string    `json:"device_id,omitempty"`
	TokenType string    `json:"token_type,omitempty"`
	Email     string    `json:"email,omitempty"`
	Nickname  string    `json:"nickname,omitempty"`
	jwt.RegisteredClaims
}

// LoadSSOPublicKey 加载 SSO 公钥
func LoadSSOPublicKey(path string) error {
	keyData, err := os.ReadFile(path)
	if err != nil {
		return err
	}

	block, _ := pem.Decode(keyData)
	if block == nil {
		return errors.New("无法解析 PEM 格式的公钥")
	}

	publicKey, err := x509.ParsePKIXPublicKey(block.Bytes)
	if err != nil {
		return err
	}

	rsaPublicKey, ok := publicKey.(*rsa.PublicKey)
	if !ok {
		return errors.New("不是 RSA 公钥")
	}

	SSOPublicKey = rsaPublicKey
	global.Log.Info("✓ SSO 公钥加载成功", zap.String("path", path))
	return nil
}

// SSOJWTAuth SSO JWT 认证中间件
func SSOJWTAuth() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 获取 Authorization header
		authHeader := c.GetHeader("Authorization")

		if authHeader == "" {
			response.NoAuth("未提供认证 token", c)
			c.Abort()
			return
		}

		// 检查 Bearer 前缀
		parts := strings.SplitN(authHeader, " ", 2)
		if !(len(parts) == 2 && parts[0] == "Bearer") {
			response.NoAuth("token 格式错误", c)
			c.Abort()
			return
		}

		token := parts[1]

		// 解析 SSO 颁发的 AccessToken
		claims, err := ParseSSOAccessToken(token)
		if err != nil {
			if errors.Is(err, TokenExpired) {
				// Token 过期，尝试自动刷新
				newToken, refreshErr := autoRefreshToken(c)
				if refreshErr != nil {
					global.Log.Warn("Token 刷新失败", zap.Error(refreshErr))
					response.NoAuth("token 已过期且刷新失败，请重新登录", c)
					c.Abort()
					return
				}

				// 刷新成功，在响应头返回新 token
				c.Header("X-New-Access-Token", newToken.AccessToken)
				c.Header("X-Token-Expires-In", strconv.Itoa(newToken.ExpiresIn))

				// 替换请求头中的 Authorization，供后续 handler 使用
				c.Request.Header.Set("Authorization", "Bearer "+newToken.AccessToken)

				// 重新解析新 token
				claims, err = ParseSSOAccessToken(newToken.AccessToken)
				if err != nil {
					response.NoAuth("新 token 解析失败", c)
					c.Abort()
					return
				}

				global.Log.Info("✓ Token 自动刷新成功", zap.String("user_uuid", claims.UserUUID.String()))
			} else {
				response.NoAuth("token 无效: "+err.Error(), c)
				c.Abort()
				return
			}
		}

		// 检查应用 ID 是否匹配
		if claims.AppID != global.Config.SSO.ClientID {
			response.NoAuth("token 不适用于此应用", c)
			c.Abort()
			return
		}

		// 将用户 UUID 存入上下文
		c.Set("user_uuid", claims.UserUUID)

		c.Next()
	}
}

// ParseSSOAccessToken 解析 SSO 颁发的 AccessToken（使用 RSA 公钥验证）
func ParseSSOAccessToken(tokenString string) (*SSOClaims, error) {
	if SSOPublicKey == nil {
		return nil, errors.New("SSO 公钥未加载")
	}

	token, err := jwt.ParseWithClaims(tokenString, &SSOClaims{}, func(token *jwt.Token) (interface{}, error) {
		// 验证签名方法必须是 RSA
		if _, ok := token.Method.(*jwt.SigningMethodRSA); !ok {
			return nil, TokenInvalid
		}
		return SSOPublicKey, nil
	})

	if err != nil {
		if ve, ok := err.(*jwt.ValidationError); ok {
			switch {
			case ve.Errors&jwt.ValidationErrorMalformed != 0:
				return nil, TokenMalformed
			case ve.Errors&jwt.ValidationErrorExpired != 0:
				return nil, TokenExpired
			case ve.Errors&jwt.ValidationErrorNotValidYet != 0:
				return nil, TokenNotValidYet
			default:
				return nil, TokenInvalid
			}
		}
		return nil, TokenInvalid
	}

	if claims, ok := token.Claims.(*SSOClaims); ok && token.Valid {
		return claims, nil
	}

	return nil, TokenInvalid
}

// autoRefreshToken 自动刷新 Token（从 session 获取 refresh_token）
func autoRefreshToken(c *gin.Context) (*api.TokenResponse, error) {
	session := sessions.Default(c)
	refreshToken := session.Get("refresh_token")
	if refreshToken == nil {
		return nil, errors.New("未找到 refresh_token，请重新登录")
	}

	refreshTokenStr, ok := refreshToken.(string)
	if !ok {
		return nil, errors.New("refresh_token 格式错误")
	}

	// 向 SSO 刷新 token
	tokenResp, err := api.RefreshAccessTokenFromSSO(refreshTokenStr)
	if err != nil {
		global.Log.Error("自动刷新 token 失败", zap.Error(err))
		// 刷新失败，清除 session
		session.Delete("refresh_token")
		session.Delete("refresh_token_expires_at")
		session.Save()
		return nil, err
	}

	// 更新 session 中的 refresh_token（如果 SSO 返回了新的）
	if tokenResp.RefreshToken != "" && tokenResp.RefreshToken != refreshTokenStr {
		session.Set("refresh_token", tokenResp.RefreshToken)
		session.Save()
	}

	global.Log.Info("✓ 自动刷新 token 成功")
	return tokenResp, nil
}

// GetUserUUID 从上下文获取用户 UUID
// Deprecated: 请使用 utils.GetUUID
// func GetUserUUID(c *gin.Context) uuid.UUID {
// 	if val, exists := c.Get("user_uuid"); exists {
// 		if userUUID, ok := val.(uuid.UUID); ok {
// 			return userUUID
// 		}
// 	}
// 	return uuid.Nil
// }

// GetUserEmail 从上下文获取用户邮箱
func GetUserEmail(c *gin.Context) string {
	if val, exists := c.Get("user_email"); exists {
		if email, ok := val.(string); ok {
			return email
		}
	}
	return ""
}

// GetClaims 从上下文获取 Claims
func GetClaims(c *gin.Context) *SSOClaims {
	if val, exists := c.Get("claims"); exists {
		if claims, ok := val.(*SSOClaims); ok {
			return claims
		}
	}
	return nil
}
