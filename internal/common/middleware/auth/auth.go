package auth

import (
	"errors"
	"github.com/joho/godotenv"
	"os"
	"scaffold/internal/common/jwt"
	"scaffold/internal/common/server/response"
	"scaffold/internal/feature/domain/user"
	"strings"

	"github.com/gin-gonic/gin"
)

const (
	authHeaderKey = "Authorization"
	bearerPrefix  = "Bearer "
)

type JwtPayload struct {
	ID        uint           `json:"id"`
	LoginType user.LoginType `json:"login_type"`
}

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

var secret string

func init() {
	_ = godotenv.Load()
	secret = os.Getenv("JWT_SECRET")
	if secret == "" {
		panic("加载JWT_SECRET环境变量失败")
	}
}

func Validate() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 1. 从请求头解析 Token
		tokenStr, err := parseTokenFromHeader(c)
		if err != nil {
			response.Error(c, response.NewAppError(response.CodeTokenInvalid, err))
			return
		}

		// 2. 解析 Token
		claims, err := jwt.ParseToken[JwtPayload](tokenStr, secret)
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

		// 3. 将用户 相关信息存入上下文
		c.Set("user_id", claims.PayLoad.ID)
		c.Set("login_type", claims.PayLoad.LoginType)
		c.Next()
	}
}
