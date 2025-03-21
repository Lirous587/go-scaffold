package cmd

import (
	"github.com/gin-gonic/gin"
)

func RegisterRouter(router *gin.Engine) {
	// 测试路由，确认服务正在运行
	router.GET("/", func(c *gin.Context) {
		c.JSON(200, gin.H{"status": "ok"})
	})

	router.POST("/test", func(c *gin.Context) {
		m := new(mock)
		if err := c.ShouldBind(m); err != nil {
			c.Error(err) // 将错误添加到Gin的错误链
			return
		}

		c.JSON(200, gin.H{
			"msg": m,
		})
	})
}

type mock struct {
	// 基本验证
	Username string `json:"username" binding:"required,min=3,max=20"`

	// 数字范围验证
	Age int `json:"age" binding:"required,gte=18,lte=120"`

	// 条件验证 (如果有email字段，则必须符合email格式)
	Email string `json:"email" binding:"omitempty,email"`

	// 枚举验证
	Role string `json:"role" binding:"required,oneof=admin user guest"`

	// 数组/切片验证
	Tags []string `json:"tags" binding:"dive,min=2,max=10"`

	// 嵌套结构体验证
	Address Address `json:"address" binding:"required"`

	// 自定义验证
	Phone string `json:"phone" binding:"required,mobile_cn" label:"手机号"`
}

type Address struct {
	Street string `json:"street" binding:"required"`
	City   string `json:"city" binding:"required"`
	ZIP    string `json:"zip" binding:"required,len=6"`
}
