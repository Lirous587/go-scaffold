package response

var (
	errCodeMsgMap = map[code]string{

		// 服务端错误消息映射
		CodeServerError:      "服务器错误",
		CodeIllegalOperation: "非法操作",

		// 认证错误
		CodeAuthFailed:     "认证失败",
		CodeTokenInvalid:   "无效的令牌",
		CodeTokenExpired:   "令牌已过期",
		CodeRefreshInvalid: "无效的refreshToken",

		// admin模块
		CodeAdminExist: "管理员已初始化",
	}
)
