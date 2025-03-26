package bind

import (
	"reflect"
	"strconv"
)

// 设置字段值，根据字段类型转换数据
func setFieldValue(fieldValue reflect.Value, strValue string) {
	if !fieldValue.CanSet() || strValue == "" {
		return
	}

	switch fieldValue.Kind() {
	case reflect.String:
		fieldValue.SetString(strValue)
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		if val, err := strconv.ParseInt(strValue, 10, 64); err == nil {
			fieldValue.SetInt(val)
		}
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		if val, err := strconv.ParseUint(strValue, 10, 64); err == nil {
			fieldValue.SetUint(val)
		}
	case reflect.Bool:
		if val, err := strconv.ParseBool(strValue); err == nil {
			fieldValue.SetBool(val)
		}
	case reflect.Float32, reflect.Float64:
		if val, err := strconv.ParseFloat(strValue, 64); err == nil {
			fieldValue.SetFloat(val)
		}
	}
}
