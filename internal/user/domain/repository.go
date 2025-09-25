package domain

type UserRepository interface {
	// 基础 CRUD
	FindByID(id int64) (*User, error)
	FindByEmail(email string) (*User, error)
	Create(user *User) (*User, error)
	Update(user *User) (*User, error)

	// OAuth 相关
	FindByOAuthID(provider, oauthID string) (*User, error)
	UpdateLastLogin(id int64) error

	// 辅助方法
	EmailExists(email string) (bool, error)
}

type TokenCache interface {
	GenRefreshToken(payload *JwtPayload) (string, error)
	ValidateRefreshToken(refreshToken string) (*JwtPayload, error)
	RemoveRefreshToken(refreshToken string) error
}
