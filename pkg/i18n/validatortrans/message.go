package validatortrans

import (
	ut "github.com/go-playground/universal-translator"
	"github.com/go-playground/validator/v10"
	"github.com/pkg/errors"
)

// 注册自定义翻译
func registerCustomMessages(v *validator.Validate) error {
	// 注册中文自定义翻译
	if t, exists := Trans["zh"]; exists {
		// 手机号验证
		err := registerMobileCNTranslation(v, t, true)
		if err != nil {
			return errors.WithMessage(err, "registerMobileCNTranslation failed")
		}

	}

	// 注册英文自定义翻译
	if t, exists := Trans["en"]; exists {
		// 手机号验证
		err := registerMobileCNTranslation(v, t, false)
		if err != nil {
			return errors.WithMessage(err, "registerMobileCNTranslation failed")
		}

		// 其他验证规则...
	}
	return nil
}

// 注册中国手机号验证翻译
func registerMobileCNTranslation(v *validator.Validate, t ut.Translator, isChinese bool) error {
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
