package bind

import (
	"errors"
	"reflect"
	"scaffold/pkg/i18n"
	"strings"
)

func (b *bind) bindHeader(req interface{}) {
	reqValue := reflect.ValueOf(req)
	if reqValue.Kind() != reflect.Ptr || reqValue.IsNil() {
		return
	}
	reqElem := reqValue.Elem()
	reqType := reqElem.Type()

	// 记录必需的头部字段
	requiredHeaders := make(map[string]string)

	// 遍历所有字段
	for i := 0; i < reqType.NumField(); i++ {
		field := reqType.Field(i)

		// 检查标记为header的字段
		if inTag := field.Tag.Get("in"); inTag == "header" {
			// 获取请求头名称 - 优先使用header标签，否则使用字段名
			headerName := field.Name
			if headerTag := field.Tag.Get("header"); headerTag != "" {
				headerName = headerTag
			}

			// 获取请求头的值
			headerValue := b.ctx.GetHeader(headerName)

			// 设置字段的值
			if headerValue != "" {
				setFieldValue(reqElem.Field(i), headerValue)
			}
		}
	}

	// 处理缺失的必需头部
	if len(requiredHeaders) > 0 {
		missing := make([]string, 0, len(requiredHeaders))
		for fieldName, headerName := range requiredHeaders {
			missing = append(missing, headerName+" ("+fieldName+")")
		}

		lang := i18n.GetLanguageFromHeader(b.ctx.GetHeader("Accept-Language"))

		var errMsg string
		if lang == "en" {
			errMsg = "Required headers missing: " + strings.Join(missing, ", ")
		} else {
			errMsg = "缺少必需的请求头: " + strings.Join(missing, ", ")
		}
		b.ctx.Error(errors.New(errMsg))
		return
	}
}
