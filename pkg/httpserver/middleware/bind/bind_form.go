package bind

import (
	"errors"
	"reflect"
	"scaffold/pkg/i18n"
	"strings"
)

func (b *bind) bindForm(req interface{}) {
	contentType := b.ctx.ContentType()
	if contentType != "application/x-www-form-urlencoded" &&
		!strings.Contains(contentType, "multipart/form-data") {
		b.ctx.Error(errors.New("请求内容类型不是表单: " + contentType))
		b.ctx.Abort()
		return
	}

	reqValue := reflect.ValueOf(req)
	if reqValue.Kind() != reflect.Ptr || reqValue.IsNil() {
		return
	}

	reqElem := reqValue.Elem()
	reqType := reqElem.Type()

	// 根据Content-Type选择合适的解析方法
	if strings.Contains(contentType, "multipart/form-data") {
		// 对于multipart/form-data类型
		if err := b.ctx.Request.ParseMultipartForm(32 << 20); err != nil { // 32MB 最大尺寸
			errorMsg, errorDetail := i18n.TranslateParseFormError(err, b.lang)
			b.ctx.Error(errors.New(errorMsg + ":" + errorDetail))
			b.ctx.Abort()
			return
		}
	} else {
		// 对于x-www-form-urlencoded类型
		if err := b.ctx.Request.ParseForm(); err != nil {
			errorMsg, errorDetail := i18n.TranslateParseFormError(err, b.lang)
			b.ctx.Error(errors.New(errorMsg + ":" + errorDetail))
			b.ctx.Abort()
			return
		}
	}

	// 遍历结构体字段
	for i := 0; i < reqType.NumField(); i++ {
		field := reqType.Field(i)
		fieldVal := reqElem.Field(i)

		if !fieldVal.CanSet() {
			continue
		}

		// 只处理明确标记为表单字段的参数
		inTag := field.Tag.Get("in")
		formTag := field.Tag.Get("form")

		// 只有同时指定了in:"form"和form:"xxx"的才会被处理
		if inTag == "form" && formTag != "" {
			// FormValue可以同时处理普通表单和multipart表单
			value := b.ctx.Request.FormValue(formTag)

			if value == "" {
				continue
			}

			// 设置字段值
			setFieldValue(fieldVal, value)
		}
	}
}
