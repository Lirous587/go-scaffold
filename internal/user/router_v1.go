package user

import (
	"github.com/gin-gonic/gin"
	"scaffold/internal/common/middleware/auth"
	"scaffold/internal/user/handler"
)

func RegisterRoutes(r *gin.RouterGroup, h *handler.HttpHandler) {
	userGroup := r.Group("/user")
	{
		// 认证相关路由（无需token）
		userGroup.POST("/auth/github", h.GithubAuth)

		// 需要token的路由
		protected := userGroup.Group("")
		protected.Use(auth.Validate())
		{
			protected.POST("/refresh", h.RefreshToken)
			protected.GET("/profile", h.GetProfile)
			protected.PUT("/profile", h.UpdateProfile)
		}
	}
}
