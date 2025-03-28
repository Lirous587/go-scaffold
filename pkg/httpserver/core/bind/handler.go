package bind

import (
	"net/http"
	"reflect"
	"scaffold/response"

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
			c.AbortWithStatusJSON(http.StatusBadRequest, response.Meta{
				Code:    response.CodeValidationError,
				Message: err.Error(),
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
		resValue := results[0]  // api对应的res
		codeValue := results[1] // 状态码

		// 获取状态码
		code := response.Code(codeValue.Int())

		res := response.GetResponse(code)

		if code != response.CodeSuccess {
			if res.HttpStatus == 0 {
				res.HttpStatus = 500
			}
			c.AbortWithStatusJSON(res.HttpStatus, response.Meta{
				Code:    code,
				Message: res.Message,
			})
		} else {
			c.JSON(res.HttpStatus, response.Success{
				Code:    code,
				Message: res.Message,
				Data:    resValue.Interface(),
			})
		}
	}
}
