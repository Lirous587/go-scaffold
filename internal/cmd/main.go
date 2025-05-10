package cmd

import (
	"scaffold/internal/domain/admin"
	"scaffold/pkg/httpserver"
)

func Main() {
	// 创建服务器
	s := httpserver.New(8080)
	r := s.Router

	// 创建 /api 分组
	api := r.Group("/api")

	var err error

	if err = admin.InitV1(api); err != nil {
		panic(err)
	}
	{
		adminWorker := admin.InitWorker()
		adminWorker.Start()
		s.RegisterStopHandler(func() {
			adminWorker.Stop()
		})
	}

	s.Run()
}
