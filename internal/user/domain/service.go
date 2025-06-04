package domain

type UserService interface {
	AuthenticateWithOAuth(provider string, userInfo *OAuthUserInfo) (*User2Token, error)
	RefreshUserToken(payload JwtPayload, refreshToken string) (*User2Token, error)

	GetUser(userID string) (*User, error)
	UpdateUserProfile(userID string, updates *UserProfileUpdate) (*User, error)
}

type TokenService interface {
	GenerateAccessToken(payload JwtPayload) (string, error)
	ValidateAccessToken(token string) (payload JwtPayload, isExpire bool, err error)
	RefreshAccessToken(domain JwtPayload, refreshToken string) (string, error)

	GenerateRefreshToken(payload JwtPayload) (string, error)
	ResetRefreshTokenExpiry(domain JwtPayload) error
}
