package middleware

// 响应状态码常量
const (
	// 成功状态码
	CodeSuccess = 1000

	// 客户端错误 (4xxx)
	CodeValidationError = 4001 // 验证错误
	CodeJSONError       = 4002 // JSON解析错误
	CodeEmptyBodyError  = 4003 // 空请求体错误

	// 服务端错误 (5xxx)
	CodeServerError = 5000 // 服务器内部错误
)

// Response 统一API响应结构
type Response struct {
	Code    int         `json:"code"`
	Message interface{} `json:"message"`
	Data    interface{} `json:"data"`
}
