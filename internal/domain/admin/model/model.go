package model

import (
	"gorm.io/gorm"
)

type Admin struct {
	gorm.Model
	Email        string `gorm:"uniqueIndex;not null;size:255"`
	HashPassword []byte `gorm:"type:varchar(255);not null"`
}

type InitReq struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,max=70"`
}

type LoginReq struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"`
}

type LoginRes struct {
	Token        string     `json:"token"`
	RefreshToken string     `json:"refresh_token"`
	Payload      JwtPayload `json:"payload"`
}

type JwtPayload struct {
	ID uint `json:"id" binding:"required"`
}

type RefreshTokenReq = JwtPayload

type RefreshTokenRes struct {
	Token string `json:"token"`
}
