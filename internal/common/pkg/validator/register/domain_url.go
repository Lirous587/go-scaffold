package register

import (
	ut "github.com/go-playground/universal-translator"
	"github.com/go-playground/validator/v10"
	"github.com/pkg/errors"
	"net"
	"net/url"
	"regexp"
)

// 自定义域名url
func validateDomainURL(fl validator.FieldLevel) bool {
	value := fl.Field().String()
	u, err := url.Parse(value)
	if err != nil || u.Scheme == "" || u.Host == "" {
		return false
	}

	// 检查是否为IP地址
	if net.ParseIP(u.Host) != nil {
		return false
	}

	pattern := `^([a-zA-Z0-9]([a-zA-Z0-9\-]{0,61}[a-zA-Z0-9])?\.)+[a-zA-Z]{2,}$`

	return regexp.MustCompile(pattern).MatchString(u.Host)
}

// 域名url验证翻译
func (r *RTrans) registerDomainURLTranslation(v *validator.Validate, t ut.Translator, isChinese bool) error {
	message := "{0} must be a valid domain URL,such as https://lirous.com"
	if isChinese {
		message = "{0}必须是有效的域名url,例如https://lirous.com"
	}

	err := v.RegisterTranslation("domain_url", t, func(ut ut.Translator) error {
		return ut.Add("domain_url", message, true)
	}, func(ut ut.Translator, fe validator.FieldError) string {
		t, _ := ut.T("domain_url", fe.Field())
		return t
	})
	if err != nil {
		return errors.WithMessage(err, "v.RegisterTranslation failed")
	}
	return nil
}
