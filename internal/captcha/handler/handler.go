package handler

import (
	"scaffold/internal/captcha/domain"
	"scaffold/internal/captcha/service"
	"scaffold/internal/common/reskit/codes"
	"scaffold/internal/common/reskit/response"
	"github.com/gin-gonic/gin"
	"github.com/pkg/errors"
	"strconv"
	"strings"
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

func (h *HttpHandler) GenWithAnswer(ctx *gin.Context) {
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
	verifyWayHeaderKey	= "verify-way"
	verifyHeaderKey		= "verify"
)

func parseFromHeader(c *gin.Context) (domain.VerifyWay, int64, string, error) {
	var way domain.VerifyWay
	var id int64
	var value string

	wayHeader := c.GetHeader(verifyWayHeaderKey)
	if wayHeader == "" {
		return "", 0, "", errors.New("验证方式为空")
	}
	way = domain.VerifyWay(wayHeader)

	verifyKey := c.GetHeader(verifyHeaderKey)
	if verifyKey == "" {
		return "", 0, "", errors.New("验证信息为空")
	}
	slice := strings.Split(verifyKey, ":")
	if len(slice) != 2 {
		return "", 0, "", errors.New("验证信息错误")
	}

	id, err := strconv.ParseInt(slice[0], 10, 64)
	if err != nil {
		return "", 0, "", errors.New("无效的id")
	}
	value = slice[1]

	return way, id, value, nil
}

func (h *HttpHandler) Verify() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 解析
		way, k, v, err := parseFromHeader(c)
		if err != nil {
			response.Error(c, codes.ErrCaptchaFormatInvalid)
			return
		}

		//fmt.Println(way)
		//fmt.Println(k)
		//fmt.Println(v)

		//	2验证
		if err := h.service.Verify(way, k, v); err != nil {
			response.Error(c, codes.ErrCaptchaVerifyFailed)
			return
		}

		c.Next()
	}
}

func (h *HttpHandler) VerifyEndpoint(ctx *gin.Context) {
	response.Success(ctx)
}
