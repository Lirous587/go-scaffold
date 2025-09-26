package mock

import (
    "scaffold/internal/common/middleware/auth"
    "scaffold/internal/mock/handler"
	"github.com/gin-gonic/gin"
)

func RegisterV1(r *gin.RouterGroup, handler *handler.HttpHandler) func() {
	g := r.Group("/v1/mock")
	{
		g.GET("/:id",handler.Read)
		g.GET("", handler.List)
	}

    protect := g.Use(auth.Validate())
    {
        protect.POST("", handler.Create)
        protect.DELETE("/:id", handler.Delete)
        protect.PUT("/:id", handler.Update)
    }
	return nil
}
