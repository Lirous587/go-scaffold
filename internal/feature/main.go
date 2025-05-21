package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/pkg/errors"
	"golang.org/x/net/context"
	"scaffold/internal/common/logger"
	"scaffold/internal/common/metrics"
	"scaffold/internal/common/server"
	"scaffold/internal/feature/ports"
	"scaffold/internal/feature/service"
)

func main() {
	var err error

	if err = logger.Init(); err != nil {
		panic(errors.WithMessage(err, "logger模块初始化失败"))
	}

	ctx := context.Background()

	metricsClient := metrics.NewPrometheusClient()
	metrics.StartPrometheusServer()

	//if err = validator.Init(); err != nil {
	//	panic(errors.WithMessage(err, "validator模块初始化失败"))
	//}

	application, cleanup := service.NewApplication(ctx, metricsClient)
	defer cleanup()

	//serverType := strings.ToLower(os.Getenv("SERVER_TO_RUN"))
	serverType := "http"
	switch serverType {
	case "http":
		server.RunHttpServer(8000, func(r *gin.RouterGroup) {
			httpServer := ports.NewHttpServer(application)
			ports.RegisterRouter(r, httpServer)
		})
	case "grpc":

	default:
		panic(fmt.Sprintf("server type '%s' is not supported", serverType))
	}
}
