package handler

import (
	"net/url"
	"os"
	"scaffold/internal/img/domain"
)

var r2PublicUrlPrefix = ""

func init() {
	r2PublicUrlPrefix = os.Getenv("R2_PUBLIC_URL_PREFIX")
	if r2PublicUrlPrefix == "" {
		panic("加载 R2_PUBLIC_URL_PREFIX 环境变量失败")
	}
}

type ImgResponse struct {
	ID          int64  `json:"id"`
	Url         string `json:"url"`
	Description string `json:"description,omitempty"`
	CreatedAt   int64  `json:"created_at"`
	UpdatedAt   int64  `json:"updated_at"`
}

type UploadRequest struct {
	Path        string `form:"path"`
	Description string `form:"description" binding:"max=60"`
	CategoryID  int64  `form:"category_id"`
}

type DeleteRequest struct {
	Hard bool `form:"hard,default=false"`
}

type ListRequest struct {
	Page       int    `form:"page,default=1" binding:"min=1"`
	PageSize   int    `form:"page_size,default=5" binding:"min=5,max=50"`
	KeyWord    string `form:"keyword" binding:"max=20"`
	Deleted    bool   `form:"deleted,default=false"`
	CategoryID int64  `form:"category_id"`
}

type ImgListResponse struct {
	Total int64          `json:"total"`
	List  []*ImgResponse `json:"list"`
}

func domainImgToResponse(img *domain.Img) *ImgResponse {
	if img == nil {
		return nil
	}

	encodedPath := url.PathEscape(img.Path)

	// 默认访问public
	resp := &ImgResponse{
		ID:          img.ID,
		Url:         r2PublicUrlPrefix + "/" + encodedPath,
		Description: img.Description,
		CreatedAt:   img.CreatedAt.Unix(),
		UpdatedAt:   img.UpdatedAt.Unix(),
	}

	// 如果是要访问已删除文件
	if img.IsDelete() {
		resp.Url = img.Path
	}

	return resp
}

func domainImgsToResponse(imgs []*domain.Img) []*ImgResponse {
	if len(imgs) == 0 {
		return nil
	}
	list := make([]*ImgResponse, 0, len(imgs))

	for _, img := range imgs {
		if img != nil {
			list = append(list, domainImgToResponse(img))
		}
	}

	return list
}

func domainImgListToResponse(data *domain.ImgList) *ImgListResponse {
	if data == nil {
		return nil
	}

	return &ImgListResponse{
		List:  domainImgsToResponse(data.List),
		Total: data.Total,
	}
}
