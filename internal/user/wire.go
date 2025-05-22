//go:build wireinject
// +build wireinject

package user

import (
	"github.com/gin-gonic/gin"
	"github.com/google/wire"
	"scaffold/internal/user/adapters"
	"scaffold/internal/user/infrastructure"
)

func InitV1(r *gin.RouterGroup) bool {
	wire.Build(
		RegisterV1,
		NewController,
		NewService,
		adapters.NewPSQLRepository,
		adapters.NewRedisCache,
		wire.Bind(new(infrastructure.UserRepository), new(*adapters.PSQLRepository)),
		wire.Bind(new(infrastructure.UserCache), new(*adapters.RedisCache)),
	)
	return true
}
