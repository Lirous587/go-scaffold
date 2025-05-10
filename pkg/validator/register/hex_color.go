package register

import (
	ut "github.com/go-playground/universal-translator"
	"github.com/go-playground/validator/v10"
	"github.com/pkg/errors"
	"regexp"
)

// 自定义十六进制颜色
func validateHexColor(fl validator.FieldLevel) bool {
	value := fl.Field().String()
	pattern := `^#([A-Fa-f0-9]{6}|[A-Fa-f0-9]{3})$`
	matched, _ := regexp.MatchString(pattern, value)
	return matched
}

// 十六进制颜色验证翻译
func (r *RTrans) registerHexColorTranslation(v *validator.Validate, t ut.Translator, isChinese bool) error {
	message := "{0} must be a valid hex color,such as #FFF or #F5A942"
	if isChinese {
		message = "{0}必须是有效的十六进颜色,例如 #FFF 或 #F5A942"
	}

	err := v.RegisterTranslation("hex_color", t, func(ut ut.Translator) error {
		return ut.Add("hex_color", message, true)
	}, func(ut ut.Translator, fe validator.FieldError) string {
		t, _ := ut.T("hex_color", fe.Field())
		return t
	})
	if err != nil {
		return errors.WithMessage(err, "v.RegisterTranslation failed")
	}
	return nil
}
