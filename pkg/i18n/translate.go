package i18n

import (
	"errors"
	"github.com/go-playground/validator/v10"
)

// Translate 翻译验证错误
func TranslateError(err error, lang string) map[string]string {
	var validationErrors validator.ValidationErrors
	if !errors.As(err, &validationErrors) {
		return map[string]string{"error": err.Error()}
	}

	// 如果未指定语言或不支持，则默认使用中文
	if _, exists := trans[lang]; !exists {
		lang = "zh"
	}

	translator := trans[lang]
	errorMessages := make(map[string]string)

	for _, err := range validationErrors {
		errorMessages[err.Field()] = err.Translate(translator)
	}

	return errorMessages
}
