package validator

import (
	"github.com/pkg/errors"
	"reflect"
	"regexp"
	"scaffold/pkg/i18n"
	"strings"

	"github.com/go-playground/validator/v10"
)

// V 全局验证器实例
var V = validator.New()

// Init 初始化验证器
func Init() error {
	// 注册自定义验证规则
	_ = V.RegisterValidation("mobile_cn", validateChineseMobile)

	// 注册结构体标签别名 - 使用更友好的字段名
	V.RegisterTagNameFunc(func(fld reflect.StructField) string {
		name := fld.Tag.Get("label")
		if name == "" {
			// 如果没有label标签，使用蛇形命名转换
			name = toSnakeCase(fld.Name)
		}
		return name
	})

	// 为验证器设置翻译
	err := i18n.SetupValidator(V)
	if err != nil {
		return errors.WithMessage(err, "i18n.SetupValidator failed")
	}

	return nil
}

// 自定义中国手机号验证
func validateChineseMobile(fl validator.FieldLevel) bool {
	value := fl.Field().String()
	pattern := `^1[3-9]\d{9}$`
	matched, _ := regexp.MatchString(pattern, value)
	return matched
}

// toSnakeCase 将字符串转为蛇形 -> snake_case
func toSnakeCase(str string) string {
	var result []rune
	for i, r := range str {
		if i > 0 && r >= 'A' && r <= 'Z' {
			result = append(result, '_')
		}
		result = append(result, r)
	}
	return strings.ToLower(string(result))
}
