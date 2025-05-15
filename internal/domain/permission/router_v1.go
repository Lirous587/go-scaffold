package permission

import (
	"comment/internal/middleware"
	"comment/pkg/config"
	"github.com/gin-gonic/gin"
)

func RegisterV1(r *gin.RouterGroup, ctrl IController) error {
	g := r.Group("/v1/permission")
	authMiddleware := middleware.NewUserAuth([]byte(config.Cfg.JWT.Secret))
	{
		g.POST("/", authMiddleware.Validate(), ctrl.CreatePermission)
		g.PUT("/:id", authMiddleware.Validate(), ctrl.UpdatePermission)
		g.DELETE("/:id", authMiddleware.Validate(), ctrl.DeletePermission)
	}
	return nil
}
