package handler

import (
	"github.com/gin-gonic/gin"
	"github.com/pkg/errors"
	"scaffold/internal/common/reqkit/bind"
	"scaffold/internal/common/reskit/response"
	"strconv"
	"scaffold/internal/mock/domain"
)

type HttpHandler struct {
    service domain.MockService
}

func NewHttpHandler(service domain.MockService) *HttpHandler {
    return &HttpHandler{
        service: service,
    }
}

func (h *HttpHandler) getID(ctx *gin.Context) (int64, error) {
    idStr := ctx.Param("id")
    idInt, err := strconv.Atoi(idStr)
    if err != nil {
        response.InvalidParams(ctx,err)
        return 0, err
    }
    if idInt == 0 {
        err := errors.New("无效的id")
        response.InvalidParams(ctx,err)
        return 0, err
    }
    return int64(idInt), err
}

func (h *HttpHandler) Create(ctx *gin.Context) {
    req := new(CreateRequest)

	if err := bind.BindingRegularAndResponse(ctx,req); err != nil {
		return
	}

    data, err := h.service.Create(&domain.Mock{
        Title:    req.Title,
        Description:  req.Description,
    })

    if err != nil {
        response.Error(ctx, err)
        return
    }

    response.Success(ctx, domainMockToResponse(data))
}


func (h *HttpHandler) Read(ctx *gin.Context) {
	id, err := h.getID(ctx)
	if err != nil {
		return
	}

	data, err := h.service.Read(id)

	if err != nil {
		response.Error(ctx, err)
		return
	}

	response.Success(ctx, domainMockToResponse(data))
}

func (h *HttpHandler) Update(ctx *gin.Context) {
    id, err := h.getID(ctx)
    if err != nil {
        return
    }

    req := new(UpdateRequest)

	if err := bind.BindingRegularAndResponse(ctx,req); err != nil {
		return
	}

    data, err := h.service.Update(&domain.Mock{
        ID:           id,
        Title:        req.Title,
        Description:  req.Description,
    })

    if err != nil {
        response.Error(ctx, err)
        return
    }

    response.Success(ctx, domainMockToResponse(data))
}

func (h *HttpHandler) Delete(ctx *gin.Context) {
    id, err := h.getID(ctx)
    if err != nil {
        return
    }

    if err := h.service.Delete(id); err != nil {
        response.Error(ctx, err)
        return
    }

    response.Success(ctx)
}

func (h *HttpHandler) List(ctx *gin.Context) {
    req := new(ListRequest)

	if err := bind.BindingRegularAndResponse(ctx,req); err != nil {
		return
	}

    data, err := h.service.List(&domain.MockQuery{
        Keyword:  req.KeyWord,
        Page:     req.Page,
        PageSize: req.PageSize,
    })

    if err != nil {
        response.Error(ctx, err)
        return
    }

    response.Success(ctx, domainMockListToResponse(data))
}