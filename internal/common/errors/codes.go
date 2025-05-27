package errors

// 全局错误代码常量
const (
	// 通用错误代码
	CodeValidationFailed  = "VALIDATION_FAILED"
	CodeInternalError     = "INTERNAL_ERROR"
	CodeExternalError     = "EXTERNAL_ERROR"
	CodeRateLimitExceeded = "RATE_LIMIT_EXCEEDED"

	// 认证和授权错误
	CodeUnauthorized = "UNAUTHORIZED"
	CodeTokenInvalid = "TOKEN_INVALID"
	CodeTokenExpired = "TOKEN_EXPIRED"
	CodeForbidden    = "FORBIDDEN"

	// 资源错误
	CodeNotFound      = "NOT_FOUND"
	CodeAlreadyExists = "ALREADY_EXISTS"
	CodeConflict      = "CONFLICT"
)

// 预定义的通用错误
var (
	ErrValidationFailed = NewValidation(CodeValidationFailed, "Validation failed")
	ErrInternalError    = NewInternal(CodeInternalError, "Internal server error")
	ErrUnauthorized     = NewUnauthorized(CodeUnauthorized, "Unauthorized")
	ErrForbidden        = NewForbidden(CodeForbidden, "Forbidden")
	ErrNotFound         = NewNotFound(CodeNotFound, "Resource not found")
	ErrAlreadyExists    = NewAlreadyExists(CodeAlreadyExists, "Resource already exists")
)
