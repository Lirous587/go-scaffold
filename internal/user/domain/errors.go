package domain

import (
	commonErrors "scaffold/internal/common/errors"
)

// 用户服务错误代码
const (
	// 用户相关
	CodeUserNotFound          = "USER_NOT_FOUND"
	CodeUserAlreadyExists     = "USER_ALREADY_EXISTS"
	CodeEmailAlreadyExists    = "EMAIL_ALREADY_EXISTS"
	CodeUsernameAlreadyExists = "USERNAME_ALREADY_EXISTS"

	// OAuth相关
	CodeInvalidOAuthCode     = "INVALID_OAUTH_CODE"
	CodeInvalidOAuthProvider = "INVALID_OAUTH_PROVIDER"
	CodeOAuthUserInfoMissing = "OAUTH_USER_INFO_MISSING"

	// Token相关
	CodeTokenGenerationFailed = "TOKEN_GENERATION_FAILED"
	CodeRefreshTokenInvalid   = "REFRESH_TOKEN_INVALID"
	CodeRefreshTokenExpired   = "REFRESH_TOKEN_EXPIRED"

	// 外部服务相关
	CodeGitHubAPIError = "GITHUB_API_ERROR"
	CodeGoogleAPIError = "GOOGLE_API_ERROR"
)

// 预定义错误变量
var (
	// 用户相关错误
	ErrUserNotFound          = commonErrors.NewNotFound(CodeUserNotFound, "用户不存在")
	ErrUserAlreadyExists     = commonErrors.NewAlreadyExists(CodeUserAlreadyExists, "用户已存在")
	ErrEmailAlreadyExists    = commonErrors.NewAlreadyExists(CodeEmailAlreadyExists, "邮箱已被使用")
	ErrUsernameAlreadyExists = commonErrors.NewAlreadyExists(CodeUsernameAlreadyExists, "用户名已被使用")

	// OAuth相关错误
	ErrInvalidOAuthCode     = commonErrors.NewValidation(CodeInvalidOAuthCode, "无效的OAuth授权码")
	ErrInvalidOAuthProvider = commonErrors.NewValidation(CodeInvalidOAuthProvider, "不支持的OAuth提供商")
	ErrOAuthUserInfoMissing = commonErrors.NewValidation(CodeOAuthUserInfoMissing, "OAuth用户信息缺失")

	// Token相关错误
	ErrTokenGenerationFailed = commonErrors.NewInternal(CodeTokenGenerationFailed, "Token生成失败")
	ErrTokenInvalid          = commonErrors.NewInternal(CodeTokenGenerationFailed, "Token无效")
	ErrTokenExpired          = commonErrors.NewInternal(CodeTokenGenerationFailed, "Token已过期")
	ErrRefreshTokenInvalid   = commonErrors.NewUnauthorized(CodeRefreshTokenInvalid, "无效的RefreshToken")
	ErrRefreshTokenExpired   = commonErrors.NewUnauthorized(CodeRefreshTokenExpired, "RefreshToken已过期")

	// 外部服务错误
	ErrGitHubAPIError = commonErrors.NewExternal(CodeGitHubAPIError, "GitHub API调用失败")
	ErrGoogleAPIError = commonErrors.NewExternal(CodeGoogleAPIError, "Google API调用失败")
)

// 错误构造函数（用于动态错误）
func NewUserNotFoundError(userID string) *commonErrors.DomainError {
	return commonErrors.NewNotFound(CodeUserNotFound, "用户不存在").
		WithDetail("user_id", userID)
}

func NewEmailExistsError(email string) *commonErrors.DomainError {
	return commonErrors.NewAlreadyExists(CodeEmailAlreadyExists, "邮箱已被使用").
		WithDetail("email", email)
}

func NewUsernameExistsError(username string) *commonErrors.DomainError {
	return commonErrors.NewAlreadyExists(CodeUsernameAlreadyExists, "用户名已被使用").
		WithDetail("username", username)
}

func NewGitHubAPIError(operation string, cause error) *commonErrors.DomainError {
	return commonErrors.NewExternal(CodeGitHubAPIError, "GitHub API调用失败").
		WithDetail("operation", operation).
		WithCause(cause)
}
