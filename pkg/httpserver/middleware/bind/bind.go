package bind

import (
	"context"
	"net/http"
	"reflect"
	"scaffold/pkg/i18n"
	"scaffold/pkg/validator"
	"strings"

	"github.com/gin-gonic/gin"
)

type bind struct {
	ctx     *gin.Context
	lang    string //响应语言
	logLang string //日志语言
}

func (b *bind) smartBind(req interface{}) {
	// 1. 绑定URI参数 (路径参数)
	b.bindUri(req)

	// 2. 绑定请求头参数
	b.bindHeader(req)

	// 3. 绑定查询参数 (对所有HTTP方法)
	b.bindQuery(req)

	// 4. 对于非GET/DELETE方法，处理请求体
	if b.ctx.Request.Method != http.MethodGet && b.ctx.Request.Method != http.MethodDelete {
		// 根据Content-Type处理请求体
		contentType := b.ctx.ContentType()

		if contentType == "application/json" {
			b.bindJson(req)
		} else if contentType == "application/x-www-form-urlencoded" ||
			strings.Contains(contentType, "multipart/form-data") {
			b.bindForm(req)
		}
	}
}

// AutoBind 自动将请求绑定到处理函数的参数，并返回处理结果
func AutoBind(handler interface{}) gin.HandlerFunc {
	handlerValue := reflect.ValueOf(handler)
	handlerType := handlerValue.Type()

	validateHandlerSignature(handlerType)

	return func(c *gin.Context) {
		// 创建请求参数实例
		reqType := handlerType.In(1)
		reqValue := reflect.New(reqType.Elem())
		req := reqValue.Interface()

		lang := i18n.GetLanguageFromHeader(c.GetHeader("Accept-Language"))
		// 自动绑定
		b := &bind{
			ctx:     c,
			lang:    lang,
			logLang: "zh",
		}

		b.smartBind(req)

		//验证参数
		if err := validator.V.Struct(req); err != nil {
			errMsg := i18n.TranslateValidatorError(err, lang)
			c.JSON(400, gin.H{
				"message": errMsg,
			})
			//c.Error(errMsg)
			c.Abort()
			return
		}

		// 调用处理函数
		response := callHandler(c, handlerValue, reqValue)

		// 默认返回成功响应
		c.JSON(http.StatusOK, gin.H{
			"data": response,
		})
	}
}

// 验证处理函数签名
func validateHandlerSignature(handlerType reflect.Type) {
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
}

// 调用处理函数并处理返回值
func callHandler(c *gin.Context, handlerValue reflect.Value, reqValue reflect.Value) interface{} {
	// 准备参数并调用
	args := []reflect.Value{
		reflect.ValueOf(c),
		reqValue,
	}
	results := handlerValue.Call(args)

	// 获取返回值
	resValue := results[0]
	errValue := results[1]

	// 处理错误
	if !errValue.IsNil() {
		err := errValue.Interface().(error)
		c.Error(err)
		c.Abort()
		return nil
	}

	return resValue.Interface()
}
