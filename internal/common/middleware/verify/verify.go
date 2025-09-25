package verify

import (
	"scaffold/internal/captcha"
	"github.com/gin-gonic/gin"
)

func Verify() gin.HandlerFunc {
	// 通过 Wire 生成的函数获取 handler
	handler := captcha.NewVerifyMiddleware()
	// 调用 handler 的 Verify 方法
	return handler.Verify()
}
