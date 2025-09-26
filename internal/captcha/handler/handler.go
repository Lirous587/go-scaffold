package handler

import (
	"github.com/gin-gonic/gin"
	"github.com/pkg/errors"
	"os"
	"scaffold/internal/captcha/domain"
	"scaffold/internal/captcha/service"
	"scaffold/internal/common/reskit/codes"
	"scaffold/internal/common/reskit/response"
	"strconv"
)

type HttpHandler struct {
	service *service.CaptchaServiceFactor
}

func NewHttpHandler(service *service.CaptchaServiceFactor) *HttpHandler {
	return &HttpHandler{
		service: service,
	}
}

func (h *HttpHandler) getID(ctx *gin.Context) (int64, error) {
	idStr := ctx.Param("id")
	idInt, err := strconv.Atoi(idStr)
	if err != nil {
		return 0, err
	}
	if idInt == 0 {
		return 0, errors.New("无效的id")
	}
	return int64(idInt), err
}

// Gen godoc
// @Summary      生成验证码
// @Description  创建新的验证码
// @Tags         captcha
// @Accept       json
// @Produce      json
// @Param        request query handler.GenRequest true "创建验证码请求"
// @Success      200  {object}  response.successResponse{data=handler.CaptchaResponse} "成功创建验证码"
// @Failure      400  {object}  response.invalidParamsResponse "参数错误"
// @Failure      500  {object}  response.errorResponse "服务器错误"
// @Router       /captcha [post]
func (h *HttpHandler) Gen(ctx *gin.Context) {
	req := new(GenRequest)

	if err := ctx.ShouldBindQuery(req); err != nil {
		response.InvalidParams(ctx, err)
		return
	}

	res, err := h.service.Generate(req.Way)
	if err != nil {
		response.Error(ctx, err)
		return
	}

	response.Success(ctx, domainCaptchaToResponse(res))
}

// GenWithAnswer godoc
// @Summary      生成带答案的验证码
// @Description  创建新的验证码并返回答案（仅用于测试或开发环境）
// @Tags         captcha
// @Accept       json
// @Produce      json
// @Param        request query handler.GenRequest true "创建验证码请求"
// @Success      200  {object}  response.successResponse{data=handler.CaptchaAnswerResponse} "成功创建验证码并返回答案"
// @Failure      400  {object}  response.invalidParamsResponse "参数错误"
// @Failure      500  {object}  response.errorResponse "服务器错误"
// @Router       /captcha/with-answer [get]
func (h *HttpHandler) GenWithAnswer(ctx *gin.Context) {
	mode := os.Getenv("SERVER_MODE")
	if mode != "dev" {
		response.Error(ctx, codes.ErrAPIForbidden)
	}

	req := new(GenRequest)

	if err := ctx.ShouldBindQuery(req); err != nil {
		response.InvalidParams(ctx, err)
		return
	}

	res, err := h.service.GenWithAnswer(req.Way)
	if err != nil {
		response.Error(ctx, err)
		return
	}

	response.Success(ctx, domainCaptchaAnswerToResponse(res))
}

const (
	verifyWayHeaderKey   = "captcha-verify-way"
	verifyIDHeaderKey    = "captcha-verify-id"
	verifyValueHeaderKey = "captcha-verify-value"
)

// parseFromHeader 从请求头中获取验证方式
func parseFromHeader(c *gin.Context) (way domain.VerifyWay, id int64, value string, err error) {
	wayFromHeader := c.GetHeader(verifyWayHeaderKey)
	if wayFromHeader == "" {
		return "", 0, "", errors.New("验证方式为空")
	}
	way = domain.VerifyWay(wayFromHeader)

	verifyID := c.GetHeader(verifyIDHeaderKey)
	if verifyID == "" {
		return "", 0, "", errors.New("无效的验证id")
	}

	id, err = strconv.ParseInt(verifyID, 10, 64)
	if err != nil {
		return "", 0, "", errors.New("无效的id")
	}

	value = c.GetHeader(verifyValueHeaderKey)
	if value == "" {
		return "", 0, "", errors.New("验证信息为空")
	}

	return
}

// Verify 作为中间件
func (h *HttpHandler) Verify() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		// 1.解析
		way, k, v, err := parseFromHeader(ctx)
		if err != nil {
			response.Error(ctx, codes.ErrCaptchaFormatInvalid.WithCause(err))
			return
		}
		// 2.验证
		if err := h.service.Verify(way, k, v); err != nil {
			response.Error(ctx, codes.ErrCaptchaVerifyFailed.WithCause(err))
			return
		}

		ctx.Next()
	}
}
