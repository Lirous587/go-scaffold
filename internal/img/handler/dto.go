package handler

import (
	"os"
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

type CategoryResponse struct {
	ID        int64  `json:"id"`
	Title     string `json:"title"`
	Prefix    string `json:"prefix"`
	CreatedAt int64  `json:"created_at"`
}

type CreateCategoryRequest struct {
	Title  string `json:"title" binding:"required,max=10"`
	Prefix string `json:"prefix" binding:"required,max=20,slug"`
}

type UpdateCategoryRequest struct {
	Title  string `json:"title" binding:"max=10"`
	Prefix string `json:"prefix" binding:"max=20"`
}
