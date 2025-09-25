package codes

// 时刻相关错误 (2600-2799)
var (
	ErrMomentNotFound = ErrCode{Msg: "时刻不存在", Type: ErrorTypeNotFound, Code: 2600}
)
