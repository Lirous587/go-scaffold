package httpserver

import (
	"fmt"
	"net/http"
	"scaffold/pkg/config"
	"syscall"

	"scaffold/pkg/logger"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type Server struct {
	engine *gin.Engine
	config *config.ServerConfig
}

func New(port int) *gin.Engine {
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

	// 配置CORS中间件
	setupCORS(r)

	// 配置404路由
	r.NoRoute(func(c *gin.Context) {
		c.JSONP(404, gin.H{"msg": "404"})
	})

	return r
}

func validateConfig(port int) (*config.ServerConfig, bool) {
	for _, server := range config.Cfg.Server {
		if port == server.Port {
			return &server, true
		}
	}
	return nil, false
}

func (s *Server) Run() {
	// 创建HTTP服务器
	srv := &http.Server{
		Addr:    fmt.Sprintf(":%d", s.config.Port),
		Handler: s.engine,
	}

	// 启动服务器
	startServer(srv, s.config.Port)

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
