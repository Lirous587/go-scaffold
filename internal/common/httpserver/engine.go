package httpserver

import (
	"context"
	"fmt"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/pkg/errors"
	"log"
	"net/http"
	"os"
	"os/exec"
	"os/signal"
	"syscall"
	"time"
)

type Server struct {
	Router       *gin.Engine
	server       *http.Server
	stopHandlers []func() // 存储关闭处理函数
	port         int
}

func New(port int, mode string) *Server {
	// 设置运行模式
	if mode == gin.ReleaseMode {
		gin.SetMode(gin.ReleaseMode)
	} else {
		gin.SetMode(gin.DebugMode)
	}

	// 创建Gin引擎并配置基础中间件
	r := gin.Default()

	r.Use(ErrorHandler())

	// 创建服务器实例
	server := &Server{
		Router: r,
		server: &http.Server{
			Addr:    fmt.Sprintf(":%d", port),
			Handler: r,
		},
		port: port,
	}

	// 配置CORS中间件
	server.setupCORS()

	// 配置404路由
	r.NoRoute(func(c *gin.Context) {
		c.JSONP(404, gin.H{"msg": "404"})
	})

	return server
}

func (s *Server) Run() {
	// 创建HTTP服务器
	s.server = &http.Server{
		Addr:    fmt.Sprintf(":%d", s.port),
		Handler: s.Router,
	}

	// 启动服务器
	s.startServer()

	// 等待终止信号
	sig := s.waitForSignal()
	log.Printf("接收到信号:%v\n", sig.String())

	// 执行关闭前的处理函数
	log.Println("正在关闭服务器...")

	s.executeStopHandlers()

	// 处理重启信号
	if sig == syscall.SIGHUP {
		log.Println("正在重启服务器...")
		s.restartServer()
	}

	// 优雅关闭服务
	s.shutdownServer()
}

func (s *Server) RegisterStopHandler(handler func()) {
	s.stopHandlers = append(s.stopHandlers, handler)
}

func (s *Server) executeStopHandlers() {
	for i := len(s.stopHandlers) - 1; i >= 0; i-- {
		s.stopHandlers[i]()
	}
}

func (s *Server) startServer() {
	go func() {
		log.Printf("服务器启动,端口:%d\n", s.port)

		if err := s.server.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.Fatalf("服务器启动失败,err:%#v\n", err)
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
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := s.server.Shutdown(ctx); err != nil {
		log.Fatalf("服务器关闭失败,err:%#v\n", err)
	}
	log.Println("服务器已退出")
}

// 重启服务器
func (s *Server) restartServer() {
	execPath, err := os.Executable()
	if err != nil {
		log.Fatalf("获取可执行文件路径失败,err:%#v\n", err)
		return
	}

	cmd := exec.Command(execPath, os.Args[1:]...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	if err := cmd.Start(); err != nil {
		log.Fatalf("启动新进程失败,err:%#v\n", err)
		return
	}
	log.Printf("新进程已启动,PID:%#v\n", cmd.Process.Pid)
}

func (s *Server) setupCORS() {
	corsCfg := cors.DefaultConfig()
	allow := os.Getenv("SERVER_ALLOW_ORIGINS")
	if allow == "" {
		panic(errors.New("httpserver加载SERVER_ALLOW_ORIGINS环境变量失败"))
	}
	corsCfg.AllowOrigins = []string{allow}
	corsCfg.AllowMethods = []string{"GET", "POST", "PUT", "DELETE", "PATCH"}
	corsCfg.AllowHeaders = []string{"Origin", "Content-Type", "Authorization", "Refresh-Token"}
	s.Router.Use(cors.New(corsCfg))
}
