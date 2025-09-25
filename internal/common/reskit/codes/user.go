package codes

import (
	"fmt"
)

// 用户模块错误码 (1000-1199)
var (
	// 基础认证错误 (1000-1019)
	ErrUnauthorized      = ErrCode{Msg: "未授权访问", Type: ErrorTypeUnauthorized, Code: 1001}
	ErrUserNotFound      = ErrCode{Msg: "用户不存在", Type: ErrorTypeNotFound, Code: 1002}
	ErrUserAlreadyExists = ErrCode{Msg: "用户已存在", Type: ErrorTypeAlreadyExists, Code: 1003}

	// 用户信息冲突错误 (1020-1039)
	ErrEmailAlreadyExists    = ErrCode{Msg: "邮箱已被使用", Type: ErrorTypeAlreadyExists, Code: 1020}
	ErrUsernameAlreadyExists = ErrCode{Msg: "用户名已被使用", Type: ErrorTypeAlreadyExists, Code: 1021}

	// OAuth相关错误 (1040-1059)
	ErrOAuthInvalidCode     = ErrCode{Msg: "无效的OAuth授权码", Type: ErrorTypeValidation, Code: 1040}
	ErrOAuthInvalidProvider = ErrCode{Msg: "不支持的OAuth提供商", Type: ErrorTypeValidation, Code: 1041}
	ErrOAuthUserInfoMissing = ErrCode{Msg: "OAuth用户信息缺失", Type: ErrorTypeValidation, Code: 1042}

	// Token相关错误 (1060-1079)
	ErrTokenGenerationFailed = ErrCode{Msg: "Token生成失败", Type: ErrorTypeInternal, Code: 1060}
	ErrTokenInvalid          = ErrCode{Msg: "Token无效", Type: ErrorTypeUnauthorized, Code: 1061}
	ErrTokenFormatInvalid    = ErrCode{Msg: "Token格式无效", Type: ErrorTypeValidation, Code: 1062}
	ErrTokenExpired          = ErrCode{Msg: "Token已过期", Type: ErrorTypeUnauthorized, Code: 1063}

	ErrRefreshTokenMissingInHeader = ErrCode{Msg: "请求头中缺少RefreshToken参数", Type: ErrorTypeValidation, Code: 1070}
	ErrRefreshTokenNotFound        = ErrCode{
		Msg: "登录凭证过期", Type: ErrorTypeUnauthorized, Code: 1071,
	}.WithCause(
		fmt.Errorf("%s", "RefreshToken不存在"),
	)

	// 外部服务错误 (1080-1099)
	ErrGitHubAPIError = ErrCode{Msg: "GitHub API调用失败", Type: ErrorTypeExternal, Code: 1080}
	ErrGoogleAPIError = ErrCode{Msg: "Google API调用失败", Type: ErrorTypeExternal, Code: 1081}
)
