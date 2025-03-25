package i18n

import (
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

// 初始化翻译器
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
func SetupValidator(v *validator.Validate) {
	validatortrans.SetupValidator(v)
}

// TranslateValidatorError 翻译验证错误
func TranslateValidatorError(err error, lang string) map[string]string {
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
