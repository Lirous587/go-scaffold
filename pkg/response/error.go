package response

import "github.com/pkg/errors"

type AppError struct {
	Code code
	Err  error
}

func (e AppError) Error() string {
	if e.Err != nil {
		return e.Err.Error()
	}
	return "未知错误"
}

func (e AppError) Unwrap() error {
	return e.Err
}

func NewAppError(code code, err error) *AppError {
	return &AppError{
		Code: code,
		Err:  errors.WithStack(err),
	}
}
