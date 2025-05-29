package response

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/pkg/errors"
	"scaffold/internal/common/reskit/codes"

	"scaffold/internal/common/validator/i18n"
)

type successResponse struct {
	Code    int         `json:"code"`
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
}

// Success 返回成功响应
func Success(c *gin.Context, data interface{}) {
	// 如果已经响应过，直接返回
	if c.Writer.Written() {
		return
	}

	c.JSON(200, successResponse{
		Code:    2000,
		Message: "Success",
		Data:    data,
	})
}

// Error 返回错误响应
func Error(c *gin.Context, err error) {
	// 如果已经响应过，直接返回
	if c.Writer.Written() {
		return
	}

	// 映射错误
	httpErr := codes.MapToHTTP(err)

	// 将需要日志记录的错误到Gin的错误列表 让后续中间件去记录
	if httpErr.Cause != nil {
		msg := httpErr.Response.Message
		if httpErr.Response.Details != nil {
			msg = fmt.Sprintf("%s | details: %v", msg, httpErr.Response.Details)
		}
		_ = c.Error(errors.WithMessage(httpErr.Cause, msg))
	}

	c.AbortWithStatusJSON(httpErr.StatusCode, httpErr.Response)
}

// ValidationError 返回验证错误响应
func ValidationError(c *gin.Context, err error) {
	// 如果已经响应过，直接返回
	if c.Writer.Written() {
		return
	}

	// 记录错误
	_ = c.Error(err)

	// 翻译验证错误
	validationErrors := i18n.TranslateError(err)

	c.AbortWithStatusJSON(400, codes.HTTPErrorResponse{
		Code:    4000,
		Message: "Validation failed",
		Details: map[string]interface{}{
			"errors": validationErrors,
		},
	})
}
