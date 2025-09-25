package codes

// 文章相关错误 (2200-2399)
var (
	// 文章基础操作 (2200-2219)
	ErrArticleNotFound = ErrCode{Msg: "文章不存在", Type: ErrorTypeNotFound, Code: 2200}

	// 文章统计相关 (2220-2239)
	ErrArticleVtNotFound = ErrCode{Msg: "文章浏览量不存在", Type: ErrorTypeNotFound, Code: 2220}

	// 文章缓存相关 (2240-2259)
	ErrArticleArchiveCacheMiss = ErrCode{Msg: "文章归档缓存未命中", Type: ErrorTypeCacheMiss, Code: 2240}
)
