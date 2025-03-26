package bind

import (
	"reflect"
)

func (b *bind) bindQuery(req interface{}) {
	reqValue := reflect.ValueOf(req)
	if reqValue.Kind() != reflect.Ptr || reqValue.IsNil() {
		return
	}
	reqElem := reqValue.Elem()
	reqType := reqElem.Type()

	// 获取查询参数
	// type queryValues map[string][]string
	queryValues := b.ctx.Request.URL.Query()

	for i := 0; i < reqType.NumField(); i++ {
		field := reqType.Field(i)
		fieldVal := reqElem.Field(i)
		if !fieldVal.CanSet() {
			continue
		}

		// 只处理明确标记为查询参数的字段
		inTag := field.Tag.Get("in")
		queryTag := field.Tag.Get("query")

		// 只有同时指定了in:"query"和query:"xxx"的才会被处理
		if inTag == "query" && queryTag != "" {
			values, exists := queryValues[queryTag]
			if !exists || len(values) == 0 {
				continue
			}

			// values是一个字符串数组，因为HTTP可能接收同名参数多次，如 ?id=1&id=2
			// values[0] 取第一个值（大多数API场景只使用第一个值）
			setFieldValue(fieldVal, values[0])
		}
	}
}
