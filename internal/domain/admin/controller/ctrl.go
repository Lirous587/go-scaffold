package controller

import (
	"github.com/gin-gonic/gin"
	"github.com/pkg/errors"
	"scaffold/internal/domain/admin/model"
	"scaffold/internal/domain/admin/service"
	"scaffold/pkg/response"
)

type Controller interface {
	IfInit(ctx *gin.Context)
	Init(ctx *gin.Context)
	Login(ctx *gin.Context)
	RefreshToken(ctx *gin.Context)
}

type controller struct {
	server service.Service
}

func NewController(service service.Service) Controller {
	return &controller{
		server: service,
	}
}

func (c *controller) IfInit(ctx *gin.Context) {
	have, err := c.server.IfInit()
	if err != nil {
		response.Error(ctx, err)
		return
	}
	if have {
		err = response.NewAppError(response.CodeAdminExist, errors.New("用户已经初始化"))
		return
	}

	response.Success(ctx)
}

func (c *controller) Init(ctx *gin.Context) {
	req := new(model.InitReq)
	if err := ctx.ShouldBindJSON(req); err != nil {
		response.ErrorParameterInvalid(ctx, err)
		return
	}
	appErr := c.server.Init(req)
	if appErr != nil {
		response.Error(ctx, appErr)
		return
	}
	response.Success(ctx)
}

func (c *controller) Login(ctx *gin.Context) {
	req := new(model.LoginReq)
	if err := ctx.ShouldBindJSON(req); err != nil {
		response.ErrorParameterInvalid(ctx, err)
		return
	}

	res, appErr := c.server.Auth(req.Email, req.Password)
	if appErr != nil {
		response.Error(ctx, appErr)
		return
	}

	response.Success(ctx, res)
}

func (c *controller) RefreshToken(ctx *gin.Context) {
	refreshToken := ctx.GetHeader("refresh-token")
	if refreshToken == "" {
		appError := response.NewAppError(response.CodeAuthFailed, errors.New("refresh-token请求头为空"))
		response.Error(ctx, appError)
		return
	}

	req := new(model.RefreshTokenReq)
	if err := ctx.ShouldBindJSON(req); err != nil {
		response.ErrorParameterInvalid(ctx, err)
		return
	}

	res, appErr := c.server.RefreshToken(req, refreshToken)
	if appErr != nil {
		response.Error(ctx, appErr)
		return
	}

	response.Success(ctx, res)
}
