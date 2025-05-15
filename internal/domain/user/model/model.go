package model

import (
	"comment/internal/domain/role/model"
	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	UserName     string       `gorm:"size:32;not null"`
	Email        string       `gorm:"uniqueIndex;not null;size:255"`
	HashPassword []byte       `gorm:"type:varchar(255)"`
	GithubID     int          `gorm:"uniqueIndex;type:int;not null"`
	Roles        []model.Role `gorm:"many2many:user_roles"`
}
