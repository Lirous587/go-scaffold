package handler

import (
	"scaffold/internal/user/domain"
)

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
