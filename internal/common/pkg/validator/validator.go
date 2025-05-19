package validator

import (
	"github.com/gin-gonic/gin/binding"
	"github.com/pkg/errors"
	"reflect"
	"scaffold/internal/common/pkg/validator/i18n"
	"scaffold/internal/common/pkg/validator/register"
	"strings"

	"github.com/go-playground/validator/v10"
)

// V 全局验证器实例
var v = validator.New()

// Init 初始化验证器
func Init() error {
	// 1. 注册自定义验证规则
	if err := register.Register(v); err != nil {
		return errors.WithMessage(err, "register.Register(v) failed")
	}

	// 2. 注册结构体标签别名 - 使用蛇形
	v.RegisterTagNameFunc(func(fld reflect.StructField) string {
		return toSnakeCase(fld.Name)
	})

	// 3. 为验证器设置翻译
	err := i18n.SetupValidator(v)
	if err != nil {
		return errors.WithMessage(err, "i18n.SetupValidator failed")
	}

	v.SetTagName("binding")

	if ginV, ok := binding.Validator.Engine().(*validator.Validate); ok {
		*ginV = *v
	}

	return nil
}

func toSnakeCase(str string) string {
	var result []rune
	for i, r := range str {
		// 如果是大写字母
		if r >= 'A' && r <= 'Z' {
			// 不是第一个字符，并且前一个字符不是大写字母，才添加下划线
			if i > 0 && (str[i-1] < 'A' || str[i-1] > 'Z') {
				result = append(result, '_')
			}
		}
		result = append(result, r)
	}
	return strings.ToLower(string(result))
}

var TranslateError = i18n.TranslateError

var GetTranslateLang = i18n.GetTranslateLang
