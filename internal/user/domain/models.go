package domain

import "time"

type User struct {
	ID           int64      `json:"id"`
	Email        string     `json:"email"`
	PasswordHash string     `json:"-"`
	Name         string     `json:"name"`
	AvatarURL    string     `json:"avatar_url,omitempty"`
	GithubID     string     `json:"github_id,omitempty"`
	LastLoginAt  *time.Time `json:"last_login_at,omitempty"`
}

type JwtPayload struct {
	UserID int64 `json:"user_id,string"`
}

type User2Token struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
}

type OAuthUserInfo struct {
	Provider string `json:"provider"`
	ID       string `json:"id"`
	Login    string `json:"login"`
	Name     string `json:"name"`
	Email    string `json:"email"`
	Avatar   string `json:"avatar_url"`
}
