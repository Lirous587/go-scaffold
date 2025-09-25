package handler

import (
	"scaffold/internal/captcha/domain"
)

type CaptchaResponse struct {
	ID	int64	`json:"id,string"`
	Image	string	`json:"image,omitempty"`	// 主图片
	Thumb	string	`json:"thumb,omitempty"`	// 缩略图
	Audio	string	`json:"audio,omitempty"`	// 音频验证码
	// 其他类型验证码的响应数据
}

func domainCaptchaToResponse(captcha *domain.Captcha) *CaptchaResponse {
	if captcha == nil {
		return nil
	}

	return &CaptchaResponse{
		ID:	captcha.ID,
		Image:	captcha.Image,
		Thumb:	captcha.Thumb,
		Audio:	captcha.Audio,
	}
}

type GenRequest struct {
	Way domain.VerifyWay `form:"way"`
}

type CaptchaAnswerResponse struct {
	ID	int64			`json:"id,string"`
	Value	string			`json:"value"`
	Way	domain.VerifyWay	`json:"way"`
	Image	string			`json:"image,omitempty"`	// 主图片
	Thumb	string			`json:"thumb,omitempty"`	// 缩略图
	Audio	string			`json:"audio,omitempty"`	// 音频验证码
}

func domainCaptchaAnswerToResponse(captcha *domain.CaptchaAnswer) *CaptchaAnswerResponse {
	if captcha == nil {
		return nil
	}

	return &CaptchaAnswerResponse{
		ID:	captcha.ID,
		Value:	captcha.Value,
		Way:	captcha.Way,
		Image:	captcha.Image,
		Thumb:	captcha.Thumb,
		Audio:	captcha.Audio,
	}
}
