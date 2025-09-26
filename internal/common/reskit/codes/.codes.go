package codes

// 错误码范围定义
const (
	// 系统级错误 (1-999)
	SystemErrorStart = 1
	SystemErrorEnd   = 999

	// 用户模块 (1000-1199)
	UserErrorStart = 1000
	UserErrorEnd   = 1199

	// 验证码模块 (1200-1399)
	CaptchaErrorStart = 1200
	CaptchaErrorEnd   = 1399

	// 图库模块 (1400-1599)
	ImgErrorStart = 1400
	ImgErrorEnd   = 1599
)
