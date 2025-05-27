package user

import (
	"github.com/gin-gonic/gin"
	"scaffold/internal/common/middleware/auth"
	"scaffold/internal/user/handler"
)

func RegisterV1(r *gin.RouterGroup, handler *handler.HttpHandler) func() {
	userGroup := r.Group("/user")
	{
		// 认证相关路由
		userGroup.POST("/auth/github", handler.GithubAuth)

		// 需要token的路由
		protected := userGroup.Group("")
		protected.Use(auth.Validate())
		{
			protected.POST("/refresh", handler.RefreshToken)
			protected.GET("/profile", handler.GetProfile)
			protected.PUT("/profile", handler.UpdateProfile)
		}
	}
	return nil
}
