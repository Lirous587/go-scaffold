package auth

import (
	"scaffold/internal/common/reskit/codes"
	"scaffold/internal/common/reskit/response"
	"scaffold/internal/user/adapters"
	"scaffold/internal/user/domain"
	"scaffold/internal/user/service"
	"github.com/pkg/errors"
	"strings"

	"github.com/gin-gonic/gin"
)

var tokenServer domain.TokenService

func init() {
	tokenCache := adapters.NewTokenRedisCache()
	userRepo := adapters.NewUserPSQLRepository()
	tokenServer = service.NewTokenService(tokenCache, userRepo)
}

const (
	authHeaderKey	= "Authorization"
	bearerPrefix	= "Bearer "
)

// 解析 Authorization 头部的 Token
func parseTokenFromHeader(c *gin.Context) (string, error) {
	authHeader := c.GetHeader(authHeaderKey)
	if authHeader == "" {
		return "", errors.New("token为空")
	}

	if !strings.HasPrefix(authHeader, bearerPrefix) {
		return "", errors.New("token格式错误")
	}

	return strings.TrimPrefix(authHeader, bearerPrefix), nil
}

func JWTValidate() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 1. 从请求头解析 Token
		tokenStr, err := parseTokenFromHeader(c)
		if err != nil {
			response.Error(c, codes.ErrTokenFormatInvalid)
			return
		}

		// 2. 验证token
		isExpire, err := tokenServer.ValidateAccessToken(tokenStr)
		if err != nil {
			if isExpire {
				response.Error(c, codes.ErrTokenExpired)
			} else {
				response.Error(c, codes.ErrTokenInvalid)
			}
			return
		}

		// 3. 解析 Token
		payload, err := tokenServer.ParseAccessToken(tokenStr)
		if err != nil {
			response.Error(c, codes.ErrTokenInvalid)
			return
		}

		// 3. 将用户 相关信息存入上下文
		c.Set("user_id", payload.UserID)

		c.Next()
	}
}
