package handler

import (
	"net/url"
	"scaffold/internal/img/domain"
)

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
