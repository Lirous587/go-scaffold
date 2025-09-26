package domain

import (
	"context"
	"database/sql"
)

type ImgRepository interface {
	FindByID(id int64, deleted ...bool) (*Img, error)
	ExistByPath(path string) (bool, error)

	// WithTX 开启事务,允许传递空值
	WithTX(ctx context.Context, opts *sql.TxOptions) (*sql.Tx, error)
	Create(img *Img, categoryID int64) (*Img, error)
	Delete(id int64, hard bool) error
	Restore(id int64) (*Img, error)
	List(query *ImgQuery) (*ImgList, error)

	CreateCategory(category *Category) (*Category, error)
	UpdateCategory(category *Category) (*Category, error)
	DeleteCategory(id int64) error
	ListCategories() ([]*Category, error)
	FindCategoryByID(id int64) (*Category, error)
	FindCategoryByTitle(title string) (*Category, error)
	CategoryExistByTitle(title string) (bool, error)
	CategoryExistByID(id int64) (bool, error)
	CountCategory() (int64, error)
	IsCategoryExistImg(id int64) (bool, error)
}

type ImgMsgQueue interface {
	AddToDeleteQueue(imgID int64) error
	ListenDeleteQueue(onExpire func(imgID int64))
	//SendDeleteMsg(imgID int64) error
	RemoveFromDeleteQueue(imgID int64) error
}
