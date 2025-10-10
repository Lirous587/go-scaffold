package handler

import (
	"fmt"
	"io"
	"math/rand"
	"mime/multipart"
	"net/http"
	"path/filepath"
	"scaffold/internal/common/reskit/response"
	"scaffold/internal/img/domain"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/pkg/errors"
)

type HttpHandler struct {
	service domain.ImgService
}

func NewHttpHandler(service domain.ImgService) *HttpHandler {
	return &HttpHandler{
		service: service,
	}
}

func (h *HttpHandler) getID(ctx *gin.Context) (int64, error) {
	idStr := ctx.Param("id")
	if idStr == "" {
		return 0, errors.New("请传递id参数")
	}
	idInt, err := strconv.Atoi(idStr)
	if err != nil {
		return 0, errors.WithStack(err)
	}
	if idInt == 0 {
		return 0, errors.WithStack(errors.New("无效的id"))
	}
	return int64(idInt), nil
}

func isImage(file multipart.File) (bool, string, error) {
	buf := make([]byte, 512)
	n, _ := file.Read(buf)

	// 复位文件指针
	if _, err := file.Seek(0, io.SeekStart); err != nil {
		return false, "", fmt.Errorf("file.Seek 复位文件指针失败,reason:%v", err)
	}

	contentType := http.DetectContentType(buf[:n])
	switch contentType {
	case "image/jpeg", "image/png", "image/gif", "image/webp", "image/avif", "image/bmp", "image/svg+xml":
		return true, contentType, nil
	default:
		return false, contentType, nil
	}
}

func generateImgPath(ext string) string {
	now := time.Now().Format("2006_01_02_150405.000")
	random := rand.Intn(1000000)
	return fmt.Sprintf("%s_%d%s", now, random, ext)
}

func getExtByContentType(realType string) (ext string) {
	switch realType {
	case "image/jpeg":
		ext = ".jpg"
	case "image/png":
		ext = ".png"
	case "image/gif":
		ext = ".gif"
	case "image/webp":
		ext = ".webp"
	case "image/avif":
		ext = ".avif"
	case "image/bmp":
		ext = ".bmp"
	case "image/svg+xml":
		ext = ".svg"
	default:
		ext = ""
	}
	return ext
}

// Upload godoc
// @Summary      上传图片
// @Description  上传单张图片（支持 jpeg/png/gif/webp/avif/bmp/svg）
// @Tags         img
// @Accept       multipart/form-data
// @Produce      json
// @Param        object      formData  file  true  "图片文件"
// @Param        path        formData  string false "自定义图片路径"
// @Param        description formData  string false "图片描述"
// @Param        category_id formData  int64  false "分类id"
// @Success      200 {object} response.successResponse{data=handler.ImgResponse} "请求成功"
// @Failure      400 {object} response.invalidParamsResponse "参数错误"
// @Failure      500 {object} response.errorResponse "服务器错误"
// @Security     BearerAuth
// @Router       /v1/img/upload [post]
func (h *HttpHandler) Upload(ctx *gin.Context) {
	fileHeader, _ := ctx.FormFile("object")
	if fileHeader == nil {
		response.InvalidParams(ctx, errors.New("未携带对象上传"))
		return
	}

	// 将 *multipart.FileHeader 转为 io.Reader
	file, _ := fileHeader.Open()
	defer file.Close()

	ok, realType, err := isImage(file)
	if err != nil {
		response.Error(ctx, errors.Errorf("isImage执行失败: %s", err))
		return
	}
	if !ok {
		response.InvalidParams(ctx, errors.Errorf("仅支持基本图片类型上传，实际类型: %s", realType))
		return
	}

	req := new(UploadRequest)

	if err := ctx.ShouldBind(req); err != nil {
		response.InvalidParams(ctx, err)
		return
	}

	var imgPath string

	// 无path则生成path
	if strings.TrimSpace(req.Path) != "" {
		ext := filepath.Ext(req.Path)
		if ext == "" {
			imgPath = req.Path + getExtByContentType(realType)
		} else {
			imgPath = req.Path
		}
	} else {
		imgPath = generateImgPath(getExtByContentType(realType))
	}

	res, err := h.service.Upload(file, &domain.Img{
		Path:        imgPath,
		Description: req.Description,
	},
		req.CategoryID,
	)

	if err != nil {
		response.Error(ctx, err)
		return
	}

	response.Success(ctx, domainImgToResponse(res))
}

