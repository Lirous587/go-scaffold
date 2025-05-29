package errors

import (
	"errors"
	"net/http"
)

// HTTPErrorResponse HTTP错误响应结构
type HTTPErrorResponse struct {
	Code    int                    `json:"code"`
	Message string                 `json:"message"`
	Details map[string]interface{} `json:"details,omitempty"`
}

// HTTPError HTTP错误信息
type HTTPError struct {
	StatusCode int
	Response   HTTPErrorResponse
}

// MapToHTTP 将领域错误映射为HTTP错误
func MapToHTTP(err error) HTTPError {
	if err == nil {
		return HTTPError{
			StatusCode: http.StatusOK,
			Response: HTTPErrorResponse{
				Code:    2000,
				Message: "Success",
			},
		}
	}

	var domainErr *DomainError
	ok := errors.As(err, &domainErr)
	if !ok {
		// 不是领域错误，返回通用服务器错误
		return HTTPError{
			StatusCode: http.StatusInternalServerError,
			Response: HTTPErrorResponse{
				Code:    5000,
				Message: "Internal server error",
			},
		}
	}

	return HTTPError{
		StatusCode: mapTypeToHTTPStatus(domainErr.Type),
		Response: HTTPErrorResponse{
			Code:    mapTypeToHTTPCode(domainErr.Type),
			Message: domainErr.Message,
			Details: domainErr.Details,
		},
	}
}

// mapTypeToHTTPStatus 映射错误类型到HTTP状态码
func mapTypeToHTTPStatus(errorType ErrorType) int {
	switch errorType {
	case ErrorTypeValidation:
		return http.StatusBadRequest
	case ErrorTypeNotFound:
		return http.StatusNotFound
	case ErrorTypeAlreadyExists:
		return http.StatusConflict
	case ErrorTypeUnauthorized:
		return http.StatusUnauthorized
	case ErrorTypeForbidden:
		return http.StatusForbidden
	case ErrorTypeRateLimit:
		return http.StatusTooManyRequests
	case ErrorTypeExternal:
		return http.StatusBadGateway
	default: // ErrorTypeInternal
		return http.StatusInternalServerError
	}
}

// mapTypeToHTTPCode 映射错误类型到业务错误码
func mapTypeToHTTPCode(errorType ErrorType) int {
	switch errorType {
	case ErrorTypeValidation:
		return 4000
	case ErrorTypeUnauthorized:
		return 4010
	case ErrorTypeForbidden:
		return 4030
	case ErrorTypeNotFound:
		return 4040
	case ErrorTypeAlreadyExists:
		return 4090
	case ErrorTypeRateLimit:
		return 4290
	case ErrorTypeExternal:
		return 5020
	default: // ErrorTypeInternal
		return 5000
	}
}
