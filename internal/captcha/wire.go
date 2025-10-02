//go:build wireinject
// +build wireinject

package captcha

import (
	"scaffold/internal/captcha/adapters"
	"scaffold/internal/captcha/handler"
	"scaffold/internal/captcha/service"
	"github.com/gin-gonic/gin"
	"github.com/google/wire"
)

func InitV1(r *gin.RouterGroup) func() {
	wire.Build(
		RegisterV1,
		handler.NewHttpHandler,
		service.NewCaptchaServiceFactor,
		adapters.NewCaptchaRedisCache,
	)

	return nil
}

func NewVerifyMiddleware() *handler.HttpHandler {
	wire.Build(
		handler.NewHttpHandler,
		service.NewCaptchaServiceFactor,
		adapters.NewCaptchaRedisCache,
	)

	return nil
}
