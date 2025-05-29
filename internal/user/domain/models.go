package domain

import "time"

// 业务实体（Domain Entity）
type User struct {
	ID            string     `json:"id"`
	Email         string     `json:"email"`
	PasswordHash  string     `json:"-"`
	Name          string     `json:"name"`
	Username      string     `json:"username,omitempty"`
	AvatarURL     string     `json:"avatar_url,omitempty"`
	EmailVerified bool       `json:"email_verified"`
	GithubID      string     `json:"github_id,omitempty"`
	GoogleID      string     `json:"google_id,omitempty"`
	GitlabID      string     `json:"gitlab_id,omitempty"`
	CreatedAt     time.Time  `json:"created_at"`
	UpdatedAt     time.Time  `json:"updated_at"`
	LastLoginAt   *time.Time `json:"last_login_at,omitempty"`
	Status        string     `json:"status"`
}

type JwtPayload struct {
	UserID     string `json:"user_id"`
	RandomCode string `json:"random_code"`
}

type User2Token struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
}

// OAuth 用户信息（值对象）
type OAuthUserInfo struct {
	Provider string `json:"provider"`
	ID       string `json:"id"`
	Login    string `json:"login"`
	Name     string `json:"name"`
	Email    string `json:"email"`
	Avatar   string `json:"avatar_url"`
}

// 用户资料更新（值对象）
type UserProfileUpdate struct {
	Name     *string `json:"name,omitempty"`
	Username *string `json:"username,omitempty"`
	Avatar   *string `json:"avatar,omitempty"`
}

// 团队（为后续微服务做准备）
type Team struct {
	ID          string    `json:"id"`
	OwnerID     string    `json:"owner_id"`
	Name        string    `json:"name"`
	Description string    `json:"description,omitempty"`
	Status      string    `json:"status"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

type TeamMember struct {
	TeamID   string    `json:"team_id"`
	UserID   string    `json:"user_id"`
	Role     string    `json:"role"`
	JoinedAt time.Time `json:"joined_at"`
	Status   string    `json:"status"`
}

// 团队创建请求（值对象）
type TeamCreateRequest struct {
	Name        string `json:"name" binding:"required"`
	Description string `json:"description,omitempty"`
}
