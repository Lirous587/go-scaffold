package i18n

import (
	"scaffold/internal/common/validator/register"
	"strings"

	"github.com/gin-gonic/gin"
	entranslations "github.com/go-playground/validator/v10/translations/en"
	zhtranslations "github.com/go-playground/validator/v10/translations/zh"
	"github.com/pkg/errors"

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

func init() {
	enLocale := en.New()
	zhLocale := zh.New()
	uni = ut.New(enLocale, zhLocale)

	trans = make(map[string]ut.Translator)
	trans["en"], _ = uni.GetTranslator("en")
	trans["zh"], _ = uni.GetTranslator("zh")
}

// SetupValidator 为指定的验证器设置翻译
func SetupValidator(v *validator.Validate) error {
	// 注册标准翻译
	if t, exists := trans["en"]; exists {
		err := entranslations.RegisterDefaultTranslations(v, t)
		if err != nil {
			return errors.WithMessage(err, "entranslations.RegisterDefaultTranslations(v, t) failed")
		}
	}

	if t, exists := trans["zh"]; exists {
		err := zhtranslations.RegisterDefaultTranslations(v, t)
		if err != nil {
			return errors.WithMessage(err, "zhtranslations.RegisterDefaultTranslations(v, t) failed")
		}
	}

	rTrans := register.NewTrans(trans)
	if err := rTrans.RegisterTranslation(v); err != nil {
		return errors.WithMessage(err, "rTrans.RegisterTranslation(v) failed")
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

// TranslateError 翻译验证错误
func TranslateError(err error, lang ...string) ValidatorError {
	// 如果错误为nil，返回空映射
	if err == nil {
		return ValidatorError{}
	}

	// 确定使用哪种语言
	language := "zh" // 默认中文
	if len(lang) > 0 {
		language = lang[0]
	}
	// 获取相应语言的翻译器
	t, exists := trans[language]
	if !exists {
		// 如果找不到指定语言的翻译器，默认使用中文
		t = trans["zh"]
	}

	// 尝试将错误转换为validator.ValidationErrors类型
	var validationErrors validator.ValidationErrors
	ok := errors.As(err, &validationErrors)
	if !ok {
		// 如果不是验证错误，直接返回原始错误信息
		return ValidatorError{"error": err.Error()}
	}

	// 翻译错误并构建ValidatorError
	result := ValidatorError{}
	for _, e := range validationErrors {
		// 使用翻译器翻译错误
		result[e.Field()] = e.Translate(t)
	}

	return result
}

func GetTranslateLang(ctx *gin.Context) string {
	acceptLang := ctx.GetHeader("Accept-Language")

	// 转换为小写并分割
	acceptLang = strings.ToLower(acceptLang)

	// 检查是否包含英文
	if strings.Contains(acceptLang, "en") {
		return "en"
	}

	// 默认中文
	return "zh"
}
