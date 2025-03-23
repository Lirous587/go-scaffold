package apigen

import (
	"fmt"
	"reflect"
	"strings"
)

// 参数信息
type ParamInfo struct {
	Name        string // 参数名称
	In          string // 参数位置：path, query, body, header
	Type        string // 参数类型
	Required    bool   // 是否必填
	Description string // 参数描述
}

// ExtractParameters 从结构体类型中提取参数信息
func ExtractParameters(reqType reflect.Type) []ParamInfo {
	params := []ParamInfo{}

	// 遍历所有字段
	for i := 0; i < reqType.NumField(); i++ {
		field := reqType.Field(i)

		// 跳过Meta字段
		if field.Name == "Meta" {
			continue
		}

		// 确定参数位置和名称
		paramIn := "body" // 默认在请求体
		paramName := field.Tag.Get("json")

		// 检查uri标签
		if uriTag := field.Tag.Get("uri"); uriTag != "" {
			paramIn = "path"
			paramName = uriTag
		} else if queryTag := field.Tag.Get("query"); queryTag != "" {
			paramIn = "query"
			paramName = queryTag
		} else if headerTag := field.Tag.Get("header"); headerTag != "" {
			paramIn = "header"
			paramName = headerTag
		}

		// 确定是否必填
		required := strings.Contains(field.Tag.Get("binding"), "required")

		// 添加参数信息
		params = append(params, ParamInfo{
			Name:        paramName,
			In:          paramIn,
			Type:        field.Type.String(),
			Required:    required,
			Description: fmt.Sprintf("%s字段", field.Name), // 默认描述
		})
	}

	return params
}
