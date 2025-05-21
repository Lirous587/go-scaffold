package response

import (
	"net/http"
	"scaffold/internal/common/validator"

	"github.com/gin-gonic/gin"
)

type response struct {
	Code    code        `json:"code"`
	Message string      `json:"msg"`
	Data    interface{} `json:"data,omitempty"`
}

func Success(ctx *gin.Context, data ...any) {
	res := response{
		Code:    codeSuccess,
		Message: "请求成功",
	}
	if len(data) > 0 {
		res.Data = data[0]
	}
	ctx.JSON(http.StatusOK, res)
}

func ErrorParameterInvalid(ctx *gin.Context, err error) {
	res := response{
		Code:    codeParamInvalid,
		Message: "参数无效",
	}
	lang := validator.GetTranslateLang(ctx)
	transErr := validator.TranslateError(err, lang)
	res.Data = transErr.Error()

	if err != nil {
		ctx.Error(err)
	}
	ctx.AbortWithStatusJSON(http.StatusBadRequest, res)
}

func Error(ctx *gin.Context, appErr *AppError) {
	res := response{
		Code: appErr.Code,
	}
	msg, ok := errCodeMsgMap[appErr.Code]
	if ok {
		res.Message = msg
	} else {
		res.Code = codeUnKnowError
		res.Message = "未知错误"
	}

	if appErr.Err != nil {
		ctx.Error(appErr.Err)
	}

	ctx.AbortWithStatusJSON(http.StatusBadRequest, res)
}
