package bind

import (
	"net/http"
	"reflect"
	"scaffold/pkg/httpserver/core"

	"github.com/gin-gonic/gin"
)

// AutoBind 自动将请求绑定到处理函数的参数，并返回处理结果
func AutoBind(apiMethod interface{}) gin.HandlerFunc {
	apiMethodValue := reflect.ValueOf(apiMethod)
	apiMethodType := apiMethodValue.Type()

	reqType := apiMethodType.In(1)
	reqElemType := reqType.Elem()

	return func(c *gin.Context) {
		// 创建请求参数实例
		reqValue := reflect.New(reqElemType)
		req := reqValue.Interface()

		// 创建绑定器并执行绑定
		lang := c.GetHeader("accept-language")
		binder := NewBinder(c, lang)
		if err := binder.Bind(req); err != nil {
			// 处理绑定错误
			c.AbortWithStatusJSON(http.StatusBadRequest, core.Response{
				Code:    core.CodeValidationError,
				Message: err,
				Data:    nil,
			})
			return
		}

		// 调用处理函数
		args := []reflect.Value{
			reflect.ValueOf(c),
			reqValue,
		}
		results := apiMethodValue.Call(args)

		// 获取返回值
		resValue := results[0]
		errValue := results[1]

		// 处理错误
		if !errValue.IsNil() {
			err := errValue.Interface().(error)
			c.AbortWithStatusJSON(http.StatusInternalServerError, core.Response{
				Code:    core.CodeServerError,
				Message: err.Error(),
				Data:    nil,
			})
			return
		}

		// 返回成功响应
		c.JSON(http.StatusOK, core.Response{
			Code:    core.CodeSuccess,
			Data:    resValue.Interface(),
			Message: "success",
		})
	}
}
