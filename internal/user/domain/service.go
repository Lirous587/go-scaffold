package domain

type UserService interface {
	AuthenticateWithOAuth(provider string, userInfo *OAuthUserInfo) (*User2Token, error)
	RefreshUserToken(refreshToken string) (*User2Token, error)
	GetUser(id int64) (*User, error)
}

type TokenService interface {
	GenerateAccessToken(payload *JwtPayload) (string, error)
	ValidateAccessToken(token string) (isExpire bool, err error)
	ParseAccessToken(token string) (payload *JwtPayload, err error)

	RefreshAccessToken(refreshToken string) (string, error)

	GenerateRefreshToken(payload *JwtPayload) (string, error)
	RemoveRefreshToken(refreshToken string) error
}
