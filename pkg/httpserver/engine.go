package httpserver

import (
	"fmt"
	"net/http"
	"path"
	"scaffold/pkg/apigen"
	"syscall"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"

	"scaffold/pkg/apigen/swagger"
	"scaffold/pkg/config"
	"scaffold/pkg/i18n"
	"scaffold/pkg/logger"
	"scaffold/pkg/middleware"
	"scaffold/pkg/validator"

	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

type Server struct {
	engine  *gin.Engine
	swagger *swagger.Swagger
}

func New() *Server {
	cfg := config.Cfg.App
	// 设置运行模式
	if cfg.Mode == gin.ReleaseMode {
		gin.SetMode(gin.ReleaseMode)
	} else {
		gin.SetMode(gin.DebugMode)
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

	// 2. 初始化swagger
	swg := swagger.New()

	// 3. 添加Swagger UI路由
	r.StaticFile("/swagger-docs/swagger.json", "./docs/swagger.json")
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler,
		ginSwagger.URL("/swagger-docs/swagger.json")))

	return &Server{
		engine:  r,
		swagger: swg,
	}
}

func (s *Server) Run(addr string) {
	// 创建HTTP服务器
	srv := &http.Server{
		Addr:    fmt.Sprintf(":%s", addr),
		Handler: s.engine,
	}

	// 将注册的路由保存到swagger
	if err := s.swagger.Save(); err != nil {
		zap.L().Fatal("swagger生成失败", zap.Error(err))
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

func (s *Server) Bind(controller interface{}) {
	apigen.RegisterAPI(s.engine, controller)
	s.swagger.Bind("", controller)
}

func (s *Server) Group(relativePath string, handle func(group *ServerGroup)) {
	group := s.engine.Group(relativePath)
	sg := &ServerGroup{
		server: Server{
			engine:  s.engine,
			swagger: s.swagger,
		},
		group:    group,
		basePath: relativePath,
	}
	handle(sg)
}

type ServerGroup struct {
	server   Server
	group    *gin.RouterGroup
	basePath string
}

func (sg *ServerGroup) Bind(controllers ...interface{}) {
	// 获取当前路由组的路径前缀
	pathPrefix := sg.basePath

	for _, controller := range controllers {
		apigen.RegisterAPI(sg.group, controller)
		sg.server.swagger.Bind(pathPrefix, controller)
	}
}

func (sg *ServerGroup) Group(relativePath string, handle func(group *ServerGroup)) {
	subGroup := sg.group.Group(relativePath)
	// 合并路径，确保路径格式正确（处理斜杠）
	newBasePath := path.Join(sg.basePath, relativePath, "/")

	subSg := &ServerGroup{
		server:   sg.server,
		group:    subGroup,
		basePath: newBasePath,
	}
	handle(subSg)
}

func (sg *ServerGroup) Middleware(middlewares ...gin.HandlerFunc) {
	sg.group.Use(middlewares...)
}

