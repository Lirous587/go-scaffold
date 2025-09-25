package register

import (
	ut "github.com/go-playground/universal-translator"
	"github.com/go-playground/validator/v10"
	"github.com/pkg/errors"
	"regexp"
)

// 预编译正则表达式，提升性能
var slugPattern = regexp.MustCompile(`^[a-zA-Z0-9_-]+$`)

// 自定义安全字符验证（字母、数字、下划线、连字符）
func validateSlug(fl validator.FieldLevel) bool {
	value := fl.Field().String()
	return slugPattern.MatchString(value)
}

// 注册slug验证翻译
func (r *RTrans) registerSlugTranslation(v *validator.Validate, t ut.Translator, isChinese bool) error {
	message := "{0} must only contain letters, numbers, underscores or hyphens"
	if isChinese {
		message = "{0}必须仅包含字母、数字、下划线或连字符"
	}

	err := v.RegisterTranslation("slug", t, func(ut ut.Translator) error {
		return ut.Add("slug", message, true)
	}, func(ut ut.Translator, fe validator.FieldError) string {
		t, _ := ut.T("slug", fe.Field())
		return t
	})
	if err != nil {
		return errors.WithMessage(err, "v.RegisterTranslation failed")
	}
	return nil
}
