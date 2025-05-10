//go:build wireinject
// +build wireinject

package admin

import (
	"github.com/gin-gonic/gin"
	"github.com/google/wire"
	"scaffold/internal/domain/admin/controller"
	"scaffold/internal/domain/admin/infrastructure/cache"
	"scaffold/internal/domain/admin/infrastructure/db"
	"scaffold/internal/domain/admin/router"
	"scaffold/internal/domain/admin/service"
	"scaffold/internal/domain/admin/worker"
	"scaffold/pkg/repository"
)

func InitV1(r *gin.RouterGroup) error {
	wire.Build(
		router.RegisterV1,
		controller.NewController,
		service.NewService,
		repository.GormDB,
		repository.RedisClient,
		db.NewDB,
		cache.NewCache,
	)
	return nil
}

func InitWorker() worker.Worker {
	gormDB := repository.GormDB()
	dbDB := db.NewDB(gormDB)
	client := repository.RedisClient()
	cacheCache := cache.NewCache(client)
	workerWorker := worker.NewWorker(dbDB, cacheCache)
	return workerWorker
}
