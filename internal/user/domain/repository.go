package domain

type UserRepository interface {
	// 基础 CRUD
	FindByID(userID string) (*User, error)
	FindByEmail(email string) (*User, error)
	Create(user *User) (*User, error)
	Update(user *User) (*User, error)

	// OAuth 相关
	FindByOAuthID(provider, oauthID string) (*User, error)
	UpdateLastLogin(userID string) error

	// 辅助方法
	EmailExists(email string) (bool, error)
	UsernameExists(username string) (bool, error)
	GenerateUniqueUsername(preferred string) (string, error)
}

type TeamRepository interface {
	// 团队 CRUD
	CreateTeam(team *Team) (*Team, error)
	FindTeamByID(teamID string) (*Team, error)
	FindTeamsByOwner(ownerID string) ([]*Team, error)
	UpdateTeam(team *Team) (*Team, error)

	// 团队成员管理
	AddTeamMember(member *TeamMember) error
	RemoveTeamMember(teamID, userID string) error
	FindTeamMembers(teamID string) ([]*TeamMember, error)
	FindUserTeams(userID string) ([]*Team, error)
}

type TokenCache interface {
	GenRefreshToken(domain JwtPayload) (string, error)
	ValidateRefreshToken(domain JwtPayload, refreshToken string) error
	ResetRefreshTokenExpiry(domain JwtPayload) error
}
