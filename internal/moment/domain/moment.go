package domain

import (
	"time"
)

type Moment struct {
	ID          int64
	Title       string
	Content     string
	Location    string
	Coordinates *MomentCoordinate
	CreatedAt   time.Time
	UpdatedAt   time.Time
	CoverURL    string
}

type MomentQuery struct {
	Keyword  string
	Page     int
	PageSize int
}

type MomentPages struct {
	Pages int
	List  []*Moment
}
