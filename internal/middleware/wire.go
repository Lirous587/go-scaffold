//go:build wireinject
// +build wireinject

package middleware

import (
	"comment/pkg/config"
	"github.com/google/wire"
)

// InitializeAdminAPI 初始化Admin模块的API
func InitUserAuth() (Auth, error) {
	wire.Build(
		provideJWTSecret,
		NewUserAuth,
	)
	return nil, nil
}

// provideJWTSecret 提供JWT Secret
func provideJWTSecret() []byte {
	return []byte(config.Cfg.JWT.Secret)
}
