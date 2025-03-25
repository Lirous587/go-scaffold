package cmd

import (
	"scaffold/internal/controller/user"
	"scaffold/pkg/config"
	"scaffold/pkg/httpserver"
)

func Main() {
	// 创建服务器
	s := httpserver.New()

	// 创建根路由组
	s.Group("/", func(group *httpserver.ServerGroup) {

		// 创建v1版本路由组
		group.Group("/v1", func(group *httpserver.ServerGroup) {
			//group.Middleware(func(ctx *gin.Context) {
			//	fmt.Println("fuck you")
			//})
			group.Bind(
				user.NewV1(),
			)
		})
	})

	// 启动服务器
	s.Run(config.Cfg.App.Port)
}
