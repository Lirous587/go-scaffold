package middleware

import (
	"encoding/json"
	"errors"
	"io"
	"scaffold/pkg/i18n"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

func ErrorHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Next()
		// 获取请求语言
		lang := i18n.GetLanguageFromHeader(c.GetHeader("Accept-Language"))
		logLang := "zh"

		if len(c.Errors) > 0 {
			for i, e := range c.Errors {

				// 处理空请求体错误
				if errors.Is(e.Err, io.EOF) {
					errorMsg, errorDetail := i18n.TranslateNullJSONError(e.Err, lang)
					c.JSON(400, Response{
						Code:    400,
						Message: errorMsg,
						Data:    errorDetail,
					})
					{
						logErrorMsg, _ := i18n.TranslateNullJSONError(e.Err, logLang)
						c.Errors[i] = &gin.Error{Err: errors.New(logErrorMsg)}
					}
					return
				}

				// 处理JSON解析错误
				var jsonError *json.SyntaxError
				var unmarshalTypeError *json.UnmarshalTypeError
				if errors.As(e.Err, &jsonError) || errors.As(e.Err, &unmarshalTypeError) {
					// 使用i18n包翻译JSON错误
					errorMsg, errorDetail := i18n.TranslateJSONError(e.Err, lang)
					c.JSON(400, Response{
						Code:    400,
						Message: errorMsg,
						Data:    errorDetail,
					})
					{
						logErrorMsg, _ := i18n.TranslateJSONError(e.Err, logLang)
						c.Errors[i] = &gin.Error{Err: errors.New(logErrorMsg)}
					}
					return
				}

				// 处理验证错误
				var validationErrors validator.ValidationErrors
				if errors.As(e.Err, &validationErrors) {
					// 使用i18n包翻译错误信息
					errorMsg := i18n.TranslateValidatorError(e.Err, lang)
					c.JSON(400, Response{
						Code:    400,
						Message: errorMsg,
						Data:    nil,
					})
					{
						logError := i18n.TranslateValidatorError(e.Err, logLang)
						logErrStr, _ := json.Marshal(logError)
						c.Errors[i] = &gin.Error{Err: errors.New(string(logErrStr))}
					}
					return
				}

				// 处理其他类型的错误
				errorMsg, _ := i18n.TranslateServerError(e.Err, lang)
				c.JSON(500, Response{
					Code:    500,
					Message: errorMsg,
					Data:    nil,
				})
				{
					logErrorMsg, _ := i18n.TranslateServerError(e.Err, lang)
					c.Errors[i] = &gin.Error{Err: errors.New(logErrorMsg)}
				}

				return
			}
		}
	}
}
