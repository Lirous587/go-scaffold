package response

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/pkg/errors"

	"scaffold/internal/common/validator/i18n"
)

type successResponse struct {
	Code	int		`json:"code" example:"2000"`
	Message	string		`json:"message" example:"请求成功"`
	Data	interface{}	`json:"data,omitempty"`
}

// Success 返回成功响应
func Success(ctx *gin.Context, data ...interface{}) {
	// 如果已经响应过，直接返回
	if ctx.Writer.Written() {
		return
	}

	if data != nil {
		ctx.JSON(200, successResponse{
			Code:		2000,
			Message:	"请求成功",
			Data:		data[0],
		})
		return
	}

	ctx.JSON(200, successResponse{
		Code:		2000,
		Message:	"请求成功",
	})
}

// Error 返回错误响应
func Error(ctx *gin.Context, err error) {
	// 如果已经响应过，直接返回
	if ctx.Writer.Written() {
		return
	}

	// 映射错误
	httpErr := MapToHTTP(err)
	msg := httpErr.Response.Message
	if httpErr.Response.Details != nil {
		msg = fmt.Sprintf("%s | details: %v", msg, httpErr.Response.Details)
	}

	// 将需要日志记录的错误到Gin的错误列表 让后续中间件去记录
	if httpErr.Cause != nil {
		_ = ctx.Error(errors.WithMessage(httpErr.Cause, msg))
	} else {
		_ = ctx.Error(errors.New(msg))
	}

	ctx.AbortWithStatusJSON(httpErr.StatusCode, httpErr.Response)
}

// InvalidParams 返回验证错误响应
func InvalidParams(ctx *gin.Context, err error) {
	// 如果已经响应过，直接返回
	if ctx.Writer.Written() {
		return
	}

	// 记录错误
	_ = ctx.Error(err)

	// 翻译验证错误
	validationErrors := i18n.TranslateError(err)

	ctx.AbortWithStatusJSON(400, HTTPErrorResponse{
		Code:		4000,
		Message:	"invalid params",
		Details: map[string]interface{}{
			"errors": validationErrors,
		},
	})
}
