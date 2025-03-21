package validator

import (
	"reflect"
	"regexp"

	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/validator/v10"
)

// Setup 初始化验证器
func Setup() {
	if v, ok := binding.Validator.Engine().(*validator.Validate); ok {
		// 注册自定义验证规则
		_ = v.RegisterValidation("mobile_cn", validateChineseMobile)

		// 注册结构体标签别名
		v.RegisterTagNameFunc(func(fld reflect.StructField) string {
			name := fld.Tag.Get("label")
			if name == "" {
				name = fld.Name
			}
			return name
		})
	}
}

// 自定义中国手机号验证
func validateChineseMobile(fl validator.FieldLevel) bool {
	value := fl.Field().String()
	pattern := `^1[3-9]\d{9}$`
	matched, _ := regexp.MatchString(pattern, value)
	return matched
}
