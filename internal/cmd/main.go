package cmd

import (
	"scaffold/internal/controller/user"

	"github.com/gin-gonic/gin"
)

func RegisterRouter(router *gin.Engine) {
	user.NewV1(router)
}
