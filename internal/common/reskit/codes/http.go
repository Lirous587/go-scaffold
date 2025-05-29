package codes

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
	Cause      error
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

	var errCode ErrCode
	var errCode2 ErrCodeWithDetail
	var errCode3 ErrCodeWithCause

	ok1 := errors.As(err, &errCode)
	ok2 := errors.As(err, &errCode2)
	ok3 := errors.As(err, &errCode3)

	if !ok1 && !ok2 && !ok3 {
		// 不是自定义错误，返回通用服务器错误
		return HTTPError{
			StatusCode: http.StatusInternalServerError,
			Response: HTTPErrorResponse{
				Code:    5000,
				Message: "Internal server error",
			},
		}
	}
	if ok1 {
		return HTTPError{
			StatusCode: mapTypeToHTTPStatus(errCode.Type),
			Response: HTTPErrorResponse{
				Code:    errCode.Code,
				Message: errCode.Msg,
			},
		}
	}

	if ok2 {
		return HTTPError{
			StatusCode: mapTypeToHTTPStatus(errCode2.Type),
			Response: HTTPErrorResponse{
				Code:    errCode2.Code,
				Message: errCode2.Msg,
				Details: errCode3.Detail,
			},
		}
	}

	return HTTPError{
		StatusCode: mapTypeToHTTPStatus(errCode3.Type),
		Response: HTTPErrorResponse{
			Code:    errCode3.Code,
			Message: errCode3.Msg,
			Details: errCode3.Detail,
		},
		Cause: errCode3.Cause,
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
