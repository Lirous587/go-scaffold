package domain

import (
	"io"
)

type ImgService interface {
	Upload(src io.Reader, img *Img, categoryID int64) (*Img, error)
	Delete(id int64, hard ...bool) error
	List(query *ImgQuery) (*ImgList, error)
	ClearRecycleBin(id int64) error
	ListenDeleteQueue()
	RestoreFromRecycleBin(id int64) (*Img, error)

	//	分类
	CreateCategory(category *Category) (*Category, error)
	UpdateCategory(category *Category) (*Category, error)
	DeleteCategory(id int64) error
	ListCategories() (categories []*Category, err error)
}
