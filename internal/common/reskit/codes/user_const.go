package codes

// 用户相关错误
var (
	ErrUnauthorized          = ErrCode{Msg: "未授权访问", Type: ErrorTypeUnauthorized, Code: 1001}
	ErrUserNotFound          = ErrCode{Msg: "用户不存在", Type: ErrorTypeNotFound, Code: 1002}
	ErrUserAlreadyExists     = ErrCode{Msg: "用户已存在", Type: ErrorTypeAlreadyExists, Code: 1003}
	ErrEmailAlreadyExists    = ErrCode{Msg: "邮箱已被使用", Type: ErrorTypeAlreadyExists, Code: 1004}
	ErrUsernameAlreadyExists = ErrCode{Msg: "用户名已被使用", Type: ErrorTypeAlreadyExists, Code: 1005}

	// OAuth相关错误
	ErrOAuthInvalidCode     = ErrCode{Msg: "无效的OAuth授权码", Type: ErrorTypeValidation, Code: 1011}
	ErrOAuthInvalidProvider = ErrCode{Msg: "不支持的OAuth提供商", Type: ErrorTypeValidation, Code: 1012}
	ErrOAuthUserInfoMissing = ErrCode{Msg: "OAuth用户信息缺失", Type: ErrorTypeValidation, Code: 1013}

	// Token相关错误
	ErrTokenGenerationFailed = ErrCode{Msg: "Token生成失败", Type: ErrorTypeInternal, Code: 1021}
	ErrTokenInvalid          = ErrCode{Msg: "Token无效", Type: ErrorTypeInternal, Code: 1022}
	ErrTokenExpired          = ErrCode{Msg: "Token已过期", Type: ErrorTypeInternal, Code: 1023}
	ErrRefreshTokenInvalid   = ErrCode{Msg: "无效的RefreshToken", Type: ErrorTypeUnauthorized, Code: 1024}
	ErrRefreshTokenExpired   = ErrCode{Msg: "RefreshToken已过期", Type: ErrorTypeUnauthorized, Code: 1025}

	// 外部服务错误
	ErrGitHubAPIError = ErrCode{Msg: "GitHub API调用失败", Type: ErrorTypeExternal, Code: 1031}
	ErrGoogleAPIError = ErrCode{Msg: "Google API调用失败", Type: ErrorTypeExternal, Code: 1032}
)
