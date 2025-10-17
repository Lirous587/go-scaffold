package service

import (
	"scaffold/internal/common/jwt"
	"scaffold/internal/common/utils"
	"scaffold/internal/user/domain"
	"time"

	"github.com/joho/godotenv"
	"github.com/pkg/errors"
)

var (
	secret string
	expire time.Duration
)

func init() {
	_ = godotenv.Load()
	secret = utils.GetEnv("JWT_SECRET")
	expireMinute := utils.GetEnvAsInt("JWT_EXPIRE_MINUTE")
	expire = time.Minute * time.Duration(expireMinute)
}

type tokenService struct {
	tokenCache domain.TokenCache
	userRepo   domain.UserRepository
}

func NewTokenService(tokenCache domain.TokenCache, userRepo domain.UserRepository) domain.TokenService {
	return &tokenService{
		tokenCache: tokenCache,
		userRepo:   userRepo,
	}
}

func (t *tokenService) GenerateAccessToken(payload *domain.JwtPayload) (string, error) {
	token, err := jwt.GenToken[domain.JwtPayload](payload, secret, expire)
	return token, errors.WithStack(err)
}

func (t *tokenService) ValidateAccessToken(token string) (isExpire bool, err error) {
	_, err = jwt.ParseToken[domain.JwtPayload](token, secret)
	if err != nil {
		switch {
		case errors.Is(err, jwt.ErrTokenExpired):
			return true, err
		default:
			return false, err
		}
	}

	return false, nil
}

func (t *tokenService) ParseAccessToken(token string) (payload *domain.JwtPayload, err error) {
	claims, err := jwt.ParseToken[domain.JwtPayload](token, secret)
	if err != nil {
		return nil, err
	}

	return claims.PayLoad, nil
}

func (t *tokenService) RefreshAccessToken(refreshToken string) (string, error) {
	payload, err := t.tokenCache.ValidateRefreshToken(refreshToken)
	if err != nil {
		return "", err
	}

	// 为后续扩展jwt携带的相应user字段保留空间
	user, err := t.userRepo.FindByID(payload.UserID)
	if err != nil {
		return "", err
	}

	newPayload := &domain.JwtPayload{
		UserID: user.ID,
	}

	return t.GenerateAccessToken(newPayload)
}

func (t *tokenService) GenerateRefreshToken(payload *domain.JwtPayload) (string, error) {
	return t.tokenCache.GenRefreshToken(payload)
}

func (t *tokenService) RemoveRefreshToken(refreshToken string) error {
	return t.tokenCache.RemoveRefreshToken(refreshToken)
}
