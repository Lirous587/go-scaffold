package httpserver

import (
	"context"
	"fmt"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/pkg/errors"
	"go.uber.org/zap"
	"net/http"
	"os"
	"os/exec"
	"os/signal"
	"scaffold/pkg/config"
	"scaffold/pkg/logger"
	"scaffold/pkg/repository"
	"syscall"
	"time"
)

func validateConfig(port int) (*config.ServerConfig, bool) {
	for _, server := range config.Cfg.Server {
		if port == server.Port {
			return &server, true
		}
	}
	return nil, false
}

type Server struct {
	Router       *gin.Engine
	config       *config.ServerConfig
	server       *http.Server
	stopHandlers []func() // 存储关闭处理函数
}

func (s *Server) setupCORS() {
	corsCfg := cors.DefaultConfig()
	corsCfg.AllowOrigins = []string{"https://Lirous.com", "http://localhost:3000"}
	corsCfg.AllowMethods = []string{"GET", "POST", "PUT", "DELETE"}
	corsCfg.AllowHeaders = []string{"Origin", "Content-Type", "Authorization", "Refresh-Token"}
	s.Router.Use(cors.New(corsCfg))
}

func New(port int) *Server {
	serverConfig, ok := validateConfig(port)
	if !ok {
		zap.L().Panic("此服务端口没有配置成功", zap.Int("port", port))
	}

	// 设置运行模式
	if serverConfig.Mode == gin.ReleaseMode {
		gin.SetMode(gin.ReleaseMode)
	} else {
		gin.SetMode(gin.DebugMode)
	}

	// 创建Gin引擎并配置基础中间件
	r := gin.New()
	r.Use(logger.GinLogger(), logger.GinRecovery(true))

	r.Use(ErrorHandler())

	// 创建服务器实例
	server := &Server{
		Router: r,
		config: serverConfig,
		server: &http.Server{
			Addr:    fmt.Sprintf(":%d", serverConfig.Port),
			Handler: r,
		},
	}

	// 配置CORS中间件
	server.setupCORS()

	// 配置404路由
	r.NoRoute(func(c *gin.Context) {
		c.JSONP(404, gin.H{"msg": "404"})
	})

	return server
}

func (s *Server) RegisterStopHandler(handler func()) {
	s.stopHandlers = append(s.stopHandlers, handler)
}

func (s *Server) executeStopHandlers() {
	for i := len(s.stopHandlers) - 1; i >= 0; i-- {
		s.stopHandlers[i]()
	}
}

func (s *Server) Run() {
	// 创建HTTP服务器
	s.server = &http.Server{
		Addr:    fmt.Sprintf(":%d", s.config.Port),
		Handler: s.Router,
	}

	// 启动服务器
	s.startServer()

	// 等待终止信号
	sig := s.waitForSignal()
	zap.L().Info("接收到信号", zap.String("信号", sig.String()))

	// 执行关闭前的处理函数
	zap.L().Info("正在关闭服务器...")
	s.executeStopHandlers()

	// 处理重启信号
	if sig == syscall.SIGHUP {
		zap.L().Info("正在重启服务器...")
		s.restartServer()
	}

	// 优雅关闭服务
	s.shutdownServer()
}

func (s *Server) startServer() {
	go func() {
		zap.L().Info("服务器启动", zap.Int("端口", s.config.Port))
		if err := s.server.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			zap.L().Fatal("服务器启动失败", zap.Error(err))
		}
	}()
}

func (s *Server) waitForSignal() os.Signal {
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM, syscall.SIGHUP)
	return <-quit
}

// 优雅关闭服务器
func (s *Server) shutdownServer() {
	s.cleanup()

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := s.server.Shutdown(ctx); err != nil {
		zap.L().Fatal("服务器关闭错误", zap.Error(err))
	}
	zap.L().Info("服务器已退出")
}

// 重启服务器
func (s *Server) restartServer() {
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

// 清理资源
func (s *Server) cleanup() {
	repository.RedisClose() // 关闭Redis连接
	repository.GormClose()  // 关闭数据库连接
	zap.L().Info("所有资源已成功关闭")
}
