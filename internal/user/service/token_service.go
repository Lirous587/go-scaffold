package service

import (
	"github.com/joho/godotenv"
	"github.com/pkg/errors"
	"os"
	"scaffold/internal/common/jwt"
	"scaffold/internal/common/utils"
	"scaffold/internal/user/domain"
	"strconv"
	"time"
)

var (
	secret string
	expire time.Duration
)

func init() {
	_ = godotenv.Load()
	secret = os.Getenv("JWT_SECRET")
	expireMinuteStr := os.Getenv("JWT_EXPIRE_MINUTE")
	if secret == "" || expireMinuteStr == "" {
		panic("加载环境变量失败")
	}
	expireMinute, err := strconv.Atoi(expireMinuteStr)
	if err != nil {
		panic(err)
	}
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

func (t tokenService) GenerateAccessToken(payload domain.JwtPayload) (string, error) {
	token, err := jwt.GenToken[domain.JwtPayload](payload, secret, expire)
	return token, errors.WithStack(err)
}

func (t tokenService) ValidateAccessToken(token string) (claim domain.JwtPayload, isExpire bool, err error) {
	claims, err := jwt.ParseToken[domain.JwtPayload](token, secret)
	if err != nil {
		switch {
		case errors.Is(err, jwt.ErrTokenExpired):
			return domain.JwtPayload{}, true, err
		default:
			return domain.JwtPayload{}, false, err
		}
	}

	return claims.PayLoad, false, nil
}

func (t tokenService) RefreshAccessToken(payload domain.JwtPayload, refreshToken string) (string, error) {
	if err := t.tokenCache.ValidateRefreshToken(payload, refreshToken); err != nil {
		return "", err
	}
	// 为后续扩展jwt字段保留空间
	user, err := t.userRepo.FindByID(payload.UserID)
	if err != nil {
		return "", err
	}
	newPayload := domain.JwtPayload{
		UserID:     user.ID,
		RandomCode: utils.GenRandomCodeForJWT(),
	}
	return t.GenerateAccessToken(newPayload)
}

func (t tokenService) GenerateRefreshToken(payload domain.JwtPayload) (string, error) {
	return t.tokenCache.GenRefreshToken(payload)
}

func (t tokenService) ResetRefreshTokenExpiry(payload domain.JwtPayload) error {
	return t.tokenCache.ResetRefreshTokenExpiry(payload)
}