// Delete godoc
// @Summary      删除图片
// @Tags         img
// @Accept       json
// @Produce      json
// @Param        id    path   int64  true  "id"
// @Param        hard  query  bool   false "是否硬删除"
// @Success      200 {object} response.successResponse "请求成功"
// @Failure      400 {object} response.invalidParamsResponse "参数错误"
// @Failure      500 {object} response.errorResponse "服务器错误"
// @Security     BearerAuth
// @Router       /v1/img/{id} [delete]
func (h *HttpHandler) Delete(ctx *gin.Context) {
	req := new(DeleteRequest)

	if err := ctx.ShouldBindQuery(req); err != nil {
		response.InvalidParams(ctx, err)
		return
	}

	id, err := h.getID(ctx)
	if err != nil {
		response.InvalidParams(ctx, err)
		return
	}

	if req.Hard {
		if err := h.service.Delete(id, true); err != nil {
			response.Error(ctx, err)
			return
		}
	} else {
		if err := h.service.Delete(id, false); err != nil {
			response.Error(ctx, err)
			return
		}
	}

	response.Success(ctx)
}

// List godoc
// @Summary      图片列表
// @Tags         img
// @Accept       json
// @Produce      json
// @Param        page        query int    false "页码"
// @Param        page_size   query int    false "每页数量"
// @Param        keyword     query string false "关键词"
// @Param        deleted     query bool   false "是否查询回收站图片"
// @Param        category_id query int64  false "分类id"
// @Success      200 {object} response.successResponse{data=handler.ImgListResponse} "请求成功"
// @Failure      400 {object} response.invalidParamsResponse "参数错误"
// @Failure      500 {object} response.errorResponse "服务器错误"
// @Security     BearerAuth
// @Router       /v1/img [get]
func (h *HttpHandler) List(ctx *gin.Context) {
	req := new(ListRequest)

	if err := ctx.ShouldBindQuery(req); err != nil {
		response.InvalidParams(ctx, err)
		return
	}

	list, err := h.service.List(&domain.ImgQuery{
		Keyword:    req.KeyWord,
		Page:       req.Page,
		PageSize:   req.PageSize,
		Deleted:    req.Deleted,
		CategoryID: req.CategoryID,
	})

	if err != nil {
		response.Error(ctx, err)
		return
	}

	response.Success(ctx, domainImgListToResponse(list))
}

// ClearRecycleBin godoc
// @Summary      清空回收站图片
// @Tags         img
// @Accept       json
// @Produce      json
// @Param        id path int64 true "id"
// @Success      200 {object} response.successResponse "请求成功"
// @Failure      400 {object} response.invalidParamsResponse "参数错误"
// @Failure      500 {object} response.errorResponse "服务器错误"
// @Security     BearerAuth
// @Router       /v1/img/recycle/{id} [delete]
func (h *HttpHandler) ClearRecycleBin(ctx *gin.Context) {
	id, err := h.getID(ctx)
	if err != nil {
		response.InvalidParams(ctx, err)
		return
	}

	if err := h.service.ClearRecycleBin(id); err != nil {
		response.Error(ctx, err)
		return
	}

	response.Success(ctx)
}

// RestoreFromRecycleBin godoc
// @Summary      恢复回收站图片
// @Tags         img
// @Accept       json
// @Produce      json
// @Param        id path int64 true "id"
// @Success      200 {object} response.successResponse{data=handler.ImgResponse} "请求成功"
// @Failure      400 {object} response.invalidParamsResponse "参数错误"
// @Failure      500 {object} response.errorResponse "服务器错误"
// @Security     BearerAuth
// @Router       /v1/img/recycle/{id} [put]
func (h *HttpHandler) RestoreFromRecycleBin(ctx *gin.Context) {
	id, err := h.getID(ctx)
	if err != nil {
		response.InvalidParams(ctx, err)
		return
	}
	res, err := h.service.RestoreFromRecycleBin(id)
	if err != nil {
		response.Error(ctx, err)
		return
	}

	response.Success(ctx, res)
}

