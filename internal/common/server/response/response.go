package response

import (
	"github.com/gin-gonic/gin"
	"github.com/pkg/errors"
	"net/http"
)

type code int

const (
	codeSuccess      code = 2000
	codeParamInvalid code = 4000
	CodeServerError  code = 5000
	codeUnKnowError  code = 9999

	// CodeAuthFailed 认证相关
	CodeAuthFailed code = 4100 + iota
	CodeTokenInvalid
	CodeTokenExpired
	CodeRefreshInvalid
)

var errCodeMsgMap = map[code]string{
	CodeServerError:    "服务器错误",
	CodeAuthFailed:     "认证失败",
	CodeTokenInvalid:   "无效的令牌",
	CodeTokenExpired:   "令牌已过期",
	CodeRefreshInvalid: "无效的refreshToken",
}

type AppError struct {
	Code code
	Err  error
}

func (e *AppError) Error() string {
	if e.Err != nil {
		return e.Err.Error()
	}
	return "未知错误"
}

func NewAppError(code code, err error) *AppError {
	return &AppError{
		Code: code,
		Err:  errors.WithStack(err),
	}
}

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
		Data:    err.Error(),
	}
	ctx.Error(err)
	ctx.AbortWithStatusJSON(http.StatusBadRequest, res)
}

func Error(ctx *gin.Context, err error) {
	var appErr *AppError
	ok := errors.As(err, &appErr)
	res := response{}
	if ok {
		res.Code = appErr.Code
		msg, exist := errCodeMsgMap[appErr.Code]
		if exist {
			res.Message = msg
		} else {
			res.Code = codeUnKnowError
			res.Message = "未知错误"
		}
		ctx.Error(appErr.Err)
	} else {
		res.Code = CodeServerError
		res.Message = err.Error()
		ctx.Error(err)
	}
	ctx.JSON(http.StatusInternalServerError, res)
}
