package cmd

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"scaffold/pkg/httpserver"
	"scaffold/utility/bind"
)

type TestJSON struct {
	Name string `json:"name" validate:"required"`
	Age  int    `json:"age" validate:"required"`
}

func Main() {
	// 创建服务器
	s := httpserver.New(8080)

	s.POST("/test", func(c *gin.Context) {
		var test TestJSON
		if err := bind.Bind(c, &test); err != nil {
			c.AbortWithStatusJSON(400, gin.H{
				"msg":   "参数错误",
				"error": err.Error(),
			})
			fmt.Printf("%#v", err)
			return
		}
		c.JSON(200, gin.H{
			"msg": "hello simple",
		})
	})

	err := s.Run()

	if err != nil {
		panic(err)
	}
}
