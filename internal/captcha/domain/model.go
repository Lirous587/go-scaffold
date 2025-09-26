package domain

import (
	"time"
)

type Captcha struct {
	ID    int64
	Image string // 主图片
	Thumb string // 缩略图
	Audio string // 音频验证码
	// 其他类型验证码的响应数据
}

type CaptchaAnswer struct {
	ID    int64
	Image string // 主图片
	Thumb string // 缩略图
	Audio string // 音频验证码
	Way   VerifyWay
	Value string
}

type VerifyWay string

const (
	WayImageClick VerifyWay = "image:click"
)

const (
	defaultExpire      = time.Minute * 1
	captchaClickExpire = time.Minute * 2
)

func (v VerifyWay) GetKey() string {
	return "captcha:" + string(v)
}

func (v VerifyWay) GetExpire() time.Duration {
	switch v {
	case WayImageClick:
		return captchaClickExpire
	default:
		return defaultExpire
	}
}
