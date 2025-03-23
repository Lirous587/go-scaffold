package middleware

import (
	"context"
	"net/http"
	"reflect"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

// 创建验证器
var validate = validator.New()

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

		// 绑定所有来源的参数，忽略错误
		c.ShouldBindUri(req)

		if c.Request.Method == "GET" {
			c.ShouldBindQuery(req)
		} else if c.ContentType() == "application/json" {
			c.ShouldBindJSON(req)
		} else {
			c.ShouldBind(req)
		}
		// 最后再进行一次完整验证
		if err := validate.Struct(req); err != nil {
			c.Error(err)
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
