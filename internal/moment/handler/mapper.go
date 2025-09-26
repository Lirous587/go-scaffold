package handler

import (
	"blog-v4/internal/moment/domain"
	"time"
)

type MomentResponse struct {
	ID          int64  `json:"id"`
	Title       string `json:"title"`
	Description string `json:"description,omitempty"`
	CreatedAt   string `json:"created_at"`
	UpdatedAt   string `json:"updated_at"`
	DeletedAt   string `json:"deleted_at,omitempty"`
}

func domainMomentToResponse(moment *domain.Moment) *MomentResponse {
	if moment == nil {
		return nil
	}

	return &MomentResponse{
		ID:        moment.ID,
		Title:     moment.Title,
		Content:   moment.Content,
		Location:  moment.Location,
		CreatedAt: moment.CreatedAt.Format(time.DateTime),
		UpdatedAt: moment.UpdatedAt.Format(time.DateTime),
	}
}

func domainMomentsToResponse(moments []*domain.Moment) []*MomentResponse {
	if len(moments) == 0 {
		return nil
	}

	ret := make([]*MomentResponse, 0, len(moments))

	for _, moment := range moments {
		if moment != nil {
			ret = append(ret, domainMomentToResponse(moment))
		}
	}
	return ret
}

type CreateRequest struct {
	Title    string `json:"title" binding:"required,max=30"`
	Content  string `json:"content" binding:"required"`
	CoverURL string `json:"cover_url" binding:"required,max=255"`
	Location string `json:"location" binding:"max=30"`
}

type UpdateRequest struct {
	Title    string `json:"title" binding:"required,max=30"`
	Content  string `json:"content" binding:"required"`
	CoverURL string `json:"cover_url" binding:"required,max=255"`
	Location string `json:"location" binding:"max=30"`
}

type ListRequest struct {
	Page     int    `form:"page,default=1" binding:"min=1"`
	PageSize int    `form:"page_size,default=5" binding:"min=5,max=20"`
	KeyWord  string `form:"keyword" binding:"max=20"`
}

type RandomNRequest struct {
	Count int8 `form:"count,default=20" binding:"max=100"`
}

type MomentPagesResponse struct {
	Pages int               `json:"pages"`
	List  []*MomentResponse `json:"list"`
}

func domainMomentWithPagesToResponse(data *domain.MomentPages) *MomentPagesResponse {
	if data == nil {
		return nil
	}

	return &MomentPagesResponse{
		Pages: data.Pages,
		List:  domainMomentsToResponse(data.List),
	}
}
