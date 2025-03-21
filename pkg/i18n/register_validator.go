package i18n

import (
	ut "github.com/go-playground/universal-translator"
	"github.com/go-playground/validator/v10"

)


// 注册自定义翻译
func registerCustomTranslations(v *validator.Validate) {
	// 注册中文自定义翻译
	if t, exists := trans["zh"]; exists {
		v.RegisterTranslation("mobile_cn", t, func(ut ut.Translator) error {
			return ut.Add("mobile_cn", "{0}必须是有效的中国大陆手机号", true)
		}, func(ut ut.Translator, fe validator.FieldError) string {
			t, _ := ut.T("mobile_cn", fe.Field())
			return t
		})
	}

	// 注册英文自定义翻译
	if t, exists := trans["en"]; exists {
		v.RegisterTranslation("mobile_cn", t, func(ut ut.Translator) error {
			return ut.Add("mobile_cn", "{0} must be a valid Chinese mainland mobile number", true)
		}, func(ut ut.Translator, fe validator.FieldError) string {
			t, _ := ut.T("mobile_cn", fe.Field())
			return t
		})
	}
}