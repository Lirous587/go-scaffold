package bind

import (
	"encoding/json"
	"errors"
	"fmt"
	"strings"
)

// translateJSONError 翻译JSON解析错误
func translateJSONError(err error, lang string) (string, string) {
	// 默认消息
	errorMsg := "JSON 格式错误"
	errorDetail := err.Error()

	if lang == "en" {
		errorMsg = "JSON Format Error"
	}

	// 提取错误类型信息
	var jsonError *json.SyntaxError
	var unmarshalTypeError *json.UnmarshalTypeError

	if errors.As(err, &unmarshalTypeError) {
		// 从路径提取字段名
		parts := strings.Split(unmarshalTypeError.Field, ".")
		fieldName := parts[len(parts)-1]
		expectedType := unmarshalTypeError.Type.String()

		// 根据语言和错误类型提供友好消息
		if lang == "en" {
			if strings.Contains(expectedType, "[]") && strings.Contains(errorDetail, "string") {
				errorDetail = fmt.Sprintf("Field '%s' requires an array, but a string was provided", fieldName)
			} else if strings.Contains(expectedType, "int") && strings.Contains(errorDetail, "string") {
				errorDetail = fmt.Sprintf("Field '%s' requires a number, but a string was provided", fieldName)
			} else if strings.Contains(expectedType, "string") && strings.Contains(errorDetail, "number") {
				errorDetail = fmt.Sprintf("Field '%s' requires a string, but a number was provided", fieldName)
			} else if strings.Contains(expectedType, "bool") {
				errorDetail = fmt.Sprintf("Field '%s' requires a boolean value (true/false)", fieldName)
			} else {
				errorDetail = fmt.Sprintf("Field '%s' has incorrect type, expected %s",
					fieldName, formatTypeNameI18n(expectedType, lang))
			}
		} else {
			// 中文翻译
			if strings.Contains(expectedType, "[]") && strings.Contains(errorDetail, "string") {
				errorDetail = fmt.Sprintf("字段 '%s' 需要数组格式，但提供了字符串", fieldName)
			} else if strings.Contains(expectedType, "int") && strings.Contains(errorDetail, "string") {
				errorDetail = fmt.Sprintf("字段 '%s' 需要数字，但提供了字符串", fieldName)
			} else if strings.Contains(expectedType, "string") && strings.Contains(errorDetail, "number") {
				errorDetail = fmt.Sprintf("字段 '%s' 需要字符串，但提供了数字", fieldName)
			} else if strings.Contains(expectedType, "bool") {
				errorDetail = fmt.Sprintf("字段 '%s' 需要布尔值(true/false)", fieldName)
			} else {
				errorDetail = fmt.Sprintf("字段 '%s' 类型不正确，期望 %s 类型",
					fieldName, formatTypeNameI18n(expectedType, lang))
			}
		}
	} else if errors.As(err, &jsonError) {
		// 处理JSON语法错误
		if lang == "en" {
			errorDetail = "JSON syntax error, please check your format (brackets, quotes, commas, etc.)"
		} else {
			errorDetail = "JSON语法错误，请检查格式是否正确，如括号、引号、逗号等"
		}
	}

	return errorMsg, errorDetail
}

// formatTypeNameI18n 根据语言将Go类型名转换为更友好的名称
func formatTypeNameI18n(typeName string, lang string) string {
	if lang == "en" {
		switch typeName {
		case "[]string":
			return "string array"
		case "[]int", "[]int64":
			return "integer array"
		case "[]float64":
			return "decimal array"
		case "string":
			return "string"
		case "int", "int64":
			return "integer"
		case "float64":
			return "decimal"
		case "bool":
			return "boolean"
		case "map[string]interface {}":
			return "object"
		case "[]interface {}":
			return "array"
		default:
			return typeName
		}
	} else {
		// 中文
		switch typeName {
		case "[]string":
			return "字符串数组"
		case "[]int", "[]int64":
			return "整数数组"
		case "[]float64":
			return "浮点数数组"
		case "string":
			return "字符串"
		case "int", "int64":
			return "整数"
		case "float64":
			return "浮点数"
		case "bool":
			return "布尔值"
		case "map[string]interface {}":
			return "对象"
		case "[]interface {}":
			return "数组"
		default:
			return typeName
		}
	}
}
