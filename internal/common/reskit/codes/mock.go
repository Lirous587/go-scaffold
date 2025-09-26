package codes

// Mock相关错误 (xx00-xx99)
var (
	ErrMockNotFound = ErrCode{Msg: "Mock不存在", Type: ErrorTypeNotFound, Code: 0000}
)
