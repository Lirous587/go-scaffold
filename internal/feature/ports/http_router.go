package ports

import (
	"github.com/gin-gonic/gin"
	"os"
)

func RegisterRouter(r *gin.RouterGroup, server HttpServer) {
	g := r.Group("/v1/user")
	secret := os.Getenv("JWT_SECRET")
	if secret == "" {
		panic("JWT_SECRET为空")
	}
	g.POST("/login", server.Login)

	//authMiddleware := middleware.NewUserAuth([]byte(secret))
	//{
	//	g.POST("/auth", authMiddleware.Validate())
	//
	//	//g.POST("/refresh_token", ctrl.RefreshToken)
	//}
	//return nil
}
