package utils

import (
	"github.com/gin-gonic/gin"
	"github.com/gofrs/uuid"
)

// GetUUID 从 Gin 的 Context 中获取 JWT 解析出来的用户 UUID
func GetUUID(c *gin.Context) uuid.UUID {
	if val, exists := c.Get("user_uuid"); exists {
		if userUUID, ok := val.(uuid.UUID); ok {
			return userUUID
		}
	}
	return uuid.Nil
}
