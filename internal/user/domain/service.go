package domain

// 纯业务逻辑，不依赖传输层
type UserService interface {
	AuthenticateWithOAuth(provider string, userInfo *OAuthUserInfo) (*User2Token, error)
	RefreshUserToken(payload JwtPayload, refreshToken string) (*User2Token, error)

	GetUser(userID string) (*User, error)
	UpdateUserProfile(userID string, updates *UserProfileUpdate) (*User, error)

	// 团队管理（为微服务做准备）
	CreateTeam(ownerID string, teamInfo *TeamCreateRequest) (*Team, error)
	GetUserTeams(userID string) ([]*Team, error)
	JoinTeam(userID, teamID string) error
}

// 令牌服务接口
type TokenService interface {
	GenerateAccessToken(payload JwtPayload) (string, error)
	ValidateAccessToken(token string) (payload JwtPayload, isExpire bool, err error)
	RefreshAccessToken(domain JwtPayload, refreshToken string) (string, error)

	GenerateRefreshToken(payload JwtPayload) (string, error)
	ResetRefreshTokenExpiry(domain JwtPayload) error
}
