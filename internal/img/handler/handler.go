package handler

import (
	"scaffold/internal/common/reskit/response"
	"scaffold/internal/img/domain"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/pkg/errors"
	"io"
	"math/rand"
	"mime/multipart"
	"net/http"
	"path/filepath"
	"strconv"
	"strings"
	"time"
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
	idInt, err := strconv.Atoi(idStr)
	if err != nil {
		return 0, err
	}
	if idInt == 0 {
		return 0, errors.New("无效的id")
	}
	return int64(idInt), err
}

func isImage(file multipart.File) (bool, string) {
	buf := make([]byte, 512)
	n, _ := file.Read(buf)
	file.Seek(0, io.SeekStart)	// 复位文件指针
	contentType := http.DetectContentType(buf[:n])
	switch contentType {
	case "image/jpeg", "image/png", "image/gif", "image/webp", "image/avif", "image/bmp", "image/svg+xml":
		return true, contentType
	default:
		return false, contentType
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

func (h *HttpHandler) Upload(ctx *gin.Context) {
	fileHeader, _ := ctx.FormFile("object")
	if fileHeader == nil {
		response.InvalidParams(ctx, errors.New("未携带对象上传"))
		return
	}

	// 将 *multipart.FileHeader 转为 io.Reader
	file, _ := fileHeader.Open()
	defer file.Close()

	ok, realType := isImage(file)
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
		Path:		imgPath,
		Description:	req.Description,
	},
		req.CategoryID,
	)

	if err != nil {
		response.Error(ctx, err)
		return
	}

	response.Success(ctx, domainImgToResponse(res))
}

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

func (h *HttpHandler) List(ctx *gin.Context) {
	req := new(ListRequest)

	if err := ctx.ShouldBindQuery(req); err != nil {
		response.InvalidParams(ctx, err)
		return
	}

	list, err := h.service.List(&domain.ImgQuery{
		Keyword:	req.KeyWord,
		Page:		req.Page,
		PageSize:	req.PageSize,
		Deleted:	req.Deleted,
		CategoryID:	req.CategoryID,
	})

	if err != nil {
		response.Error(ctx, err)
		return
	}

	response.Success(ctx, domainImgWithPagesToResponse(list))
}

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

func (h *HttpHandler) ListenDeleteQueue() {
	h.service.ListenDeleteQueue()
}

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

// --- 分类管理 ---

func (h *HttpHandler) CreateCategory(ctx *gin.Context) {
	req := new(CreateCategoryRequest)
	if err := ctx.ShouldBindJSON(req); err != nil {
		response.InvalidParams(ctx, err)
		return
	}

	res, err := h.service.CreateCategory(&domain.Category{
		Title:	req.Title,
		Prefix:	req.Prefix,
	})

	if err != nil {
		response.Error(ctx, err)
		return
	}

	response.Success(ctx, domainCategoryToResponse(res))
}

func (h *HttpHandler) UpdateCategory(ctx *gin.Context) {
	id, err := h.getID(ctx)
	if err != nil {
		response.InvalidParams(ctx, err)
		return
	}

	req := new(UpdateCategoryRequest)
	if err := ctx.ShouldBindJSON(req); err != nil {
		response.InvalidParams(ctx, err)
		return
	}

	res, err := h.service.UpdateCategory(&domain.Category{
		ID:	id,
		Title:	req.Title,
		Prefix:	req.Prefix,
	})

	if err != nil {
		response.Error(ctx, err)
		return
	}

	response.Success(ctx, domainCategoryToResponse(res))
}

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

func (h *HttpHandler) ListCategories(ctx *gin.Context) {
	res, err := h.service.ListCategories()
	if err != nil {
		response.Error(ctx, err)
		return
	}

	response.Success(ctx, domainCategoriesToResponse(res))
}
