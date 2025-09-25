package domain

type CaptchaService interface {
	GetVerifyWay() VerifyWay
	Generate() (*Captcha, string, error)
}
