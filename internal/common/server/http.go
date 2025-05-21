package server

import (
	"context"
	"fmt"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"github.com/pkg/errors"
	"log"
	"net/http"
	"os"
	"os/signal"
	"scaffold/internal/common/validator"
	"strings"
	"syscall"
	"time"
)

func RunHttpServer(port string, registerRouter func(r *gin.RouterGroup)) {
	if port == "" {
		panic(errors.New("RunHttpServer中的port无效"))
	}

	_ = godotenv.Load()
	mode := os.Getenv("SERVER_MODE")
	if mode == "" {
		panic("读取SERVER_MODE环境变量失败")
	}

	if mode == "dev" {
		gin.SetMode(gin.DebugMode)
	} else {
		gin.SetMode(gin.ReleaseMode)
	}

	engine := gin.Default()

	engine.Use(errorHandler())

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
	allow := os.Getenv("SERVER_ALLOW_ORIGINS")
	if allow == "" {
		panic(errors.New("httpserver加载SERVER_ALLOW_ORIGINS环境变量失败"))
	}
	corsCfg.AllowOrigins = []string{allow}
	corsCfg.AllowMethods = []string{"GET", "POST", "PUT", "DELETE", "PATCH"}
	corsCfg.AllowHeaders = []string{"Origin", "Content-Type", "Authorization", "Refresh-Token"}
	r.Use(cors.New(corsCfg))
}

// 错误链追踪 用于开发环境
func errorHandler() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		ctx.Next()
		// 处理错误
		if len(ctx.Errors) > 0 {
			for _, e := range ctx.Errors {
				// 记录详细错误日志
				//log.Printf("Error: %+v\n", e.Err)

				// 使用自定义格式化错误栈
				printBusinessStack(e.Err)
			}
		}
	}
}

func printBusinessStack(err error) {
	// 获取完整错误栈
	stackTrace := fmt.Sprintf("%+v", err)
	lines := strings.Split(stackTrace, "\n")

	// 错误消息
	if len(lines) > 0 {
		log.Printf("\n\n")
		log.Printf("业务逻辑错误: %s\n", lines[0])
	}

	// 记录已打印的栈帧数量
	framePrinted := 0
	maxBusinessFrames := 3 // 最多打印栈帧条数

	// 逐行检查并不做任何修改，保持原始格式
	for i := 0; i < len(lines)-1 && framePrinted < maxBusinessFrames; i++ {
		currentLine := lines[i]
		nextLine := lines[i+1]

		// 只检查是否为业务相关行，但完全保持原始格式
		if strings.Contains(currentLine, "internal") &&
			!strings.Contains(currentLine, "github.com/gin-gonic") &&
			!strings.Contains(currentLine, "net/http") &&
			!strings.Contains(currentLine, "net/http") &&
			!strings.Contains(currentLine, "/common/decorator") &&
			strings.Contains(nextLine, ".go:") {
			log.Println(currentLine)
			log.Println(nextLine)
			framePrinted++
		}
	}

	// 如果还有更多栈帧但已达到限制
	totalBusinessFrames := countBusinessFrames(lines)
	if framePrinted == maxBusinessFrames && framePrinted < totalBusinessFrames {
		log.Printf("一共%d条栈帧,实际打印%d条 (更多栈帧已省略)\n", totalBusinessFrames, maxBusinessFrames)
	}
}

// 计算业务栈帧总数
func countBusinessFrames(lines []string) int {
	count := 0
	for i := 0; i < len(lines)-1; i++ {
		currentLine := lines[i]
		nextLine := lines[i+1]

		if strings.Contains(currentLine, "internal") &&
			!strings.Contains(currentLine, "github.com/gin-gonic") &&
			!strings.Contains(currentLine, "net/http") &&
			strings.Contains(nextLine, ".go:") {
			count++
		}
	}
	return count
}
