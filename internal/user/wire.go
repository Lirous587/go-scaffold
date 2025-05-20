//go:build wireinject
// +build wireinject

package user

import (
	"github.com/gin-gonic/gin"
	"github.com/google/wire"
	"scaffold/internal/user/infrastructure"
)

func InitV1(r *gin.RouterGroup) error {
	wire.Build(
		RegisterV1,
		NewController,
		NewService,
		repository.GormDB,
		repository.RedisClient,
		infrastructure.NewDB,
		infrastructure.NewCache,
	)
	return nil
}
