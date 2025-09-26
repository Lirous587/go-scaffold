package domain

import (
	"time"
)

type Img struct {
	ID          int64
	Path        string
	Description string
	CreatedAt   time.Time
	UpdatedAt   time.Time
	DeletedAt   time.Time
	deleted     bool
}

func (img *Img) SetDeletedStatus(deleted bool) {
	img.deleted = deleted
}

func (img *Img) IsDelete() bool {
	return img.deleted
}

type ImgQuery struct {
	Keyword    string
	Page       int
	PageSize   int
	Deleted    bool
	CategoryID int64
}

type ImgList struct {
	List  []*Img
	Total int64
}

type Category struct {
	ID        int64
	Title     string
	Prefix    string
	CreatedAt time.Time
}
