package permission

import (
	"comment/internal/domain/role/model"
	"comment/pkg/response"
	"comment/utils"
	"github.com/gin-gonic/gin"
)

type IController interface {
	CreatePermission(ctx *gin.Context)
	UpdatePermission(ctx *gin.Context)
	DeletePermission(ctx *gin.Context)
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

func (c *controller) CreatePermission(ctx *gin.Context) {
	req := new(model.CreateRoleReq)
	if err := ctx.ShouldBindJSON(req); err != nil {
		response.ErrorParameterInvalid(ctx, err)
		return
	}
	if err := c.server.CreatePermission(req); err != nil {
		response.Error(ctx, err)
		return
	}
	response.Success(ctx)
}

func (c *controller) UpdatePermission(ctx *gin.Context) {
	id, err := utils.GetID(ctx)
	if err != nil {
		response.ErrorParameterInvalid(ctx, err)
		return
	}
	if err := c.server.DeletePermission(id); err != nil {
		response.Error(ctx, err)
		return
	}
	response.Success(ctx)
}

func (c *controller) DeletePermission(ctx *gin.Context) {
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
	if err := c.server.UpdatePermission(id, req); err != nil {
		response.Error(ctx, err)
		return
	}
	response.Success(ctx)
}
