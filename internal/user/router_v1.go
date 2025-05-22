package user

import (
	"github.com/gin-gonic/gin"
	"scaffold/internal/common/middleware/auth"
)

func RegisterV1(r *gin.RouterGroup, ctrl *Controller) bool {
	g := r.Group("/v1/user")

	g.POST("/auth", auth.Validate())
	g.POST("/login", ctrl.Login)
	g.POST("/refresh_token", ctrl.RefreshToken)

	return true
}
