package domain

// 纯业务逻辑，不依赖传输层
type UserService interface {
	// 认证相关
	AuthenticateWithOAuth(provider string, userInfo *OAuthUserInfo) (*UserSession, error)
	RefreshUserSession(userID string, refreshToken string) (*UserSession, error)

	// 用户管理
	GetUser(userID string) (*User, error)
	UpdateUserProfile(userID string, updates *UserProfileUpdate) (*User, error)

	// 团队管理（为微服务做准备）
	CreateTeam(ownerID string, teamInfo *TeamCreateRequest) (*Team, error)
	GetUserTeams(userID string) ([]*Team, error)
	JoinTeam(userID, teamID string) error
}

// 令牌服务接口
type TokenService interface {
	GenerateAccessToken(payload *JwtPayload) (string, error)
	ValidateAccessToken(token string) (*JwtPayload, error)
	GenerateRefreshToken(userID string) (string, error)
	ValidateRefreshToken(userID, token string) error
	RefreshAccessToken(userID, refreshToken string) (string, error)
}
