package handler

import (
	"blog-v4/internal/common/reskit/response"
	"blog-v4/internal/moment/domain"
	"github.com/gin-gonic/gin"
	"github.com/pkg/errors"
	"strconv"
)

type HttpHandler struct {
	service domain.MomentService
}

func NewHttpHandler(service domain.MomentService) *HttpHandler {
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

func (h *HttpHandler) Create(ctx *gin.Context) {
	req := new(CreateRequest)

	req.Coordinates = MomentCoordinate{
		X: domain.InvalidCoordinate,
		Y: domain.InvalidCoordinate,
	}

	if err := ctx.ShouldBind(req); err != nil {
		response.InvalidParams(ctx, err)
		return
	}

	data, err := h.service.Create(&domain.Moment{
		Title:    req.Title,
		Content:  req.Content,
		Location: req.Location,
		Coordinates: &domain.MomentCoordinate{
			X: req.Coordinates.X,
			Y: req.Coordinates.Y,
		},
		CoverURL: req.CoverURL,
	})

	if err != nil {
		response.Error(ctx, err)
		return
	}

	response.Success(ctx, domainMomentToResponse(data))
}

func (h *HttpHandler) Update(ctx *gin.Context) {
	id, err := h.getID(ctx)
	if err != nil {
		response.InvalidParams(ctx, err)
		return
	}

	req := new(UpdateRequest)
	req.Coordinates = MomentCoordinate{
		X: domain.InvalidCoordinate,
		Y: domain.InvalidCoordinate,
	}

	if err := ctx.ShouldBindJSON(req); err != nil {
		response.InvalidParams(ctx, err)
		return
	}

	data, err := h.service.Update(&domain.Moment{
		ID:       id,
		Title:    req.Title,
		Content:  req.Content,
		Location: req.Location,
		Coordinates: &domain.MomentCoordinate{
			X: req.Coordinates.X,
			Y: req.Coordinates.Y,
		},
		CoverURL: req.CoverURL,
	})

	if err != nil {
		response.Error(ctx, err)
		return
	}

	response.Success(ctx, domainMomentToResponse(data))
}

func (h *HttpHandler) Delete(ctx *gin.Context) {
	id, err := h.getID(ctx)
	if err != nil {
		response.InvalidParams(ctx, err)
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

	if err := ctx.ShouldBindQuery(req); err != nil {
		response.InvalidParams(ctx, err)
		return
	}

	data, err := h.service.ListMoments(&domain.MomentQuery{
		Keyword:  req.KeyWord,
		Page:     req.Page,
		PageSize: req.PageSize,
	})

	if err != nil {
		response.Error(ctx, err)
		return
	}

	response.Success(ctx, domainMomentWithPagesToResponse(data))
}

func (h *HttpHandler) RandomN(ctx *gin.Context) {
	req := new(RandomNRequest)
	if err := ctx.ShouldBindQuery(req); err != nil {
		response.InvalidParams(ctx, err)
		return
	}
	data, err := h.service.RandomN(req.Count)
	if err != nil {
		response.Error(ctx, err)
		return
	}
	response.Success(ctx, domainMomentsToResponse(data))
}
