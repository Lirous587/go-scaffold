package handler

import (
	"scaffold/internal/img/domain"
)

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

func domainCategoryToResponse(category *domain.Category) *CategoryResponse {
	if category == nil {
		return nil
	}

	resp := &CategoryResponse{
		ID:        category.ID,
		Title:     category.Title,
		Prefix:    category.Prefix,
		CreatedAt: category.CreatedAt.Unix(),
	}

	return resp
}

func domainCategoriesToResponse(categories []*domain.Category) []*CategoryResponse {
	if len(categories) == 0 {
		return nil
	}
	list := make([]*CategoryResponse, 0, len(categories))

	for _, category := range categories {
		if category != nil {
			list = append(list, domainCategoryToResponse(category))
		}
	}

	return list
}
