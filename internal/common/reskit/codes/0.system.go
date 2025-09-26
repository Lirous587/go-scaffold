package codes

// 系统级错误码 (0-999)
var (
	ErrAPIForbidden = ErrCode{Msg: "当前接口禁止访问", Type: ErrorTypeForbidden, Code: 1}
)
