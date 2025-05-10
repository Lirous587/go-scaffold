package middleware

import (
	"errors"
	"scaffold/internal/domain/admin/model"
	"scaffold/pkg/jwt"
	"scaffold/pkg/response"
	"strings"

	"github.com/gin-gonic/gin"
)

const (
	authHeaderKey = "Authorization"
	bearerPrefix  = "Bearer "
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

type Auth interface {
	Validate() gin.HandlerFunc
}

type adminAuth struct {
	secret []byte
}

func NewAdminAuth(secret []byte) Auth {
	return &adminAuth{
		secret: secret,
	}
}

// Validate 验证管理员 Token
func (auth *adminAuth) Validate() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 1. 从请求头解析 Token
		tokenStr, err := parseTokenFromHeader(c)
		if err != nil {
			response.Error(c, response.NewAppError(response.CodeTokenInvalid, err))
			return
		}

		// 2. 解析 Token
		secret := auth.secret
		claims, err := jwt.ParseToken[model.JwtPayload](tokenStr, secret)
		if err != nil {
			switch {
			case errors.Is(err, jwt.ErrTokenExpired):
				response.Error(c, response.NewAppError(response.CodeTokenExpired, err))
				break
			default:
				response.Error(c, response.NewAppError(response.CodeTokenInvalid, err))
			}

			return
		}

		// 3. 将用户 ID 存入上下文
		c.Set("admin_id", claims.PayLoad.ID)
		c.Next()
	}
}
