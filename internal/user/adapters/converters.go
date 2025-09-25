package adapters

import (
	"scaffold/internal/common/orm"
	"scaffold/internal/user/domain"
	"github.com/volatiletech/null/v8"
)

func DomainUserToORM(user *domain.User) *orm.User {
	if user == nil {
		return nil
	}

	ormUser := &orm.User{
		ID:	user.ID,
		Email:	user.Email,
		Name:	user.Name,
	}

	if user.PasswordHash != "" {
		ormUser.PasswordHash = null.StringFrom(user.PasswordHash)
	}

	if user.AvatarURL != "" {
		ormUser.PasswordHash = null.StringFrom(user.AvatarURL)
	}

	if user.GithubID != "" {
		ormUser.GithubID = null.StringFrom(user.GithubID)
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
		ID:	ormUser.ID,
		Email:	ormUser.Email,
		Name:	ormUser.Name,
	}

	if ormUser.PasswordHash.Valid {
		user.PasswordHash = ormUser.PasswordHash.String
	}

	if ormUser.GithubID.Valid {
		user.GithubID = ormUser.GithubID.String
	}

	if ormUser.LastLoginAt.Valid {
		user.LastLoginAt = &ormUser.LastLoginAt.Time
	}

	return user
}
