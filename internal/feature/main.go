package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/pkg/errors"
	"golang.org/x/net/context"
	"scaffold/internal/common/httpserver"
	"scaffold/internal/common/logger"
	"scaffold/internal/common/metrics"
	"scaffold/internal/feature/app/command"
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

	application := service.NewApplication(ctx, metricsClient)

	//serverType := strings.ToLower(os.Getenv("SERVER_TO_RUN"))
	serverType := "http"
	switch serverType {
	case "http":
		// 使用application去注册路由
		loginT := command.LoginWithType{
			LoginType: "haha",
		}

		r := httpserver.New(8080, "dev")
		r.Router.GET("/user", func(c *gin.Context) {
			application.Commands.LoginWithType.Handle(ctx, loginT)
			c.JSON(200, gin.H{
				"msg": "ok",
			})
		})
		r.Run()
		//server.RunHTTPServer(func(router chi.Router) http.Handler {
		//	return ports.HandlerFromMux(
		//		ports.NewHttpServer(application),
		//		router,
		//	)
		//})
	case "grpc":

	default:
		panic(fmt.Sprintf("server type '%s' is not supported", serverType))
	}

}
