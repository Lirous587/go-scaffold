package cmd

import (
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"log"
	"scaffold/pkg/config"
	"scaffold/pkg/engine"
	"scaffold/pkg/logger"
)

func Main() {
	if err := config.Init(); err != nil {
		log.Fatalf("配置初始化失败: %v", err)
	}

	if err := logger.InitZap(&config.Cfg.Log); err != nil {
		log.Fatalf("zap日志初始化失败: %v", err)
		return
	}

	defer func() {
		logger.FinishLog()
	}()

	router := engine.Init(&config.Cfg.App)

	router.GET("/", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"msg": "hello",
		})
	})

	router.GET("/warning", func(c *gin.Context) {
		c.JSON(404, gin.H{
			"msg": "warning",
		})
	})
	router.GET("/error", func(c *gin.Context) {
		c.JSON(500, gin.H{
			"msg": "error",
		})
	})

	router.Run(":8080")
}
