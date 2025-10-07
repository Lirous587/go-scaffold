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
}

func (img *Img) IsDeleted() bool {
	// deleted_time非零值则为软删除记录
	return !img.DeletedAt.IsZero()
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
