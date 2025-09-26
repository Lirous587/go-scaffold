package codes

// 验证码模块错误码 (1200-1399)
var (
	// 生成类错误 (1200-1219)
	ErrCaptchaGenerateFailed = ErrCode{Msg: "验证码生成失败", Type: ErrorTypeInternal, Code: 1201}
	ErrCaptchaImageEmpty     = ErrCode{Msg: "验证码图片为空", Type: ErrorTypeInternal, Code: 1202}

	// 验证类错误 (1220-1239)
	ErrCaptchaVerifyFailed = ErrCode{Msg: "验证码验证失败", Type: ErrorTypeUnauthorized, Code: 1220}
	ErrCaptchaInvalid      = ErrCode{Msg: "验证码无效", Type: ErrorTypeUnauthorized, Code: 1221}
	ErrCaptchaExpired      = ErrCode{Msg: "验证码已过期", Type: ErrorTypeUnauthorized, Code: 1222}
	ErrCaptchaNotFound     = ErrCode{Msg: "验证码不存在", Type: ErrorTypeNotFound, Code: 1223}

	// 参数类错误 (1240-1259)
	ErrCaptchaFormatInvalid = ErrCode{Msg: "验证码格式错误", Type: ErrorTypeValidation, Code: 1240}
	ErrCaptchaHeaderMissing = ErrCode{Msg: "缺少验证码请求头", Type: ErrorTypeValidation, Code: 1241}

	// 存储类错误 (1260-1279)
	ErrCaptchaCacheError = ErrCode{Msg: "验证码缓存操作失败", Type: ErrorTypeInternal, Code: 1260}

	// 频率限制类错误 (1280-1299)
	ErrCaptchaRateLimit = ErrCode{Msg: "验证码请求过于频繁", Type: ErrorTypeRateLimit, Code: 1280}
)
