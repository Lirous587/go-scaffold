package adapters

import (
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
