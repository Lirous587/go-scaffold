package register

import (
	ut "github.com/go-playground/universal-translator"
	"github.com/go-playground/validator/v10"
	"github.com/pkg/errors"
	"regexp"
)

// 自定义中国手机号验证
func validateChineseMobile(fl validator.FieldLevel) bool {
	value := fl.Field().String()
	pattern := `^1[3-9]\d{9}$`
	matched, _ := regexp.MatchString(pattern, value)
	return matched
}

// 注册中国手机号验证翻译
func (r *RTrans) registerMobileCNTranslation(v *validator.Validate, t ut.Translator, isChinese bool) error {
	message := "{0} must be a valid Chinese mainland mobile number"
	if isChinese {
		message = "{0}必须是有效的中国大陆手机号"
	}

	err := v.RegisterTranslation("mobile_cn", t, func(ut ut.Translator) error {
		return ut.Add("mobile_cn", message, true)
	}, func(ut ut.Translator, fe validator.FieldError) string {
		t, _ := ut.T("mobile_cn", fe.Field())
		return t
	})
	if err != nil {
		return errors.WithMessage(err, "v.RegisterTranslation failed")
	}
	return nil
}
