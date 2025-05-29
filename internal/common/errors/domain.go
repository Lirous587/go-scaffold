package errors

import (
	"errors"
	"fmt"
)

// ErrorType 错误类型枚举
type ErrorType string

const (
	ErrorTypeValidation    ErrorType = "VALIDATION"
	ErrorTypeNotFound      ErrorType = "NOT_FOUND"
	ErrorTypeAlreadyExists ErrorType = "ALREADY_EXISTS"
	ErrorTypeUnauthorized  ErrorType = "UNAUTHORIZED"
	ErrorTypeForbidden     ErrorType = "FORBIDDEN"
	ErrorTypeInternal      ErrorType = "INTERNAL"
	ErrorTypeExternal      ErrorType = "EXTERNAL"
	ErrorTypeRateLimit     ErrorType = "RATE_LIMIT"
)

// DomainError 领域错误结构
type DomainError struct {
	Type    ErrorType              `json:"type"`
	Code    string                 `json:"code"`
	Message string                 `json:"message"`
	Details map[string]interface{} `json:"details,omitempty"`
	Cause   error                  `json:"-"`
}

func (e *DomainError) Error() string {
	if e.Cause != nil {
		return fmt.Sprintf("%s: %s (caused by: %v)", e.Code, e.Message, e.Cause)
	}
	return fmt.Sprintf("%s: %s", e.Code, e.Message)
}

func (e *DomainError) Unwrap() error {
	return e.Cause
}

func (e *DomainError) WithDetail(key string, value interface{}) *DomainError {
	if e.Details == nil {
		e.Details = make(map[string]interface{})
	}
	e.Details[key] = value
	return e
}

func (e *DomainError) WithCause(cause error) *DomainError {
	e.Cause = cause
	return e
}

func New(errorType ErrorType, code, message string) *DomainError {
	return &DomainError{
		Type:    errorType,
		Code:    code,
		Message: message,
	}
}

func NewValidation(code, message string) *DomainError {
	return New(ErrorTypeValidation, code, message)
}

func NewNotFound(code, message string) *DomainError {
	return New(ErrorTypeNotFound, code, message)
}

func NewAlreadyExists(code, message string) *DomainError {
	return New(ErrorTypeAlreadyExists, code, message)
}

func NewUnauthorized(code, message string) *DomainError {
	return New(ErrorTypeUnauthorized, code, message)
}

func NewForbidden(code, message string) *DomainError {
	return New(ErrorTypeForbidden, code, message)
}

func NewInternal(code, message string) *DomainError {
	return New(ErrorTypeInternal, code, message)
}

func NewExternal(code, message string) *DomainError {
	return New(ErrorTypeExternal, code, message)
}

func NewRateLimit(code, message string) *DomainError {
	return New(ErrorTypeRateLimit, code, message)
}

// 类型检查函数
func IsType(err error, errorType ErrorType) bool {
	var domainErr *DomainError
	ok := errors.As(err, &domainErr)
	return ok && domainErr.Type == errorType
}

func IsValidation(err error) bool {
	return IsType(err, ErrorTypeValidation)
}

func IsNotFound(err error) bool {
	return IsType(err, ErrorTypeNotFound)
}

func IsUnauthorized(err error) bool {
	return IsType(err, ErrorTypeUnauthorized)
}
