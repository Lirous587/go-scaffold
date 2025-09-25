package codes

var (
	ErrFriendlinkNotFound            = ErrCode{Msg: "友链不存在", Type: ErrorTypeNotFound, Code: 1501}
	ErrFriendlinkMissingRejectReason = ErrCode{
		Msg: "拒绝申请时必须提供拒绝理由", Type: ErrorTypeExternal, Code: 1502,
	}
	ErrFriendlinkInvalidStatusTransition = ErrCode{Msg: "非法的友链状态转换", Type: ErrorTypeExternal, Code: 1503}
)
