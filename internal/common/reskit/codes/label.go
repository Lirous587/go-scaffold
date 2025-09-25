package codes

// 标签相关错误 (2000-2199)
var (
	ErrLabelNotFound      = ErrCode{Msg: "标签不存在", Type: ErrorTypeNotFound, Code: 2000}
	ErrLabelAlreadyExists = ErrCode{Msg: "标签已存在", Type: ErrorTypeAlreadyExists, Code: 2001}
	ErrLabelHasArticle    = ErrCode{Msg: "该标签已关联文章", Type: ErrorTypeExternal, Code: 2002}
)
