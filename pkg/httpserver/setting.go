package httpserver

import (
	"context"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/pkg/errors"
	"go.uber.org/zap"
	"net/http"
	"os"
	"os/exec"
	"os/signal"
	"scaffold/pkg/repository/db"
	"scaffold/pkg/repository/redis"
	"syscall"
	"time"
)

// 配置CORS中间件
func setupCORS(r *gin.Engine) {
	corsCfg := cors.DefaultConfig()
	corsCfg.AllowOrigins = []string{"https://Lirous.com", "http://localhost:3000"}
	corsCfg.AllowMethods = []string{"GET", "POST", "PUT", "DELETE"}
	corsCfg.AllowHeaders = []string{"Origin", "Content-Type", "Authorization", "Refresh-Token"}
	r.Use(cors.New(corsCfg))
}

// 启动HTTP服务器
func startServer(srv *http.Server, port int) {
	go func() {
		zap.L().Info("服务器启动", zap.Int("端口", port))
		if err := srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
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
