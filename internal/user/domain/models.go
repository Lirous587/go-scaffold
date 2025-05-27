package domain

import "time"

// 业务实体（Domain Entity）
type User struct {
    ID            string    `json:"id"`
    Email         string    `json:"email"`
    PasswordHash  string    `json:"-"` // 不在JSON中暴露
    Name          string    `json:"name"`
    Username      string    `json:"username,omitempty"`
    AvatarURL     string    `json:"avatar_url,omitempty"`
    EmailVerified bool      `json:"email_verified"`
    GithubID      string    `json:"github_id,omitempty"`
    GoogleID      string    `json:"google_id,omitempty"`
    GitlabID      string    `json:"gitlab_id,omitempty"`
    CreatedAt     time.Time `json:"created_at"`
    UpdatedAt     time.Time `json:"updated_at"`
    LastLoginAt   *time.Time `json:"last_login_at,omitempty"`
    Status        string    `json:"status"`
}

// 业务方法
func (u *User) IsActive() bool {
    return u.Status == "active"
}

func (u *User) HasPassword() bool {
    return u.PasswordHash != ""
}

func (u *User) IsOAuthUser() bool {
    return u.GithubID != "" || u.GoogleID != "" || u.GitlabID != ""
}

func (u *User) GetOAuthProvider() string {
    if u.GithubID != "" {
        return "github"
    }
    if u.GoogleID != "" {
        return "google"
    }
    if u.GitlabID != "" {
        return "gitlab"
    }
    return ""
}

// JWT 载荷（业务概念）
type JwtPayload struct {
    UserID    string    `json:"user_id"`
    LoginType string    `json:"login_type"`
    IssuedAt  time.Time `json:"issued_at"`
}

// 用户会话（业务概念）
type UserSession struct {
    User         *User     `json:"user"`
    AccessToken  string    `json:"access_token"`
    RefreshToken string    `json:"refresh_token"`
    ExpiresAt    time.Time `json:"expires_at"`
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