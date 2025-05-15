package role

import (
	"comment/internal/domain/role/model"
	"comment/pkg/response"
	"comment/utils"
	"github.com/gin-gonic/gin"
)

type IController interface {
	CreateRole(ctx *gin.Context)
	DeleteRole(ctx *gin.Context)
	UpdateRole(ctx *gin.Context)
}

type controller struct {
	server IService
}

func NewController(svc IService) IController {
	ctrl := &controller{
		server: svc,
	}
	return ctrl
}

func (c *controller) CreateRole(ctx *gin.Context) {
	req := new(model.CreateRoleReq)
	if err := ctx.ShouldBindJSON(req); err != nil {
		response.ErrorParameterInvalid(ctx, err)
		return
	}
	if err := c.server.CreateRole(req); err != nil {
		response.Error(ctx, err)
		return
	}
	response.Success(ctx)
}

func (c *controller) DeleteRole(ctx *gin.Context) {
	id, err := utils.GetID(ctx)
	if err != nil {
		response.ErrorParameterInvalid(ctx, err)
		return
	}
	if err := c.server.DeleteRole(id); err != nil {
		response.Error(ctx, err)
		return
	}
	response.Success(ctx)
}

func (c *controller) UpdateRole(ctx *gin.Context) {
	id, err := utils.GetID(ctx)
	if err != nil {
		response.ErrorParameterInvalid(ctx, err)
		return
	}

	req := new(model.UpdateRoleReq)

	if err = ctx.ShouldBindJSON(req); err != nil {
		response.ErrorParameterInvalid(ctx, err)
		return
	}
	if err := c.server.UpdateRole(id, req); err != nil {
		response.Error(ctx, err)
		return
	}
	response.Success(ctx)
}
