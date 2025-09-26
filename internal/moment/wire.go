//go:build wireinject
// +build wireinject

package moment

import (
	"blog-v4/internal/moment/adapters"
	"blog-v4/internal/moment/handler"
	"blog-v4/internal/moment/service"
	"github.com/gin-gonic/gin"
	"github.com/google/wire"
)

func InitV1(r *gin.RouterGroup) func() {
	wire.Build(
		RegisterV1,
		handler.NewHttpHandler,
		service.NewMomentService,
		adapters.NewPSQLMomentRepository,
	)

	return nil
}
