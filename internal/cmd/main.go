package cmd

import (
	"github.com/gin-gonic/gin"
	"scaffold/pkg/httpserver"
	"scaffold/pkg/validator"
)

type TestJSON struct {
	Name  string `json:"name" binding:"required"`
	Age   int    `json:"age" binding:"required"`
	Phone string `json:"phone" binding:"mobile_cn"`
}

func Main() {
	// 创建服务器
	s := httpserver.New(8080)

	r := s.Router

	r.POST("/test", func(c *gin.Context) {
		var test TestJSON
		if err := c.ShouldBindJSON(&test); err != nil {
			// 翻译错误
			errMsg := validator.TranslateError(err, "zh")
			c.AbortWithStatusJSON(400, gin.H{
				"msg":   "参数错误",
				"error": errMsg,
			})
			return
		}
		c.JSON(200, gin.H{
			"msg": "hello",
		})
	})

	s.Run()
}
