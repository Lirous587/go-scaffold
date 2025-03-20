package engine

import (
	"context"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
	"os"
	"os/signal"
	"scaffold/pkg/config"
	"scaffold/pkg/logger"
	"syscall"
	"time"
)

func Init(cfg *config.AppConfig) *gin.Engine {
	if cfg.Mode == gin.ReleaseMode {
		gin.SetMode(gin.ReleaseMode)
	}


	// r := gin.Default()
	r := gin.New()
	// 创建CORS中间件
	corsCfg := cors.DefaultConfig()
	corsCfg.AllowOrigins = []string{"https://Lirous.com", "http://localhost:3000"}
	corsCfg.AllowMethods = []string{"GET", "POST", "PUT", "DELETE"}
	corsCfg.AllowHeaders = []string{"Origin", "Content-Type", "Authorization", "Refresh-Token"}

	r.Use(cors.New(corsCfg))
	r.Use(logger.GinLogger(), logger.GinRecovery(true))

	r.NoRoute(func(c *gin.Context) {
		c.JSONP(404, gin.H{
			"msg": "404",
		})
	})

	return r
}

func test() {
	router := gin.Default()

	srv := &http.Server{
		Addr:    ":8080",
		Handler: router,
	}

	go func() {
		// 开启一个goroutine启动服务
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("listen: %s\n", err)
		}
	}()

	// 等待中断信号来优雅地关闭服务器，为关闭服务器操作设置一个5秒的超时
	quit := make(chan os.Signal, 1) // 创建一个接收信号的通道
	// kill 默认会发送 syscall.SIGTERM 信号
	// kill -2 发送 syscall.SIGINT 信号，我们常用的Ctrl+C就是触发系统SIGINT信号
	// kill -9 发送 syscall.SIGKILL 信号，但是不能被捕获，所以不需要添加它
	// signal.Notify把收到的 syscall.SIGINT或syscall.SIGTERM 信号转发给quit
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM) // 此处不会阻塞

	<-quit // 阻塞在此，当接收到上述两种信号时才会往下执行
	log.Println("Shutdown Server ...")
	// 创建一个5秒超时的context
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	// 5秒内优雅关闭服务（将未处理完的请求处理完再关闭服务），超过5秒就超时退出
	if err := srv.Shutdown(ctx); err != nil {
		log.Fatal("Server Shutdown: ", err)
	}

	log.Println("Server exiting")
}
