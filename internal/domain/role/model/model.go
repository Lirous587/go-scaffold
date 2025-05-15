package model

import "gorm.io/gorm"

type Role struct {
	gorm.Model
	Name        string `gorm:"size:32;not null"`
	Description string `gorm:"size:32"`
}

type CreateRoleReq struct {
	Name        string `json:"name" binding:"required,max=32"`
	Description string `json:"description" binding:"max=32"`
}

type UpdateRoleReq struct {
	Name        string `json:"name" binding:"required,max=32"`
	Description string `json:"description" binding:"max=32"`
}
