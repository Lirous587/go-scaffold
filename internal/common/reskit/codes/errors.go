package codes

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

type ErrCode struct {
	Msg  string
	Type ErrorType
	Code int
}

func (e ErrCode) Error() string {
	return e.Msg
}

func (e ErrCode) WithSlug(slug string) ErrCode {
	return ErrCode{
		Msg:  e.Msg + " " + slug,
		Type: e.Type,
		Code: e.Code,
	}
}

type ErrCodeWithDetail struct {
	Msg    string
	Type   ErrorType
	Code   int
	Detail map[string]any `json:"detail,omitempty"`
}

func (e ErrCode) WithDetail(detail map[string]any) ErrCodeWithDetail {
	return ErrCodeWithDetail{
		Msg:    e.Msg,
		Type:   e.Type,
		Code:   e.Code,
		Detail: detail,
	}
}

func (e ErrCodeWithDetail) Error() string {
	return e.Msg
}

type ErrCodeWithCause struct {
	Msg    string
	Type   ErrorType
	Code   int
	Detail map[string]any `json:"detail,omitempty"`
	Cause  error          `json:"-"`
}

func (e ErrCode) WithCause(err error) ErrCodeWithCause {
	return ErrCodeWithCause{
		Msg:   e.Msg,
		Type:  e.Type,
		Code:  e.Code,
		Cause: err,
	}
}

func (e ErrCodeWithDetail) WithCause(err error) ErrCodeWithCause {
	return ErrCodeWithCause{
		Msg:    e.Msg,
		Type:   e.Type,
		Code:   e.Code,
		Detail: e.Detail,
		Cause:  err,
	}
}

func (e ErrCodeWithCause) Error() string {
	return e.Msg
}
