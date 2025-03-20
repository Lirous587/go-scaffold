package engine

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/exec"
	"os/signal"
	"scaffold/pkg/config"
	"scaffold/pkg/logger"
	"scaffold/pkg/repository/redis"
	"syscall"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

func Init(cfg *config.AppConfig) *gin.Engine {
	if cfg.Mode == gin.ReleaseMode {
		gin.SetMode(gin.ReleaseMode)
	}

	r := gin.New()
	r.Use(logger.GinLogger(), logger.GinRecovery(true))

	// 创建CORS中间件
	corsCfg := cors.DefaultConfig()
	corsCfg.AllowOrigins = []string{"https://Lirous.com", "http://localhost:3000"}
	corsCfg.AllowMethods = []string{"GET", "POST", "PUT", "DELETE"}
	corsCfg.AllowHeaders = []string{"Origin", "Content-Type", "Authorization", "Refresh-Token"}
	r.Use(cors.New(corsCfg))

	r.NoRoute(func(c *gin.Context) {
		c.JSONP(404, gin.H{
			"msg": "404",
		})
	})

	return r
}

func Run(router *gin.Engine, addr string) {
	srv := &http.Server{
		Addr:    fmt.Sprintf(":%s", addr),
		Handler: router,
	}

	// 在goroutine中启动服务器
	go func() {
		zap.L().Info("服务器启动", zap.String("地址", addr))
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			zap.L().Fatal("服务器启动失败", zap.Error(err))
		}
	}()

	// 创建信号通道
	quit := make(chan os.Signal, 1)
	// 监听SIGINT和SIGTERM信号
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM, syscall.SIGHUP)

	// 等待信号
	sig := <-quit
	zap.L().Info("接收到信号", zap.String("信号", sig.String()))

	if sig == syscall.SIGHUP {
		// 处理SIGHUP为重启信号
		zap.L().Info("正在重启服务器...")
		restartServer()
	}

	zap.L().Info("正在关闭服务器...")

	cleanup()

	// 创建5秒超时上下文
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// 优雅关闭
	if err := srv.Shutdown(ctx); err != nil {
		zap.L().Fatal("服务器关闭错误", zap.Error(err))
	}

	zap.L().Info("服务器已退出")
}

// restartServer 重启服务器
func restartServer() {
	// 获取当前可执行文件路径
	execPath, err := os.Executable()
	if err != nil {
		zap.L().Error("获取可执行文件路径失败", zap.Error(err))
		return
	}

	// 使用相同的命令行参数启动新进程
	cmd := exec.Command(execPath, os.Args[1:]...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	// 启动新进程
	if err := cmd.Start(); err != nil {
		zap.L().Error("启动新进程失败", zap.Error(err))
		return
	}

	zap.L().Info("新进程已启动", zap.Int("PID", cmd.Process.Pid))
}

func cleanup() {
	// 关闭Redis连接
	redis.Close()

	// 关闭数据库连接
	// db.Close()

	zap.L().Info("所有资源已成功关闭")
}
