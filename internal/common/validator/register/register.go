package register

import (
	ut "github.com/go-playground/universal-translator"
	"github.com/go-playground/validator/v10"
	"github.com/pkg/errors"
)

func Register(v *validator.Validate) error {
	var err error
	if err = v.RegisterValidation("mobile_cn", validateChineseMobile); err != nil {
		return errors.WithMessage(err, "register mobile_cn failed")
	}
	if err = v.RegisterValidation("hex_color", validateHexColor); err != nil {
		return errors.WithMessage(err, "register hex_color failed")
	}
	if err = v.RegisterValidation("domain_url", validateDomainURL); err != nil {
		return errors.WithMessage(err, "register domain_url failed")
	}

	return nil
}

type RTrans struct {
	trans map[string]ut.Translator
}

func NewTrans(trans map[string]ut.Translator) *RTrans {
	return &RTrans{
		trans: trans,
	}
}

// RegisterTranslation 注册自定义翻译
func (r *RTrans) RegisterTranslation(v *validator.Validate) error {
	// 注册中文自定义翻译
	if t, exists := r.trans["zh"]; exists {
		// 手机号验证
		if err := r.registerMobileCNTranslation(v, t, true); err != nil {
			return errors.WithMessage(err, "registerMobileCNTranslation failed")
		}
		// 十六进颜色
		if err := r.registerHexColorTranslation(v, t, true); err != nil {
			return errors.WithMessage(err, "registerHexColorTranslation failed")
		}
		// 域名url
		if err := r.registerDomainURLTranslation(v, t, true); err != nil {
			return errors.WithMessage(err, "registerHexColorTranslation failed")
		}
	}

	// 注册英文自定义翻译
	if t, exists := r.trans["en"]; exists {
		// 手机号验证
		if err := r.registerMobileCNTranslation(v, t, false); err != nil {
			return errors.WithMessage(err, "registerMobileCNTranslation failed")
		}
		if err := r.registerHexColorTranslation(v, t, false); err != nil {
			return errors.WithMessage(err, "registerHexColorTranslation failed")
		}
		// 域名url
		if err := r.registerDomainURLTranslation(v, t, false); err != nil {
			return errors.WithMessage(err, "registerHexColorTranslation failed")
		}
	}
	return nil
}
