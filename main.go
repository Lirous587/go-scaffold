package main

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/aarondl/sqlboiler/v4/boil"
	"github.com/gin-gonic/gin"
	"github.com/pkg/errors"
	"github.com/subosito/gotenv"
	swaggerfiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"go.uber.org/zap"
	"os"
	"os/signal"
	_ "scaffold/api/openapi"
	"scaffold/internal/captcha"
	"scaffold/internal/common/logger"
	"scaffold/internal/common/metrics"
	"scaffold/internal/common/server"
	"scaffold/internal/common/uid"
	"scaffold/internal/img"
	"scaffold/internal/user"
	"syscall"
)

func setGDB() {
	host := os.Getenv("PSQL_HOST")
	port := os.Getenv("PSQL_PORT")
	username := os.Getenv("PSQL_USERNAME")
	password := os.Getenv("PSQL_PASSWORD")
	dbname := os.Getenv("PSQL_DB_NAME")
	sslmode := os.Getenv("PSQL_SSL_MODE")

	dsn := fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
		host, port, username, password, dbname, sslmode,
	)

	db, err := sql.Open("postgres", dsn)
	if err != nil {
		panic(err)
	}

	// 测试连接
	if err := db.Ping(); err != nil {
		panic(fmt.Sprintf("无法连接到数据库: %v", err))
	}

	boil.SetDB(db)

	boil.DebugMode = true

	logMode := os.Getenv("LOG_MODE")
	if logMode != "dev" {
		if err := os.MkdirAll("./logs", 0755); err != nil {
			panic(fmt.Sprintf("创建日志目录失败:%v", err))
		}
		fh, err := os.OpenFile("./logs/sqlboiler.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
		if err != nil {
			panic(fmt.Sprintf("打开debug日志错误:%v", err))
		}
		boil.DebugWriter = fh
	}
}

func syncWorker(ctx context.Context) {

}

func sync(ctx context.Context, cancel context.CancelFunc) {
	// 监听系统信号
	go func() {
		c := make(chan os.Signal, 1)
		signal.Notify(c, os.Interrupt, syscall.SIGTERM)
		<-c
		zap.L().Info("收到退出信号，开始优雅关闭")
		cancel()
	}()

	go syncWorker(ctx)
}

// @title           自定义title
// @version         1.0
// @description     自定义描述
// @termsOfService  https://Lirous.com

// @contact.name   Lirous
// @contact.url    https://Lirous.com
// @contact.email  lirous@lirous.com

// @license.name  MIT
// @license.url   https://github.com/Lirous587/go-scaffold/main/LICENSE

// @host      localhost:8080
// @BasePath  /api

// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
// @description Type "Bearer" followed by a space and JWT token.

// @externalDocs.description  OpenAPI
// @externalDocs.url          https://swagger.io/resources/open-api/
// swag init -g main.go -o ./api/openapi
func main() {
	var err error

	if err := gotenv.Load(); err != nil {
		panic(err)
	}

	uid.Init()

	setGDB()

	if err = logger.Init(); err != nil {
		panic(errors.WithMessage(err, "logger模块初始化失败"))
	}

	ctx, cancel := context.WithCancel(context.Background())
	go sync(ctx, cancel)

	metricsClient := metrics.NewPrometheusClient()
	metrics.StartPrometheusServer()

	server.RunHttpServer(os.Getenv("SERVER_PORT"), metricsClient, func(r *gin.RouterGroup) {
		r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerfiles.Handler,
			ginSwagger.PersistAuthorization(true)))

		user.InitV1(r)
		captcha.InitV1(r)
		img.InitV1(r)
	})
}
