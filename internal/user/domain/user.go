package domain

import "time"

type User struct {
	ID           int64
	Email        string
	PasswordHash string
	Nickname     string
	GithubID     string
	Avatar       string
	CreatedAt    time.Time
	UpdatedAt    time.Time
	LastLoginAt  time.Time
}

type JwtPayload struct {
	UserID int64 `json:"user_id"`
}

type User2Token struct {
	AccessToken  string
	RefreshToken string
}

type OAuthUserInfo struct {
	Provider string
	ID       string
	Login    string
	Nickname string
	Email    string
	Avatar   string
}
