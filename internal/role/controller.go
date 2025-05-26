package role

import (
	"github.com/gin-gonic/gin"
	"scaffold/internal/common/server/response"
	"scaffold/internal/common/utils"
	"scaffold/internal/role/model"
)

type Controller struct {
	server Service
}

func NewController(svc Service) *Controller {
	return &Controller{
		server: svc,
	}
}

func (c *Controller) CreateRole(ctx *gin.Context) {
	req := new(model.CreateRoleReq)

	if err := ctx.ShouldBindJSON(req); err != nil {
		response.ErrorParameterInvalid(ctx, err)
		return
	}

	role, err := c.server.CreateRole(req)
	if err != nil {
		response.Error(ctx, err)
		return
	}

	response.Success(ctx, role)
}

func (c *Controller) DeleteRole(ctx *gin.Context) {
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

func (c *Controller) UpdateRole(ctx *gin.Context) {
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

func (c *Controller) AllRole(ctx *gin.Context) {
	roles, err := c.server.AllRole()
	if err != nil {
		response.Error(ctx, err)
		return
	}
	response.Success(ctx, roles)
}
