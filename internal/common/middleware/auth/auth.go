package auth

import (
	"github.com/joho/godotenv"
	"github.com/pkg/errors"
	"os"
	"scaffold/internal/common/jwt"
	"scaffold/internal/common/server/response"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
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

const (
	authHeaderKey = "Authorization"
	bearerPrefix  = "Bearer "
)

type JwtPayload struct {
	ID        uint   `json:"id"`
	LoginType string `json:"login_type"`
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

func GenUserToken[T any](payload T) (string, error) {
	token, err := jwt.GenToken[T](payload, secret, expire)
	return token, errors.WithStack(err)
}
