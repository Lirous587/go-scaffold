package user

import (
	"scaffold/internal/common/middleware/auth"
	"scaffold/internal/user/handler"
	"github.com/gin-gonic/gin"
)

func RegisterV1(r *gin.RouterGroup, handler *handler.HttpHandler) func() {
	userGroup := r.Group("/v1/user")

	{
		// 认证相关路由
		userGroup.POST("/auth/github", handler.GithubAuth)
		userGroup.POST("/refresh_token", handler.RefreshToken)

		// 需要token的路由
		protected := userGroup.Group("")
		protected.Use(auth.Validate())
		{
			protected.POST("/auth", handler.ValidateAuth)
			protected.GET("/profile", handler.GetProfile)
		}
	}
	return nil
}
