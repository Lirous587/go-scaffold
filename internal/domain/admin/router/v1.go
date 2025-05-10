package router

import (
	"scaffold/internal/domain/admin/controller"

	"github.com/gin-gonic/gin"
)

func RegisterV1(r *gin.RouterGroup, ctrl controller.Controller) error {
	g := r.Group("/v1/admin")
	{
		g.GET("/ifInit", ctrl.IfInit)
		g.POST("/init", ctrl.Init)
		g.POST("/login", ctrl.Login)
		g.POST("/refresh_token", ctrl.RefreshToken)
	}
	return nil
}
