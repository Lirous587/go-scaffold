package response

type code int

const (
	codeSuccess          code = 2000
	codeParamInvalid     code = 4000
	CodeServerError      code = 5000
	CodeIllegalOperation code = 5001
	codeUnKnowError      code = 9999
)

// 认证授权错误 4100-4199
const (
	CodeAuthFailed code = 4100 + iota
	CodeTokenInvalid
	CodeTokenExpired
	CodeRefreshInvalid
)
