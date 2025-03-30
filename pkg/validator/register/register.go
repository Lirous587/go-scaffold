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
		err := r.registerMobileCNTranslation(v, t, true)
		if err != nil {
			return errors.WithMessage(err, "registerMobileCNTranslation failed")
		}
	}

	// 注册英文自定义翻译
	if t, exists := r.trans["en"]; exists {
		// 手机号验证
		err := r.registerMobileCNTranslation(v, t, false)
		if err != nil {
			return errors.WithMessage(err, "registerMobileCNTranslation failed")
		}
	}
	return nil
}
