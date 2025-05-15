package user

import (
	"comment/internal/middleware"
	"comment/pkg/config"
	"github.com/gin-gonic/gin"
)

func RegisterV1(r *gin.RouterGroup, ctrl IController) error {
	g := r.Group("/v1/user")
	authMiddleware := middleware.NewUserAuth([]byte(config.Cfg.JWT.Secret))
	{
		g.POST("/auth", authMiddleware.Validate())
		g.POST("/login", ctrl.Login)
		g.POST("/refresh_token", ctrl.RefreshToken)
	}
	return nil
}
