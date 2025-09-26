package domain

import (
	"time"
)

type Mock struct {
	ID          int64
	Title       string
	Description string
	CreatedAt   time.Time
	UpdatedAt   time.Time
	DeletedAt   time.Time
}

type MockQuery struct {
	Keyword  string
	Page     int
	PageSize int
}

type MockList struct {
	Total int64
	List  []*Mock
}
