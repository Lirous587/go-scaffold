package user

import (
	"github.com/gin-gonic/gin"
	"os"
	"scaffold/internal/common/middleware"
)

func RegisterV1(r *gin.RouterGroup, ctrl IController) error {
	g := r.Group("/v1/user")
	secret := os.Getenv("JWT_SECRET")
	if secret == "" {
		panic("JWT_SECRET为空")
	}

	authMiddleware := middleware.NewUserAuth([]byte(secret))
	{
		g.POST("/auth", authMiddleware.Validate())
		g.POST("/login", ctrl.Login)
		g.POST("/refresh_token", ctrl.RefreshToken)
	}
	return nil
}
