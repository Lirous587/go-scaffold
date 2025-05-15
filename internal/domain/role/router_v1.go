package role

import (
	"github.com/gin-gonic/gin"
	"scaffold/internal/middleware"
	"scaffold/pkg/config"
)

func RegisterV1(r *gin.RouterGroup, ctrl IController) error {
	g := r.Group("/v1/role")
	authMiddleware := middleware.NewUserAuth([]byte(config.Cfg.JWT.Secret))
	{
		g.POST("/auth", authMiddleware.Validate())
	}
	return nil
}
