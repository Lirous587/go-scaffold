package httpserver

import (
	"fmt"
	"net/http"
	"path"
	"syscall"

	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"go.uber.org/zap"
	"scaffold/pkg/config"
	"scaffold/pkg/httpserver/core/apigen"
	"scaffold/pkg/httpserver/core/apigen/swagger"
	"scaffold/pkg/logger"
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

	// 创建Gin引擎并配置基础中间件
	r := gin.New()
	r.Use(logger.GinLogger(), logger.GinRecovery(true))

	// 配置CORS中间件
	setupCORS(r)

	// 配置404路由
	r.NoRoute(func(c *gin.Context) {
		c.JSONP(404, gin.H{"msg": "404"})
	})

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

func (s *Server) Bind(apiInterfaces ...interface{}) {
	for _, apiInterface := range apiInterfaces {
		apiInfos := apigen.GetApiInfo(apiInterface)
		apigen.RegisterAPI(s.engine, apiInfos)
		s.swagger.GenerateDocs("", apiInfos)
	}
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

func (s *Server) Middleware(middlewares ...gin.HandlerFunc) {
	s.engine.Use(middlewares...)
}

type ServerGroup struct {
	server   Server
	group    *gin.RouterGroup
	basePath string
}

func (sg *ServerGroup) Bind(apiInterfaces ...interface{}) {
	// 获取当前路由组的路径前缀
	pathPrefix := sg.basePath
	for _, apiInterface := range apiInterfaces {
		apiInfos := apigen.GetApiInfo(apiInterface)
		apigen.RegisterAPI(sg.group, apiInfos)
		sg.server.swagger.GenerateDocs(pathPrefix, apiInfos)
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
