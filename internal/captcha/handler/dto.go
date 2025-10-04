package handler

import (
	"scaffold/internal/captcha/domain"
)

type CaptchaResponse struct {
	ID    int64  `json:"id,string"`
	Image string `json:"image,omitempty"` // 主图片
	Thumb string `json:"thumb,omitempty"` // 缩略图
	Audio string `json:"audio,omitempty"` // 音频验证码
	// 其他类型验证码的响应数据
}

type GenRequest struct {
	Way domain.VerifyWay `form:"way"`
}

type CaptchaAnswerResponse struct {
	ID    int64            `json:"id,string"`
	Value string           `json:"value"`
	Way   domain.VerifyWay `json:"way"`
	Image string           `json:"image,omitempty"` // 主图片
	Thumb string           `json:"thumb,omitempty"` // 缩略图
	Audio string           `json:"audio,omitempty"` // 音频验证码
}
