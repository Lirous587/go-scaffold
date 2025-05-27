package adapters

import (
	"time"

	"github.com/volatiletech/null/v8"
	"scaffold/internal/common/orm"
	"scaffold/internal/user/domain"
)

// Domain User <-> ORM User 转换
func DomainUserToORM(user *domain.User) *orm.User {
	if user == nil {
		return nil
	}

	ormUser := &orm.User{
		UserID:        user.ID,
		Email:         user.Email,
		Name:          user.Name,
		EmailVerified: user.EmailVerified,
		Status:        user.Status,
		CreatedAt:     user.CreatedAt,
		UpdatedAt:     user.UpdatedAt,
	}

	if user.PasswordHash != "" {
		ormUser.PasswordHash = null.StringFrom(user.PasswordHash)
	}

	if user.Username != "" {
		ormUser.Username = null.StringFrom(user.Username)
	}

	if user.AvatarURL != "" {
		ormUser.AvatarURL = null.StringFrom(user.AvatarURL)
	}

	if user.GithubID != "" {
		ormUser.GithubID = null.StringFrom(user.GithubID)
	}

	if user.GoogleID != "" {
		ormUser.GoogleID = null.StringFrom(user.GoogleID)
	}

	if user.GitlabID != "" {
		ormUser.GitlabID = null.StringFrom(user.GitlabID)
	}

	if user.LastLoginAt != nil {
		ormUser.LastLoginAt = null.TimeFrom(*user.LastLoginAt)
	}

	return ormUser
}

func ORMUserToDomain(ormUser *orm.User) *domain.User {
	if ormUser == nil {
		return nil
	}

	user := &domain.User{
		ID:            ormUser.UserID,
		Email:         ormUser.Email,
		Name:          ormUser.Name,
		EmailVerified: ormUser.EmailVerified,
		Status:        ormUser.Status,
		CreatedAt:     ormUser.CreatedAt,
		UpdatedAt:     ormUser.UpdatedAt,
	}

	if ormUser.PasswordHash.Valid {
		user.PasswordHash = ormUser.PasswordHash.String
	}

	if ormUser.Username.Valid {
		user.Username = ormUser.Username.String
	}

	if ormUser.AvatarURL.Valid {
		user.AvatarURL = ormUser.AvatarURL.String
	}

	if ormUser.GithubID.Valid {
		user.GithubID = ormUser.GithubID.String
	}

	if ormUser.GoogleID.Valid {
		user.GoogleID = ormUser.GoogleID.String
	}

	if ormUser.GitlabID.Valid {
		user.GitlabID = ormUser.GitlabID.String
	}

	if ormUser.LastLoginAt.Valid {
		user.LastLoginAt = &ormUser.LastLoginAt.Time
	}

	return user
}

// Domain User <-> HTTP Response 转换
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

// HTTP Request -> Domain 转换
type AuthRequest struct {
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

func HTTPUserUpdateToDomain(req *UserProfileUpdateRequest) *domain.UserProfileUpdate {
	return &domain.UserProfileUpdate{
		Name:     req.Name,
		Username: req.Username,
		Avatar:   req.Avatar,
	}
}

// HTTP Response 模型
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

func DomainSessionToAuthResponse(session *domain.UserSession) *AuthResponse {
	return &AuthResponse{
		User:         DomainUserToResponse(session.User),
		AccessToken:  session.AccessToken,
		RefreshToken: session.RefreshToken,
		ExpiresAt:    session.ExpiresAt,
	}
}
