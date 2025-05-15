//go:build wireinject
// +build wireinject

package permission

import (
	"comment/internal/domain/permission/infrastructure"
	"comment/pkg/repository"
	"github.com/gin-gonic/gin"
	"github.com/google/wire"
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
