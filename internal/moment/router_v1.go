package moment

import (
	"blog-v4/internal/common/middleware/auth"
	"blog-v4/internal/moment/handler"
	"github.com/gin-gonic/gin"
)

func RegisterV1(r *gin.RouterGroup, handler *handler.HttpHandler) func() {
	g := r.Group("/v1/moment")
	{
		g.GET("", handler.List)
		g.GET("/random", handler.RandomN)
	}

	protect := g.Use(auth.Validate())
	{
		protect.POST("", handler.Create)
		protect.DELETE("/:id", handler.Delete)
		protect.PUT("/:id", handler.Update)
	}
	return nil
}
