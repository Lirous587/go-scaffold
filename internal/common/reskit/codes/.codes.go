package codes

// 错误码范围定义
const (
	// 系统级错误 (1-999)
	SystemErrorStart = 1
	SystemErrorEnd   = 999

	// 用户模块 (1000-1199)
	UserErrorStart = 1000
	UserErrorEnd   = 1199

	// 验证码模块 (1200-1399)
	CaptchaErrorStart = 1200
	CaptchaErrorEnd   = 1399

	// 图库模块 (1400-1599)
	ImgErrorStart = 1400
	ImgErrorEnd   = 1599

	// 标签模块 (2000-2199)
	LabelErrorStart = 2000
	LabelErrorEnd   = 2199

	// 文章模块 (2200-2399)
	ArticleErrorStart = 2200
	ArticleErrorEnd   = 2399

	// 格言模块 (2400-2599)
	MaximErrorStart = 2400
	MaximErrorEnd   = 2599

	// 时刻/瞬间模块 (2200-2399)
	MomentErrorStart = 2600
	MomentErrorEnd   = 2799

	// 友链模块
	FriendlinkErrorStart = 2800
	FriendlinkErrorEnd   = 2999
)
