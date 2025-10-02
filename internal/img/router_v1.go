package img

import (
	"scaffold/internal/common/middleware/auth"
	"scaffold/internal/img/handler"
	"github.com/gin-gonic/gin"
)

func RegisterV1(r *gin.RouterGroup, handler *handler.HttpHandler) func() {
	g := r.Group("/v1/img")
	{
	}

	protect := g.Use(auth.JWTValidate())
	{
		// 如果上传文件过大 可能导致连接重置 后端解决方案如下
		//g.POST("/upload",middlewares.FullRequest() ,auth.Validate(), handler.Upload)
		// 关于连接reset的原因: 上传文件为流式操作，不同于简单的crud 如果上传时token过期 返回错误会导致连接重置 对前端及其不友好
		// 当前前端解决方案为上传时刷新token 这样可以有效避免服务端的资源浪费
		protect.POST("/upload", handler.Upload)

		protect.DELETE("/:id", handler.Delete)
		protect.GET("", handler.List)

		// 回收站
		protect.DELETE("/recycle/:id", handler.ClearRecycleBin)
		protect.PUT("/recycle/:id", handler.RestoreFromRecycleBin)

		// 分类
		protect.POST("/category", handler.CreateCategory)
		protect.DELETE("/category/:id", handler.DeleteCategory)
		protect.PUT("/category/:id", handler.UpdateCategory)
		protect.GET("/categories", handler.ListCategories)
	}

	go func() {
		handler.ListenDeleteQueue()
	}()

	return nil
}
