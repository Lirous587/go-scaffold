package cmd

import (
	"scaffold/internal/controller/mock1"
	"scaffold/pkg/httpserver"
)

func Main() {
	// 创建服务器
	s := httpserver.New(8080)

	// 创建根路由组
	s.Group("/", func(group *httpserver.ServerGroup) {
		group.Group("/mock", func(group *httpserver.ServerGroup) {
			group.Bind(
				mock1.NewV1(),
			)
		})
	})

	// 启动服务器
	s.Run()
}
