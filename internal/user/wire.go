//go:build wireinject
// +build wireinject

package user

import (
	"github.com/gin-gonic/gin"
	"github.com/google/wire"
	"scaffold/internal/user/adapters"
	"scaffold/internal/user/handler"
	"scaffold/internal/user/infrastructure"
	"scaffold/internal/user/service"
)

func InitV1(r *gin.RouterGroup) func() {
	wire.Build(
		RegisterV1,
		handler.NewHttpHandler,
		service.NewService,
		adapters.NewPSQLRepository,
		adapters.NewRedisCache,
		wire.Bind(new(infrastructure.UserRepository), new(*adapters.PSQLRepository)),
		wire.Bind(new(infrastructure.UserCache), new(*adapters.RedisCache)),
	)
	return nil
}
