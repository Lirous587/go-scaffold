package captcha

import (
	"scaffold/internal/captcha/handler"
	"scaffold/internal/common/middleware/auth"
	"github.com/gin-gonic/gin"
)

func RegisterV1(r *gin.RouterGroup, handler *handler.HttpHandler) func() {
	g := r.Group("/v1/captcha")
	{
		g.POST("", handler.Gen)
		//测试端点
		g.POST("/verify", handler.Verify(), handler.VerifyEndpoint)

		// 测试路由：生成验证码并返回图片+验证答案
		g.POST("/test", auth.Validate(), handler.GenWithAnswer)
	}

	return nil
}
