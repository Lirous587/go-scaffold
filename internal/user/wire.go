//go:build wireinject
// +build wireinject

package user

import (
	"github.com/gin-gonic/gin"
	"github.com/google/wire"
	"scaffold/internal/user/adapters"
	"scaffold/internal/user/handler"
	"scaffold/internal/user/service"
)

func InitV1(r *gin.RouterGroup) func() {
	wire.Build(
		RegisterV1,
		handler.NewHttpHandler,
		service.NewTokenService,
		service.NewUserService,
		adapters.NewPSQLUserRepository,
		adapters.NewRedisTokenCache,
	)
	return nil
}
