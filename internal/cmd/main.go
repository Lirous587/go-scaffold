package cmd

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"scaffold/internal/controller/mock1"
	mock2 "scaffold/internal/controller/mock2"
	"scaffold/pkg/config"
	"scaffold/pkg/httpserver"
)

func Main() {
	// 创建服务器
	s := httpserver.New()

	// 创建根路由组
	s.Group("/", func(group *httpserver.ServerGroup) {
		group.Group("/mock", func(group *httpserver.ServerGroup) {
			group.Bind(
				mock1.NewV1(),
				mock2.NewV1(),
			)
		})
		group.Middleware()
	})
	s.Middleware(
		func(context *gin.Context) {
			context.Next()
			fmt.Println(1)
		},
		func(context *gin.Context) {
			fmt.Println(2)
		},
	)
	s.Bind(
		mock1.NewV1(),
		mock2.NewV1(),
	)

	// 启动服务器
	s.Run(config.Cfg.App.Port)
}
