package response

import (
	"github.com/gin-gonic/gin"

	commonErrors "scaffold/internal/common/errors"
	"scaffold/internal/common/validator/i18n"
)

// successResponse 成功响应结构
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
	// 记录错误到 Gin 的错误列表（用于日志中间件）
	c.Error(err)

	// 如果已经响应过，直接返回
	if c.Writer.Written() {
		return
	}

	// 映射错误
	httpErr := commonErrors.MapToHTTP(err)

	c.AbortWithStatusJSON(httpErr.StatusCode, httpErr.Response)
}

// ValidationError 返回验证错误响应
func ValidationError(c *gin.Context, err error) {
	// 记录错误
	c.Error(err)

	// 如果已经响应过，直接返回
	if c.Writer.Written() {
		return
	}

	// 翻译验证错误
	validationErrors := i18n.TranslateError(err)

	c.AbortWithStatusJSON(400, commonErrors.HTTPErrorResponse{
		Code:    4000,
		Message: "Validation failed",
		Details: map[string]interface{}{
			"validation_errors": validationErrors,
		},
	})
}
