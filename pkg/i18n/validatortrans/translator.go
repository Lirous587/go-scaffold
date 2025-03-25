package validatortrans

import (
    "github.com/go-playground/validator/v10"
    ut "github.com/go-playground/universal-translator"
    en_translations "github.com/go-playground/validator/v10/translations/en"
    zh_translations "github.com/go-playground/validator/v10/translations/zh"
)

// 导出从父包访问的翻译器映射
var Trans map[string]ut.Translator

// Setup 初始化验证翻译器
func Setup(translators map[string]ut.Translator) {
    Trans = translators
}

// SetupValidator 为指定的验证器设置翻译
func SetupValidator(v *validator.Validate) {
    // 注册标准翻译
    if t, exists := Trans["en"]; exists {
        en_translations.RegisterDefaultTranslations(v, t)
    }
    
    if t, exists := Trans["zh"]; exists {
        zh_translations.RegisterDefaultTranslations(v, t)
    }

    // 注册自定义翻译
    registerCustomMessages(v)
}

// TranslateError 翻译验证错误
func TranslateError(err error, lang string) map[string]string {
    return TranslateValidatorError(err, lang)
}