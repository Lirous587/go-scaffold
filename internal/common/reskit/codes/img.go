package codes

// 图库相关错误 (1400-1599)
var (
	// 图库基础操作 (1400-1419)
	ErrImgNotFound   = ErrCode{Msg: "图片不存在", Type: ErrorTypeNotFound, Code: 1400}
	ErrImgPathRepeat = ErrCode{Msg: "图片路径重复,请检查图库和回收站", Type: ErrorTypeAlreadyExists, Code: 1401}

	ErrImgCategoryNotFound    = ErrCode{Msg: "图片分类不存在", Type: ErrorTypeNotFound, Code: 1402}
	ErrImgCategoryTitleRepeat = ErrCode{Msg: "图片分类名重复", Type: ErrorTypeAlreadyExists, Code: 1403}
	ErrImgCategoryToMany      = ErrCode{Msg: "图片分类过多", Type: ErrorTypeExternal, Code: 1404}
	ErrImgCategoryExistImg    = ErrCode{
		Msg: "当前图片分类下存在图片,请检查图库和回收站", Type: ErrorTypeExternal, Code: 1405,
	}

	// 图片处理 (1420-1439)
	ErrImgCompress         = ErrCode{Msg: "压缩图片失败", Type: ErrorTypeInternal, Code: 1420}
	ErrImgUploadToR3Failed = ErrCode{Msg: "上传图片到R3失败", Type: ErrorTypeInternal, Code: 1421}
)
