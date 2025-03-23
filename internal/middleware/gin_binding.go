package middleware

import (
	"context"
	"net/http"
	"reflect"

	"github.com/gin-gonic/gin"
)

// AutoBind 自动将请求绑定到处理函数的参数，并返回处理结果
func AutoBind(handler interface{}) gin.HandlerFunc {
	handlerValue := reflect.ValueOf(handler)
	handlerType := handlerValue.Type()

	// 验证handler是函数
	if handlerType.Kind() != reflect.Func {
		panic("handler must be a function")
	}

	// 验证函数参数格式：func(ctx context.Context, req *SomeReq) (res *SomeRes, err error)
	if handlerType.NumIn() != 2 || handlerType.NumOut() != 2 {
		panic("handler must have signature: func(ctx context.Context, req *SomeReq) (*SomeRes, error)")
	}

	// 检查第一个参数是否为context.Context
	if !handlerType.In(0).Implements(reflect.TypeOf((*context.Context)(nil)).Elem()) {
		panic("first parameter must be context.Context")
	}

	// 检查第二个参数是否为指针类型
	reqType := handlerType.In(1)
	if reqType.Kind() != reflect.Ptr {
		panic("second parameter must be a pointer")
	}

	// 检查第二个返回值是否为error
	if !handlerType.Out(1).Implements(reflect.TypeOf((*error)(nil)).Elem()) {
		panic("second return value must be error")
	}

	return func(c *gin.Context) {
		// 创建请求参数实例
		reqValue := reflect.New(reqType.Elem())
		req := reqValue.Interface()

		// 尝试绑定适合的参数来源
		var bindErr error

		// 根据请求方法和Content-Type选择绑定方式
		if c.Request.Method == "GET" {
			// GET请求优先绑定查询参数
			bindErr = c.ShouldBindQuery(req)
			if bindErr == nil {
				// 再尝试绑定URI参数
				bindErr = c.ShouldBindUri(req)
			}
		} else if c.Request.Method == "POST" || c.Request.Method == "PUT" || c.Request.Method == "PATCH" {
			// POST/PUT/PATCH请求优先绑定JSON
			contentType := c.ContentType()
			if contentType == "application/json" {
				bindErr = c.ShouldBindJSON(req)
			} else if contentType == "application/x-www-form-urlencoded" || contentType == "multipart/form-data" {
				bindErr = c.ShouldBind(req)
			} else {
				// 默认尝试JSON绑定
				bindErr = c.ShouldBindJSON(req)
			}
		} else {
			// 其他方法使用通用绑定
			bindErr = c.ShouldBind(req)
		}

		// 无论哪种绑定方式，始终绑定URL参数
		if uriErr := c.ShouldBindUri(req); uriErr != nil && bindErr == nil {
			// 只有在主要绑定成功，但URI参数绑定失败时才报告URI错误
			bindErr = uriErr
		}

		// 处理绑定错误
		if bindErr != nil {
			c.Error(bindErr)
			return
		}

		// 调用处理函数
		args := []reflect.Value{
			reflect.ValueOf(c), // 上下文
			reqValue,           // 请求参数
		}
		results := handlerValue.Call(args)

		// 获取返回值
		resValue := results[0]
		errValue := results[1]

		// 处理错误
		if !errValue.IsNil() {
			err := errValue.Interface().(error)
			c.Error(err)
			return
		}

		// 处理成功响应
		if !resValue.IsNil() {
			res := resValue.Interface()
			c.JSON(http.StatusOK, res)
		} else {
			c.Status(http.StatusNoContent)
		}
	}
}
