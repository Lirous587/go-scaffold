package cmd

import (
	"scaffold/internal/domain/user"
	"scaffold/pkg/httpserver"
)

func Main() {
	// 创建服务器
	s := httpserver.New(8080)
	r := s.Router

	// 创建 /api 分组
	api := r.Group("/api")

	var err error

	if err = user.InitV1(api); err != nil {
		panic(err)
	}

	s.Run()
}
