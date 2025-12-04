package router

import (
	"go-mcp-context/internal/api"

	"github.com/gin-gonic/gin"
)

type UserRouter struct{}

func (u *UserRouter) InitUserRouter(PublicRouter *gin.RouterGroup) {
	userRouter := PublicRouter.Group("user")
	userApi := api.ApiGroupApp.UserApi
	{
		// 获取用户信息（从 SSO 获取）
		userRouter.GET("info", userApi.GetUserInfo)
	}
}
