package response

type Code int

// 成功响应码 (1xxx)
const (
	CodeSuccess Code = 1000 + iota
	CodeCreated
	CodeAccepted
)

// 客户端错误码 (4xxx)
const (
	CodeValidationError Code = 4000 + iota
	CodeJSONError
	CodeEmptyBodyError
	CodeUnauthorized // 未授权
	CodeForbidden    // 禁止访问
	CodeNotFound     // 资源不存在
)

// 服务器错误码 (5xxx)
const (
	CodeServerError Code = 5000 + iota
	CodeDBError          // 数据库错误
	CodeCacheError       // 缓存错误
)
