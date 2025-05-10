package jwt

import (
	"errors"
	"github.com/golang-jwt/jwt/v5"
	"time"
)

type JWTTokenParams struct {
	Payload  interface{}
	Secret   []byte
	Duration time.Duration
}

// MyClaims 自定义声明结构体并内嵌 jwt.RegisteredClaims
type MyClaims[T any] struct {
	PayLoad T `json:"payLoad"`
	jwt.RegisteredClaims
}

func GenToken[T any](data *JWTTokenParams) (string, error) {
	payload, ok := data.Payload.(T)
	if !ok {
		return "", errors.New("invalid payload type")
	}

	claims := &MyClaims[T]{
		PayLoad: payload,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(data.Duration)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			NotBefore: jwt.NewNumericDate(time.Now()),
			Issuer:    "Lirous-Go-Scaffold",
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(data.Secret)
}

var (
	ErrTokenExpired = errors.New("token已过期")
	ErrInvalidToken = errors.New("无效的token")
	// ErrTokenNotValidYet     = errors.New("token尚未生效")
	// ErrTokenMalformed       = errors.New("token格式错误")
	// ErrTokenInvalidIssuer   = errors.New("token颁发者无效")
	// ErrTokenInvalidAudience = errors.New("token接收者无效")
)

func ParseToken[T any](tokenString string, secret []byte) (myClaims *MyClaims[T], err error) {
	token, err := jwt.ParseWithClaims(tokenString, &MyClaims[T]{}, func(token *jwt.Token) (interface{}, error) {
		return secret, nil
	})

	if err != nil {
		// 对于 JWT v5，直接判断错误类型
		switch {
		case errors.Is(err, jwt.ErrTokenExpired):
			return nil, ErrTokenExpired
		default:
			return nil, ErrInvalidToken
		}
	}

	if claims, ok := token.Claims.(*MyClaims[T]); ok && token.Valid {
		return claims, nil
	}

	return nil, ErrInvalidToken
}
