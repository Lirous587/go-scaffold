package validatortrans

import (
	ut "github.com/go-playground/universal-translator"
	"github.com/go-playground/validator/v10"
	entranslations "github.com/go-playground/validator/v10/translations/en"
	zhtranslations "github.com/go-playground/validator/v10/translations/zh"
	"github.com/pkg/errors"
)

// Trans 导出从父包访问的翻译器映射
var Trans map[string]ut.Translator

// Setup 初始化验证翻译器
func Setup(translators map[string]ut.Translator) {
	Trans = translators
}

// SetupValidator 为指定的验证器设置翻译
func SetupValidator(v *validator.Validate) error {
	// 注册标准翻译
	if t, exists := Trans["en"]; exists {
		err := entranslations.RegisterDefaultTranslations(v, t)
		if err != nil {
			return errors.WithMessage(err, "entranslations.RegisterDefaultTranslations(v, t) failed")
		}
	}

	if t, exists := Trans["zh"]; exists {
		err := zhtranslations.RegisterDefaultTranslations(v, t)
		if err != nil {
			return errors.WithMessage(err, "nil.RegisterDefaultTranslations(v, t) failed")
		}
	}

	// 注册自定义翻译
	err := registerCustomMessages(v)
	if err != nil {
		return errors.WithMessage(err, "registerCustomMessages failed")
	}
	return nil
}

// TranslateError 翻译验证错误
func TranslateError(err error, lang string) map[string]string {
	return TranslateValidatorError(err, lang)
}
