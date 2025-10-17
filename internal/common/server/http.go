package server

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"scaffold/internal/common/metrics"
	"scaffold/internal/common/utils"
	"scaffold/internal/common/validator"
	"strings"
	"syscall"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"github.com/pkg/errors"
)

func RunHttpServer(port string, metricsClient metrics.Client, registerRouter func(r *gin.RouterGroup), clearFunc ...func()) {
	if port == "" {
		panic(errors.New("RunHttpServer中的port无效"))
	}

	_ = godotenv.Load()
	mode := utils.GetEnv("SERVER_MODE")

	if mode == "dev" {
		gin.SetMode(gin.DebugMode)
	} else {
		gin.SetMode(gin.ReleaseMode)
	}

	engine := gin.Default()

	engine.Use(errorHandler(), logHandler(), metricsHandler(metricsClient))

	// 注册验证器
	if err := validator.Init(); err != nil {
		panic(errors.WithMessage(err, "validator模块初始化失败"))
	}

	// 配置CORS中间件
	setCORS(engine)

	// 配置404路由
	engine.NoRoute(func(c *gin.Context) {
		c.JSONP(404, gin.H{"msg": "404"})
	})

	routerGroup := engine.Group("/api")

	registerRouter(routerGroup)

	// 创建HTTP服务器
	server := &http.Server{
		Addr:    fmt.Sprintf(":%s", port),
		Handler: engine,
	}

	// 启动服务器
	go func() {
		log.Printf("服务器启动,端口:%v\n", port)

		if err := server.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.Fatalf("服务器启动失败,err:%#v\n", err)
		}
	}()

	// 等待终止信号
	sig := waitForSignal()
	log.Printf("接收到信号:%v\n", sig.String())

	log.Println("正在关闭服务器...")

	if len(clearFunc) > 0 {
		log.Println("正在执行资源清理")
		clearFunc[0]()
	}

	// 优雅关闭服务
	shutdownServer(server)
}

func waitForSignal() os.Signal {
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM, syscall.SIGHUP)
	return <-quit
}

// 优雅关闭服务器
func shutdownServer(server *http.Server) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := server.Shutdown(ctx); err != nil {
		log.Fatalf("服务器关闭失败,err:%#v\n", err)
	}
	log.Println("服务器已退出")
}

func setCORS(r *gin.Engine) {
	corsCfg := cors.DefaultConfig()
	allowsStr := utils.GetEnv("SERVER_ALLOW_ORIGINS")
	allows := strings.Split(allowsStr, ",")

	corsCfg.AllowOrigins = allows
	corsCfg.AllowMethods = []string{"GET", "POST", "PUT", "DELETE", "PATCH"}
	corsCfg.AllowHeaders = []string{"Origin", "Content-Type", "Authorization", "X-Refresh-Token"}
	r.Use(cors.New(corsCfg))
}
