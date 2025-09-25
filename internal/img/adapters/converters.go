package adapters

import (
	"scaffold/internal/common/orm"
	"scaffold/internal/img/domain"
	"github.com/volatiletech/null/v8"
)

func domainImgToORM(img *domain.Img) *orm.Img {
	if img == nil {
		return nil
	}

	ormImg := &orm.Img{
		ID:		img.ID,
		Path:		img.Path,
		UpdatedAt:	img.UpdatedAt,
	}

	// 处理null项
	if img.Description != "" {
		ormImg.Description = null.StringFrom(img.Description)
	}

	return ormImg
}

func ormImgToDomain(ormImg *orm.Img, isDeleted ...bool) *domain.Img {
	if ormImg == nil {
		return nil
	}

	deleted := false
	if len(isDeleted) > 0 {
		deleted = isDeleted[0]
	}

	img := &domain.Img{
		ID:		ormImg.ID,
		Path:		ormImg.Path,
		CreatedAt:	ormImg.CreatedAt,
		UpdatedAt:	ormImg.UpdatedAt,
	}

	img.SetDeletedStatus(deleted)

	// 处理null项
	if ormImg.Description.Valid {
		img.Description = ormImg.Description.String
	}

	if ormImg.DeletedAt.Valid {
		img.DeletedAt = ormImg.DeletedAt.Time
	}

	return img
}

func ormImgsToDomain(ormImgs []*orm.Img, isDeleted ...bool) []*domain.Img {
	if len(ormImgs) == 0 {
		return nil
	}

	deleted := false
	if len(isDeleted) > 0 {
		deleted = isDeleted[0]
	}

	imgs := make([]*domain.Img, 0, len(ormImgs))
	for _, ormImg := range ormImgs {
		if ormImg == nil {
			continue
		}
		img := ormImgToDomain(ormImg, deleted)
		imgs = append(imgs, img)
	}

	return imgs
}

func domainCategoryToORM(category *domain.Category) *orm.ImgCategory {
	if category == nil {
		return nil
	}

	ormImg := &orm.ImgCategory{
		ID:	category.ID,
		Title:	category.Title,
		Prefix:	category.Prefix,
	}

	// 处理null项

	return ormImg
}

func ormCategoryToDomain(ormCategory *orm.ImgCategory) *domain.Category {
	if ormCategory == nil {
		return nil
	}

	img := &domain.Category{
		ID:		ormCategory.ID,
		Title:		ormCategory.Title,
		Prefix:		ormCategory.Prefix,
		CreatedAt:	ormCategory.CreatedAt,
	}

	// 处理null项
	return img
}

func ormCategoriesToDomain(ormCategories []*orm.ImgCategory) []*domain.Category {
	if len(ormCategories) == 0 {
		return nil
	}

	categories := make([]*domain.Category, 0, len(ormCategories))
	for _, ormCategory := range ormCategories {
		if ormCategory == nil {
			continue
		}
		category := ormCategoryToDomain(ormCategory)
		categories = append(categories, category)
	}

	return categories
}
