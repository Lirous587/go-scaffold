package handler

import (
	"time"

	"scaffold/internal/user/domain"
)

// HTTP 请求/响应模型
type GithubAuthRequest struct {
	Code string `json:"code" binding:"required"`
}

type RefreshTokenRequest struct {
	RefreshToken string `json:"refresh_token" binding:"required"`
}

type UserProfileUpdateRequest struct {
	Name     *string `json:"name,omitempty"`
	Username *string `json:"username,omitempty"`
	Avatar   *string `json:"avatar,omitempty"`
}

type UserResponse struct {
	ID            string     `json:"id"`
	Email         string     `json:"email"`
	Name          string     `json:"name"`
	Username      string     `json:"username,omitempty"`
	AvatarURL     string     `json:"avatar_url,omitempty"`
	EmailVerified bool       `json:"email_verified"`
	Status        string     `json:"status"`
	CreatedAt     time.Time  `json:"created_at"`
	UpdatedAt     time.Time  `json:"updated_at"`
	LastLoginAt   *time.Time `json:"last_login_at,omitempty"`
}

type AuthResponse struct {
	User         *UserResponse `json:"user"`
	AccessToken  string        `json:"access_token"`
	RefreshToken string        `json:"refresh_token"`
	ExpiresAt    time.Time     `json:"expires_at"`
}

type RefreshTokenResponse struct {
	AccessToken  string    `json:"access_token"`
	RefreshToken string    `json:"refresh_token"`
	ExpiresAt    time.Time `json:"expires_at"`
}

// GitHub API 响应模型
type GithubUser struct {
	ID        int64  `json:"id"`
	Login     string `json:"login"`
	Name      string `json:"name"`
	Email     string `json:"email"`
	AvatarURL string `json:"avatar_url"`
}

type GithubAccessTokenResponse struct {
	AccessToken string `json:"access_token"`
	TokenType   string `json:"token_type"`
	Scope       string `json:"scope"`
}

// 转换函数
func DomainUserToResponse(user *domain.User) *UserResponse {
	if user == nil {
		return nil
	}

	return &UserResponse{
		ID:            user.ID,
		Email:         user.Email,
		Name:          user.Name,
		Username:      user.Username,
		AvatarURL:     user.AvatarURL,
		EmailVerified: user.EmailVerified,
		Status:        user.Status,
		CreatedAt:     user.CreatedAt,
		UpdatedAt:     user.UpdatedAt,
		LastLoginAt:   user.LastLoginAt,
	}
}

func Domain2TokenToAuthResponse(token2 *domain.User2Token) *AuthResponse {
	return &AuthResponse{
		AccessToken:  token2.AccessToken,
		RefreshToken: token2.RefreshToken,
	}
}

func DomainSessionToRefreshResponse(token2 *domain.User2Token) *RefreshTokenResponse {
	return &RefreshTokenResponse{
		AccessToken:  token2.AccessToken,
		RefreshToken: token2.RefreshToken,
	}
}

func HTTPUserUpdateToDomain(req *UserProfileUpdateRequest) *domain.UserProfileUpdate {
	return &domain.UserProfileUpdate{
		Name:     req.Name,
		Username: req.Username,
		Avatar:   req.Avatar,
	}
}
