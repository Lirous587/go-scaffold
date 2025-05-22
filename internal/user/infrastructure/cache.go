package infrastructure

import (
	"scaffold/internal/user/model"
)

type UserCache interface {
	GenRefreshToken(payload *model.JwtPayload) (string, error)
	ValidateRefreshToken(payload *model.JwtPayload, refreshToken string) error
	ResetRefreshTokenExpiry(payload *model.JwtPayload) error
}
