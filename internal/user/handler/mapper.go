package handler

import (
	"scaffold/internal/user/domain"
)

type GithubAuthRequest struct {
	Code string `json:"code" binding:"required"`
}

type UserResponse struct {
	ID            int64  `json:"id"`
	Email         string `json:"email"`
	NickName      string `json:"username"`
	Avatar        string `json:"avatar_url,omitempty"`
	EmailVerified bool   `json:"email_verified"`
	CreatedAt     int64  `json:"created_at"`
	UpdatedAt     int64  `json:"updated_at"`
	LastLoginAt   int64  `json:"last_login_at"`
}

type AuthResponse struct {
	User         *UserResponse `json:"user"`
	AccessToken  string        `json:"access_token"`
	RefreshToken string        `json:"refresh_token"`
}

type RefreshTokenResponse struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
}

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

func domainUserToResponse(user *domain.User) *UserResponse {
	if user == nil {
		return nil
	}

	return &UserResponse{
		ID:          user.ID,
		Email:       user.Email,
		NickName:    user.Nickname,
		Avatar:      user.Avatar,
		CreatedAt:   user.CreatedAt.Unix(),
		UpdatedAt:   user.UpdatedAt.Unix(),
		LastLoginAt: user.LastLoginAt.Unix(),
	}
}

func domain2TokenToAuthResponse(token2 *domain.User2Token) *AuthResponse {
	return &AuthResponse{
		AccessToken:  token2.AccessToken,
		RefreshToken: token2.RefreshToken,
	}
}

func domainSessionToRefreshResponse(token2 *domain.User2Token) *RefreshTokenResponse {
	return &RefreshTokenResponse{
		AccessToken:  token2.AccessToken,
		RefreshToken: token2.RefreshToken,
	}
}
