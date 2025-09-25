package middlewares

import (
	"bytes"
	"github.com/gin-gonic/gin"
	"io"
	"net/http"
)

func FullRequest() func(c *gin.Context) {
	return func(c *gin.Context) {
		bodyBytes, err := io.ReadAll(c.Request.Body)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
				"msg": "无法读取请求体",
			})
			return
		}
		c.Request.Body = io.NopCloser(bytes.NewBuffer(bodyBytes))
	}
}
