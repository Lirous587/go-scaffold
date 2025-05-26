package role

import (
	"github.com/gin-gonic/gin"
)

func RegisterV1(r *gin.RouterGroup, ctrl *Controller) bool {
	g := r.Group("/v1/role")
	{
		g.POST("/", ctrl.CreateRole)
		g.DELETE("/:id", ctrl.DeleteRole)
		g.PUT("/:id", ctrl.UpdateRole)
		g.GET("/all", ctrl.AllRole)
	}
	return true
}
