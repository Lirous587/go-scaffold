package role

import (
	"comment/internal/middleware"
	"comment/pkg/config"
	"github.com/gin-gonic/gin"
)

func RegisterV1(r *gin.RouterGroup, ctrl IController) error {
	g := r.Group("/v1/role")
	authMiddleware := middleware.NewUserAuth([]byte(config.Cfg.JWT.Secret))
	{
		g.POST("/auth", authMiddleware.Validate())
	}
	return nil
}
