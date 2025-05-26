//go:build wireinject
// +build wireinject

package role

import (
	"github.com/gin-gonic/gin"
	"github.com/google/wire"
	"scaffold/internal/role/adapters"
	"scaffold/internal/role/infrastructure"
)

func InitV1(r *gin.RouterGroup) bool {
	wire.Build(
		RegisterV1,
		NewController,
		NewService,
		adapters.NewPSQLRepository,
		adapters.NewRedisCache,
		wire.Bind(new(infrastructure.RoleRepository), new(*adapters.PSQLRepository)),
		wire.Bind(new(infrastructure.RoleCache), new(*adapters.RedisCache)),
	)
	return true
}
