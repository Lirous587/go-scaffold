package i18n

import (
	"github.com/pkg/errors"
	"scaffold/pkg/i18n/validatortrans"
	"strings"

	"github.com/go-playground/locales/en"
	"github.com/go-playground/locales/zh"
	ut "github.com/go-playground/universal-translator"
	"github.com/go-playground/validator/v10"
)

var (
	// 全局翻译器映射
	trans map[string]ut.Translator
	uni   *ut.UniversalTranslator
)

// Init 初始化翻译器
func Init() {
	enLocale := en.New()
	zhLocale := zh.New()
	uni = ut.New(enLocale, zhLocale)

	trans = make(map[string]ut.Translator)
	trans["en"], _ = uni.GetTranslator("en")
	trans["zh"], _ = uni.GetTranslator("zh")

	// 设置验证翻译器
	validatortrans.Setup(trans)
}

// SetupValidator 为指定的验证器设置翻译
func SetupValidator(v *validator.Validate) error {
	err := validatortrans.SetupValidator(v)
	if err != nil {
		return errors.WithMessage(err, "validatortrans.SetupValidator(v) failed")
	}
	return nil
}

type ValidatorError map[string]string

// ValidatorError 实现 error 接口
func (v ValidatorError) Error() string {
	if len(v) == 0 {
		return "validation failed"
	}

	// 构建包含所有错误的格式化字符串
	var sb strings.Builder
	first := true
	for field, msg := range v {
		if !first {
			sb.WriteString("; ")
		}
		sb.WriteString(field)
		sb.WriteString(": ")
		sb.WriteString(msg)
		first = false
	}
	return sb.String()
}

// TranslateValidatorError 翻译验证错误
func TranslateValidatorError(err error, lang string) ValidatorError {
	// 转发到子包的验证错误翻译函数
	return validatortrans.TranslateError(err, lang)
}

// GetLanguageFromHeader 从请求头获取语言代码
func GetLanguageFromHeader(acceptLanguage string) string {
	if strings.HasPrefix(acceptLanguage, "en") {
		return "en"
	}
	return "zh"
}