func (h *HttpHandler) ListenDeleteQueue() {
	h.service.ListenDeleteQueue()
}

// --- 分类管理 ---

// CreateCategory godoc
// @Summary      创建图片分类
// @Tags         img-category
// @Accept       json
// @Produce      json
// @Param        request body handler.CreateCategoryRequest true "请求参数"
// @Success      200 {object} response.successResponse{data=handler.CategoryResponse} "请求成功"
// @Failure      400 {object} response.invalidParamsResponse "参数错误"
// @Failure      500 {object} response.errorResponse "服务器错误"
// @Security     BearerAuth
// @Router       /v1/img/category [post]
func (h *HttpHandler) CreateCategory(ctx *gin.Context) {
	req := new(CreateCategoryRequest)
	if err := ctx.ShouldBindJSON(req); err != nil {
		response.InvalidParams(ctx, err)
		return
	}

	res, err := h.service.CreateCategory(&domain.Category{
		Title:  req.Title,
		Prefix: req.Prefix,
	})

	if err != nil {
		response.Error(ctx, err)
		return
	}

	response.Success(ctx, domainCategoryToResponse(res))
}

// UpdateCategory godoc
// @Summary      更新图片分类
// @Tags         img-category
// @Accept       json
// @Produce      json
// @Param        id      path   int64  true  "分类id"
// @Param        request body   handler.UpdateCategoryRequest true "请求参数"
// @Success      200 {object} response.successResponse{data=handler.CategoryResponse} "请求成功"
// @Failure      400 {object} response.invalidParamsResponse "参数错误"
// @Failure      500 {object} response.errorResponse "服务器错误"
// @Security     BearerAuth
// @Router       /v1/img/category/{id} [put]
func (h *HttpHandler) UpdateCategory(ctx *gin.Context) {
	req := new(UpdateCategoryRequest)
	if err := ctx.ShouldBindJSON(req); err != nil {
		response.InvalidParams(ctx, err)
		return
	}

	res, err := h.service.UpdateCategory(&domain.Category{
		ID:     req.ID,
		Title:  req.Title,
		Prefix: req.Prefix,
	})

	if err != nil {
		response.Error(ctx, err)
		return
	}

	response.Success(ctx, domainCategoryToResponse(res))
}

// DeleteCategory godoc
// @Summary      删除图片分类
// @Tags         img-category
// @Accept       json
// @Produce      json
// @Param        id path int64 true "分类id"
// @Success      200 {object} response.successResponse "请求成功"
// @Failure      400 {object} response.invalidParamsResponse "参数错误"
// @Failure      500 {object} response.errorResponse "服务器错误"
// @Security     BearerAuth
// @Router       /v1/img/category/{id} [delete]
func (h *HttpHandler) DeleteCategory(ctx *gin.Context) {
	id, err := h.getID(ctx)
	if err != nil {
		response.InvalidParams(ctx, err)
		return
	}

	if err := h.service.DeleteCategory(id); err != nil {
		response.Error(ctx, err)
		return
	}

	response.Success(ctx)
}

// ListCategories godoc
// @Summary      分类列表
// @Tags         img-category
// @Accept       json
// @Produce      json
// @Success      200 {object} response.successResponse{data=[]handler.CategoryResponse} "请求成功"
// @Failure      500 {object} response.errorResponse "服务器错误"
// @Security     BearerAuth
// @Router       /v1/img/categories [get]
func (h *HttpHandler) ListCategories(ctx *gin.Context) {
	res, err := h.service.ListCategories()
	if err != nil {
		response.Error(ctx, err)
		return
	}

	response.Success(ctx, domainCategoriesToResponse(res))
}
