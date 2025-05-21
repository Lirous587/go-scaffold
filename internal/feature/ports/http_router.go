package ports

import (
	"github.com/gin-gonic/gin"
	"scaffold/internal/common/middleware/auth"
)

func RegisterRouter(r *gin.RouterGroup, server HttpServer) {
	g := r.Group("/v1/user")
	g.POST("/login", server.Login)
	g.POST("/auth", auth.Validate())
	//g.POST("/refresh_token", ctrl.RefreshToken)
}
