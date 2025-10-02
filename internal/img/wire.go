//go:build wireinject
// +build wireinject

package img

import (
	"scaffold/internal/img/adapters"
	"scaffold/internal/img/handler"
	"scaffold/internal/img/service"
	"github.com/gin-gonic/gin"
	"github.com/google/wire"
)

func InitV1(r *gin.RouterGroup) func() {
	wire.Build(
		RegisterV1,
		handler.NewHttpHandler,
		service.NewImgService,
		adapters.NewImgPSQLRepository,
		adapters.NewImgRedisCache,
	)

	return nil
}
