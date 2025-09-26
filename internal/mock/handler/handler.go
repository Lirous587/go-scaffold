package handler

import (
	"github.com/gin-gonic/gin"
	"github.com/pkg/errors"
	"scaffold/internal/common/reqkit/bind"
	"scaffold/internal/common/reskit/response"
	"scaffold/internal/mock/domain"
	"strconv"
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
		response.InvalidParams(ctx, err)
		return 0, err
	}
	if idInt == 0 {
		err := errors.New("无效的id")
		response.InvalidParams(ctx, err)
		return 0, err
	}
	return int64(idInt), err
}

// Create godoc
// @Summary      创建 Mock
// @Description  创建新的 Mock
// @Tags         mock
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        request body handler.CreateRequest true "创建 Mock 请求"
// @Success      200  {object}  response.successResponse{data=handler.MockResponse} "成功创建 Mock"
// @Failure      400  {object}  response.invalidParamsResponse "参数错误"
// @Failure      500  {object}  response.errorResponse "服务器错误"
// @Router       /mock [post]
func (h *HttpHandler) Create(ctx *gin.Context) {
	req := new(CreateRequest)

	if err := bind.BindingRegularAndResponse(ctx, req); err != nil {
		return
	}

	data, err := h.service.Create(&domain.Mock{
		Title:       req.Title,
		Description: req.Description,
	})

	if err != nil {
		response.Error(ctx, err)
		return
	}

	response.Success(ctx, domainMockToResponse(data))
}

// Read godoc
// @Summary      读取单条 Mock
// @Description  读取单条 Mock
// @Tags         mock
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        id   path int true "Mock ID"
// @Success      200  {object}  response.successResponse{data=handler.MockResponse} "成功创建 Mock"
// @Failure      400  {object}  response.invalidParamsResponse "参数错误"
// @Failure      500  {object}  response.errorResponse "服务器错误"
// @Router       /mock/{id} [get]
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

// Update godoc
// @Summary      更新 Mock
// @Description  根据ID更新 Mock 信息
// @Tags         mock
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        id   path int true "Mock ID"
// @Param        request body handler.UpdateRequest true "更新 Mock 请求"
// @Success      200  {object}  response.successResponse{data=handler.MockResponse} "成功更新 Mock"
// @Failure      400  {object}  response.invalidParamsResponse "参数错误"
// @Failure      500  {object}  response.errorResponse "服务器错误"
// @Router       /mock/{id} [put]
func (h *HttpHandler) Update(ctx *gin.Context) {
	id, err := h.getID(ctx)
	if err != nil {
		return
	}

	req := new(UpdateRequest)

	if err := bind.BindingRegularAndResponse(ctx, req); err != nil {
		return
	}

	data, err := h.service.Update(&domain.Mock{
		ID:          id,
		Title:       req.Title,
		Description: req.Description,
	})

	if err != nil {
		response.Error(ctx, err)
		return
	}

	response.Success(ctx, domainMockToResponse(data))
}

// Delete godoc
// @Summary      删除 Mock
// @Description  根据ID删除 Mock
// @Tags         mock
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        id   path int true "Mock ID"
// @Success      200  {object}  response.successResponse "成功删除 Mock"
// @Failure      400  {object}  response.invalidParamsResponse "参数错误"
// @Failure      500  {object}  response.errorResponse "服务器错误"
// @Router       /mock/{id} [delete]
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

// List godoc
// @Summary      获取 Mock 列表
// @Description  根据查询参数获取Mock列表，返回当前页数据和total数量
// @Tags         mock
// @Accept       json
// @Produce      json
// @Param        keyword    query     string  false  "关键词搜索"
// @Param        page       query     int     false  "页码" default(1)
// @Param        page_size  query     int     false  "每页数量" default(10)
// @Success      200  {object}  response.successResponse{data=handler.MockListResponse} "Mock列表"
// @Failure      400  {object}  response.invalidParamsResponse "参数错误"
// @Failure      500  {object}  response.errorResponse "服务器错误"
// @Router       /mock [get]
func (h *HttpHandler) List(ctx *gin.Context) {
	req := new(ListRequest)

	if err := bind.BindingRegularAndResponse(ctx, req); err != nil {
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
