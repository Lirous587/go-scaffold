package response

// Meta 定义返回的元数据
type Meta struct {
	Code       Code   `json:"code"`
	Message    string `json:"message"`
	HttpStatus int    `json:"-"`
}

type Success struct {
	Code    Code   `json:"code"`
	Message string `json:"message"`
	Data    any    `json:"data"`
}

// 全局响应映射表
var codeAndResponseMap map[Code]Meta

func GetResponse(code Code) Meta {
	res, ok := codeAndResponseMap[code]
	if !ok {
		res.Message = "未知错误"
	}
	return res
}

func init() {
	codeAndResponseMap = map[Code]Meta{
		CodeSuccess:         {Message: "操作成功"},
		CodeValidationError: {Message: "参数验证失败", HttpStatus: 400},
		CodeJSONError:       {Message: "JSON解析错误"},
		CodeEmptyBodyError:  {Message: "请求体为空"},
		CodeServerError:     {Message: "服务器内部错误"},
	}
}
