package bind

import (
	"reflect"
)

func (b *bind) bindUri(req interface{}) {
	reqValue := reflect.ValueOf(req)
	if reqValue.Kind() != reflect.Ptr || reqValue.IsNil() {
		return
	}
	reqElem := reqValue.Elem()
	reqType := reqElem.Type()

	// 获取参数路径
	params := b.ctx.Params

	for i := 0; i < reqType.NumField(); i++ {
		field := reqType.Field(i)
		fieldVal := reqElem.Field(i)

		if !fieldVal.CanSet() {
			continue
		}

		// 拿到uri标签
		uriTag := field.Tag.Get("uri")
		if uriTag == "" {
			continue
		}

		var paramValue string
		for _, param := range params {
			if param.Key == uriTag {
				paramValue = param.Value
				break
			}
		}

		if paramValue == "" {
			continue
		}

		setFieldValue(fieldVal, paramValue)
	}
}
