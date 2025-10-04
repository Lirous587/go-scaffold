package handler

import (
	"scaffold/internal/captcha/domain"
)

func domainCaptchaToResponse(captcha *domain.Captcha) *CaptchaResponse {
	if captcha == nil {
		return nil
	}

	return &CaptchaResponse{
		ID:    captcha.ID,
		Image: captcha.Image,
		Thumb: captcha.Thumb,
		Audio: captcha.Audio,
	}
}

func domainCaptchaAnswerToResponse(captcha *domain.CaptchaAnswer) *CaptchaAnswerResponse {
	if captcha == nil {
		return nil
	}

	return &CaptchaAnswerResponse{
		ID:    captcha.ID,
		Value: captcha.Value,
		Way:   captcha.Way,
		Image: captcha.Image,
		Thumb: captcha.Thumb,
		Audio: captcha.Audio,
	}
}
