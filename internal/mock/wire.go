//go:build wireinject
// +build wireinject

package mock

import (
	"github.com/gin-gonic/gin"
	"github.com/google/wire"
	"scaffold/internal/mock/adapters"
	"scaffold/internal/mock/handler"
	"scaffold/internal/mock/service"
)

func InitV1(r *gin.RouterGroup) func() {
	wire.Build(
		RegisterV1,
		handler.NewHttpHandler,
		service.NewMockService,
		adapters.NewPSQLMockRepository,
	)

	return nil
}
