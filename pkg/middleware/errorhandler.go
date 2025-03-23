package middleware

import (
	"encoding/json"
	"errors"
	"io"
	"scaffold/pkg/i18n"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

// ErrorHandler 返回统一的错误处理中间件
func ErrorHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Next()

		if len(c.Errors) > 0 {
			for _, e := range c.Errors {
				// 获取请求语言
				lang := i18n.GetLanguageFromHeader(c.GetHeader("Accept-Language"))

				// 处理验证错误
				var validationErrors validator.ValidationErrors
				if errors.As(e.Err, &validationErrors) {
					// 使用i18n包翻译错误信息
					errorMessages := i18n.TranslateError(e.Err, lang)

					c.JSON(400, gin.H{
						"code":    400,
						"message": "请求参数错误",
						"errors":  errorMessages,
					})
					return
				}

				// 处理空请求体错误
				if errors.Is(e.Err, io.EOF) {
					c.JSON(400, gin.H{
						"code":    400,
						"message": "请求体为空或格式错误",
					})
					return
				}

				// 处理JSON解析错误
				var jsonError *json.SyntaxError
				var unmarshalTypeError *json.UnmarshalTypeError
				if errors.As(e.Err, &jsonError) || errors.As(e.Err, &unmarshalTypeError) {
					// 使用i18n包翻译JSON错误
					lang := i18n.GetLanguageFromHeader(c.GetHeader("Accept-Language"))
					errorMsg, errorDetail := i18n.TranslateJSONError(e.Err, lang)
					
					c.JSON(400, gin.H{
						"code":    400,
						"message": errorMsg,
						"error":   errorDetail,
					})
					return
				}

				// 处理其他类型的错误
				c.JSON(500, gin.H{
					"code":    500,
					"message": "服务器内部错误",
					"error":   e.Err.Error(),
				})
				return
			}
		}
	}
}
