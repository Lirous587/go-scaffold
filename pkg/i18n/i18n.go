package i18n

import (
	"strings"

	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/locales/en"
	"github.com/go-playground/locales/zh"
	ut "github.com/go-playground/universal-translator"
	"github.com/go-playground/validator/v10"
	en_translations "github.com/go-playground/validator/v10/translations/en"
	zh_translations "github.com/go-playground/validator/v10/translations/zh"
)

var (
	trans map[string]ut.Translator
	uni   *ut.UniversalTranslator
)

// Setup 初始化翻译器
func Setup() {
	enLocale := en.New()
	zhLocale := zh.New()
	uni = ut.New(enLocale, zhLocale)

	trans = make(map[string]ut.Translator)

	// 注册英语翻译器
	trans["en"], _ = uni.GetTranslator("en")

	// 注册中文翻译器
	trans["zh"], _ = uni.GetTranslator("zh")

	// 获取验证器并注册翻译
	if v, ok := binding.Validator.Engine().(*validator.Validate); ok {
		en_translations.RegisterDefaultTranslations(v, trans["en"])
		zh_translations.RegisterDefaultTranslations(v, trans["zh"])

		// 注册自定义翻译
		registerCustomTranslations(v)
	}
}

// GetLanguageFromHeader 从请求头获取语言代码
func GetLanguageFromHeader(acceptLanguage string) string {
	if strings.HasPrefix(acceptLanguage, "en") {
		return "en"
	}
	return "zh"
}

