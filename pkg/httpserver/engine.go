package httpserver

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/exec"
	"os/signal"
	"syscall"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"

	"scaffold/pkg/config"
	"scaffold/pkg/i18n"
	"scaffold/pkg/logger"
	"scaffold/pkg/middleware"
	"scaffold/pkg/repository/db"
	"scaffold/pkg/repository/redis"
	"scaffold/pkg/validator"
)

// Init 初始化并配置Gin引擎
func Init(cfg *config.AppConfig) *gin.Engine {
	// 设置运行模式
	if cfg.Mode == gin.ReleaseMode {
		gin.SetMode(gin.ReleaseMode)
	}

	// 初始化依赖模块
	validator.Setup() // 设置自定义验证器
	i18n.Setup()      // 设置国际化

	// 创建Gin引擎并配置基础中间件
	r := gin.New()
	r.Use(logger.GinLogger(), logger.GinRecovery(true))

	// 配置CORS中间件
	setupCORS(r)

	// 配置404路由
	r.NoRoute(func(c *gin.Context) {
		c.JSONP(404, gin.H{"msg": "404"})
	})

	// 配置错误处理中间件
	r.Use(middleware.ErrorHandler())

	return r
}

// 配置CORS中间件
func setupCORS(r *gin.Engine) {
	corsCfg := cors.DefaultConfig()
	corsCfg.AllowOrigins = []string{"https://Lirous.com", "http://localhost:3000"}
	corsCfg.AllowMethods = []string{"GET", "POST", "PUT", "DELETE"}
	corsCfg.AllowHeaders = []string{"Origin", "Content-Type", "Authorization", "Refresh-Token"}
	r.Use(cors.New(corsCfg))
}

// Run 启动HTTP服务器并处理优雅关闭
func Run(router *gin.Engine, addr string) {
	// 创建HTTP服务器
	srv := &http.Server{
		Addr:    fmt.Sprintf(":%s", addr),
		Handler: router,
	}

	// 启动服务器
	startServer(srv, addr)

	// 等待终止信号
	sig := waitForSignal()
	zap.L().Info("接收到信号", zap.String("信号", sig.String()))

	// 处理重启信号
	if sig == syscall.SIGHUP {
		zap.L().Info("正在重启服务器...")
		restartServer()
	}

	// 优雅关闭服务
	shutdownServer(srv)
}

// 启动HTTP服务器
func startServer(srv *http.Server, addr string) {
	go func() {
		zap.L().Info("服务器启动", zap.String("地址", addr))
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			zap.L().Fatal("服务器启动失败", zap.Error(err))
		}
	}()
}

// 等待终止信号
func waitForSignal() os.Signal {
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM, syscall.SIGHUP)
	return <-quit
}

// 优雅关闭服务器
func shutdownServer(srv *http.Server) {
	zap.L().Info("正在关闭服务器...")
	cleanup()

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		zap.L().Fatal("服务器关闭错误", zap.Error(err))
	}
	zap.L().Info("服务器已退出")
}

// restartServer 重启服务器
func restartServer() {
	execPath, err := os.Executable()
	if err != nil {
		zap.L().Error("获取可执行文件路径失败", zap.Error(err))
		return
	}

	cmd := exec.Command(execPath, os.Args[1:]...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	if err := cmd.Start(); err != nil {
		zap.L().Error("启动新进程失败", zap.Error(err))
		return
	}
	zap.L().Info("新进程已启动", zap.Int("PID", cmd.Process.Pid))
}

// cleanup 清理资源
func cleanup() {
	redis.Close() // 关闭Redis连接
	db.Close()    // 关闭数据库连接
	zap.L().Info("所有资源已成功关闭")
}
