package cmd

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"scaffold/pkg/httpserver"
)

type TestUri interface {
}

type TestHeader interface {
}

type TestQuery interface {
}

type TestJSON struct {
}

type TestForm struct {
}

func Main() {
	// 创建服务器
	s := httpserver.New(8080)

	s.GET("/test", func(c *gin.Context) {
		fmt.Println("simple is good")
		c.JSON(200, gin.H{
			"msg": "hello simple",
		})
	})

	err := s.Run()

	if err != nil {
		panic(err)
	}
}
